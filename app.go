package ogpapp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// App ogp.app
type App struct {
	Config        *Config
	KoruriBold    *truetype.Font
	OgpPagePath   string
	IndexPagePath string
	OgpPageTmpl   *template.Template
	IndexPageTmpl *template.Template
}

// NewApp create app
func NewApp(cfg *Config) (*App, error) {
	fontBytes, err := ioutil.ReadFile(cfg.KoruriBoldFontPath)
	if err != nil {
		return nil, err
	}
	ft, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	pf, err := os.Open(path.Join("public", "p.html"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse p.html")
	}
	defer pf.Close()

	pbuf, err := ioutil.ReadAll(pf)
	if err != nil {
		return nil, err
	}
	ogpPageTmpl, err := template.New("page").Parse(string(pbuf))
	if err != nil {
		return nil, err
	}
	idxf, err := os.Open(path.Join("public", "index.html"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse index.html")
	}
	defer idxf.Close()

	idxbuf, err := ioutil.ReadAll(idxf)
	if err != nil {
		return nil, err
	}
	idxPageTmpl, err := template.New("index").Parse(string(idxbuf))
	if err != nil {
		return nil, err
	}

	return &App{
		Config:        cfg,
		KoruriBold:    ft,
		OgpPageTmpl:   ogpPageTmpl,
		IndexPageTmpl: idxPageTmpl,
	}, nil
}

// OgpPage display ogp page
func (app *App) OgpPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data := map[string]string{
		"id":      id,
		"file":    fmt.Sprintf("%s.png", id),
		"baseURL": app.Config.BaseURL,
	}
	w.WriteHeader(http.StatusOK)
	if err := app.OgpPageTmpl.Execute(w, data); err != nil {
		return
	}
	return
}

// IndexPage display index page
func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"baseURL": app.Config.BaseURL,
	}
	w.WriteHeader(http.StatusOK)
	if err := app.IndexPageTmpl.Execute(w, data); err != nil {
		return
	}
	return
}

// CreateImage create ogp image API
func (app *App) CreateImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	words := r.Form.Get("words")
	id := uuid.New()
	log.Printf("%s, %s", words, id)
	filename := fmt.Sprintf("%s.png", id.String())
	filepath := path.Join("data", filename)
	wi, he, fs := app.Config.DefaultImageWidth, app.Config.DefaultImageHeight, app.Config.DefaultFontSize
	if err := createImage(wi, he, fs, app.KoruriBold, words, filepath); err != nil {
		return
	}
	redirectURL := fmt.Sprintf("%s/p/%s", app.Config.BaseURL, id.String())
	log.Print(redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return
}
