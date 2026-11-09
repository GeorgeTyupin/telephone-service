//go:build integration

package cases

import (
	"os"
	"testing"
	"time"

	api "telephone-service/codegen/oas"
	"telephone-service/internal/model"
	"telephone-service/internal/tests"
	"telephone-service/internal/tests/builders"
	"telephone-service/internal/tests/fixtures"

	"github.com/stretchr/testify/suite"
)

type CallsSuite struct {
	tests.AcceptanceSuite
}

func TestCalls_Acceptance(t *testing.T) {
	t.Parallel()
	suite.Run(t, &CallsSuite{AcceptanceSuite: tests.AcceptanceSuite{
		DatabaseURL: databaseURL,
		Migrations:  os.DirFS("../../../migrations"),
	}})
}

func (s *CallsSuite) TestCalls_Start_CreatesActiveCall() {
	// arrange
	first := builders.NewSubscriber()
	second := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	request := builders.NewStartCallRequest(builders.WithStartCallParticipants(first, second))
	expectedParticipants := map[string]model.Subscriber{
		request.Participants[0].PhoneNumber: first,
		request.Participants[1].PhoneNumber: second,
	}

	// act
	response, err := s.APIClient.StartCall(s.Ctx, request)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Call{}, response)
	call := response.(*api.Call)
	s.Require().Positive(call.ID)
	s.Require().Equal(api.CallStatusActive, call.Status)
	s.Require().True(call.StartedAt.Set)
	s.Require().False(call.FinishedAt.Set)
	s.Require().Len(call.Participants, 2)
	s.Require().Len(call.Participations, 2)
	s.Require().Less(call.Participants[0].ID, call.Participants[1].ID)
	duration, err := time.ParseDuration(call.Duration)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(duration, time.Duration(0))

	stored, err := s.Repository.GetCall(s.Ctx, call.ID)
	s.Require().NoError(err)
	s.Require().Equal(model.CallStatusActive, stored.Status)
	s.Require().Equal(int64(1), stored.Version)
	s.Require().Len(stored.Participations, 2)
	s.Require().Len(stored.Participants, 2)
	s.Require().Equal(time.Time(stored.CreatedAt).String(), call.CreatedAt)
	s.Require().Equal(time.Time(*stored.StartedAt).Truncate(time.Second), call.StartedAt.Value)
	for i, participant := range call.Participants {
		expected, exists := expectedParticipants[participant.PhoneNumber]
		s.Require().True(exists)
		expected.Id = participant.ID
		s.Require().Equal(expected.Name, participant.Name)
		s.Require().Equal(expected, stored.Participants[i])
		s.Require().Positive(call.Participations[i].ID)
		s.Require().Equal(participant.ID, call.Participations[i].SubscriberID)
		s.Require().Equal(call.StartedAt.Value, call.Participations[i].JoinedAt)
		s.Require().False(call.Participations[i].LeftAt.Set)
		s.Require().Equal(stored.Participations[i].Id, call.Participations[i].ID)
		s.Require().Equal(*stored.StartedAt, stored.Participations[i].JoinedAt)
		s.Require().Nil(stored.Participations[i].LeftAt)
	}
}

func (s *CallsSuite) TestCalls_List_FiltersActiveCalls() {
	// arrange
	first, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall())
	s.Require().NoError(err)
	_, err = fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall(builders.WithFinished(time.Minute)))
	s.Require().NoError(err)
	second, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall())
	s.Require().NoError(err)
	params := api.GetCallsParams{Status: api.NewOptCallStatus(api.CallStatusActive)}

	// act
	response, err := s.APIClient.GetCalls(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.GetCallsOKApplicationJSON{}, response)
	calls := *response.(*api.GetCallsOKApplicationJSON)
	s.Require().Len(calls, 2)
	s.Require().Equal([]int{first.Id, second.Id}, []int{calls[0].ID, calls[1].ID})
	for _, call := range calls {
		s.Require().Equal(api.CallStatusActive, call.Status)
		s.Require().Len(call.Participants, 1)
		s.Require().Len(call.Participations, 1)
	}
}

