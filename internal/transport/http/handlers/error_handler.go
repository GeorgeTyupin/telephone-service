package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	api "telephone-service/codegen/oas"
)

const (
	statusClientClosedRequest = 499
	internalErrorMessage      = "Внутренняя ошибка сервера"
	clientClosedMessage       = "Запрос отменён клиентом"
)

func NewErrorHandler(logger *slog.Logger) api.ErrorHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		status := http.StatusInternalServerError
		message := internalErrorMessage
		ogenErrorsCode := ogenerrors.ErrorCode(err)
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			status = statusClientClosedRequest
			message = clientClosedMessage
		case ogenErrorsCode != http.StatusInternalServerError:
			status = ogenErrorsCode
			message = err.Error()
		}

		level := slog.LevelWarn
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		}
		logger.Log(ctx, level, "Ошибка обработки HTTP-запроса",
			slog.String("path", r.URL.Path),
			slog.Int("status", status),
			slog.Any("error", err),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(api.Error{ErrorMessage: message})
	}
}
