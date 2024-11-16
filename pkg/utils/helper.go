package utils

import (
	"encoding/json"
	"net/http"
)

type GetExchangeDTO struct {
	Bid string `json:"bid"`
}

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
