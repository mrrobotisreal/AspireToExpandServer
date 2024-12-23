package chatsHandlers

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"net/http"
	"time"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := createUser(req)
	if err != nil {
		http.Error(w, "Error creating the user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUser(req types.CreateUserRequest) (types.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.UsersCollection)
	_, err := collection.InsertOne(ctx, bson.M{
		"userId":            req.UserId,
		"userType":          req.UserType,
		"preferredName":     req.PreferredName,
		"firstName":         req.FirstName,
		"lastName":          req.LastName,
		"profilePictureUrl": req.ProfilePictureURL,
	})
	if err != nil {
		fmt.Println("Error inserting user into the database: ", err)
		return types.CreateUserResponse{
			IsCreated: false,
		}, err
	}

	return types.CreateUserResponse{
		IsCreated: true,
	}, nil
}
