package wsd

import (
	"io"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"github.com/recursionexcursion/dd-go-api/internal/wsd/core"
)

func WsdRoutes(mwChain []api.Middleware) []api.RouteHandler {
	var postWsdHome = api.RouteHandler{
		MethodAndPath: "POST /wsd/build",
		Handler:       postWsdBuildHandler,
		Middleware:    mwChain,
	}

	getSupportedOs := api.RouteHandler{
		MethodAndPath: "GET /wsd/os",
		Handler:       getSupportedOsHandler,
		Middleware:    mwChain,
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
		Middleware: mwChain,
	}

	return []api.RouteHandler{
		postWsdHome,
		getSupportedOs,
		routes,
	}
}

var postWsdBuildHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.Response.ServerError(w, "Failed to read body")
		return
	}

	params, err := lib.Map[core.CreateExeParams](bodyBytes)
	if err != nil {
		api.Response.ServerError(w, "Failed to map body")
		return
	}

	// Validate body
	if params.Arch == "" {
		api.Response.BadRequest(w, "No arch provided")
		return
	}

	if len(params.Commands) == 0 {
		api.Response.BadRequest(w, "No commands provided")
		return
	}

	binPath, name, err := core.CreateGoExe(params)
	if err != nil {
		panic(err)
	}

	api.Response.StreamFile(w, 200, binPath, name)
}

var getWsdTestHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	testParams := core.CreateExeParams{
		Arch: "win",
		Commands: []string{
			"url:www.facebook.com",
			"url:www.chatgpt.com",
			"cmd:code C:/Users/rloup/dev/workspaces/vsc/xpres",
		},
	}

	bin, name, err := core.CreateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	api.Response.StreamFile(w, 200, bin, name)
}

var getSupportedOsHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	keys := []string{}
	for k := range core.SupportedArchitecture {
		keys = append(keys, k)
	}

	api.Response.Ok(w, keys)
}
