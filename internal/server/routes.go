package server

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// API versioning for backwards compatibility in case the functionality will change
	mux.HandleFunc("POST /api/v1/calculate-packs", calculatePacksHandler)

	return mux
}
