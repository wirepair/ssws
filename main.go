package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

var (
	staticPath string
	hostname   string
	certPath   string
	email      string
)

func init() {
	flag.StringVar(&staticPath, "files", "/opt/ssws/static", "path to static files to serve")
	flag.StringVar(&hostname, "host", "", "hostname to use for serving files from")
	flag.StringVar(&certPath, "certs", "/opt/ssws/certs", "path to autocert cache")
	flag.StringVar(&email, "email", "", "email to register with lets encrypt")

}

func main() {
	flag.Parse()

	if hostname == "" {
		log.Fatal("you did not configure a valid hostname")
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostname),
		Cache:      autocert.DirCache(certPath),
	}

	if email != "" {
		m.Email = email
	}

	mux := http.NewServeMux()
	addRoutes(mux)

	httpsServer := &http.Server{
		Addr:         ":443",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
		Handler:      mux,
	}

	httpServer := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      m.HTTPHandler(nil),
	}

	go func() {
		log.Printf("Starting HTTP server\n")
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("error from HTTP server: %s\n", err)
		}
	}()

	log.Printf("Starting HTTPS server\n")
	err := httpsServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("error from HTTPS server: %s\n", err)
	}
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/", serve())
	// add custom handlers here if necessary
}

func serve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		http.StripPrefix("/", http.FileServer(http.Dir(staticPath))).ServeHTTP(w, r)
		log.Printf("%s [%s] %s", r.RemoteAddr, r.URL, time.Since(start))
	})
}
