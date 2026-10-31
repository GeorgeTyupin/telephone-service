package handlers

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/require"
)

func TestNewErrorHandler_InternalError_LogsCauseAndReturnsSafeResponse(t *testing.T) {
	t.Parallel()

	// arrange
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := NewErrorHandler(logger)
	request := httptest.NewRequest(http.MethodGet, "/calls", nil)
	response := httptest.NewRecorder()
	cause := errors.New("database password=secret: connection failed")

	// act
	handler(context.Background(), response, request, cause)

	// assert
	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.JSONEq(t, `{"error_message":"Внутренняя ошибка сервера"}`, response.Body.String())
	require.NotContains(t, response.Body.String(), cause.Error())
	require.Contains(t, logs.String(), `"level":"ERROR"`)
	require.Contains(t, logs.String(), cause.Error())
}

func TestNewErrorHandler_BadRequest_PreservesResponseAndLogsWarning(t *testing.T) {
	t.Parallel()

	// arrange
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := NewErrorHandler(logger)
	request := httptest.NewRequest(http.MethodPost, "/call", nil)
	response := httptest.NewRecorder()
	cause := &ogenerrors.DecodeRequestError{
		OperationContext: ogenerrors.OperationContext{Name: "CreateCall"},
		Err:              errors.New("invalid request body"),
	}
	expected := httptest.NewRecorder()
	ogenerrors.DefaultErrorHandler(context.Background(), expected, request, cause)

	// act
	handler(context.Background(), response, request, cause)

	// assert
	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Equal(t, expected.Header(), response.Header())
	require.JSONEq(t, expected.Body.String(), response.Body.String())
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), cause.Error())
}

func TestNewErrorHandler_CanceledContext_Returns499AndLogsWarning(t *testing.T) {
	t.Parallel()

	// arrange
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := NewErrorHandler(logger)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/call", nil).WithContext(ctx)
	response := httptest.NewRecorder()
	cause := &ogenerrors.DecodeRequestError{
		OperationContext: ogenerrors.OperationContext{Name: "CreateCall"},
		Err:              context.Canceled,
	}

	// act
	handler(ctx, response, request, cause)

	// assert
	require.Equal(t, 499, response.Code)
	require.JSONEq(t, `{"error_message":"Запрос отменён клиентом"}`, response.Body.String())
	require.Contains(t, logs.String(), `"level":"WARN"`)
	require.Contains(t, logs.String(), `"status":499`)
	require.Contains(t, logs.String(), `"path":"/call"`)
	require.Contains(t, logs.String(), cause.Error())
}
