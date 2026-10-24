package record

type Abonent struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
}
