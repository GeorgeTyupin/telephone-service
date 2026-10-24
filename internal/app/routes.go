package app

import (
	"net/http"
)

func (a *App) registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /calls", a.callsHandler.GetCalls)
	mux.HandleFunc("GET /abonets", a.abonentsHandler.GetAbonents)
	mux.HandleFunc("POST /call", a.callsHandler.CreateCall)

	return mux
}
