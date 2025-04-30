package handler

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Cep string `json:"cep"`
}

func ValidateRequest(w http.ResponseWriter, r *http.Request) {
	var cep Request
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if cep.Cep == "" {
		http.Error(w, "CEP is required", http.StatusBadRequest)
		return
	}

	if len(cep.Cep) < 9 {
		http.Error(w, "CEP must be at least 9 characters long", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	
}
