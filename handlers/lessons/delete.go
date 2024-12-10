package lessonsHandlers

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

func DeleteLessonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.DeleteLessonRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := deleteLesson(req)
	if err != nil {
		http.Error(w, "Error deleting the lesson", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteLesson(req types.DeleteLessonRequest) (types.DeleteLessonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.LessonsCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"lessonID": req.LessonID})
	if err != nil {
		fmt.Errorf("Error deleting the lesson from the database: %v", err)
		return types.DeleteLessonResponse{
			IsDeleted: false,
		}, err
	}

	return types.DeleteLessonResponse{
		IsDeleted: true,
	}, nil
}
