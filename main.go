package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidatePost)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidatePost(w http.ResponseWriter, r *http.Request) {
	type parammeter struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parammeter{}
	err := decoder.Decode(&params)

	if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		type returnVal struct {
			Error string `json:"error"`
		}

		fmt.Println("longer that 140 chars post")
		w.WriteHeader(http.StatusBadRequest)

		returnVal_ := returnVal{
			Error: "Chirp is too long",
		}

		json, err := json.Marshal(returnVal_)

		if err != nil {
			fmt.Printf("Error encoding return value: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	returnVal_ := returnVal{
		CleanedBody: cleanProfane(params.Body),
	}

	json, err := json.Marshal(returnVal_)

	if err != nil {
		fmt.Printf("Error encoding return value: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
		<body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		</body>

		</html>`, cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func cleanProfane(s string) string {
	words := strings.Split(s, " ")
	profanes := []string{"kerfuffle",
		"sharbert",
		"fornax",
	}

	for index, word := range words {
		for _, profane := range profanes {
			if strings.ToLower(word) == profane {
				words[index] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}
