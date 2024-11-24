package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"net/http"
	"time"
)

func CreateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateRegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("CreateRegistration request incoming...")
	fmt.Println(req.RegistrationCode)
	fmt.Println(req.FirstName)
	fmt.Println(req.LastName)
	fmt.Println(req.EmailAddress)
	fmt.Println(req.IsValid)

	// Validate registration code
	result, err := createRegistration(req)
	if err != nil {
		fmt.Println("Error Creating registration result...")
		http.Error(w, "Error Creating registration result...", http.StatusInternalServerError)
		return
	}

	response := types.CreateRegistrationResponse{
		RegistrationCode: result.RegistrationCode,
		FirstName:        result.FirstName,
		LastName:         result.LastName,
		EmailAddress:     result.EmailAddress,
		IsValid:          result.IsValid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createRegistration(req types.CreateRegistrationRequest) (types.CreateRegistrationResponse, error) {
	response := types.CreateRegistrationResponse{
		RegistrationCode: "abcd-1234",
		FirstName:        "Bob",
		LastName:         "Smith",
		EmailAddress:     "bob@gmail.com",
		IsValid:          true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.RegistrationCollection)

	_, err := collection.InsertOne(ctx, bson.M{
		"registrationcode": req.RegistrationCode,
		"firstname":        req.FirstName,
		"lastname":         req.LastName,
		"emailaddress":     req.EmailAddress,
	})
	if err != nil {
		fmt.Println("Error inserting registration code...")
		fmt.Println("Error is: " + err.Error())
		return response, err
	}

	fmt.Println("Successfully inserted new registration...")
	response.RegistrationCode = req.RegistrationCode
	response.FirstName = req.FirstName
	response.LastName = req.LastName
	response.EmailAddress = req.EmailAddress
	response.IsValid = req.IsValid
	return response, nil
}

// HTTP handler for validating registration code
func ValidateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON body
	var req types.RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Registration request incoming...")
	fmt.Println(req.RegistrationCode)

	// Validate registration code
	result, err := validateRegistrationCode(req.RegistrationCode)
	if err != nil {
		fmt.Println("Error validating registration result...")
		http.Error(w, "Error validating registration result...", http.StatusInternalServerError)
		return
	}

	response := types.ValidationResponse{
		IsValid:          result.IsValid,
		FirstName:        result.Result.FirstName,
		LastName:         result.Result.LastName,
		EmailAddress:     result.Result.EmailAddress,
		RegistrationCode: result.Result.RegistrationCode,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Function to validate the registration code
func validateRegistrationCode(code string) (types.RegistrationValidationResult, error) {
	// For testing purposes, auto-validate the test code
	if code == "test-registration-code" {
		return types.RegistrationValidationResult{
			IsValid: false,
			Result: types.Registration{
				FirstName:        "Student",
				LastName:         "Pupil",
				EmailAddress:     "student@alina.edu",
				RegistrationCode: "Student-12345678",
			},
		}, nil
	}

	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.RegistrationCollection)

	// Query MongoDB
	var result types.Registration
	err := collection.FindOne(ctx, bson.M{"registrationcode": code}).Decode(&result)
	if err != nil {
		fmt.Println("Error finding registration... returning error.")
		fmt.Println("Error: " + err.Error())
		return types.RegistrationValidationResult{
			IsValid: false,
			Result: types.Registration{
				FirstName:        "",
				LastName:         "",
				EmailAddress:     "",
				RegistrationCode: "",
			},
		}, err
	}

	// If no error, code exists, so it's valid
	return types.RegistrationValidationResult{
		IsValid: true,
		Result:  result,
	}, nil
}
