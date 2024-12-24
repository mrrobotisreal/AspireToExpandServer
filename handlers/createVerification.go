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

func CreateVerificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateVerificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := createVerification(req)
	if err != nil {
		http.Error(w, "Error creating verification object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createVerification(req types.CreateVerificationRequest) (types.CreateVerificationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.VerificationsCollection)
	_, err := collection.InsertOne(ctx, bson.M{
		"token":            req.Token,
		"email":            req.Email,
		"isVerified":       req.IsVerified,
		"registrationCode": req.RegistrationCode,
		"isRegistered":     req.IsRegistered,
	})
	if err != nil {
		fmt.Println("Error attempting to insert verification object into the database: ", err)
		return types.CreateVerificationResponse{
			IsCreated: false,
		}, err
	}

	return types.CreateVerificationResponse{
		IsCreated: true,
	}, nil
}
