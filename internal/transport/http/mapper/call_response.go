package mapper

import (
	"telephone-service/codegen/oas"
	"telephone-service/internal/domain"
)

func CallToResponse(call domain.Call) api.Call {
	return api.Call{
		ID:               call.Id,
		FirstSubscriber:  call.FirstSubscriber.Name,
		SecondSubscriber: call.SecondSubscriber.Name,
		Duration:         call.Duration.DurationToString(),
		CreatedAt:        call.CreatedAt.TimeToString(),
	}
}

func CallsToResponse(calls []domain.Call) []api.Call {
	response := make([]api.Call, len(calls))
	for i, call := range calls {
		response[i] = CallToResponse(call)
	}
	return response
}
