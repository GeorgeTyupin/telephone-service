package mapper

import (
	"log/slog"
	"telephone-service/internal/domain"
	"telephone-service/internal/dto"
	"telephone-service/internal/record"
)

func (m *Mapper) CallRecordToDomain(
	call record.Call,
	first record.Abonent,
	second record.Abonent,
) domain.Call {
	return domain.Call{
		Id:            int(call.Id),
		FirstAbonent:  m.AbonentRecordToDomain(first),
		SecondAbonent: m.AbonentRecordToDomain(second),
		Duration:      domain.Duration(call.Duration),
		CreatedAt:     domain.Time(call.CreatedAt),
	}
}

func (*Mapper) CallDomainToDTO(call domain.Call) dto.Call {
	return dto.Call{
		Id:            call.Id,
		FirstAbonent:  call.FirstAbonent.Name,
		SecondAbonent: call.SecondAbonent.Name,
		Duration:      call.Duration.DurationToString(),
		CreatedAt:     call.CreatedAt.TimeToString(),
	}
}

func (m *Mapper) CallDTOToDomain(
	call dto.Call,
	first dto.Abonent,
	second dto.Abonent,
) domain.Call {
	const op = "mapper.CallDTOToDomain"
	logger := m.logger.With(slog.String("operation", op))

	duration, err := domain.StringToDuration(call.Duration)
	if err != nil {
		logger.Error(
			"Ошибка маппинга dto продолжительности времени к домену",
			slog.String("error", err.Error()),
		)
	}

	createdAt, err := domain.StringToTime(call.CreatedAt)
	if err != nil {
		logger.Error("Ошибка маппинга dto времени к домену",
			slog.String("error", err.Error()),
		)
	}

	return domain.Call{
		Id:            call.Id,
		FirstAbonent:  m.AbonentDTOToDomain(first),
		SecondAbonent: m.AbonentDTOToDomain(second),
		Duration:      duration,
		CreatedAt:     createdAt,
	}
}
