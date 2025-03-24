package api

import (
	"net/http"
)

type RouteHandler struct {
	MethodAndPath string
	Handler       HandlerFn
	Middleware    []Middleware
}

func (rh *RouteHandler) handleHttp() (string, func(http.ResponseWriter, *http.Request)) {

	if rh.Handler == nil {
		panic("handler is nil for route " + rh.MethodAndPath)
	}

	mwPipe := middlewarePipe(rh.Middleware...)
	return rh.MethodAndPath, mwPipe(rh.Handler)
}
