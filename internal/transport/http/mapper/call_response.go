package mapper

import (
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

func CallToResponse(call model.Call) api.Call {
	return api.Call{
		ID:               call.Id,
		FirstSubscriber:  call.FirstSubscriber.Name,
		SecondSubscriber: call.SecondSubscriber.Name,
		Duration:         call.Duration.DurationToString(),
		CreatedAt:        call.CreatedAt.TimeToString(),
	}
}

func CallsToResponse(calls []model.Call) []api.Call {
	response := make([]api.Call, len(calls))
	for i, call := range calls {
		response[i] = CallToResponse(call)
	}
	return response
}
