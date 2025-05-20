package api

import (
	"fmt"
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

// addr = :PORT
func NewApiServer(addr string) *APIServer {
	return &APIServer{
		Addr:       addr,
		initalized: false,
	}
}

func (s *APIServer) Init(routes []RouteHandler) {
	defer func() {
		s.initalized = true
	}()

	s.Router = http.NewServeMux()

	for _, r := range routes {
		s.Router.HandleFunc(r.handleHttp())
	}

	s.Server = http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}
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
	return rh.MethodAndPath, pipeMiddleware(rh.Middleware...)(rh.Handler)
}

/* Middleware */

/* pipeMiddleware passes request through middleware fns from left to right */
func pipeMiddleware(mws ...Middleware) Middleware {
	return func(final HandlerFn) HandlerFn {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

type HTTPMethods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

func HttpMethodGenerator(base string) func(path ...string) HTTPMethods {
	return func(paths ...string) HTTPMethods {

		pathStr := ""

		numArgs := len(paths)

		if numArgs != 0 {
			if numArgs == 1 {
				pathStr = "/" + paths[0]
			} else {
				for _, arg := range paths {
					pathStr += "/" + arg
				}
			}
		}

		route := base + pathStr

		assign := func(s string) string {
			return fmt.Sprintf("%v %v", s, route)
		}

		routes := HTTPMethods{
			GET:    assign("GET"),
			POST:   assign("POST"),
			PUT:    assign("PUT"),
			PATCH:  assign("PATCH"),
			DELETE: assign("DELETE"),
		}

		return routes
	}
}
