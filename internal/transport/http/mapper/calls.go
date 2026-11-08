package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
	"time"
)

func CallFromRequest(call api.StartCallRequest) (model.Call, error) {
	participants := make([]model.Subscriber, len(call.Participants))
	for i, participant := range call.Participants {
		subscriber, err := SubscriberFromRequest(participant)
		if err != nil {
			return model.Call{}, fmt.Errorf("не удалось преобразовать участника %d: %w", i+1, err)
		}

		participants[i] = subscriber
	}

	return model.Call{Participants: participants}, nil
}

func CallToResponse(call model.Call) api.Call {
	participants := make([]api.Subscriber, len(call.Participants))
	for i, participant := range call.Participants {
		participants[i] = SubscriberToResponse(participant)
	}

	participations := make([]api.CallParticipation, len(call.Participations))
	for i, participation := range call.Participations {
		participations[i] = api.CallParticipation{
			ID:           participation.Id,
			SubscriberID: participation.SubscriberId,
			JoinedAt:     time.Time(participation.JoinedAt),
			LeftAt:       timeToResponse(participation.LeftAt),
		}
	}

	return api.Call{
		ID:             call.Id,
		Participants:   participants,
		Duration:       call.Duration.DurationToString(),
		CreatedAt:      call.CreatedAt.TimeToString(),
		Status:         api.CallStatus(call.Status),
		StartedAt:      timeToResponse(call.StartedAt),
		FinishedAt:     timeToResponse(call.FinishedAt),
		Participations: participations,
	}
}

func CallsToResponse(calls []model.Call) []api.Call {
	response := make([]api.Call, len(calls))
	for i, call := range calls {
		response[i] = CallToResponse(call)
	}

	return response
}

func timeToResponse(value *model.Time) api.OptDateTime {
	if value == nil {
		return api.OptDateTime{}
	}

	return api.NewOptDateTime(time.Time(*value))
}

func CallsFilterFromParams(params api.GetCallsParams) model.CallsFilter {
	return model.CallsFilter{Status: model.CallStatus(params.Status.Value)}
}
