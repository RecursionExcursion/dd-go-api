package main

import (
	"log"
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/app"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func main() {
	go beeGeesProtocol()
	app.App()
}

/* Keeps server from spinning down on Render's free tier */
func beeGeesProtocol() {
	self := lib.EnvGet("SELF_URL")
	for {
		time.Sleep(10 * time.Minute)
		http.Get(self)
		log.Printf("BeeGees protocol complete!")
	}
}
