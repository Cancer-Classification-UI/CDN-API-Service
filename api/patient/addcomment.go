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

// PostAddComment godoc
// @Summary      Adds a comment to a patient
// @Description  Finds a matching id in the database and adds a comment to the patient
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        id             query int    true "id of the patient"
// @Param        comment        query string true "comment to add"
// @Success      200  {array}   mAPI.AddCommentResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /add-comment [post]
func PostAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	comment := r.URL.Query().Get("comment")

	if id == "" {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	if comment == "" {
		api.Respond(w, "Invalid Comment Parameter", http.StatusBadRequest)
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
			"comments": comment,
		},
	}

	// Update the document in the collection
	_, err = coll.UpdateOne(context.Background(), bson.D{{Key: "id", Value: id_int}}, update)
	if err != nil {
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In add comment handler -------------------------")
	response := mAPI.AddCommentResponse{
		Success:          true,
	}

	api.RespondOK(w, response)
}