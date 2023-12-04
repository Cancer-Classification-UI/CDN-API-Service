package auth

import (
	"context"
	"net/http"
	"strconv"

	"ccu/api"
	mAPI "ccu/model/api"
	db "ccu/db"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// PostAddSample godoc
// @Summary      Adds a base64 encoded sample image to a patient
// @Description  Finds a matching id in the database and adds a sample to the patient
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        id            query int    true "id of the patient"
// @Param        sample        query string true "sample to add"
// @Success      200  {array}   mAPI.AddSampleResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /add-sample [post]
func PostAddSample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	sample := r.URL.Query().Get("sample")

	if id == "" {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	if sample == "" {
		api.Respond(w, "Invalid Sample Parameter", http.StatusBadRequest)
		return
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	// Checks for a specific username in the cdn Database
	coll := db.CLIENT.Database("cdn-api-db").Collection("patients")

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id_int}}).Decode(&result)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$push": bson.M{
			"samples": sample,
		},
	}

	// Update the document in the collection
	_, err = coll.UpdateOne(context.Background(), bson.D{{Key: "id", Value: id_int}}, update)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In add sample handler -------------------------")
	response := mAPI.AddSampleResponse{
		Success:          true,
	}

	api.RespondOK(w, response)
}