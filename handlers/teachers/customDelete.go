package teachersHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"net/http"
	"time"
)

type CustomDeleteRequest struct {
	TeacherID         string `json:"teacherID"`
	PreferredName     string `json:"preferred_name"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	EmailAddress      string `json:"email_address"`
}

type CustomDeleteResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

func CustomDeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var req CustomDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	response, err := customDeleteTeacher(req)
	if err != nil {
		http.Error(w, "Error deleting teacher via custom delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func customDeleteTeacher(req CustomDeleteRequest) (CustomDeleteResponse, error) {
	update := bson.M{}

	if req.TeacherID != "" {
		update["teacherID"] = req.TeacherID
	}
	if req.ProfilePictureURL != "" {
		update["profilepictureurl"] = req.ProfilePictureURL
	}
	if req.PreferredName != "" {
		update["preferredname"] = req.PreferredName
	}
	if req.EmailAddress != "" {
		update["emailaddress"] = req.EmailAddress
	}
	if req.FirstName != "" {
		update["firstname"] = req.FirstName
	}
	if req.LastName != "" {
		update["lastname"] = req.LastName
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)
	_, err := collection.DeleteOne(ctx, update)
	if err != nil {
		fmt.Println("Error attempting to custom delete teacher from database:", err)
		return CustomDeleteResponse{
			IsDeleted: false,
		}, err
	}

	return CustomDeleteResponse{
		IsDeleted: true,
	}, nil
}
