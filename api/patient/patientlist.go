package auth

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"ccu/api"
	"ccu/db"
	mAPI "ccu/model/api/patientlist"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPatientList godoc
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
// @Router       /patient-list [get]
func GetPatientList(w http.ResponseWriter, r *http.Request) {
	log.Info("In patient list handler -------------------------")
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	coll := db.CLIENT.Database("cdn-api-db").Collection("patientlist")

	// Search a database for a certain document
	cursor, err := coll.Find(context.Background(), bson.D{{Key: "username", Value: username}})
	if err == mongo.ErrNoDocuments {
		log.Debug("No documents bound to username", username)
	} else if err != nil {
		log.Error("Error while getting patientlist for username", username, err)
		api.Respond(w, "Error while getting patientlist for username", http.StatusInternalServerError)
		return
	}

	// Iterate through the cursor and decode each document
	patients := []mAPI.PatientListEntry{}

	var result bson.M
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&result)
		if err != nil {
			log.Warning("Error while decoding a patientlist entry and skipping for username", username, err)
		}

		// Add the patient to the list
		patients = append(patients, mAPI.PatientListEntry{
			Ref_ID:     result["ref_id"].(string),
			Name:       result["name"].(string),
			Patient_ID: result["patient_id"].(string),
			Samples:    result["samples"].(string),
			Date:       result["date"].(string),
		})

		log.Debug("Found document for username", username)
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
