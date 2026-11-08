package model

type CallStatus string

const (
	CallStatusActive   CallStatus = "active"
	CallStatusFinished CallStatus = "finished"
)

type Call struct {
	Id             int
	Participants   []Subscriber
	Duration       Duration
	CreatedAt      Time
	Status         CallStatus
	StartedAt      *Time
	FinishedAt     *Time
	Participations []CallParticipation
	Version        int64
	ObservedAt     Time
}

type CallParticipation struct {
	Id           int
	SubscriberId int
	JoinedAt     Time
	LeftAt       *Time
}

type CallsFilter struct {
	Status CallStatus
}

type CallUpdate struct {
	Call                         Call
	ChangedAt                    Time
	ConnectedParticipant         *Subscriber
	DisconnectedParticipationIds []int
}
