package model

type PatientListEntry struct {
	Ref_ID     string `json:"ref_id"`
	Name       string `json:"name"`
	Patient_ID string `json:"patient_id"`
	Samples    string `json:"samples"`
	Date       string `json:"date"`
}

type PatientList struct {
	Patients []PatientListEntry `json:"patients"`
}
