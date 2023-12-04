package model

type CreatePatientResponse struct {
	Id          int64     `json:"id"`
	Success     bool      `json:"success"`
}

type CreatePatientDatabase struct {
	ID           int64    `bson:"id"`
	Username     string   `bson:"username"`
	Name         string   `bson:"name"`
	Sex          string   `bson:"sex"`
	DOB          string   `bson:"dob"`
	Samples      []string `bson:"samples"`
	Comments     []string `bson:"comments"`
}