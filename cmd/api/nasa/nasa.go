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

type Asteroid struct {
	Name              string   `json:"name"`
	InfoLink          string   `json:"nasa_jpl_url"`
	IsHazardous       bool     `json:"is_potentially_hazardous_asteroid"`
	Diameter          Diameter `json:"estimated_diameter"`
	CloseApproachDate string   `json:"close_approach_date"`
}

type AsteroidResponse struct {
	Count        int                   `json:"element_count"`
	AsteroidData map[string][]Asteroid `json:"asteroids"`
}

type Diameter struct {
	Feet struct {
		Min float64 `json:"estimated_diameter_min"`
		Max float64 `json:"estimated_diameter_max"`
	} `json:"feet"`
}

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

	if start_date == "" {
		start_date = time.Now().AddDate(0, 0, -30).Format("Jan 2, 2006")
	} else {
		sDate, err := time.Parse("2006-01-02", start_date)
		if err != nil {
			return SolarFlareResponse{}, fmt.Errorf("error processing date. Err: %v", err)
		}
		start_date = sDate.Format("Jan 2, 2006")
	}

	if end_date == "" {
		end_date = time.Now().Format("Jan 2, 2006")
	} else {
		eDate, err := time.Parse("2006-01-02", end_date)
		if err != nil {
			return SolarFlareResponse{}, fmt.Errorf("error processing date. Err: %v", err)

		}
		end_date = eDate.Format("Jan 2, 2006")
	}

	count := len(solarFlares)
	solarFlareResponse := SolarFlareResponse{
		SolarFlares: solarFlares,
		Count:       count,
		StartDate:   start_date,
		EndDate:     end_date,
	}

	return solarFlareResponse, nil
}

func GetAsteroidData() (AsteroidResponse, error) {
	api_key := os.Getenv("NASA_API_KEY")
	baseUrl, err := url.Parse("https://api.nasa.gov/neo/rest/v1/feed")
	if err != nil {
		log.Fatalf("error parsing URL. Err: %v", err)
	}

	params := url.Values{}
	params.Add("api_key", api_key)
	baseUrl.RawQuery = params.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return AsteroidResponse{}, fmt.Errorf("error making request to NASA API. Err: %v", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return AsteroidResponse{}, fmt.Errorf("error making request to NASA API. Rate limit exceeded")
	}

	var asteroidResponse struct {
		NearEarthObjects map[string][]Asteroid `json:"near_earth_objects"`
		ElementCount     int                   `json:"element_count"`
	}
	err = json.NewDecoder(resp.Body).Decode(&asteroidResponse)
	if err != nil {
		return AsteroidResponse{}, fmt.Errorf("error decoding JSON. Err: %v", err)
	}

	externalResponse := AsteroidResponse{
		AsteroidData: asteroidResponse.NearEarthObjects,
		Count:        asteroidResponse.ElementCount,
	}

	return externalResponse, nil
}
