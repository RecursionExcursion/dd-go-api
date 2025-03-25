package api

import (
	"log"
	"net/http"
)

type APIServer struct {
	Addr   string
	Server http.Server
	Router *http.ServeMux

	initalized bool
}

type HandlerFn = func(http.ResponseWriter, *http.Request)
type Middleware = func(HandlerFn) HandlerFn

type RouteHandler struct {
	MethodAndPath string
	Handler       HandlerFn
	Middleware    []Middleware
}

func NewApiServer(addr string) *APIServer {
	return &APIServer{
		Addr:       addr,
		initalized: false,
	}
}

func (s *APIServer) Init(routes []RouteHandler) {
	s.Router = http.NewServeMux()

	for _, r := range routes {
		s.Router.HandleFunc(r.handleHttp())
	}

	s.Server = http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}

	s.initalized = true

}

func (s *APIServer) ListenAndServe() error {
	if !s.initalized {
		panic("server not initalzed")
	}

	log.Printf("Server is listening on %s", s.Addr)
	return s.Server.ListenAndServe()
}

func (rh *RouteHandler) handleHttp() (string, func(http.ResponseWriter, *http.Request)) {

	if rh.Handler == nil {
		panic("handler is nil for route " + rh.MethodAndPath)
	}

	mwPipe := pipeMiddleware(rh.Middleware...)
	return rh.MethodAndPath, mwPipe(rh.Handler)
}

/* Middleware */

/* pipeMiddleware passes request through middleware fn from left to right */
func pipeMiddleware(mws ...Middleware) Middleware {
	return func(final HandlerFn) HandlerFn {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}
