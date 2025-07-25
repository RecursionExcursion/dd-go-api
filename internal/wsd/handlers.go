package wsd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/gogen/gogen"
	"github.com/RecursionExcursion/gouse/gouse"
)

func WsdRoutes(mwChain []gouse.Middleware) []gouse.RouteHandler {
	var postWsdHome = gouse.RouteHandler{
		MethodAndPath: "POST /wsd/build",
		Handler:       postWsdBuildHandler,
		Middleware:    mwChain,
	}

	getSupportedOs := gouse.RouteHandler{
		MethodAndPath: "GET /wsd/os",
		Handler:       getSupportedOsHandler,
		Middleware:    mwChain,
	}

	routes := gouse.RouteHandler{
		MethodAndPath: "GET /wsd/routes",
		Handler: func(w http.ResponseWriter, r *http.Request) {

			routeMap := map[string]string{
				"getOs":     "/wsd/os",
				"postBuild": "/wsd/build",
			}

			gouse.Response.Ok(w, routeMap)
		},
		Middleware: mwChain,
	}

	return []gouse.RouteHandler{
		postWsdHome,
		getSupportedOs,
		routes,
	}
}

var postWsdBuildHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		gouse.Response.ServerError(w, "Failed to read body")
		return
	}

	params, err := core.Map[gogen.CreateExeParams](bodyBytes)
	if err != nil {
		gouse.Response.ServerError(w, "Failed to map body")
		return
	}

	// Validate body
	if params.Arch == "" {
		gouse.Response.BadRequest(w, "No arch provided")
		return
	}

	if len(params.Commands) == 0 {
		gouse.Response.BadRequest(w, "No commands provided")
		return
	}

	ret, err := gogen.GenerateGoExe(params)
	// binPath, name, err := core.CreateGoExe(params)
	if err != nil {
		fmt.Println(err)
		gouse.Response.ServerError(w)
		return
	}

	gouse.Response.StreamFile(w, 200, ret.BinPath, ret.ExeName)
	ret.Cleanup()
}

var getWsdTestHandler gouse.HandlerFn = func(w http.ResponseWriter, _ *http.Request) {

	testParams := gogen.CreateExeParams{
		Arch: "win",
		Commands: []string{
			"url:www.facebook.com",
			"url:www.chatgpt.com",
			"cmd:code C:/Users/rloup/dev/workspaces/vsc/xpres",
		},
	}

	ret, err := gogen.GenerateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	gouse.Response.StreamFile(w, 200, ret.BinPath, ret.ExeName)
	ret.Cleanup()
}

var getSupportedOsHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	keys := []string{}
	for k := range gogen.SupportedArchitecture {
		keys = append(keys, k)
	}

	gouse.Response.Ok(w, keys)
}
