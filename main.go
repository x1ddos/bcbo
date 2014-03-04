// +build !appengine

// A stand-alone HTTP server providing web UI for Beancounter Office.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/crhym3/bcbo/bc"
	"github.com/crhym3/bcbo/bo"
)

const (
	PATH_API_V1 = "/api/v1/"
	PATH_STATIC = "/static/"
)

var (
	// address:port the server will be listening on
	bindAddr string

	// Beancounter API base URL
	bcApiUrl string

	// static assets directory
	assetsDir string

	// temporary using a cmd-line provided API key.
	// This will be removed in the near future.
	bcApiKey string
)

// main is the program entry point when running in stand-alone (local) mode.
func main() {
	procDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	flag.StringVar(&bindAddr, "a", ":9090", "Address to listen on.")
	flag.StringVar(&bcApiUrl, "bcapi", "http://localhost:8080/beancounter-platform/rest",
		"Base URL of Beancounter Platform API.")
	flag.StringVar(&assetsDir, "assets", filepath.Join(procDir, "static"),
		"Static assets directory.")
	flag.StringVar(&bcApiKey, "apikey", "",
		"BC API key. This will be removed in the near future.")
	flag.Parse()

	var bcClient bc.Client
	var api *bo.Api

	log.Printf("Beancounter API base URL: %s", bcApiUrl)
	bcClient, err = bc.NewClient(bcApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	api, err = bo.NewApi(PATH_API_V1, bcClient, bcApiKey)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.DefaultServeMux
	mux.Handle(PATH_API_V1, api)

	if assetsDir != "" {
		log.Printf("Serving static assets from: %s", assetsDir)
		fileServer := http.FileServer(http.Dir(assetsDir))
		mux.Handle(PATH_STATIC, http.StripPrefix(PATH_STATIC, fileServer))
		// TODO: make a favicon
		mux.HandleFunc("/favicon.ico", http.NotFound)
		mux.HandleFunc("/", indexHandler)
	}

	log.Printf("Listening on %s", bindAddr)
	err = http.ListenAndServe(bindAddr, logHandler(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// indexHandler handles homepage requests.
// It currently serves static HTML file from static/index.html.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(assetsDir, "index.html"))
}

// logHandler logs inflight request and hands it over to the handler.
func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		log.Printf("%s %s %s\n", host, r.Method, r.URL.RequestURI())
		handler.ServeHTTP(w, r)
	})
}
