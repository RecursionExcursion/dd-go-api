package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"slices"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/jwt"
	"github.com/RecursionExcursion/gouse/gouse"
	"golang.org/x/time/rate"
)

func KeyAuthMW(key string) gouse.Middleware {
	return func(next gouse.HandlerFn) gouse.HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(("Authorization"))
			parts := strings.SplitN(token, " ", 2)

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] != key {
				gouse.Response.Unauthorized(w, "Invalid token")
				return
			}

			next(w, r)
		}
	}
}

func JWTAuthMW(key string) gouse.Middleware {
	return func(next gouse.HandlerFn) gouse.HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(("Authorization"))
			parts := strings.SplitN(token, " ", 2)

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				gouse.Response.Unauthorized(w, "Malformed token")
				return
			}
			if isValid, _, _ := jwt.ParseJWT(parts[1], key); !isValid {
				gouse.Response.Unauthorized(w, "Invalid token")
				return
			}

			next(w, r)
		}
	}
}

func RateLimitMW(next gouse.HandlerFn) gouse.HandlerFn {
	// refil rate 5/sec, total bucket size is 10
	var limiter = rate.NewLimiter(5, 10)
	return func(w http.ResponseWriter, r *http.Request) {
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

	BlacklistZipCodes []string
}

func GeoLimitMW(params GeoLimitParams) gouse.Middleware {
	return func(next gouse.HandlerFn) gouse.HandlerFn {
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
			addr := r.RemoteAddr

			data, res, err := core.FetchAndMap[GeoLimitData](func() (resp *http.Response, err error) {
				return http.Get(fmt.Sprintf("http://ip-api.com/json/%v", addr))
			})
			if err != nil {
				if res.StatusCode == 429 {
					gouse.Response.ServerError(w, "too many requests, please try again later")
				} else {
					log.Println(addr, err, "GeoLimitMW", "FetchAndMap")
					gouse.Response.ServerError(w, "something went wrong, please try again later")
				}
				return
			}

			isBlackListed := func() bool {
				return slices.Contains(params.BlacklistCountryCodes, data.CountryCode) ||
					slices.Contains(params.BlacklistContinentCodes, data.ContinentCode)
			}

			isWhitelisted := func() bool {
				if params.WhitelistContinentCodes != nil {
					if !slices.Contains(params.WhitelistContinentCodes, data.ContinentCode) {
						return false
					}
				}

				if params.WhitelistCountryCodes != nil {
					if !slices.Contains(params.WhitelistCountryCodes, data.CountryCode) {
						return false
					}
				}

				return true
			}

			if isBlackListed() || !isWhitelisted() {
				gouse.Response.Forbidden(w, "")
				return
			}

			next(w, r)
		}
	}
}
