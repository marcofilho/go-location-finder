package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Cep string `json:"cep"`
}

type Response struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

func main() {
	port := ":9000"
	http.HandleFunc("/validate", validateRequest)

	log.Printf("Listening on port %s", port)
	http.ListenAndServe(port, nil)
}

func validateRequest(w http.ResponseWriter, r *http.Request) {
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

	callExternalAPI(w, cep.Cep)
}
func callExternalAPI(w http.ResponseWriter, cep string) {
	apiURL := "http://localhost:8080/cep"
	payload := map[string]string{"cep": cep}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to call external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(errBody)
		return
	}

	var apiResp Response
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		http.Error(w, "Failed to decode external API response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResp)
}
