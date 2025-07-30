package server

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/achere/homework-pack-sizes/internal/pack"
)

//go:embed templates
var content embed.FS

type calculatePacksRequestV1 struct {
	Sizes []int `json:"sizes"`
	Order int   `json:"order"`
}

type calculatePacksResponseV1 struct {
	Packs map[int]int `json:"packs,omitempty"`
	Error string      `json:"error,omitempty"`
}

type calculatePacksRequest struct {
	Order int `json:"order"`
}

type calculatePacksResponse struct {
	Packs map[int]int `json:"packs,omitempty"`
	Sizes []int       `json:"sizes,omitempty"`
	Error string      `json:"error,omitempty"`
}

type storePackSizesRequest struct {
	Sizes []int `json:"sizes"`
}

type storePackSizesResponse struct {
	Error string `json:"error"`
}

type retrievePackSizesResponse struct {
	Sizes []int  `json:"sizes,omitempty"`
	Error string `json:"error,omitempty"`
}

// calculatePacksHandlerV1 provides an JSON interface to calculate pack sizes from the request
func (a *App) calculatePacksHandlerV1(w http.ResponseWriter, r *http.Request) {
	var req calculatePacksRequestV1

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(calculatePacksResponseV1{Error: err.Error()})
		return
	}

	packs, err := pack.CalculatePacks(req.Sizes, req.Order)
	if err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(calculatePacksResponseV1{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculatePacksResponseV1{Packs: packs})
}

// calculatePacksHandler provides an JSON interface to calculate pack sizes from SizeRepo
func (a *App) calculatePacksHandler(w http.ResponseWriter, r *http.Request) {
	var req calculatePacksRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(calculatePacksResponse{Error: err.Error()})
		return
	}

	packs, sizes, err := pack.CalculatePacksWithRepo(r.Context(), a.SizeRepo, req.Order)
	if err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, pack.ErrInvalidArg) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(calculatePacksResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculatePacksResponse{Packs: packs, Sizes: sizes})
}

// storePackSizesHandler allows to store pack sizes in the SizeRepo
func (a *App) storePackSizesHandler(w http.ResponseWriter, r *http.Request) {
	var req storePackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(storePackSizesResponse{Error: err.Error()})
		return
	}

	err := pack.SavePackSizes(context.Background(), a.SizeRepo, req.Sizes)
	if err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, pack.ErrInvalidArg) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(storePackSizesResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// retrievePackSizesHandler allows to retrieve pack sizes from SizeRepo
func (a *App) retrievePackSizesHandler(w http.ResponseWriter, r *http.Request) {
	sizes, err := a.SizeRepo.GetPackSizes(r.Context())
	if err != nil {
		a.logger.Error(err.Error(), "url", r.RequestURI)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(retrievePackSizesResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retrievePackSizesResponse{Sizes: sizes})
}

// uiHandler handles displating HTML UI
func (a *App) uiHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sizes, err := a.SizeRepo.GetPackSizes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Order string
		Sizes []int
	}{
		Order: a.Config.Order,
		Sizes: sizes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
