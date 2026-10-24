package domain

type Abonent struct {
	Id          int
	Name        string
	PhoneNumber PhoneNumber
}

var WrongAbonent = Abonent{
	Id:          -1,
	Name:        "Wrong abonent",
	PhoneNumber: WrongPhoneNumber,
}
