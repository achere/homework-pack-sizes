package server

import (
	"net/http"
)

func (a *App) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.uiHandler)

	// API versioning for backwards compatibility in case the functionality will change
	mux.HandleFunc("POST /api/v1/calculate-packs", a.calculatePacksHandler)

	return mux
}
