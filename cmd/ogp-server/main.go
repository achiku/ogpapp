package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/achiku/ogpapp"
)

var (
	configFile = flag.String("c", "", "configuration file path")
)

func main() {
	flag.Parse()

	s, err := ogpapp.NewServer(*configFile)
	if err != nil {
		log.Fatalf("NewServer failed: %s", err)
	}

	p := fmt.Sprintf("localhost:%s", s.Config.APIServerPort)
	switch s.Config.TLS {
	case false:
		if err := http.ListenAndServe(p, s.Mux); err != nil {
			log.Fatalf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		if err := http.ListenAndServeTLS(p, s.Config.ServerCertPath, s.Config.ServerKeyPath, s.Mux); err != nil {
			log.Fatalf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
