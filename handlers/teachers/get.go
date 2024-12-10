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

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	teacherID := r.URL.Query().Get("teacherID")
	fmt.Println("?teacherID=", teacherID)

	response, err := getTeacher(teacherID)
	if err != nil {
		http.Error(w, "Error getting teacher info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTeacher(teacherID string) (types.GetTeacherResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)

	var teacherResult types.Teacher
	err := collection.FindOne(ctx, bson.M{"teacherID": teacherID}).Decode(&teacherResult)
	if err != nil {
		fmt.Println("Error getting teacher info from the database:", err)
		return types.GetTeacherResponse{}, err
	}

	return types.GetTeacherResponse{
		Teacher: teacherResult,
	}, nil
}
