package teachersHandlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"log"
	"net/http"
	"time"
)

func ListTeachersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	response, err := listTeachers()
	if err != nil {
		http.Error(w, "Error listing all teachers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listTeachers() (types.ListTeachersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)
	pipeline := mongo.Pipeline{
		{{"$sort", bson.D{{"preferredname", 1}}}},
		{{"$project", bson.M{
			"teacherid":         1,
			"firstname":         1,
			"preferredname":     1,
			"lastname":          1,
			"emailaddress":      1,
			"profilepictureurl": 1,
			"_id":               0,
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error fetching all teachers from database: %v", err.Error())
		return types.ListTeachersResponse{
			Teachers: nil,
		}, err
	}
	defer cursor.Close(ctx)

	var results []types.TeacherInfo
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Error compiling all teachers results from cursor: %v", err.Error())
		return types.ListTeachersResponse{
			Teachers: nil,
		}, err
	}

	return types.ListTeachersResponse{
		Teachers: results,
	}, nil
}
