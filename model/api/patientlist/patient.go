package model

type Patient struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Samples     int64     `json:"samples"`
}
