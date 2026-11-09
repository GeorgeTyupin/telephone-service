package calls

import (
	"testing"
	"time"

	"telephone-service/internal/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCalls_ReturnsFilteredCalls(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	filter := model.CallsFilter{Status: model.CallStatusActive}
	expected := []model.Call{{Id: 42, Status: model.CallStatusActive}}
	repository.EXPECT().GetCalls(ctx, filter).Return(expected, nil)

	// act
	calls, err := service.GetCalls(ctx, filter)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, calls)
}

func TestGetCall_ReturnsCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	expected := model.Call{Id: 42, Status: model.CallStatusActive}
	repository.EXPECT().GetCall(ctx, 42).Return(expected, nil)

	// act
	call, err := service.GetCall(ctx, 42)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestGetSubscriberCalls_ReturnsExistingSubscriberCalls(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	subscriber := model.Subscriber{Id: 7, Name: "Иван"}
	expected := []model.Call{{Id: 42, Status: model.CallStatusActive, Participants: []model.Subscriber{subscriber}}}
	gomock.InOrder(
		repository.EXPECT().GetSubscriber(ctx, 7).Return(subscriber, nil),
		repository.EXPECT().GetSubscriberCalls(ctx, 7).Return(expected, nil),
	)

	// act
	calls, err := service.GetSubscriberCalls(ctx, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, calls)
}

func TestStartCall_CreatesActiveCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	input := model.Call{Participants: []model.Subscriber{
		{Name: "Иван", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}},
		{Name: "Пётр", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9997654321}},
	}}
	toCreate := model.Call{Participants: input.Participants, Status: model.CallStatusActive}
	expected := toCreate
	expected.Id = 42
	repository.EXPECT().StartCall(ctx, toCreate).Return(expected, nil)

	// act
	call, err := service.StartCall(ctx, input)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestStartCall_NullInName_ReturnsInvalidCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	input := model.Call{Participants: []model.Subscriber{
		{Name: "Иван\x00", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9991234567}},
	}}

	// act
	call, err := service.StartCall(ctx, input)

	// assert
	require.ErrorIs(t, err, model.ErrInvalidCall)
	require.Empty(t, call)
}

func TestConnectCall_AddsParticipant(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	changedAt := model.Time(time.Time(startedAt).Add(time.Minute))
	participant := model.Subscriber{Id: 7, Name: "Пётр", PhoneNumber: model.PhoneNumber{CountryCode: 7, Number: 9997654321}}
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: changedAt,
		Participants:   []model.Subscriber{{Id: 1, Name: "Иван"}},
		Participations: []model.CallParticipation{{Id: 11, SubscriberId: 1, JoinedAt: startedAt}},
	}
	update := model.CallUpdate{Call: current, ChangedAt: changedAt, ConnectedParticipant: &participant}
	expected := current
	expected.Version = 2
	expected.Participants = []model.Subscriber{current.Participants[0], participant}
	expected.Participations = []model.CallParticipation{
		current.Participations[0],
		{Id: 12, SubscriberId: 7, JoinedAt: changedAt},
	}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().GetSubscriber(ctx, 7).Return(participant, nil),
		repository.EXPECT().UpdateCall(ctx, update).Return(expected, nil),
	)

	// act
	call, err := service.ConnectCall(ctx, 42, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestDisconnectCall_DisconnectsParticipantAndKeepsCallActive(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	changedAt := model.Time(time.Time(startedAt).Add(time.Minute))
	participant := model.Subscriber{Id: 7, Name: "Пётр"}
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: changedAt,
		Participants: []model.Subscriber{{Id: 1, Name: "Иван"}, participant},
		Participations: []model.CallParticipation{
			{Id: 11, SubscriberId: 1, JoinedAt: startedAt},
			{Id: 12, SubscriberId: 7, JoinedAt: startedAt},
		},
	}
	update := model.CallUpdate{Call: current, ChangedAt: changedAt, DisconnectedParticipationIds: []int{12}}
	expected := current
	expected.Version = 2
	expected.Participations = []model.CallParticipation{
		current.Participations[0],
		{Id: 12, SubscriberId: 7, JoinedAt: startedAt, LeftAt: &changedAt},
	}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().GetSubscriber(ctx, 7).Return(participant, nil),
		repository.EXPECT().UpdateCall(ctx, update).Return(expected, nil),
	)

	// act
	call, err := service.DisconnectCall(ctx, 42, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestDisconnectCall_LastParticipant_FinishesCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	changedAt := model.Time(time.Time(startedAt).Add(time.Minute))
	participant := model.Subscriber{Id: 7, Name: "Иван"}
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: changedAt,
		Participants:   []model.Subscriber{participant},
		Participations: []model.CallParticipation{{Id: 11, SubscriberId: 7, JoinedAt: startedAt}},
	}
	finished := current
	finished.Status = model.CallStatusFinished
	finished.FinishedAt = &changedAt
	finished.Duration = model.Duration(time.Minute)
	update := model.CallUpdate{Call: finished, ChangedAt: changedAt, DisconnectedParticipationIds: []int{11}}
	expected := finished
	expected.Version = 2
	expected.Participations = []model.CallParticipation{{Id: 11, SubscriberId: 7, JoinedAt: startedAt, LeftAt: &changedAt}}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().GetSubscriber(ctx, 7).Return(participant, nil),
		repository.EXPECT().UpdateCall(ctx, update).Return(expected, nil),
	)

	// act
	call, err := service.DisconnectCall(ctx, 42, 7)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestFinishCall_ClosesActiveParticipations(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	changedAt := model.Time(time.Time(startedAt).Add(2 * time.Minute))
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: changedAt,
		Participants: []model.Subscriber{{Id: 1, Name: "Иван"}, {Id: 7, Name: "Пётр"}},
		Participations: []model.CallParticipation{
			{Id: 11, SubscriberId: 1, JoinedAt: startedAt},
			{Id: 12, SubscriberId: 7, JoinedAt: startedAt},
		},
	}
	finished := current
	finished.Status = model.CallStatusFinished
	finished.FinishedAt = &changedAt
	finished.Duration = model.Duration(2 * time.Minute)
	update := model.CallUpdate{Call: finished, ChangedAt: changedAt, DisconnectedParticipationIds: []int{11, 12}}
	expected := finished
	expected.Version = 2
	expected.Participations = []model.CallParticipation{
		{Id: 11, SubscriberId: 1, JoinedAt: startedAt, LeftAt: &changedAt},
		{Id: 12, SubscriberId: 7, JoinedAt: startedAt, LeftAt: &changedAt},
	}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().UpdateCall(ctx, update).Return(expected, nil),
	)

	// act
	call, err := service.FinishCall(ctx, 42)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestFinishCall_ClockMovedBack_UsesLatestDepartureTime(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	leftAt := model.Time(time.Time(startedAt).Add(time.Minute))
	observedAt := model.Time(time.Time(startedAt).Add(30 * time.Second))
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: observedAt,
		Participants: []model.Subscriber{{Id: 1, Name: "Иван"}, {Id: 7, Name: "Пётр"}},
		Participations: []model.CallParticipation{
			{Id: 11, SubscriberId: 1, JoinedAt: startedAt, LeftAt: &leftAt},
			{Id: 12, SubscriberId: 7, JoinedAt: startedAt},
		},
	}
	finished := current
	finished.Status = model.CallStatusFinished
	finished.FinishedAt = &leftAt
	finished.Duration = model.Duration(time.Minute)
	update := model.CallUpdate{Call: finished, ChangedAt: leftAt, DisconnectedParticipationIds: []int{12}}
	expected := finished
	expected.Version = 2
	expected.Participations = []model.CallParticipation{
		current.Participations[0],
		{Id: 12, SubscriberId: 7, JoinedAt: startedAt, LeftAt: &leftAt},
	}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().UpdateCall(ctx, update).Return(expected, nil),
	)

	// act
	call, err := service.FinishCall(ctx, 42)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

