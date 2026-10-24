package domain

type Call struct {
	Id            int
	FirstAbonent  Abonent
	SecondAbonent Abonent
	Duration      Duration
	CreatedAt     Time
}
