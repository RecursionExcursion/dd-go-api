package app

import (
	"log"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
)

type Server struct {
	repo *cfbr.CfbrRepo
}

func (s *Server) handleCfbr(w http.ResponseWriter, r *http.Request) {

	cfh := cfbr.CfbrHandler{
		Repo: s.repo,
	}

	switch r.Method {
	case http.MethodGet:
		pipe(globalMw...)(cfh.CFBRGet)(w, r)
	case http.MethodPost:
		pipe(globalMw...)(cfh.CFBRPost)(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

var globalMw = []middleware{
	LoggerMW(log.Default()),
	// GeoLimitMW(geoParams),
	RateLimitMW(5, 10),
}
