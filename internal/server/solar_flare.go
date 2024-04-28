package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SolarFlare struct {
	Begin string `json:"beginTime"`
	Peak  string `json:"peakTime"`
	End   string `json:"endTime"`
	Class string `json:"classType"`
	Link  string `json:"link"`
}

type SolarFlareResponse struct {
	Count       int          `json:"count"`
	SolarFlares []SolarFlare `json:"solarFlares"`
}

func (s *Server) SolarFlareHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.getSolarFlareData()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error getting solar flare data: %v", err)
	}

	if resp.Count == 0 {
		http.Error(w, "No solar flare data found", http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getSolarFlareData() (SolarFlareResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/DONKI/FLR?api_key=%s", s.nasa_api_key))
	if err != nil {
		return SolarFlareResponse{}, fmt.Errorf("error making request to NASA API. Err: %v", err)
	}

	var solarFlares []SolarFlare
	err = json.NewDecoder(resp.Body).Decode(&solarFlares)
	if err != nil {
		return SolarFlareResponse{}, fmt.Errorf("error decoding JSON. Err: %v", err)
	}

	count := len(solarFlares)
	solarFlareResponse := SolarFlareResponse{
		SolarFlares: solarFlares,
		Count:       count,
	}

	return solarFlareResponse, nil
}
