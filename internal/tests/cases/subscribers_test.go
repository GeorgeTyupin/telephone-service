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

type SubscribersSuite struct {
	tests.AcceptanceSuite
}

func TestSubscribers_Acceptance(t *testing.T) {
	t.Parallel()
	suite.Run(t, &SubscribersSuite{AcceptanceSuite: tests.AcceptanceSuite{
		DatabaseURL: databaseURL,
		Migrations:  os.DirFS("../../../migrations"),
	}})
}

func (s *SubscribersSuite) TestSubscribers_List_ReturnsOrderedSubscribers() {
	// arrange
	subscribers := []model.Subscriber{
		builders.NewSubscriber(
			builders.WithName("Пётр"),
			builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
		),
		builders.NewSubscriber(),
	}
	for i := range subscribers {
		subscriber := &subscribers[i]
		err := s.DB.QueryRow(s.Ctx,
			"insert into abonents (name, phone_number) values ($1, $2) returning id",
			subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
		).Scan(&subscriber.Id)
		s.Require().NoError(err)
	}

	// act
	response, err := s.APIClient.GetSubscribers(s.Ctx)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.GetSubscribersOKApplicationJSON{}, response)
	result := *response.(*api.GetSubscribersOKApplicationJSON)
	s.Require().Equal([]api.Subscriber{
		{ID: subscribers[0].Id, Name: subscribers[0].Name, PhoneNumber: model.PhoneNumberToString(subscribers[0].PhoneNumber)},
		{ID: subscribers[1].Id, Name: subscribers[1].Name, PhoneNumber: model.PhoneNumberToString(subscribers[1].PhoneNumber)},
	}, []api.Subscriber(result))
}

func (s *SubscribersSuite) TestSubscribers_Get_ReturnsSubscriber() {
	// arrange
	subscriber := builders.NewSubscriber()
	err := s.DB.QueryRow(s.Ctx,
		"insert into abonents (name, phone_number) values ($1, $2) returning id",
		subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
	).Scan(&subscriber.Id)
	s.Require().NoError(err)
	params := api.GetSubscriberParams{ID: subscriber.Id}

	// act
	response, err := s.APIClient.GetSubscriber(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.Subscriber{}, response)
	s.Require().Equal(&api.Subscriber{
		ID: subscriber.Id, Name: subscriber.Name, PhoneNumber: model.PhoneNumberToString(subscriber.PhoneNumber),
	}, response)
}

func (s *SubscribersSuite) TestSubscribers_Calls_ReturnsOnlySubscriberCalls() {
	// arrange
	subscriber := builders.NewSubscriber()
	err := s.DB.QueryRow(s.Ctx,
		"insert into abonents (name, phone_number) values ($1, $2) returning id",
		subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
	).Scan(&subscriber.Id)
	s.Require().NoError(err)
	other := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	active, err := fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall(), fixtures.WithParticipants(subscriber, other))
	s.Require().NoError(err)
	finished, err := fixtures.CreateCall(s.Ctx, s.DB,
		builders.NewCall(builders.WithFinished(time.Minute)), fixtures.WithParticipants(subscriber),
	)
	s.Require().NoError(err)
	_, err = fixtures.CreateCall(s.Ctx, s.DB,
		builders.NewCall(builders.WithFinished(2*time.Minute)), fixtures.WithParticipants(other),
	)
	s.Require().NoError(err)
	params := api.GetSubscriberCallsParams{ID: subscriber.Id}

	// act
	response, err := s.APIClient.GetSubscriberCalls(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.GetSubscriberCallsOKApplicationJSON{}, response)
	calls := *response.(*api.GetSubscriberCallsOKApplicationJSON)
	s.Require().Len(calls, 2)
	s.Require().Equal([]int{active.Id, finished.Id}, []int{calls[0].ID, calls[1].ID})
	s.Require().Equal(api.CallStatusActive, calls[0].Status)
	s.Require().Equal(api.CallStatusFinished, calls[1].Status)
	s.Require().Len(calls[0].Participants, 2)
	s.Require().Equal(active.Participants[0].Id, calls[0].Participants[0].ID)
	s.Require().Equal(active.Participants[1].Id, calls[0].Participants[1].ID)
	s.Require().Len(calls[1].Participants, 1)
	s.Require().Equal(subscriber.Id, calls[1].Participants[0].ID)
	s.Require().Equal("1m0s", calls[1].Duration)
}

func (s *SubscribersSuite) TestSubscribers_Stats_ReturnsFinishedCallAggregates() {
	// arrange
	subscriber := builders.NewSubscriber()
	err := s.DB.QueryRow(s.Ctx,
		"insert into abonents (name, phone_number) values ($1, $2) returning id",
		subscriber.Name, model.PhoneNumberToString(subscriber.PhoneNumber),
	).Scan(&subscriber.Id)
	s.Require().NoError(err)
	other := builders.NewSubscriber(
		builders.WithName("Пётр"),
		builders.WithPhoneNumber(model.PhoneNumber{CountryCode: 7, Number: 9997654321}),
	)
	_, err = fixtures.CreateCall(s.Ctx, s.DB,
		builders.NewCall(builders.WithFinished(time.Minute)), fixtures.WithParticipants(subscriber, other),
	)
	s.Require().NoError(err)
	_, err = fixtures.CreateCall(s.Ctx, s.DB,
		builders.NewCall(builders.WithFinished(3*time.Minute)), fixtures.WithParticipants(subscriber),
	)
	s.Require().NoError(err)
	_, err = fixtures.CreateCall(s.Ctx, s.DB, builders.NewCall(), fixtures.WithParticipants(subscriber))
	s.Require().NoError(err)
	_, err = fixtures.CreateCall(s.Ctx, s.DB,
		builders.NewCall(builders.WithFinished(4*time.Minute)), fixtures.WithParticipants(other),
	)
	s.Require().NoError(err)
	params := api.GetSubscriberStatsParams{ID: subscriber.Id}

	// act
	response, err := s.APIClient.GetSubscriberStats(s.Ctx, params)

	// assert
	s.Require().NoError(err)
	s.Require().IsType(&api.SubscriberStats{}, response)
	s.Require().Equal(&api.SubscriberStats{
		CallsCount: 2, TotalDuration: "4m0s", AverageDuration: "2m0s",
	}, response)
}
