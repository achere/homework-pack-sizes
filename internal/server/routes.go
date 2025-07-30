package server

import (
	"net/http"
)

func (a *App) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.uiHandler)

	// API versioning for backwards compatibility in case the functionality will change
	mux.HandleFunc("POST /api/v1/calculate-packs", a.calculatePacksHandlerV1)

	// V2
	mux.HandleFunc("POST /api/v2/calculate-packs", a.calculatePacksHandler)
	mux.HandleFunc("POST /api/v2/sizes", a.storePackSizesHandler)
	mux.HandleFunc("GET /api/v2/sizes", a.retrievePackSizesHandler)

	return mux
}
