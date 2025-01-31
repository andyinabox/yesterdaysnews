package assetshandler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
)

type Route interface {
	Path() string
	HandlerFunc() http.HandlerFunc
}

type Config struct {
	AssetsUrlPath     string // url path ("/assets")
	AssetsFS          fs.FS  // if using embed.FS, you will probably need to use `fs.Sub`
	StripAssetsPrefix bool   // probably true
	RedirectToIndex   bool
}

type Handler struct {
	mux           *http.ServeMux
	cfg           *Config
	assetsHandler http.Handler
}

func New(cfg *Config) *Handler {

	handler := http.FileServer(http.FS(cfg.AssetsFS))

	if cfg.StripAssetsPrefix {
		handler = http.StripPrefix(cfg.AssetsUrlPath, handler)
	}

	return &Handler{
		mux:           http.NewServeMux(),
		cfg:           cfg,
		assetsHandler: handler,
	}
}

func (h *Handler) AddRoute(path string, handler http.HandlerFunc) {
	h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			log.Infof("attempt to access path %q", r.URL.Path)
			if h.cfg.RedirectToIndex {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.NotFound(w, r)
			}
			return
		}
		handler(w, r)
	})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if the url contains the assets path, use the assets handler
	if strings.HasPrefix(r.URL.Path, h.cfg.AssetsUrlPath) {
		log.Debugf("serve using file server: %s", r.URL)
		h.assetsHandler.ServeHTTP(w, r)
		return
	}

	log.Debugf("serve using ServeMux: %s", r.URL)
	h.mux.ServeHTTP(w, r)
}
