package domain

type Subscriber struct {
	Id          int
	Name        string
	PhoneNumber PhoneNumber
}

var WrongSubscriber = Subscriber{
	Id:          -1,
	Name:        "Wrong subscriber",
	PhoneNumber: WrongPhoneNumber,
}
