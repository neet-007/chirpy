package api

import (
	"encoding/json"
	"fmt"
	"github.com/neet-007/chirpy/database"
	"net/http"
	"strings"
)

func NewApiConfig() (ApiConfig, error) {
	db, err := database.NewDB("./database/database.json")

	if err != nil {
		return ApiConfig{}, err
	}

	return ApiConfig{
		fileserverHits: 0,
		db:             db,
	}, nil

}

type ApiConfig struct {
	fileserverHits int
	db             *database.DB
}

func (cfg *ApiConfig) HandlerValidatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		chirps, err := cfg.db.GetChirps()
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		json_, err := json.Marshal(chirps)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json_)
		return
	}

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

	CleanedBody := cleanProfane(params.Body)

	newData, err := cfg.db.CreateChirp(CleanedBody)
	fmt.Println(newData)
	if err != nil {
		fmt.Printf("Error creating chirp value: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(newData)
	if err != nil {
		fmt.Printf("Error encoding return value: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
		<body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		</body>

		</html>`, cfg.fileserverHits)))
}

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
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
