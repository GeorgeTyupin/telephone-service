package record

import "time"

type Call struct {
	Id            int64         `db:"id"`
	FirstAbonent  int64         `db:"first_name"`
	SecondAbonent int64         `db:"second_name"`
	Duration      time.Duration `db:"description"`
	CreatedAt     time.Time     `db:"created_at"`
}
