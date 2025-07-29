package server

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/achere/homework-pack-sizes/internal/pack"
)

//go:embed templates
var content embed.FS

type calculatePacksRequest struct {
	Sizes []int `json:"sizes"`
	Order int   `json:"order"`
}

type calculatePacksResponse struct {
	Packs map[int]int `json:"packs"`
	Error string      `json:"error"`
}

// calculatePacksHandler provides an JSON interface to calculate pack sizes.
func (a *App) calculatePacksHandler(w http.ResponseWriter, r *http.Request) {
	var req calculatePacksRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(calculatePacksResponse{Error: err.Error()})
		return
	}

	packs, err := pack.CalculatePacks(req.Sizes, req.Order)
	if err != nil {
		a.logger.Error(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(calculatePacksResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculatePacksResponse{Packs: packs})
}

// uiHandler handles displating HTML UI
func (a *App) uiHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Order string
		Sizes []int
	}{
		Order: a.Config.Order,
		Sizes: a.Config.Sizes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
