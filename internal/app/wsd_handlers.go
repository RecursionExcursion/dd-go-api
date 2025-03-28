package app

import (
	"io"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"github.com/recursionexcursion/dd-go-api/internal/wsd"
)

var postWsdBuildHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.Response.ServerError(w, "Failed to read body")
		return
	}

	params, err := lib.Map[wsd.CreateExeParams](bodyBytes)
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

	binBytes, name, err := wsd.CreateGoExe(params)
	if err != nil {
		panic(err)
	}

	api.Response.StreamBytes(w, 200, binBytes, name)
}

var getWsdTestHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	testParams := wsd.CreateExeParams{
		Arch: "win",
		Commands: []string{
			"url:www.facebook.com",
			"url:www.chatgpt.com",
			"cmd:code C:/Users/rloup/dev/workspaces/vsc/xpres",
		},
	}

	binBytes, name, err := wsd.CreateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	api.Response.StreamBytes(w, 200, binBytes, name)
}

var getSupportedOsHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	keys := []string{}
	for k := range wsd.SupportedArchitecture {
		keys = append(keys, k)
	}

	api.Response.Ok(w, keys)
}
