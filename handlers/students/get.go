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

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Incoming GetStudent request...")

	studentID := r.URL.Query().Get("studentID")
	if studentID == "" {
		http.Error(w, "Invalid request query, \"studentID\" cannot be empty", http.StatusBadRequest)
		return
	}
	fmt.Println("For StudentID:", studentID)

	response, err := getStudent(studentID)
	if err != nil {
		http.Error(w, "Error getting student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getStudent(studentID string) (types.GetStudentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Getting student", studentID, "from the database...")
	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)
	var studentResult types.StudentInfo
	err := collection.FindOne(ctx, bson.M{"studentid": studentID}).Decode(&studentResult)
	if err != nil {
		fmt.Println("Error finding student in the database:", err)
		return types.GetStudentResponse{}, err
	}
	fmt.Println("Successfully retrieved student", studentResult.FirstName, studentResult.LastName, "from the database!")

	return types.GetStudentResponse{
		Student: studentResult,
	}, nil
}
