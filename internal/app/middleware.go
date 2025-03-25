package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"slices"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"golang.org/x/time/rate"
)

func LoggerMW(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {

		remote := r.RemoteAddr
		accessedPath := r.Host + r.RequestURI
		time := time.Now().Format("2006-01-02 15:04:05")

		logMsg := fmt.Sprintf("%v accessed %v at %v", remote, accessedPath, time)

		log.Println(logMsg)

		next(w, r)
	}
}

func KeyAuthMW(key string) api.Middleware {
	return func(next api.HandlerFn) api.HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(("Authorization"))
			parts := strings.SplitN(token, " ", 2)

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] != key {
				api.Response.Unauthorized(w, "Invalid token")
				return
			}

			next(w, r)
		}
	}
}

// refil rate 5/sec, total bucket size is 10
var limiter = rate.NewLimiter(5, 10)

func RateLimitMW(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("RL MW")
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

type GeoLimitParams struct {
	WhitelistCountryCodes []string
	BlacklistCountryCodes []string

	WhitelistContinentCodes []string
	BlacklistContinentCodes []string
}

func GeoLimitMW(params GeoLimitParams) api.Middleware {
	return func(next api.HandlerFn) api.HandlerFn {
		type GeoLimitData struct {
			Query         string `json:"query"`
			Status        string `json:"status"`
			Message       string `json:"message"`
			Region        string `json:"region"`
			CountryCode   string `json:"countryCode"`
			City          string `json:"city"`
			RegionName    string `json:"regionName"`
			Country       string `json:"country"`
			Zip           string `json:"zip"`
			Isp           string `json:"isp"`
			Continent     string `json:"continent"`
			ContinentCode string `json:"continentCode"`
		}

		return func(w http.ResponseWriter, r *http.Request) {
			addr := "24.39.0.1"
			// addr := r.RemoteAddr

			data, err := lib.FetchAndMap[GeoLimitData](func() (resp *http.Response, err error) {
				return http.Get(fmt.Sprintf("http://ip-api.com/json/%v", addr))
			})
			if err != nil {
				panic(err)
			}

			//Handle blacklists (Contains -> 403)
			if slices.Contains(params.BlacklistCountryCodes, data.CountryCode) ||
				slices.Contains(params.BlacklistContinentCodes, data.ContinentCode) {
				api.Response.Forbidden(w)
				return
			}

			if params.WhitelistContinentCodes != nil || params.WhitelistCountryCodes != nil {
				//Handle whitelists (Does not contain -> 403)
				if !slices.Contains(params.WhitelistCountryCodes, data.CountryCode) ||
					!slices.Contains(params.WhitelistContinentCodes, data.ContinentCode) {
					api.Response.Forbidden(w)
					return
				}
			}

			next(w, r)
		}
	}
}
