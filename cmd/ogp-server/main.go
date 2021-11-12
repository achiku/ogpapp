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

	p := fmt.Sprintf("127.0.0.1:%s", s.Config.APIServerPort)
	switch s.Config.TLS {
	case false:
		if err := http.ListenAndServe(p, s.Mux); err != nil {
			log.Fatalf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		srv := http.Server{
			Addr:      p,
			Handler:   s.Mux,
			TLSConfig: cfg,
		}
		if err := srv.ListenAndServeTLS(s.Config.ServerCertPath, s.Config.ServerKeyPath); err != nil {
			log.Fatalf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
