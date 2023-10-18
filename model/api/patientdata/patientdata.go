package model

import (
	"time"
)

type PatientData struct {
	Id          string        `json:"id"`
	DateCreated time.Time     `json:"date_created"`
	Name        string        `json:"name"`
	Sex         string        `json:"sex"`
	DOB         string        `json:"date_of_birth"`
	Samples     []Sample `json:"samples"`
    Comments    []string      `json:"comments"`
}
