package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/neet-007/chirpy/api"
)

func main() {
	godotenv.Load()
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
	mux.HandleFunc("GET /api/chirps/{chat_id}", apiCfg.HandlerGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chat_id}", apiCfg.HandlerDeleteChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerValidatePost)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUpdateUser)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevokeToken)
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
