package main

import (
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Addr:    ":8080", // Specify the port number here
		Handler: serverMux,
	}

	server.ListenAndServe()
}
