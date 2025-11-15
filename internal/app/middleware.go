package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"slices"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/jwt"
	"github.com/RecursionExcursion/gouse/gouse"
	"golang.org/x/time/rate"
)

type handler = func(w http.ResponseWriter, r *http.Request)
type middleware = func(handler) handler

func pipe(mws ...middleware) middleware {
	return func(hndlr handler) handler {
		for i := len(mws) - 1; i >= 0; i-- {
			hndlr = mws[i](hndlr)
		}
		return hndlr
	}
}

func LoggerMW(logger *log.Logger) middleware {
	return func(next handler) handler {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			logger.Printf("%v %v accessed %v in %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		}
	}
}

func KeyAuthMW(key string) middleware {
	return func(next handler) handler {
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

func JWTAuthMW(key string) middleware {
	return func(next handler) handler {
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

// refil rate /sec, total bucket size
func RateLimitMW(refillRate int, size int) middleware {
	return func(next handler) handler {
		var limiter = rate.NewLimiter(rate.Limit(refillRate), 10)
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}

type GeoLimitParams struct {
	WhitelistCountryCodes []string
	BlacklistCountryCodes []string

	WhitelistContinentCodes []string
	BlacklistContinentCodes []string

	BlacklistZipCodes []string
}

func GeoLimitMW(params GeoLimitParams) middleware {
	return func(next handler) handler {
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

		isBlackListed := func(p GeoLimitParams, d GeoLimitData) bool {
			return slices.Contains(p.BlacklistCountryCodes, d.CountryCode) ||
				slices.Contains(p.BlacklistContinentCodes, d.ContinentCode)
		}

		isWhitelisted := func(p GeoLimitParams, d GeoLimitData) bool {
			if params.WhitelistContinentCodes != nil {
				if !slices.Contains(p.WhitelistContinentCodes, d.ContinentCode) {
					return false
				}
			}

			if params.WhitelistCountryCodes != nil {
				if !slices.Contains(p.WhitelistCountryCodes, d.CountryCode) {
					return false
				}
			}

			return true
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

			if isBlackListed(params, data) || !isWhitelisted(params, data) {
				gouse.Response.Forbidden(w, "")
				return
			}

			next(w, r)
		}
	}
}
