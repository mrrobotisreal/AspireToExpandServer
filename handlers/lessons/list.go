package lessonsHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"net/http"
	"strconv"
	"time"
)

func ListLessonsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var isCanceled bool
	var isCompleted bool
	isCanceledStr := r.URL.Query().Get("isCanceled")
	isCompletedStr := r.URL.Query().Get("isCompleted")

	if isCanceledStr == "" {
		http.Error(w, "Invalid request query, \"isCanceled\" cannot be empty", http.StatusBadRequest)
		return
	}
	if isCompletedStr == "" {
		http.Error(w, "Invalid request query, \"isCompleted\" cannot be empty", http.StatusBadRequest)
		return
	}

	if isCanceledResult, err := strconv.ParseBool(isCanceledStr); err == nil {
		isCanceled = isCanceledResult
	} else {
		http.Error(w, "Invalid request query, \"isCanceled\" must be either \"true\" or \"false\"", http.StatusBadRequest)
		return
	}
	if isCompletedResult, err := strconv.ParseBool(isCompletedStr); err == nil {
		isCompleted = isCompletedResult
	} else {
		http.Error(w, "Invalid request query, \"isCompleted\" must be either \"true\" or \"false\"", http.StatusBadRequest)
		return
	}

	response, err := listLessons(isCanceled, isCompleted)
	if err != nil {
		http.Error(w, "Error listing the lessons", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listLessons(isCanceled, isCompleted bool) (types.ListLessonsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.LessonsCollection)
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"iscanceled":  isCanceled,
			"iscompleted": isCompleted,
		}}},
		{{"$sort", bson.D{{"scheduleddatetime", -1}}}},
		{{"$project", bson.M{
			"lessonid":          1,
			"teacherid":         1,
			"studentid":         1,
			"subject":           1,
			"scheduleddatetime": 1,
			"room":              1,
			"iscanceled":        1,
			"iscompleted":       1,
			"timesrescheduled":  1,
			"isstudentlate":     1,
			"isteacherlate":     1,
			"isconnectionlost":  1,
			"_id":               0, // Exclude MongoDB internal _id
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Errorf("Error aggregating lessons from the database: %v", err)
		return types.ListLessonsResponse{
			Lessons: nil,
			Page:    1, // TODO: figure out pagination for this
		}, err
	}
	defer cursor.Close(ctx)

	var results []types.Lesson
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Errorf("Error compiling all lessons into results: %v", err)
		return types.ListLessonsResponse{
			Lessons: nil,
			Page:    1, // TODO: figure out pagination for this
		}, err
	}

	return types.ListLessonsResponse{
		Lessons: results,
		Page:    1, // TODO: figure out pagination for this
	}, nil
}
