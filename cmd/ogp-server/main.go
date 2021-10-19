package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/achiku/ogpapp"
	"github.com/gorilla/mux"
)

var (
	configFile = flag.String("c", "", "configuration file path")
)

func main() {
	flag.Parse()

	cfg, err := ogpapp.NewConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to create config: %s", err)
	}
	app, err := ogpapp.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/").HandlerFunc(app.IndexPage)
	r.Methods(http.MethodGet).Path("/p/{id}").HandlerFunc(app.OgpPage)
	r.Methods(http.MethodGet).PathPrefix("/image/").Handler(
		http.StripPrefix("/image/", http.FileServer(http.Dir("data"))))

	// static asset
	r.Methods(http.MethodGet).PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir(path.Join("client", "dist", "js")))))
	r.Methods(http.MethodGet).PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir(path.Join("client", "dist", "css")))))
	r.Methods(http.MethodGet).PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(http.Dir(path.Join("client", "dist", "img")))))

	// API
	r.Methods(http.MethodPost).Path("/api/image").Handler(
		ogpapp.LoggingMiddleware(http.HandlerFunc(app.CreateImage)))

	switch cfg.TLS {
	case false:
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIServerPort), r); err != nil {
			log.Fatalf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", cfg.APIServerPort), cfg.ServerCertPath, cfg.ServerKeyPath, r); err != nil {
			log.Fatalf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
