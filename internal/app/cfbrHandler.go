package app

import (
	"log"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
)

func handleCfbrGet(w http.ResponseWriter, r *http.Request) {

	season, err := cfbr.Create("fbs", 2024)
	if err != nil {
		panic(err)
	}

	s, err := season.FindSchoolById(194)
	if err != nil {
		panic(err)
	}
	// log.Println(s)
	log.Println(len(s.Games))
}
