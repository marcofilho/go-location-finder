package cmd

import (
	"net/http"

	"github.com/marcofilho/go-location-finder/go-location-validation/internal/handler"
)

func main() {
	// Initialize the HTTP server
	http.HandleFunc("/validate", handler.ValidateRequest)

}