func (s *CallsSuite) TestCalls_Get_ReturnsParticipantsAndHistory() {
	// arrange
	subscriber := builders.NewSubscriber()
	err := s.DB.QueryRow(s.Ctx,
		"insert into abonents (name, phone_number) values ($1, $2) returning id",
		subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
	).Scan(&subscriber.Id)
	s.Require().NoError(err)
	input := builders.NewCall()
	leftAt := time.Time(*input.StartedAt).Add(time.Minute)
	joinedAt := leftAt.Add(time.Minute)
	stored, err := fixtures.CreateCall(s.Ctx, s.DB, input,
		fixtures.WithParticipants(subscriber),
		fixtures.WithParticipations(
			builders.NewParticipation(input, subscriber.Id, builders.WithLeftAt(leftAt)),
			builders.NewParticipation(input, subscriber.Id, builders.WithJoinedAt(joinedAt)),
		),
	)
	s.Require().NoError(err)
	params := api.GetCallParams{ID: stored.Id}

	// act
	response, err := s.APIClient.GetCall(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Call{}, response)
	call := response.(*api.Call)
	s.Require().Equal(stored.Id, call.ID)
	s.Require().Equal(api.CallStatusActive, call.Status)
	s.Require().Equal(time.Time(stored.CreatedAt).String(), call.CreatedAt)
	s.Require().Equal(api.NewOptDateTime(time.Time(*stored.StartedAt).Truncate(time.Second)), call.StartedAt)
	s.Require().False(call.FinishedAt.Set)
	s.Require().Equal([]api.Subscriber{{
		ID: subscriber.Id, Name: subscriber.Name, PhoneNumber: model.PhoneNumberToString(subscriber.PhoneNumber),
	}}, call.Participants)
	s.Require().Equal([]api.CallParticipation{
		{ID: stored.Participations[0].Id, SubscriberID: subscriber.Id,
			JoinedAt: time.Time(*stored.StartedAt).Truncate(time.Second), LeftAt: api.NewOptDateTime(leftAt.Truncate(time.Second))},
		{ID: stored.Participations[1].Id, SubscriberID: subscriber.Id, JoinedAt: joinedAt.Truncate(time.Second)},
	}, call.Participations)
}

func (s *CallsSuite) TestCalls_Connect_PersistsParticipation() {
	// arrange
	initial, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall())
	s.Require().NoError(err)
	subscriber := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	err = s.DB.QueryRow(s.Ctx,
		"insert into abonents (name, phone_number) values ($1, $2) returning id",
		subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
	).Scan(&subscriber.Id)
	s.Require().NoError(err)
	request := builders.NewCallParticipantRequest(subscriber.Id)
	params := api.ConnectCallParams{ID: initial.Id}

	// act
	response, err := s.APIClient.ConnectCall(s.Ctx, request, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Call{}, response)
	call := response.(*api.Call)
	s.Require().Equal(initial.Id, call.ID)
	s.Require().Equal(api.CallStatusActive, call.Status)
	s.Require().Len(call.Participants, 2)
	s.Require().Equal(subscriber.Id, call.Participants[1].ID)
	s.Require().Equal(subscriber.Name, call.Participants[1].Name)
	s.Require().Equal(model.PhoneNumberToString(subscriber.PhoneNumber), call.Participants[1].PhoneNumber)
	s.Require().Len(call.Participations, 2)
	s.Require().Equal(initial.Participations[0].Id, call.Participations[0].ID)
	s.Require().Equal(subscriber.Id, call.Participations[1].SubscriberID)
	s.Require().False(call.Participations[1].JoinedAt.Before(time.Time(*initial.StartedAt)))
	s.Require().False(call.Participations[1].LeftAt.Set)

	stored, err := s.Repository.GetCall(s.Ctx, initial.Id)
	s.Require().NoError(err)
	s.Require().Equal(initial.Version+1, stored.Version)
	s.Require().Equal([]model.Subscriber{initial.Participants[0], subscriber}, stored.Participants)
	s.Require().Len(stored.Participations, 2)
	s.Require().Equal(call.Participations[1].ID, stored.Participations[1].Id)
	s.Require().Equal(call.Participations[1].JoinedAt, time.Time(stored.Participations[1].JoinedAt).Truncate(time.Second))
	s.Require().Nil(stored.Participations[1].LeftAt)
}

