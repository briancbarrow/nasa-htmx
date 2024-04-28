package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
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
	var requestData struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if (requestData.StartDate != "" && !isValidDate(requestData.StartDate)) || (requestData.EndDate != "" && !isValidDate(requestData.EndDate)) {
		http.Error(w, "Invalid date format. Use yyyy-MM-dd", http.StatusBadRequest)
		return
	}

	resp, err := s.getSolarFlareData(requestData.StartDate, requestData.EndDate)
	if err != nil {
		if err.Error() == "error making request to NASA API. Rate limit exceeded" {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatalf("Error getting solar flare data: %v", err)
		}
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

func (s *Server) getSolarFlareData(start_date, end_date string) (SolarFlareResponse, error) {
	baseUrl, err := url.Parse("https://api.nasa.gov/DONKI/FLR")
	if err != nil {
		log.Fatalf("error parsing URL. Err: %v", err)
	}

	params := url.Values{}
	if start_date != "" {
		params.Add("startDate", start_date)
	}
	if end_date != "" {
		params.Add("endDate", end_date)
	}
	params.Add("api_key", s.nasa_api_key)
	baseUrl.RawQuery = params.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return SolarFlareResponse{}, fmt.Errorf("error making request to NASA API. Err: %v", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return SolarFlareResponse{}, fmt.Errorf("error making request to NASA API. Rate limit exceeded")
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

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
