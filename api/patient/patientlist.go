package auth

import (
	"math/rand"
	"net/http"
	"time"
	"fmt"

	"ccu/api"
	mAPI "ccu/model/api/patientlist"

	log "github.com/sirupsen/logrus"
)

// PostPatientList godoc
// @Summary      Retrieves a list of patients for the specified doctor username
// @Description  Finds username in database and retrieves all patients for that user
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        username       query string    true "username of the doctor"
// @Success      200  {array}   mAPI.PatientList
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patient-list-no-auth [post]
func PostPatientList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.Form.Get("username")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In patient list handler -------------------------")
	patients := []mAPI.Patient{}

	num_patients := rand.Intn(10)
	
	for i := 0; i < num_patients; i++ {
		patients = append(patients, mAPI.Patient{
			Id:          int64(rand.Intn(1000)),
			DateCreated: time.Now(),
			Name:        RandomName(),
			Samples:     int64(rand.Intn(10)),
			Date:        RandomRegistrationDate(),
		})
	}

	response := mAPI.PatientList{Patients: patients}

	api.RespondOK(w, response)
}

func RandomRegistrationDate() string {
	year := rand.Intn(22) + 2000

	month := rand.Intn(12) + 1

	days_in_month := 31

	if month == 4 || month == 6 || month == 9 || month == 11 {
		days_in_month = 30
	} else if month == 2 {
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			days_in_month = 29
		} else {
			days_in_month = 28
		}
	}

	day := rand.Intn(days_in_month) + 1

	date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

	return date
}