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

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	teacherID := r.URL.Query().Get("teacherID")
	if teacherID == "" {
		http.Error(w, "Invalid request query. The \"teacherID\" query is required", http.StatusBadRequest)
		return
	}

	response, err := deleteTeacher(teacherID)
	if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteTeacher(teacherID string) (types.DeleteTeacherResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"teacherid": teacherID})
	if err != nil {
		fmt.Println("Error attempting to delete the teacher from the database:", err)
		return types.DeleteTeacherResponse{
			IsDeleted: false,
		}, err
	}

	return types.DeleteTeacherResponse{
		IsDeleted: true,
	}, nil
}
