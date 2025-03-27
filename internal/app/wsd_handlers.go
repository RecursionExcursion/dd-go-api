package app

import (
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/wsd"
)

var getWsdHomeHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	testParams := wsd.CreateExeParams{
		Arch: "win",
		Commands: []string{
			"www.facebook.com",
			"www.chatgpt.com",
		},
	}

	binBytes, name, err := wsd.CreateGoExe(testParams)
	if err != nil {
		panic(err)
	}
	api.Response.StreamBytes(w, 200, binBytes, name)
}
