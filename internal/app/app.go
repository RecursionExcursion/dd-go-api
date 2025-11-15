package app

import (
	"context"
	"log"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
)

func App() {
	ctx := context.Background()

	repo, err := cfbr.CfbrRepository(ctx)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer repo.Conn.Close(ctx)

	srv := &Server{
		repo: &repo,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cfbr", srv.handleCfbr)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
