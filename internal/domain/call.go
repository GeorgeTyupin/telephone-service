package domain

type Call struct {
	Id               int
	FirstSubscriber  Subscriber
	SecondSubscriber Subscriber
	Duration         Duration
	CreatedAt        Time
}
