package studentsHandlers

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

func HandleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.DeleteStudentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := deleteStudent(req.StudentId)
	if err != nil {
		http.Error(w, "Error deleting student. Student was not deleted", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteStudent(StudentId string) (types.DeleteStudentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	_, err := collection.DeleteOne(ctx, bson.M{"studentid": StudentId})
	if err != nil {
		fmt.Errorf("Error deleting student from database: %v", err)
		return types.DeleteStudentResponse{
			IsDeleted: false,
		}, err
	}

	return types.DeleteStudentResponse{
		IsDeleted: true,
	}, nil
}