func (s *CallsSuite) TestCalls_Disconnect_ClosesParticipation() {
	// arrange
	first := builders.NewSubscriber()
	second := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	initial, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall(), fixtures.WithParticipants(first, second))
	s.Require().NoError(err)
	request := builders.NewCallParticipantRequest(initial.Participants[0].Id)
	params := api.DisconnectCallParams{ID: initial.Id}

	// act
	response, err := s.APIClient.DisconnectCall(s.Ctx, request, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Call{}, response)
	call := response.(*api.Call)
	s.Require().Equal(initial.Id, call.ID)
	s.Require().Equal(api.CallStatusActive, call.Status)
	s.Require().False(call.FinishedAt.Set)
	s.Require().Len(call.Participants, 2)
	s.Require().Len(call.Participations, 2)
	s.Require().Equal(initial.Participations[0].Id, call.Participations[0].ID)
	s.Require().Equal(request.SubscriberID, call.Participations[0].SubscriberID)
	s.Require().True(call.Participations[0].LeftAt.Set)
	s.Require().False(call.Participations[0].LeftAt.Value.Before(call.Participations[0].JoinedAt))
	s.Require().False(call.Participations[1].LeftAt.Set)

	stored, err := s.Repository.GetCall(s.Ctx, initial.Id)
	s.Require().NoError(err)
	s.Require().Equal(initial.Version+1, stored.Version)
	s.Require().Equal(model.CallStatusActive, stored.Status)
	s.Require().Nil(stored.FinishedAt)
	s.Require().Equal(initial.Participants, stored.Participants)
	s.Require().Len(stored.Participations, 2)
	s.Require().NotNil(stored.Participations[0].LeftAt)
	s.Require().Equal(call.Participations[0].LeftAt.Value, time.Time(*stored.Participations[0].LeftAt).Truncate(time.Second))
	s.Require().Nil(stored.Participations[1].LeftAt)
}

func (s *CallsSuite) TestCalls_Finish_PersistsFinishedCall() {
	// arrange
	first := builders.NewSubscriber()
	second := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	initial, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall(), fixtures.WithParticipants(first, second))
	s.Require().NoError(err)
	params := api.FinishCallParams{ID: initial.Id}

	// act
	response, err := s.APIClient.FinishCall(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Call{}, response)
	call := response.(*api.Call)
	s.Require().Equal(initial.Id, call.ID)
	s.Require().Equal(api.CallStatusFinished, call.Status)
	s.Require().Equal(api.NewOptDateTime(time.Time(*initial.StartedAt).Truncate(time.Second)), call.StartedAt)
	s.Require().True(call.FinishedAt.Set)
	duration, err := time.ParseDuration(call.Duration)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(duration, time.Duration(0))
	s.Require().Len(call.Participants, 2)
	s.Require().Len(call.Participations, 2)

	stored, err := s.Repository.GetCall(s.Ctx, initial.Id)
	s.Require().NoError(err)
	s.Require().Equal(initial.Version+1, stored.Version)
	s.Require().Equal(model.CallStatusFinished, stored.Status)
	s.Require().Equal(model.Duration(duration), stored.Duration)
	s.Require().NotNil(stored.FinishedAt)
	s.Require().Equal(time.Time(*stored.FinishedAt).Sub(time.Time(*stored.StartedAt)), duration)
	s.Require().Equal(call.FinishedAt.Value, time.Time(*stored.FinishedAt).Truncate(time.Second))
	s.Require().Equal(initial.Participants, stored.Participants)
	s.Require().Len(stored.Participations, 2)
	for i, participation := range call.Participations {
		s.Require().Equal(initial.Participations[i].Id, participation.ID)
		s.Require().Equal(initial.Participations[i].SubscriberId, participation.SubscriberID)
		s.Require().Equal(call.FinishedAt, participation.LeftAt)
		s.Require().NotNil(stored.Participations[i].LeftAt)
		s.Require().Equal(*stored.FinishedAt, *stored.Participations[i].LeftAt)
		s.Require().Equal(call.FinishedAt.Value, time.Time(*stored.Participations[i].LeftAt).Truncate(time.Second))
	}
}
