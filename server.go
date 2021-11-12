package ogpapp

import (
	"net/http"

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
	r.Methods(http.MethodGet).PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("public"))))

	// API
	r.Methods(http.MethodPost).Path("/api/image").Handler(
		LoggingMiddleware(http.HandlerFunc(app.CreateImage)))

	return &Server{
		App:    app,
		Config: cfg,
		Mux:    r,
	}, nil
}
