package main

import (
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	fileServer := http.FileServer(http.Dir("."))
	serverMux.Handle("/app/", http.StripPrefix("/app", fileServer))

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()
}
