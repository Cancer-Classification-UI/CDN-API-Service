package auth

import (
	"context"
	_ "encoding/base64"
	_ "fmt"
	"math/rand"
	"net/http"
	_ "os"
	_ "time"

	"ccu/api"
	mAPI "ccu/model/api"
	db "ccu/db"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

// PostCreatePatient godoc
// @Summary      Creates a table entry for a patient
// @Description  Creates a new patient with a randomly assigned ID
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        username       query string    true "username of the doctor"
// @Param        name           query string    true "name of the patient"
// @Param        sex            query string    true "sex of the patient"
// @Param        dob            query string    true "date of birth of the patient"
// @Success      200  {array}   mAPI.CreatePatientResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /create-patient [post]
func PostCreatePatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	name := r.URL.Query().Get("name")
	sex := r.URL.Query().Get("sex")
	dob := r.URL.Query().Get("dob")

	if username == "" {
		api.Respond(w, "Invalid username Parameter", http.StatusBadRequest)
		return
	}

	if name == "" {
		api.Respond(w, "Invalid name Parameter", http.StatusBadRequest)
		return
	}

	if sex == "" {
		api.Respond(w, "Invalid sex Parameter", http.StatusBadRequest)
		return
	}

	if dob == "" {
		api.Respond(w, "Invalid DOB Parameter", http.StatusBadRequest)
		return
	}

	result := AddPatient(username, name, sex, dob)

	if result == -1 {
		api.Respond(w, "Invalid username Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In create patient handler -------------------------")
	response := mAPI.CreatePatientResponse {
		Id:          result,
		Success:     true,
	}

	api.RespondOK(w, response)
}

// Insert Credentials Code Here
func AddPatient(Username string, Name string, Sex string, DOB string) int64 {
	//checks for a specific username in the login Database
	coll_login := db.CLIENT.Database("login-api-db").Collection("users")

	//search a database for a certain document
	var result bson.M
	err := coll_login.FindOne(context.TODO(), bson.D{{Key: "username", Value: Username}}).Decode(&result)

	// Invalid doctor username
	if err == mongo.ErrNoDocuments {
		return -1
	}
	if err != nil {
		panic(err)
		return -1
	}
	id := int64(rand.Int63n(1000000000))

	coll_cdn := db.CLIENT.Database("cdn-api-db").Collection("patients")

	row := mAPI.CreatePatientDatabase{
		ID:           id,
		Username:     Username,
		Name:         Name,
		Sex:          Sex,
		DOB:          DOB,
		Samples:      []string{},
		Comments:     []string{},
	}

	_, err = coll_cdn.InsertOne(context.TODO(), row)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return id
}