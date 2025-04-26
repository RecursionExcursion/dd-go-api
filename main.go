package main

import (
	"log"
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/app"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func main() {
	go beeGeesProtocol(12)
	app.App()
}

func beeGeesProtocol(min int) {
	self := lib.EnvGetOrPanic("SELF_URL")

	for {
		<-time.After(time.Minute * time.Duration(min))
		http.Get(self)
		log.Printf("BGP %v", time.Now())
	}
}
