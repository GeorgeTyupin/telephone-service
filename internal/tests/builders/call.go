package builders

import (
	"time"

	api "telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

type CallOption func(*model.Call)
type ParticipationOption func(*model.CallParticipation)
type StartCallRequestOption func(*api.StartCallRequest)

func NewCall(options ...CallOption) model.Call {
	startedAt := model.Time(time.Now().UTC().Truncate(time.Microsecond).Add(-5 * time.Minute))
	call := model.Call{
		Participants: []model.Subscriber{NewSubscriber()},
		CreatedAt:    startedAt,
		Status:       model.CallStatusActive,
		StartedAt:    &startedAt,
		Version:      1,
	}
	for _, option := range options {
		option(&call)
	}
	return call
}

func WithFinished(duration time.Duration) CallOption {
	return func(call *model.Call) {
		finishedAt := model.Time(time.Time(*call.StartedAt).Add(duration))
		call.Status = model.CallStatusFinished
		call.Duration = model.Duration(duration)
		call.FinishedAt = &finishedAt
	}
}

func NewParticipation(call model.Call, subscriberID int, options ...ParticipationOption) model.CallParticipation {
	participation := model.CallParticipation{
		SubscriberId: subscriberID,
		JoinedAt:     *call.StartedAt,
		LeftAt:       call.FinishedAt,
	}
	for _, option := range options {
		option(&participation)
	}
	return participation
}

func WithJoinedAt(joinedAt time.Time) ParticipationOption {
	return func(participation *model.CallParticipation) {
		participation.JoinedAt = model.Time(joinedAt)
	}
}

func WithLeftAt(leftAt time.Time) ParticipationOption {
	return func(participation *model.CallParticipation) {
		value := model.Time(leftAt)
		participation.LeftAt = &value
	}
}

func NewStartCallRequest(options ...StartCallRequestOption) *api.StartCallRequest {
	request := &api.StartCallRequest{}
	WithStartCallParticipants(NewSubscriber())(request)
	for _, option := range options {
		option(request)
	}
	return request
}

func WithStartCallParticipants(participants ...model.Subscriber) StartCallRequestOption {
	return func(request *api.StartCallRequest) {
		request.Participants = make([]api.StartCallParticipant, len(participants))
		for i, participant := range participants {
			request.Participants[i] = api.StartCallParticipant{
				Name:        participant.Name,
				PhoneNumber: model.PhoneNumberToString(participant.PhoneNumber),
			}
		}
	}
}

func NewCallParticipantRequest(subscriberID int) *api.CallParticipantRequest {
	return &api.CallParticipantRequest{SubscriberID: subscriberID}
}
