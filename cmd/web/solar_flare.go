package web

import (
	"fmt"
	"log"
	"nasa-htmx/cmd/api/nasa"
	helpers "nasa-htmx/internal"
	"time"

	"net/http"
)

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04Z", timeStr)
}

func convertTimes(sft *SolarFlareTime, origSf *nasa.SolarFlare) error {
	var err error
	sft.Begin, err = parseTime(origSf.Begin)
	if err != nil {
		return fmt.Errorf("error parsing Begin time: %w", err)
	}
	sft.Peak, err = parseTime(origSf.Peak)
	if err != nil {
		return fmt.Errorf("error parsing Peak time: %w", err)
	}
	sft.End, err = parseTime(origSf.End)
	if err != nil {
		return fmt.Errorf("error parsing End time: %w", err)
	}
	return nil
}

func SolarFlareHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	start_date := r.FormValue("start_date")
	end_date := r.FormValue("end_date")
	if (start_date != "" && !helpers.IsValidDate(start_date)) || (end_date != "" && !helpers.IsValidDate(end_date)) {
		http.Error(w, "Invalid date format. Use yyyy-MM-dd", http.StatusBadRequest)
		return
	}

	sfData, err := nasa.GetSolarFlareData(start_date, end_date)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error getting solar flare data: %v", err)
	}

	sfTimeData := SolarFlareTableData{
		Count: sfData.Count,
	}
	sfTimeList := []SolarFlareTime{}
	for _, sf := range sfData.SolarFlares {
		sfTime := SolarFlareTime{}
		convertTimes(&sfTime, &sf)
		sfTime.Class = sf.Class
		sfTime.Link = sf.Link
		sfTimeList = append(sfTimeList, sfTime)
	}
	sfTimeData.SolarFlares = sfTimeList

	if start_date == "" {
		start_date = time.Now().AddDate(0, 0, -30).Format("Jan 2, 2006")
	} else {
		sDate, err := time.Parse("2006-01-02", start_date)
		if err != nil {
			http.Error(w, "error processing date", http.StatusInternalServerError)
			return
		}
		start_date = sDate.Format("Jan 2, 2006")
	}

	if end_date == "" {
		end_date = time.Now().Format("Jan 2, 2006")
	} else {
		eDate, err := time.Parse("2006-01-02", end_date)
		if err != nil {
			http.Error(w, "error processing date", http.StatusInternalServerError)
			return
		}
		end_date = eDate.Format("Jan 2, 2006")
	}

	component := SolarFlareTable(sfTimeData, start_date, end_date)

	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
