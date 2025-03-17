package main

import (
	"net/http"

	"github.com/balayher/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteSingleChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

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

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", nil)
		return
	}

	err = cfg.db.DeleteChrip(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
