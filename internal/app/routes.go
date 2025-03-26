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
			api.Response.Ok(w, "dd-api")
		}),
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

			packagedData, err := betbot.NewStatCalculator(fsd).CalcAndPackage()
			if err != nil {
				lib.LogError("", err)
				api.Response.ServerError(w)
				return
			}

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

	return routes
}

func betbotRoutes() []api.RouteHandler {
	geoParams := GeoLimitParams{
		// WhitelistCountryCodes: strings.Split(lib.EnvGet("CC_WHITELIST"), ","),
	}

	bbMwChain := []api.Middleware{
		LoggerMW,
		GeoLimitMW(geoParams),
		RateLimitMW,
		KeyAuthMW(lib.EnvGet("BB_API_KEY")),
	}

	var getBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot",
		Handler:       HandleGetBetBot,
		Middleware:    bbMwChain,
	}

	var revalidateBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot/revalidate",
		Handler:       HandleBetBotRevalidation,
		Middleware:    bbMwChain,
	}

	var loginBBUserRoute = api.RouteHandler{
		MethodAndPath: "POST /user",
		Handler:       HandleUserLogin,
		Middleware:    bbMwChain,
	}

	return []api.RouteHandler{getBetBotRoute, revalidateBetBotRoute, loginBBUserRoute}
}
