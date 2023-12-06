package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"ccu/api"
	"ccu/db"
	mAPI "ccu/model/api/patientdata"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPatientData godoc
// @Summary      Provides data about a patient
// @Description  Checks for a matching id in the databse and returns data for that id
// @Tags         Patient
// @Accept       json
// @Produce      json
// @Param        ref_id      query string  true "reference id of the samples"
// @Param        patient_id  query string  true "id of the patient"
// @Success      200  {array}   mAPI.PatientData
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patient-data [get]
func GetPatientData(w http.ResponseWriter, r *http.Request) {
	log.Info("In patient data handler -------------------------")
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ref_id := r.URL.Query().Get("ref_id")
	if ref_id == "" { // What is the default int value?
		api.Respond(w, "Invalid Reference ID Parameter", http.StatusBadRequest)
		return
	}

	patient_id := r.URL.Query().Get("patient_id")
	if patient_id == "" { // What is the default int value?
		api.Respond(w, "Invalid Patient ID Parameter", http.StatusBadRequest)
		return
	}

	// Magically grab patient data and also the reference id samples and comments
	coll_data := db.CLIENT.Database("cdn-api-db").Collection("patientdata")
	coll_samples := db.CLIENT.Database("cdn-api-db").Collection("patientsamples")

	// Search a db for patient data
	var result_data bson.M
	err := coll_data.FindOne(context.Background(), bson.D{{Key: "patient_id", Value: patient_id}}).Decode(&result_data)
	if err == mongo.ErrNoDocuments {
		log.Warning("No patient data document was found for patient_id: ", patient_id)
		api.RespondOK(w, mAPI.PatientData{Name: "", Sex: "", DOB: "", Samples: []mAPI.Sample{}, Comments: []string{}})
		return
	} else if err != nil {
		log.Error("Error while getting patient data", err)
		api.Respond(w, "Error while getting patient data", http.StatusInternalServerError)
		return
	}

	// Search a db for patient data
	var result_samples bson.M
	err = coll_samples.FindOne(context.Background(), bson.D{{Key: "ref_id", Value: ref_id}}).Decode(&result_samples)
	if err == mongo.ErrNoDocuments {
		log.Warning("No samples document was found for ref_id: ", patient_id)
		api.RespondOK(w, mAPI.PatientData{Name: "", Sex: "", DOB: "", Samples: []mAPI.Sample{}, Comments: []string{}})
		return
	} else if err != nil {
		log.Error("Error while getting reference id samples", err)
		api.Respond(w, "Error while getting reference id samples", http.StatusInternalServerError)
		return
	}

	// Iterate through samples, add them to a mAPI.Sample struct
	samples := []mAPI.Sample{}
	for _, sample := range result_samples["samples"].(bson.A) {
		samples = append(samples, mAPI.Sample{Image: sample.(string)})
	}

	// Iterate through comments, add them to a string array
	comments := []string{}
	for _, comment := range result_samples["comments"].(bson.A) {
		comments = append(comments, comment.(string))
	}

	response := mAPI.PatientData{
		Name:     result_data["name"].(string),
		Sex:      result_data["sex"].(string),
		DOB:      result_data["dob"].(string),
		Samples:  samples,
		Comments: comments,
	}

	api.RespondOK(w, response)
}

// Generate Samples with Images and Model output for samples parameter
func GenerateSamples() []mAPI.Sample {
	result := []mAPI.Sample{}
	num_samples := rand.Intn(5) + 1

	for i := 0; i < num_samples; i++ {
		imageFile := images[rand.Intn(len(images))]
		imageData, _ := os.ReadFile(imageFile)

		// Encode the image data to Base64
		base64String := base64.StdEncoding.EncodeToString(imageData)

		result = append(result, mAPI.Sample{
			Image: base64String,
		})
	}
	return result
}

// Generate random value for sex parameter
func RandomSex() string {
	result := "M"
	if rand.Intn(2) == 1 {
		result = "F"
	}
	return result
}

var (
	first_names = []string{
		"Alice", "Bob", "Charlie", "David", "Eva", "Frank", "Grace", "Hannah",
		"Ivy", "Jack", "Katie", "Leo", "Mia", "Nathan", "Olivia", "Steve",
	}
	last_names = []string{
		"Smith", "Johnson", "Davis", "Lee", "Garcia", "Wilson", "Taylor",
		"Martin", "Anderson", "Clark", "Hall", "Moore", "Young", "Walker",
	}
	images = []string{
		"sample_images/ISIC_0034525.jpg", "sample_images/ISIC_0034526.jpg", "sample_images/ISIC_0034527.jpg", "sample_images/ISIC_0034528.jpg", "sample_images/ISIC_0034529.jpg",
	}
)

// Generate random name for name parameter
func RandomName() string {
	// Select random first and last names
	first_name := first_names[rand.Intn(len(first_names))]
	last_name := last_names[rand.Intn(len(last_names))]

	return fmt.Sprintf("%s %s", first_name, last_name)
}

// Generate random date for date parameter
func RandomDate() string {
	year := rand.Intn(93) + 1930

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
