package main

import (
	"log"
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/cfbr/core"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func main() {
	// go beeGeesProtocol(12)
	// app.App()
	core.CompileSeason(2024)
}

func beeGeesProtocol(min int) {
	self := lib.EnvGetOrPanic("SELF_URL")

	for {
		<-time.After(time.Minute * time.Duration(min))
		http.Get(self)
		log.Printf("BGP %v", time.Now())
	}
}
