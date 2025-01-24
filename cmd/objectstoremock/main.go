package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/assetshandler"
)

var port int
var dir string
var verbose bool

func init() {
	flag.IntVar(&port, "p", 9000, "port to serve on")
	flag.StringVar(&dir, "d", "dist", "dir to serve")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func main() {

	// Check if the directory exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Fatalf("directory %q not found.\n", dir)
	}

	data, err := os.ReadFile(filepath.Join(dir, domain.ManifestFileName))
	if err != nil {
		log.Fatalf("problem reading manifest file: %s", err)
	}

	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("problem unmarshaling manifest file: %s", err)
	}

	if manifest.ID == "" {
		log.Fatal("no build ID found")
	}

	buildID := manifest.ID

	// Create a file server handler to serve the directory's contents
	// fileServer := http.FileServer(http.Dir(dir))

	handler := assetshandler.New(
		&assetshandler.Config{
			AssetsUrlPath:     "/" + buildID,
			AssetsFS:          os.DirFS(dir),
			StripAssetsPrefix: true,
		},
	)

	handler.AddRoute("/"+domain.CurrentBuildIDFileName, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(buildID))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: cors(handler),
	}

	// // configure cert
	// serverTLSCert, err := tls.LoadX509KeyPair(".cert/localhost.crt", ".cert/localhost.key")
	// if err != nil {
	// 	log.Fatalf("Error loading certificate and key file: %v", err)
	// }

	// // configure server
	// srv.TLSConfig = &tls.Config{
	// 	Certificates: []tls.Certificate{serverTLSCert},
	// }

	// // Start the server on port 8080
	// fmt.Printf("tls fileserver started at https://localhost:%d\n", port)
	// log.Fatal(srv.ListenAndServeTLS("", ""))

	fmt.Printf("fileserver started at http://localhost:%d\n", port)
	log.Fatal(srv.ListenAndServe())
}
