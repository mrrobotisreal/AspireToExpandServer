package studentsHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"io.winapps.aspirewithalina.aspirewithalinaserver/utils"
	"net/http"
	"time"
)

func ValidateLoginMobileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.ValidateLoginMobileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := validateLoginMobile(req)
	if err != nil {
		http.Error(w, "Error occurred while logging in", http.StatusInternalServerError)
		return
	}
	if !result.IsValid {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		// add logging here for unauthorized attempts later
		fmt.Println("Unauthorized login attempt. Invalid password for email: ", req.EmailAddress)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.StudentInfo)
}

func validateLoginMobile(req types.ValidateLoginMobileRequest) (types.ValidateLoginMobileResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)
	var result types.Student
	err := collection.FindOne(ctx, bson.M{"emailaddress": req.EmailAddress}).Decode(&result)
	if err != nil {
		fmt.Println("An error occurred while attempting to find the student in the database: ", err)
		return types.ValidateLoginMobileResult{}, err
	}

	isPasswordValid := utils.CheckPasswordHash(req.Password+result.Salt, result.Password)
	return types.ValidateLoginMobileResult{
		IsValid: isPasswordValid,
		StudentInfo: types.ValidateLoginMobileResponse{
			StudentId:          result.StudentId,
			FirstName:          result.FirstName,
			PreferredName:      result.PreferredName,
			LastName:           result.LastName,
			EmailAddress:       result.EmailAddress,
			NativeLanguage:     result.NativeLanguage,
			PreferredLanguage:  result.PreferredLanguage,
			StudentSince:       result.StudentSince,
			ProfilePictureURL:  result.ProfilePictureURL,
			ProfilePicturePath: result.ProfilePicturePath,
			ThemeMode:          result.ThemeMode,
			FontStyle:          result.FontStyle,
			LessonsRemaining:   result.LessonsRemaining,
			LessonsCompleted:   result.LessonsCompleted,
		},
	}, nil
}
