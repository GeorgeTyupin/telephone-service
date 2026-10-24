package dto

type Call struct {
	Id            int    `json:"id"`
	FirstAbonent  string `json:"first_abonent"`
	SecondAbonent string `json:"second_abonent"`
	Duration      string `json:"duration"`
	CreatedAt     string `json:"created_at"`
}
