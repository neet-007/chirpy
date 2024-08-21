package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neet-007/chirpy/api"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg, err := api.NewApiConfig()

	if err != nil {
		fmt.Printf("error happend %s", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerValidatePost)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerValidatePost)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.HandlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
