package app

import (
	"fmt"
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

			betbot.FindGameInFsd(fsd, strconv.Itoa(401705613))

			// Compile stats
			packagedData, err := betbot.NewStatCalculator(fsd).CalculateAndPackage()
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

	bbBase := pathGenerator("/betbot")

	chains := mwChainMap()
	jwtChain := chains["bb-jwt-chain"]
	apiKeyChain := chains["bb-key-chain"]

	var getBetBotRoute = api.RouteHandler{
		MethodAndPath: bbBase().GET,
		Handler:       HandleBBGet,
		Middleware:    jwtChain,
	}

	var revalidateBetBotRoute = api.RouteHandler{
		MethodAndPath: bbBase().POST,
		Handler:       HandleGetBBRevalidation,
		Middleware:    jwtChain,
	}
	var pollBetBotRoute = api.RouteHandler{
		MethodAndPath: bbBase("poll").GET,
		Handler:       handleRevalidationPolling,
		Middleware:    jwtChain,
	}

	var bbRevalidateAndZip = api.RouteHandler{
		MethodAndPath: bbBase("zip").GET,
		Handler:       HandleBBValidateAndZip,
		Middleware:    jwtChain,
	}

	var loginBBUserRoute = api.RouteHandler{
		MethodAndPath: bbBase("user", "login").POST,
		Handler:       HandleUserLogin,
		Middleware:    apiKeyChain,
	}

	var bbPing = api.RouteHandler{
		MethodAndPath: bbBase("ping").GET,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			api.Response.Ok(w)
		},
		Middleware: apiKeyChain,
	}

	return []api.RouteHandler{
		getBetBotRoute,
		revalidateBetBotRoute,
		loginBBUserRoute,
		bbRevalidateAndZip,
		bbPing,
		pollBetBotRoute,
	}
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

	routes := api.RouteHandler{
		MethodAndPath: "GET /wsd/routes",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			routeMap := map[string]string{
				"getOs":     "/wsd/os",
				"postBuild": "/wsd/build",
			}

			api.Response.Ok(w, routeMap)
		},
		Middleware: keyChain,
	}

	return []api.RouteHandler{
		postWsdHome,
		getWsdTest,
		getSupportedOs,
		routes,
	}
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

type HTTPMethods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

func pathGenerator(base string) func(path ...string) HTTPMethods {
	return func(paths ...string) HTTPMethods {

		pathStr := ""

		numArgs := len(paths)

		if numArgs != 0 {
			if numArgs == 1 {
				pathStr = "/" + paths[0]
			} else {
				for _, arg := range paths {
					pathStr += "/" + arg
				}
			}
		}

		route := base + pathStr

		assign := func(s string) string {
			return fmt.Sprintf("%v %v", s, route)
		}

		routes := HTTPMethods{
			GET:    assign("GET"),
			POST:   assign("POST"),
			PUT:    assign("PUT"),
			PATCH:  assign("PATCH"),
			DELETE: assign("DELETE"),
		}

		return routes
	}
}
