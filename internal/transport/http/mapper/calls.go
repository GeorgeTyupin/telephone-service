package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

func CallFromRequest(call api.CreateCallRequest) (model.Call, error) {
	duration, err := model.StringToDuration(call.Duration)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать длительность звонка: %w", err)
	}

	createdAt, err := model.StringToTime(call.CreatedAt)
	if err != nil {
		return model.Call{}, fmt.Errorf("не удалось преобразовать время создания звонка: %w", err)
	}

	participants := make([]model.Subscriber, len(call.Participants))
	for i, participant := range call.Participants {
		subscriber, err := SubscriberFromRequest(participant)
		if err != nil {
			return model.Call{}, fmt.Errorf("не удалось преобразовать участника %d: %w", i+1, err)
		}

		participants[i] = subscriber
	}

	return model.Call{
		Participants: participants,
		Duration:     duration,
		CreatedAt:    createdAt,
	}, nil
}

func CallToResponse(call model.Call) api.Call {
	participants := make([]api.Subscriber, len(call.Participants))
	for i, participant := range call.Participants {
		participants[i] = SubscriberToResponse(participant)
	}

	return api.Call{
		ID:           call.Id,
		Participants: participants,
		Duration:     call.Duration.DurationToString(),
		CreatedAt:    call.CreatedAt.TimeToString(),
	}
}

func CallsToResponse(calls []model.Call) []api.Call {
	response := make([]api.Call, len(calls))
	for i, call := range calls {
		response[i] = CallToResponse(call)
	}

	return response
}
