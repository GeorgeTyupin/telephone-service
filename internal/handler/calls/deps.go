package calls_handler

import (
	"net/url"
	"telephone-service/internal/dto"
)

type GetCallsUC interface {
	Execute(params url.Values) ([]dto.Call, error)
}
