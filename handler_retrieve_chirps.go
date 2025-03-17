package main

import (
	"net/http"

	"github.com/balayher/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRetrieveChirps(w http.ResponseWriter, req *http.Request) {
	authorID := req.URL.Query().Get("author_id")

	if authorID == "" {
		dbChirps, err := cfg.db.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
			return
		}
		FormatChirps(w, dbChirps)
		return
	}

	userID, err := uuid.Parse(authorID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}
	dbChirps, err := cfg.db.GetChirpsByUserID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps", err)
		return
	}
	FormatChirps(w, dbChirps)
}

func (cfg *apiConfig) handlerRetrieveSingleChirp(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}

	retrieved_chirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, retrieved_chirp)
}

func FormatChirps(w http.ResponseWriter, dbChirps []database.Chirp) {
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
