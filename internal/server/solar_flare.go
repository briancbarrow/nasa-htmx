package server

import (
	"encoding/json"
	"log"
	"nasa-htmx/cmd/api/nasa"
	helpers "nasa-htmx/internal"
	"net/http"
)

func (s *Server) SolarFlareHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	queryParams := r.URL.Query()
	if queryParams.Get("start_date") != "" {
		requestData.StartDate = queryParams.Get("start_date")
	}
	if queryParams.Get("end_date") != "" {
		requestData.EndDate = queryParams.Get("end_date")
	}

	if (requestData.StartDate != "" && !helpers.IsValidDate(requestData.StartDate)) || (requestData.EndDate != "" && !helpers.IsValidDate(requestData.EndDate)) {
		http.Error(w, "Invalid date format. Use yyyy-MM-dd", http.StatusBadRequest)
		return
	}

	resp, err := nasa.GetSolarFlareData(requestData.StartDate, requestData.EndDate)
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
