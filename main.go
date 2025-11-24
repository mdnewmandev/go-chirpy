package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	server.ListenAndServe()
}