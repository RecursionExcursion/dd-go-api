package main

import (
	"log"
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func main() {
	// go beeGeesProtocol()
	// app.App()
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

func beeGeesProtocol() {
	self := lib.EnvGetOrPanic("SELF_URL")

	for {
		<-time.After(time.Minute * 12)
		http.Get(self)
		log.Printf("BGP %v", time.Now())
	}
}
