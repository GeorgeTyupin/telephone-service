package model

type Call struct {
	Id           int
	Participants []Subscriber
	Duration     Duration
	CreatedAt    Time
}
