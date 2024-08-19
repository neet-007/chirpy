package main

import (
	"net/http"
)

func main() {
	serverMux := http.NewServeMux()
	server := http.Server{Handler: serverMux}

	server.ListenAndServe()
}
