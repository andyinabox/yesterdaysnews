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
	// AssetsRequestPathPrefix is a request prefix that will be stripped when converting
	// the URL to a filesystem path. For instance, the for request path "/assets/file.txt"
	// with `AssetsRequestPathPrefix` set to "/assets" will look in the filesystem for
	// the file path "file.txt".
	AssetsRequestPathPrefix string // url path ("/assets")
	// There are two modes of operation:
	// 1. Serving static files via an embedded filesystem
	// 2. Serving static files from the filesystem
	// The first is generaly for production, the secont for development
	AssetsEmbeddedFS        fs.FS
	AssetsDirFS             fs.FS
	UseFilesystemAssets     bool
	RedirectNotFoundToIndex bool
}

type Handler struct {
	mux           *http.ServeMux
	cfg           *Config
	assetsHandler http.Handler
	assetsFS      fs.FS
}

func New(cfg *Config) *Handler {

	h := &Handler{
		mux:      http.NewServeMux(),
		cfg:      cfg,
		assetsFS: cfg.AssetsEmbeddedFS,
	}

	// use the directory filesystem
	if cfg.UseFilesystemAssets {
		h.assetsFS = cfg.AssetsDirFS
	}

	// set the assets handler
	h.assetsHandler = http.StripPrefix(cfg.AssetsRequestPathPrefix, http.FileServer(http.FS(h.assetsFS)))

	return h
}

func (h *Handler) AddRoute(path string, handler http.HandlerFunc) {
	h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			log.Debugf("attempt to access path %q", r.URL.Path)
			if h.cfg.RedirectNotFoundToIndex {
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
	if strings.HasPrefix(r.URL.Path, h.cfg.AssetsRequestPathPrefix) {
		log.Debugf("serve using file server: %s", r.URL)
		h.assetsHandler.ServeHTTP(w, r)
		return
	}

	log.Debugf("serve using ServeMux: %s", r.URL)
	h.mux.ServeHTTP(w, r)
}
