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
	cfbr.CFBR()
}

func beeGeesProtocol() {
	self := lib.EnvGetOrPanic("SELF_URL")

	for {
		<-time.After(time.Minute * 12)
		http.Get(self)
		log.Printf("BGP %v", time.Now())
	}
}
