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
	fmt.Println("Incoming GetTeacher request...")

	teacherID := r.URL.Query().Get("teacherID")
	fmt.Println("For TeacherID:", teacherID)

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

	fmt.Println("Getting teacher", teacherID, "from the database...")
	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)
	var teacherResult types.TeacherInfo
	err := collection.FindOne(ctx, bson.M{"teacherid": teacherID}).Decode(&teacherResult)
	if err != nil {
		fmt.Println("Error getting teacher info from the database:", err)
		return types.GetTeacherResponse{}, err
	}
	fmt.Println("Successfully retrieved teacher", teacherResult.FirstName, teacherResult.LastName, "from the database!")

	return types.GetTeacherResponse{
		Teacher: teacherResult,
	}, nil
}
