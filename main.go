package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFileserverHits() int32 {
	return cfg.fileserverHits.Load()
}

func (cfg *apiConfig) resetFileserverHits() {
	cfg.fileserverHits.Store(0)
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("app")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("app/assets"))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Hits: " + fmt.Sprintf("%d", apiCfg.getFileserverHits())))
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.resetFileserverHits()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Hits counter reset"))
	})

	server.ListenAndServe()
}