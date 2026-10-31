package http

import (
	"context"
	"net/http"
	api "telephone-service/codegen/oas"
)

func newErrorHandler() api.ErrorHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {

	}
}
