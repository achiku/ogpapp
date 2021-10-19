package ogpapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// App ogp.app
type App struct {
	Config        *Config
	KoruriBold    *truetype.Font
	OgpPagePath   string
	IndexPagePath string
	OgpPageTmpl   *template.Template
	IndexPageTmpl string
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
	pf, err := os.Open(path.Join("app", "public", "p.html"))
	if err != nil {
		return nil, err
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
	idxf, err := os.Open(path.Join("app", "public", "index.html"))
	if err != nil {
		return nil, err
	}
	defer idxf.Close()

	idxbuf, err := ioutil.ReadAll(idxf)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:        cfg,
		KoruriBold:    ft,
		OgpPageTmpl:   ogpPageTmpl,
		IndexPageTmpl: string(idxbuf),
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, app.IndexPageTmpl)
	return
}

type createImageReq struct {
	Words string `json:"words"`
}

// CreateImage create ogp image API
func (app *App) CreateImage(w http.ResponseWriter, r *http.Request) {
	var d createImageReq
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		// logger.Error().Msgf("decode failed: %s", err)
		return
	}
	words := d.Words
	id := uuid.New()
	filename := fmt.Sprintf("%s.png", id.String())
	filepath := path.Join("data", filename)
	wi, he, fs := app.Config.DefaultImageWidth, app.Config.DefaultImageHeight, app.Config.DefaultFontSize
	if err := createImage(wi, he, fs, app.KoruriBold, words, filepath); err != nil {
		// logger.Error().Msgf("create image failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	data := map[string]string{
		"words":   words,
		"file":    filename,
		"id":      id.String(),
		"baseURL": app.Config.BaseURL,
	}
	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		// logger.Printf("encode failed: %s", err)
		return
	}
	return
}
