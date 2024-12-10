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

func UpdateLessonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.UpdateLessonRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := updateLesson(req)
	if err != nil {
		http.Error(w, "Error updating the lesson", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateLesson(req types.UpdateLessonRequest) (types.UpdateLessonResponse, error) {
	updatedLesson := types.Lesson{}
	update := bson.M{}

	if req.Subject != "" {
		update["subject"] = req.Subject
	}
	var scheduledDateTimeInterface interface{} = req.ScheduledDateTime
	if scheduledDateTimeResult, ok := scheduledDateTimeInterface.(int64); ok {
		fmt.Println("scheduledDateTimeResult is an int64:", scheduledDateTimeResult)
		update["scheduleddatetime"] = req.ScheduledDateTime
	}
	var roomInterface interface{} = req.Room
	if roomResult, ok := roomInterface.(int32); ok {
		fmt.Println("roomResult is an int32:", roomResult)
		update["room"] = req.Room
	}
	var timesRescheduledInterface interface{} = req.TimesRescheduled
	if timesRescheduledResult, ok := timesRescheduledInterface.(int32); ok {
		fmt.Println("timesRescheduledResult is an int32:", timesRescheduledResult)
		update["timesrescheduled"] = req.TimesRescheduled
	}
	var isCanceledExists = !(req.IsCanceled != false && req.IsCanceled != true)
	if isCanceledExists {
		fmt.Println("IsCanceled exists:", req.IsCanceled)
		update["iscanceled"] = req.IsCanceled
	}
	var isCompletedExists = !(req.IsCompleted != false && req.IsCanceled != true)
	if isCompletedExists {
		fmt.Println("IsCompleted exists:", req.IsCompleted)
		update["iscompleted"] = req.IsCompleted
	}
	var isStudentLateExists = !(req.IsStudentLate != false && req.IsStudentLate != true)
	if isStudentLateExists {
		fmt.Println("IsStudentLate exists:", req.IsStudentLate)
		update["isstudentlate"] = req.IsStudentLate
	}
	var isTeacherLateExists = !(req.IsTeacherLate != false && req.IsTeacherLate != true)
	if isTeacherLateExists {
		fmt.Println("IsTeacherLate exists:", req.IsTeacherLate)
		update["isteacherlate"] = req.IsTeacherLate
	}
	var isConnectionLostExists = !(req.IsConnectionLost != false && req.IsConnectionLost != true)
	if isConnectionLostExists {
		fmt.Println("IsConnectionLost exists:", req.IsConnectionLost)
		update["isconnectionlost"] = req.IsConnectionLost
	}

	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.LessonsCollection)

	var updateLessonResult types.Lesson
	err := collection.FindOneAndUpdate(updateCtx, bson.M{"lessonID": req.LessonID}, bson.M{
		"$set": update,
	}).Decode(&updateLessonResult)
	if err != nil {
		fmt.Errorf("Error finding and/or updating the lesson in the database: %v", err)
		return types.UpdateLessonResponse{}, err
	}

	updatedLesson.LessonID = updateLessonResult.LessonID
	updatedLesson.TeacherID = updateLessonResult.TeacherID
	updatedLesson.StudentId = updateLessonResult.StudentId
	updatedLesson.Subject = updateLessonResult.Subject
	updatedLesson.ScheduledDateTime = updateLessonResult.ScheduledDateTime
	updatedLesson.Room = updateLessonResult.Room
	updatedLesson.TimesRescheduled = updateLessonResult.TimesRescheduled
	updatedLesson.IsCanceled = updateLessonResult.IsCanceled
	updatedLesson.IsCompleted = updateLessonResult.IsCompleted
	updatedLesson.IsTeacherLate = updateLessonResult.IsTeacherLate
	updatedLesson.IsStudentLate = updateLessonResult.IsStudentLate
	updatedLesson.IsConnectionLost = updateLessonResult.IsConnectionLost

	return types.UpdateLessonResponse{
		Lesson: updatedLesson,
	}, nil
}
