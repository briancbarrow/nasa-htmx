package nasa

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
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
	StartDate   string       `json:"startDate"`
	EndDate     string       `json:"endDate"`
}

func GetSolarFlareData(start_date, end_date string) (SolarFlareResponse, error) {
	api_key := os.Getenv("NASA_API_KEY")
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
	params.Add("api_key", api_key)
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

	sort.Slice(solarFlares, func(i, j int) bool {
		iTime, _ := time.Parse("2006-01-02T15:04Z", solarFlares[i].Begin)
		jTime, _ := time.Parse("2006-01-02T15:04Z", solarFlares[j].Begin)
		return iTime.After(jTime)
	})

	count := len(solarFlares)
	solarFlareResponse := SolarFlareResponse{
		SolarFlares: solarFlares,
		Count:       count,
	}

	return solarFlareResponse, nil
}