func TestFinishCall_ConcurrentUpdate_ReloadsCall(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := t.Context()
	repository := NewMockCallsRepository(gomock.NewController(t))
	service := NewService(repository)
	startedAt := model.Time(time.Date(2026, time.November, 8, 12, 0, 0, 0, time.UTC))
	firstObservedAt := model.Time(time.Time(startedAt).Add(time.Minute))
	secondObservedAt := model.Time(time.Time(startedAt).Add(2 * time.Minute))
	current := model.Call{
		Id: 42, Status: model.CallStatusActive, Version: 1, StartedAt: &startedAt, ObservedAt: firstObservedAt,
		Participants:   []model.Subscriber{{Id: 1, Name: "Иван"}},
		Participations: []model.CallParticipation{{Id: 11, SubscriberId: 1, JoinedAt: startedAt}},
	}
	firstFinished := current
	firstFinished.Status = model.CallStatusFinished
	firstFinished.FinishedAt = &firstObservedAt
	firstFinished.Duration = model.Duration(time.Minute)
	firstUpdate := model.CallUpdate{Call: firstFinished, ChangedAt: firstObservedAt, DisconnectedParticipationIds: []int{11}}
	reloaded := current
	reloaded.Version = 2
	reloaded.ObservedAt = secondObservedAt
	reloaded.Participants = []model.Subscriber{current.Participants[0], {Id: 7, Name: "Пётр"}}
	reloaded.Participations = []model.CallParticipation{
		current.Participations[0],
		{Id: 12, SubscriberId: 7, JoinedAt: firstObservedAt},
	}
	secondFinished := reloaded
	secondFinished.Status = model.CallStatusFinished
	secondFinished.FinishedAt = &secondObservedAt
	secondFinished.Duration = model.Duration(2 * time.Minute)
	secondUpdate := model.CallUpdate{Call: secondFinished, ChangedAt: secondObservedAt, DisconnectedParticipationIds: []int{11, 12}}
	expected := secondFinished
	expected.Version = 3
	expected.Participations = []model.CallParticipation{
		{Id: 11, SubscriberId: 1, JoinedAt: startedAt, LeftAt: &secondObservedAt},
		{Id: 12, SubscriberId: 7, JoinedAt: firstObservedAt, LeftAt: &secondObservedAt},
	}
	gomock.InOrder(
		repository.EXPECT().GetCall(ctx, 42).Return(current, nil),
		repository.EXPECT().UpdateCall(ctx, firstUpdate).Return(model.Call{}, model.ErrCallChanged),
		repository.EXPECT().GetCall(ctx, 42).Return(reloaded, nil),
		repository.EXPECT().UpdateCall(ctx, secondUpdate).Return(expected, nil),
	)

	// act
	call, err := service.FinishCall(ctx, 42)

	// assert
	require.NoError(t, err)
	require.Equal(t, expected, call)
}
