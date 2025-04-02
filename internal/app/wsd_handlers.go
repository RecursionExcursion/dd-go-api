package app

import (
	"io"
	"log"
	"net/http"
	"sync/atomic"

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

	binPath, name, err := wsd.CreateGoExe(params)
	if err != nil {
		panic(err)
	}

	api.Response.StreamFile(w, 200, binPath, name)
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

	bin, name, err := wsd.CreateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	api.Response.StreamFile(w, 200, bin, name)
}

var getSupportedOsHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	keys := []string{}
	for k := range wsd.SupportedArchitecture {
		keys = append(keys, k)
	}

	api.Response.Ok(w, keys)
}

var isReady atomic.Bool

var getPipelineWarmUpHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	func() {

		log.Println("Warming up Go build pipeline")

		//TODO add darwin warmup too
		winParams := wsd.CreateExeParams{
			Name: "WIN_WARMUP",
			Arch: "win",
			Commands: []string{
				"url:www.facebook.com",
			},
		}
		darwinParams := wsd.CreateExeParams{
			Name: "DARWIN-WARMUP",
			Arch: "win",
			Commands: []string{
				"url:www.facebook.com",
			},
		}

		log.Println("Caching Win dist")
		_, _, err := wsd.CreateGoExe(winParams)
		if err != nil {
			log.Println("Win warmup failed:", err)
			return
		}

		log.Println("Caching Win dist")
		_, _, err = wsd.CreateGoExe(darwinParams)
		if err != nil {
			log.Println("Darwin warmup failed:", err)
			return
		}

		log.Println(`Pipeline warmup successful`)
		isReady.Store(true)
	}()
	api.Response.Ok(w)
}

var getStatusHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	api.Response.Ok(w, isReady.Load())
}
