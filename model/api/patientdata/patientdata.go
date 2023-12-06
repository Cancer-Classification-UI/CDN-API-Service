package model

type PatientData struct {
	Name     string   `json:"name"`
	Sex      string   `json:"sex"`
	DOB      string   `json:"date_of_birth"`
	Samples  []Sample `json:"samples"`
	Comments []string `json:"comments"`
}
