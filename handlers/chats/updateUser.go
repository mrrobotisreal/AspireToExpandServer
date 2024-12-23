package chatsHandlers

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

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := updateUser(req)
	if err != nil {
		http.Error(w, "Error updating chat user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateUser(req types.UpdateUserRequest) (types.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	if req.PreferredName != "" {
		update["preferredName"] = req.PreferredName
	}
	if req.ProfilePictureURL != "" {
		update["profilePictureUrl"] = req.ProfilePictureURL
	}

	collection := db.MongoClient.Database(db.DbName).Collection(db.UsersCollection)
	var userResult bson.M
	err := collection.FindOneAndUpdate(ctx, bson.M{
		"userId": req.UserId,
	}, bson.M{
		"$set": update,
	}).Decode(&userResult)
	if err != nil {
		fmt.Println("Error attempting to find and update user in the database:", err)
		return types.UpdateUserResponse{
			IsUpdated: false,
		}, err
	}

	return types.UpdateUserResponse{
		IsUpdated: true,
	}, nil
}
