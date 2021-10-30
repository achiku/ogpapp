package ogpapp

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Server server
type Server struct {
	App    *App
	Config *Config
	Mux    *mux.Router
}

// NewServer creates server
func NewServer(cfgFile string) (*Server, error) {
	cfg, err := NewConfig(cfgFile)
	if err != nil {
		return nil, errors.Wrap(err, "NewConfig failed")
	}
	app, err := NewApp(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "NewApp failed")
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
		LoggingMiddleware(http.HandlerFunc(app.CreateImage)))

	return &Server{
		App:    app,
		Config: cfg,
		Mux:    r,
	}, nil
}
