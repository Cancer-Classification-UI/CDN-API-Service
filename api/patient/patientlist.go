package auth

import (
	"context"
	"net/http"
	_ "time"

	"ccu/api"
	mAPI "ccu/model/api/patientlist"
	mAPIDB "ccu/model/api"
	db "ccu/db"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
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
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In patient list handler -------------------------")
	patients := []mAPI.Patient{}

    coll_cdn := db.CLIENT.Database("cdn-api-db").Collection("patients")

    // Define username filter
    filter := bson.M{"username": username}

    // Retrieve data from MongoDB and store it in the results slice
    cursor, err := coll_cdn.Find(context.Background(), filter)
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var patient mAPIDB.CreatePatientDatabase
        if err := cursor.Decode(&patient); err != nil {
            log.Fatal(err)
			api.Respond(w, "Database Error", http.StatusBadRequest)
			return
        }
        patients = append(patients, mAPI.Patient{
			Id:          patient.ID,
			Name:        patient.Name,
			Samples:     int64(len(patient.Samples)),
		})
    }

    // Check for cursor errors after iteration
    if err := cursor.Err(); err != nil {
        log.Fatal(err)
		api.Respond(w, "Database Error", http.StatusBadRequest)
		return
    }

	response := mAPI.PatientList{Patients: patients}

	api.RespondOK(w, response)
}
