package app

import (
	"net/http"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/gouse/gouse"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/hash"
	"github.com/recursionexcursion/dd-go-api/internal/pickle"
	"github.com/recursionexcursion/dd-go-api/internal/wsd"
)

// TODO rethink flow here, works but feels convoluted with repeated appends
func routes() []gouse.RouteHandler {

	var getbaseRoute = gouse.RouteHandler{
		MethodAndPath: "GET /",
		Handler: gouse.HandlerFn(func(w http.ResponseWriter, r *http.Request) {
			gouse.Response.Ok(w, "API Status: Healthy")
		}),
		Middleware: globalMWChain,
	}

	var routes = []gouse.RouteHandler{getbaseRoute}

	routes = append(routes, betbot.BetbotRoutes(struct {
		JwtChain []gouse.Middleware
		KeyChain []gouse.Middleware
	}{
		JwtChain: append(globalMWChain, JWTAuthMW(core.EnvGetOrPanic("BB_JWT_SECRET"))),
		KeyChain: append(globalMWChain, KeyAuthMW(core.EnvGetOrPanic("BB_API_KEY"))),
	})...)

	routes = append(routes, wsd.WsdRoutes(append(globalMWChain, KeyAuthMW(core.EnvGetOrPanic("WSD_API_KEY"))))...)

	routes = append(routes, cfbr.CfbrRoutes(globalMWChain)...)

	routes = append(routes, pickle.PickleRoutes(append(globalMWChain, JWTAuthMW(core.EnvGetOrPanic("PICKLE_SECRET"))))...)
	routes = append(routes, pickle.PickleLoginRoute(globalMWChain)...)

	routes = append(routes, hash.HashRoutes(globalMWChain)...)

	return routes
}

var globalMWChain = func() []gouse.Middleware {
	// geoParams := GeoLimitParams{
	// 	// WhitelistCountryCodes: strings.Split(lib.EnvGet("CC_WHITELIST"), ","),
	// }

	globalMWChain := []gouse.Middleware{
		LoggerMW,
		// GeoLimitMW(geoParams),
		RateLimitMW,
	}

	return globalMWChain
}()
