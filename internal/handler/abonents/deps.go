package abonents_handler

import (
	"net/url"
	"telephone-service/internal/dto"
)

type GetAbonentsUC interface {
	Execute(params url.Values) ([]dto.Abonent, error)
}
