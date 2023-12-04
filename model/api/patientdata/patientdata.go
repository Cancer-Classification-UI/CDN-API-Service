package model

type PatientData struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Sex         string        `json:"sex"`
	DOB         string        `json:"date_of_birth"`
	Samples     []string      `json:"samples"`
    Comments    []string      `json:"comments"`
}
