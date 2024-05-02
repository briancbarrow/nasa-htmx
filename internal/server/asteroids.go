package server

import (
	"encoding/json"
	"log"
	"nasa-htmx/cmd/api/nasa"
	"net/http"
)

func (s *Server) AsteroidHandler(w http.ResponseWriter, r *http.Request) {
	asteroidData, err := nasa.GetAsteroidData()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error getting asteroid data: %v", err)
	}

	jsonResp, err := json.Marshal(asteroidData)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
