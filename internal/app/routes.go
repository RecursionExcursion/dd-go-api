package app

import (
	"net/http"

	"github.com/RecursionExcursion/api-go/api"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/wsd"
)

// TODO rethink flow here, works but feels convoluted with repeated appends
func routes() []api.RouteHandler {

	var getbaseRoute = api.RouteHandler{
		MethodAndPath: "GET /",
		Handler: api.HandlerFn(func(w http.ResponseWriter, r *http.Request) {
			api.Response.Ok(w, "API Status: Healthy")
		}),
		Middleware: globalMWChain,
	}

	var routes = []api.RouteHandler{getbaseRoute}

	routes = append(routes, betbot.BetbotRoutes(struct {
		JwtChain []api.Middleware
		KeyChain []api.Middleware
	}{
		JwtChain: append(globalMWChain, JWTAuthMW(core.EnvGetOrPanic("BB_JWT_SECRET"))),
		KeyChain: append(globalMWChain, KeyAuthMW(core.EnvGetOrPanic("BB_API_KEY"))),
	})...)

	routes = append(routes, wsd.WsdRoutes(append(globalMWChain, KeyAuthMW(core.EnvGetOrPanic("WSD_API_KEY"))))...)

	routes = append(routes, cfbr.CfbrRoutes(globalMWChain)...)

	return routes
}

var globalMWChain = func() []api.Middleware {
	geoParams := GeoLimitParams{
		// WhitelistCountryCodes: strings.Split(lib.EnvGet("CC_WHITELIST"), ","),
	}

	globalMWChain := []api.Middleware{
		LoggerMW,
		GeoLimitMW(geoParams),
		RateLimitMW,
	}

	return globalMWChain
}()
