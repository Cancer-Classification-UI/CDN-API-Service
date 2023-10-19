package auth

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"ccu/api"
	mAPI "ccu/model/api/patientdata"

	log "github.com/sirupsen/logrus"
)

// PostPatientData godoc
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
// @Router       /patient-data-no-auth [post]
func PostPatientData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.Form.Get("id")

	if id == "" { // What is the default int value?
		api.Respond(w, "Invalid ID Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In patient data handler -------------------------")
	response := mAPI.PatientData{
		Id:          id,
		DateCreated: time.Now(),
		Name:        RandomName(),
		Sex:         RandomSex(),
		DOB:         RandomDate(),
		Samples:     GenerateSamples(),
		Comments:    []string{"Lorem ipsum dolor sit amet. 1", "Lorem ipsum dolor sit amet. 2", "Lorem ipsum dolor sit amet. 3"},
	}

	api.RespondOK(w, response)
}

// Generate Samples with Images and Model output for samples parameter
func GenerateSamples() []mAPI.Sample {
	result := []mAPI.Sample{}
	num_samples := rand.Intn(5)

	for i := 0; i < num_samples; i++ {
		imageFile := images[rand.Intn(len(images))]
		imageData, _ := os.ReadFile(imageFile)

		// Encode the image data to Base64
		base64String := base64.StdEncoding.EncodeToString(imageData)

		result = append(result, mAPI.Sample{
			Image:           base64String,
			ModelPrediction: rand.Float64(),
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
