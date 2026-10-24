package abonents_handler

import (
	"log/slog"
	"net/http"
)

type AbonentsHandler struct {
	logger *slog.Logger

	getAbonentsUC GetAbonentsUC
}

func NewAbonentsHandler(
	logger *slog.Logger,
	getAbonentsUC GetAbonentsUC,
) *AbonentsHandler {
	return &AbonentsHandler{
		logger:        logger,
		getAbonentsUC: getAbonentsUC,
	}
}

func (h *AbonentsHandler) GetAbonents(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	abonents, err := h.getAbonentsUC.Execute(params)
	if err != nil {

	}
}
