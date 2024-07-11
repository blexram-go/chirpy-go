package main

import (
	"net/http"
	"strconv"

	"github.com/gobash-blex/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	if dbChirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "Request not allowed")
		return
	}
	err = cfg.DB.DeleteChirp(dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
