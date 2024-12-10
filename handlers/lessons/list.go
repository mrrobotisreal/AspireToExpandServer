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
	"time"
)

func ListLessonsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.ListLessonsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := listLessons(req)
	if err != nil {
		http.Error(w, "Error listing the lessons", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listLessons(req types.ListLessonsRequest) (types.ListLessonsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.LessonsCollection)
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"iscanceled":  req.IsCanceled,
			"iscompleted": req.IsCompleted,
		}}},
		{{"$sort", bson.D{{"scheduleddatetime", -1}}}},
		{{"$project", bson.M{
			"lessonID":          1,
			"teacherID":         1,
			"studentID":         1,
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
