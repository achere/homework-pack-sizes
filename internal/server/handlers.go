package server

import (
	"encoding/json"
	"net/http"

	"github.com/achere/homework-pack-sizes/internal/pack"
)

type calculatePacksRequest struct {
	Sizes []int `json:"sizes"`
	Order int   `json:"order"`
}

func calculatePacksHandler(w http.ResponseWriter, r *http.Request) {
	var req calculatePacksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	packs, err := pack.CalculatePacks(req.Sizes, req.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(packs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
