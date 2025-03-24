package api

import (
	"log"
	"net/http"
)

type APIServer struct {
	addr       string
	initalized bool
	server     http.Server
	router     *http.ServeMux
}

func NewApiServer(addr string) *APIServer {
	return &APIServer{
		addr:       addr,
		initalized: false,
	}
}

func (s *APIServer) Init(routes []RouteHandler) {
	s.router = http.NewServeMux()

	for _, r := range routes {
		s.router.HandleFunc(r.handleHttp())
	}

	s.server = http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	s.initalized = true

}

func (s *APIServer) ListenAndServe() error {
	if !s.initalized {
		panic("server not initalzed")
	}

	log.Printf("Server is listening on %s", s.addr)
	return s.server.ListenAndServe()
}
