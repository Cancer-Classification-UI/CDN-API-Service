package model

import "time"

type Patient struct {
	Id          int64     `json:"id"`
	DateCreated time.Time `json:"date_created"`
	Name        string    `json:"name"`
	Samples     int64     `json:"samples"`
	Date        string    `json:"date"`
}
