package main

import (
	"crypto/tls"
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

	var p string
	if s.Config.LocalDev {
		p = fmt.Sprintf("localhost:%s", s.Config.APIServerPort)
	} else {
		p = fmt.Sprintf(":%s", s.Config.APIServerPort)
	}
	switch s.Config.TLS {
	case false:
		log.Printf("starting dev http server at %s...", p)
		if err := http.ListenAndServe(p, s.Mux); err != nil {
			log.Fatalf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		log.Printf("starting tls server at %s...", p)
		log.Printf("server cert path=%s", s.Config.ServerCertPath)
		log.Printf("server key path=%s", s.Config.ServerKeyPath)
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		srv := http.Server{
			Handler:   s.Mux,
			TLSConfig: cfg,
		}
		if err := srv.ListenAndServeTLS(s.Config.ServerCertPath, s.Config.ServerKeyPath); err != nil {
			log.Fatalf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
