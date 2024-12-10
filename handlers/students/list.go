package studentsHandlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ListStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var page int64
	var limit int64
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageResult, err := strconv.ParseInt(pageStr, 10, 64); err == nil {
		page = pageResult
	} else {
		http.Error(w, "Invalid request query, \"page\" must be a number between 1 and 1,000,000", http.StatusBadRequest)
		return
	}
	if page <= 0 {
		http.Error(w, "Invalid page size, \"page\" must be a number between 1 and 1,000,000", http.StatusBadRequest)
		return
	}
	if page > 1000000 {
		http.Error(w, "Invalid page size, \"page\" must be a number between 1 and 1,000,000", http.StatusBadRequest)
		return
	}

	if limitResult, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
		limit = limitResult
	} else {
		http.Error(w, "Invalid request query, \"limit\" must be a number between 1 and 100", http.StatusBadRequest)
		return
	}
	if limit <= 0 {
		http.Error(w, "Invalid limit size, \"limit\" must be a number between 1 and 100", http.StatusBadRequest)
		return
	}
	if limit > 100 {
		http.Error(w, "Invalid limit size, \"limit\" must be a number between 1 and 100", http.StatusBadRequest)
		return
	}

	response, err := listStudents(page, limit)
	if err != nil {
		http.Error(w, "Error listing students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listStudents(page, limit int64) (types.ListStudentsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)
	pipeline := mongo.Pipeline{
		{{"$sort", bson.D{{"preferredname", 1}}}},
		{{"$project", bson.M{
			"preferredname":     1,
			"firstname":         1,
			"lastname":          1,
			"studentid":         1,
			"emailaddress":      1,
			"profilepictureurl": 1,
			"_id":               0,
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error listing all students: %v", err.Error())
		return types.ListStudentsResponse{
			Students: nil,
			Page:     page,
		}, err
	}
	defer cursor.Close(ctx)

	var results []types.Student
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Error getting all students results from cursor: %v", err.Error())
		return types.ListStudentsResponse{
			Students: nil,
			Page:     page,
		}, err
	}

	return types.ListStudentsResponse{
		Students: results,
		Page:     page,
	}, nil
}
