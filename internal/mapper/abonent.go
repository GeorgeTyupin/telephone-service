package mapper

import (
	"log/slog"
	"telephone-service/internal/domain"
	"telephone-service/internal/dto"
	"telephone-service/internal/record"
)

func (m *Mapper) AbonentRecordToDomain(abonent record.Abonent) domain.Abonent {
	const op = "mapper.AbonentRecordToDomain"
	logger := m.logger.With(slog.String("operation", op))

	number, err := domain.StringToPhoneNumber(abonent.PhoneNumber)
	if err != nil {
		logger.Error("Ошибка маппинга записи абонента к домену", slog.String("error", err.Error()))
		return domain.WrongAbonent
	}

	return domain.Abonent{
		Id:          int(abonent.Id),
		Name:        abonent.Name,
		PhoneNumber: number,
	}
}

func (*Mapper) AbonentDomainToDTO(abonent domain.Abonent) dto.Abonent {
	return dto.Abonent{
		Id:          abonent.Id,
		Name:        abonent.Name,
		PhoneNumber: domain.PhoneNumberToString(abonent.PhoneNumber),
	}
}

func (*Mapper) AbonentDTOToDomain(abonent dto.Abonent) domain.Abonent {
	number, err := domain.StringToPhoneNumber(abonent.PhoneNumber)
	if err != nil {
		return domain.WrongAbonent
	}

	return domain.Abonent{
		Id:          abonent.Id,
		Name:        abonent.Name,
		PhoneNumber: number,
	}
}
