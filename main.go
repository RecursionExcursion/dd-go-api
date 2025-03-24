package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/app"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func main() {
	osInfo()

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

	var testRoute = api.RouteHandler{
		MethodAndPath: "GET /test",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("TEST")
			api.Response.OkPayload(w, "foo")
		},
	}

	return []api.RouteHandler{getBetBotRoute, revalidateBetBotRoute, testRoute}
}

func osInfo() {
	lib.Log("\nOS INFO:", -1)
	lib.Log(fmt.Sprintf("CPUs available: %d", runtime.NumCPU()), -1)
	lib.Log(fmt.Sprintf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)), -1)
}
