package teachersHandlers

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

type CustomGetRequest struct {
	TeacherID     string `json:"teacherID"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	PreferredName string `json:"preferred_name"`
}

func CustomGetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CustomGetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := customGetTeacher(req)
	if err != nil {
		http.Error(w, "Error getting teacher info via custom handler", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func customGetTeacher(req CustomGetRequest) (types.GetTeacherResponse, error) {
	update := bson.M{}

	if req.TeacherID != "" {
		update["teacherid"] = req.TeacherID
	}
	if req.PreferredName != "" {
		update["preferredname"] = req.PreferredName
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

	var teacherResult types.Teacher
	err := collection.FindOne(ctx, update).Decode(&teacherResult)
	if err != nil {
		fmt.Println("Error getting custom teacher info from the database:", err)
		return types.GetTeacherResponse{}, err
	}

	return types.GetTeacherResponse{
		Teacher: teacherResult,
	}, nil
}
