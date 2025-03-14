package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	clean_body := replaceBadWords(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: clean_body,
	})
}

func replaceBadWords(body string) string {
	words := strings.Split(body, " ")
	if len(words) < 1 {
		return body
	}
	for idx, word := range words {
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}
