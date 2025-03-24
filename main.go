package main

import (
	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/app"
)

func main() {
	s := api.NewApiServer(":8080")
	s.Init(createRoutes())
	s.ListenAndServe()
}

func createRoutes() []api.RouteHandler {
	mwChain := []api.Middleware{app.LoggerMW, app.RateLimitMW, app.AuthMW}

	var getBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot",
		Handler:       app.HandleGetBetBot,
		Middleware:    mwChain,
	}

	var revalidateBetBotRoute = api.RouteHandler{
		MethodAndPath: "GET /betbot/revalidate",
		Handler:       app.HandleBetBotRevalidation,
		Middleware:    mwChain,
	}

	return []api.RouteHandler{getBetBotRoute, revalidateBetBotRoute}
}
