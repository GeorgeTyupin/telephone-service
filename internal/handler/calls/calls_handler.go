package calls_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type CallsHandler struct {
	logger *slog.Logger

	getCallsUC GetCallsUC
}

func NewCallsHandler(
	logger *slog.Logger,
	getCallsUC GetCallsUC,
) *CallsHandler {
	return &CallsHandler{
		logger:     logger,
		getCallsUC: getCallsUC,
	}
}

func (h *CallsHandler) GetCalls(w http.ResponseWriter, r *http.Request) {
	const op = "calls_handler.CallsHandler.GetCalls"
	logger := h.logger.With(slog.String("operation", op))

	params := r.URL.Query()

	calls, err := h.getCallsUC.Execute(params)
	if err != nil {
		http.Error(w, "Не удалось получить звонки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(calls); err != nil {
		logger.Error("не удалось закодировать ответ", slog.String("error", err.Error()))
		return
	}
}

func (h *CallsHandler) CreateCall(w http.ResponseWriter, r *http.Request) {}
