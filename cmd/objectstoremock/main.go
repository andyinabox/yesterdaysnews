package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
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

func main() {

	// Check if the directory exists
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Fatalf("Directory '%s' not found.\n", dir)
	}

	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir(dir))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: fileServer,
	}

	// configure cert
	serverTLSCert, err := tls.LoadX509KeyPair(".cert/localhost.crt", ".cert/localhost.key")
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	// configure server
	srv.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	// Start the server on port 8080
	fmt.Printf("tls fileserver started at https://localhost:%d\n", port)
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
