package app

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func routes() []api.RouteHandler {

	/* base */

	var getbaseRoute = api.RouteHandler{
		MethodAndPath: "GET /",
		Handler: api.HandlerFn(func(w http.ResponseWriter, r *http.Request) {
			api.Response.Ok(w, "API Status: Healthy")
		}),
		Middleware: mwChainMap()["global"],
	}

	/* Test */

	var testRoute = api.RouteHandler{
		MethodAndPath: "GET /test",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			//collect data
			fsd, err := betbot.CollectData()
			if err != nil {
				api.Response.ServerError(w)
				return
			}

			log.Println("In handlers")
			betbot.FindGameInFsd(fsd, strconv.Itoa(401705613))

			// Compile stats
			packagedData, err := betbot.NewStatCalculator(fsd).CalcAndPackage()
			if err != nil {
				lib.LogError("", err)
				api.Response.ServerError(w)
				return
			}

			// Zip and return
			lib.Log("Gzipping payload", 5)
			api.Response.Gzip(w, 200,
				struct {
					Meta string
					Data []betbot.PackagedPlayer
				}{
					Meta: strconv.FormatInt(time.Now().UnixMilli(), 10),
					Data: packagedData,
				},
			)
		},
	}

	var routes = []api.RouteHandler{getbaseRoute, testRoute}
	routes = append(routes, betbotRoutes()...)
	routes = append(routes, wsdRoutes()...)

	return routes
}

func betbotRoutes() []api.RouteHandler {

	chains := mwChainMap()
	jwtChain := chains["bb-jwt-chain"]
	apiKeyChain := chains["bb-key-chain"]

	var getBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot",
		Handler:       HandleBBGet,
		Middleware:    jwtChain,
	}

	var revalidateBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot/revalidate",
		Handler:       HandleGetBBRevalidation,
		Middleware:    jwtChain,
	}

	var bbRevalidateAndZip = api.RouteHandler{
		MethodAndPath: "GET /betbot/zip",
		Handler:       HandleBBValidateAndZip,
		Middleware:    jwtChain,
	}

	var loginBBUserRoute = api.RouteHandler{
		MethodAndPath: "POST /betbot/user/login",
		Handler:       HandleUserLogin,
		Middleware:    apiKeyChain,
	}

	return []api.RouteHandler{getBetBotRoute, revalidateBetBotRoute, loginBBUserRoute, bbRevalidateAndZip}
}

func wsdRoutes() []api.RouteHandler {

	mwchains := mwChainMap()
	keyChain := mwchains["wsd-key-chain"]

	var postWsdHome = api.RouteHandler{
		MethodAndPath: "POST /wsd/build",
		Handler:       postWsdBuildHandler,
		Middleware:    keyChain,
	}

	getWsdTest := api.RouteHandler{
		MethodAndPath: "GET /wsd/test",
		Handler:       getWsdTestHandler,
		Middleware:    mwchains["global"],
	}

	getSupportedOs := api.RouteHandler{
		MethodAndPath: "GET /wsd/os",
		Handler:       getSupportedOsHandler,
		Middleware:    keyChain,
	}

	warmup := api.RouteHandler{
		MethodAndPath: "GET /wsd/warmup",
		Handler:       getPipelineWarmUpHandler,
		Middleware:    keyChain,
	}

	status := api.RouteHandler{
		MethodAndPath: "GET /wsd/status",
		Handler:       getStatusHandler,
		Middleware:    keyChain,
	}

	routes := api.RouteHandler{
		MethodAndPath: "GET /wsd/routes",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			routeMap := map[string]string{
				"getOs":     "/wsd/os",
				"postBuild": "/wsd/build",
				"getWarmUp": "/wsd/warmup",
				"getStatus": "/wsd/status",
			}

			api.Response.Ok(w, routeMap)
		},
		Middleware: keyChain,
	}

	return []api.RouteHandler{postWsdHome, getWsdTest, getSupportedOs, warmup, routes, status}
}

var mwChainMap = func() func() map[string][]api.Middleware {
	geoParams := GeoLimitParams{
		// WhitelistCountryCodes: strings.Split(lib.EnvGet("CC_WHITELIST"), ","),
	}

	globalMWChain := []api.Middleware{
		LoggerMW,
		GeoLimitMW(geoParams),
		RateLimitMW,
	}

	mwChainMap := map[string][]api.Middleware{
		"global":        globalMWChain,
		"bb-jwt-chain":  append(globalMWChain, JWTAuthMW(lib.EnvGet("BB_JWT_SECRET"))),
		"bb-key-chain":  append(globalMWChain, KeyAuthMW(lib.EnvGet("BB_API_KEY"))),
		"wsd-key-chain": append(globalMWChain, KeyAuthMW(lib.EnvGet("WSD_API_KEY"))),
	}

	return func() map[string][]api.Middleware {
		return mwChainMap
	}
}()
