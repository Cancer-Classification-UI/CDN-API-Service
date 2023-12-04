package auth

import (
	"context"
	"net/http"
	"strconv"

	"ccu/api"
	mAPI "ccu/model/api/patientdata"
	db "ccu/db"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetPatientData godoc
// @Summary      Provides data about a patient
// @Description  Checks for a matching id in the databse and returns data for that id
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        id             query int    true "id of the patient"
// @Success      200  {array}   mAPI.PatientData
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patient-data [get]
func GetPatientData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" { // What is the default int value?
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	// Checks for a specific username in the login Database
	coll := db.CLIENT.Database("cdn-api-db").Collection("patients")

	// Search a database for a certain document
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id_int}}).Decode(&result)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	// Convert primitive.A to []string
	var samples []string
	for _, v := range result["samples"].(primitive.A) {
		if s, ok := v.(string); ok {
			samples = append(samples, s)
		} else {
			// Not a string
		}
	}

	var comments []string
	for _, v := range result["comments"].(primitive.A) {
		if s, ok := v.(string); ok {
			comments = append(comments, s)
		} else {
			// Not a string
		}
	}

	log.Info("In patient data handler -------------------------")
	response := mAPI.PatientData{
		Id:          result["id"].(int64),
		Name:        result["name"].(string),
		Sex:         result["sex"].(string),
		DOB:         result["dob"].(string),
		Samples:     samples,
		Comments:    comments,
	}

	api.RespondOK(w, response)
}