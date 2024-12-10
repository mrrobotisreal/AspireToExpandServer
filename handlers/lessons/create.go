package lessonsHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"net/http"
	"time"
)

func CreateLessonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateLessonRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := createLesson(req)
	if err != nil {
		http.Error(w, "Error creating lesson", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createLesson(req types.CreateLessonRequest) (types.CreateLessonResponse, error) {
	newLesson := types.Lesson{
		LessonID:          uuid.New().String(),
		TeacherID:         req.TeacherID,
		StudentId:         req.StudentId,
		Subject:           req.Subject,
		ScheduledDateTime: req.ScheduledDateTime,
		Room:              req.Room,
		IsCanceled:        false,
		IsCompleted:       false,
		TimesRescheduled:  0,
		IsStudentLate:     false,
		IsTeacherLate:     false,
		IsConnectionLost:  false,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.LessonsCollection)

	_, err := collection.InsertOne(ctx, newLesson)

	if err != nil {
		fmt.Errorf("Error inserting new lesson into the database: %v", err)
		return types.CreateLessonResponse{}, err
	}

	return types.CreateLessonResponse{
		Lesson: newLesson,
	}, nil
}
