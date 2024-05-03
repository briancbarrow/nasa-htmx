package web

import (
	"log"
	"nasa-htmx/cmd/api/nasa"
	"net/http"
)

func AsteroidHandler(w http.ResponseWriter, r *http.Request) {
	asteroidData, err := nasa.GetAsteroidData()
	if err != nil {
		if err.Error() == "error making request to NASA API. Rate limit exceeded" {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatalf("Error getting asteroid data: %v", err)
		}
	}

	component := AsteroidsPage(asteroidData)

	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
