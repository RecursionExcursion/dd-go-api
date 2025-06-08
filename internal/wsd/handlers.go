package wsd

import (
	"io"
	"net/http"

	"github.com/RecursionExcursion/api-go/api"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/wsd-core/pkg"
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

	params, err := core.Map[pkg.CreateExeParams](bodyBytes)
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

	binPath, name, err := pkg.CreateGoExe(params)
	// binPath, name, err := core.CreateGoExe(params)
	if err != nil {
		panic(err)
	}

	api.Response.StreamFile(w, 200, binPath, name)
}

var getWsdTestHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	testParams := pkg.CreateExeParams{
		Arch: "win",
		Commands: []string{
			"url:www.facebook.com",
			"url:www.chatgpt.com",
			"cmd:code C:/Users/rloup/dev/workspaces/vsc/xpres",
		},
	}

	bin, name, err := pkg.CreateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	api.Response.StreamFile(w, 200, bin, name)
}

var getSupportedOsHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	keys := []string{}
	for k := range pkg.SupportedArchitecture {
		keys = append(keys, k)
	}

	api.Response.Ok(w, keys)
}
