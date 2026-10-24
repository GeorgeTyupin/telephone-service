package app

import "net/http"

type CallsHandler interface {
	GetCalls(w http.ResponseWriter, r *http.Request)
	CreateCall(w http.ResponseWriter, r *http.Request)
}

type AbonentsHandler interface {
	GetAbonents(w http.ResponseWriter, r *http.Request)
}
