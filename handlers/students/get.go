package studentsHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"log"
	"net/http"
	"time"
)

func HandleFetchAllStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	students, err := fetchAllStudents()
	if err != nil {
		http.Error(w, "Error fetching all students.", http.StatusInternalServerError)
		log.Println(err)
	}

	fmt.Println(students)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func fetchAllStudents() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	pipeline := mongo.Pipeline{
		{{"$sort", bson.D{{"preferredname", 1}}}},
		{{"$project", bson.M{
			"preferredname":     1,
			"studentid":         1,
			"emailaddress":      1,
			"profilepictureurl": 1,
			"_id":               0,
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error fetching all students: %v", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Error getting all students results from cursor: %v", err.Error())
		return nil, err
	}
	fmt.Println("RESULTS:")
	fmt.Println(results)

	return results, nil
}
