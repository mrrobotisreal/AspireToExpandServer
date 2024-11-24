package teachersHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"io.winapps.aspirewithalina.aspirewithalinaserver/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

func CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateTeacherRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("CreateTeacher request incoming...")

	result, err := createTeacher(req)
	if err != nil {
		http.Error(w, "Error creating teacher: "+err.Error(), http.StatusInternalServerError)
	}

	response := types.CreateTeacherResponse{
		TeacherID:          result.TeacherID,
		FirstName:          result.FirstName,
		PreferredName:      result.PreferredName,
		LastName:           result.LastName,
		NativeLanguage:     result.NativeLanguage,
		PreferredLanguage:  result.PreferredLanguage,
		EmailAddress:       result.EmailAddress,
		ProfilePictureURL:  result.ProfilePictureURL,
		ProfilePicturePath: result.ProfilePicturePath,
		ThemeMode:          result.ThemeMode,
		FontStyle:          result.FontStyle,
		TimeZone:           result.TimeZone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createTeacher(req types.CreateTeacherRequest) (types.CreateTeacherResponse, error) {
	newTeacher := types.Teacher{
		TeacherID:          req.TeacherID,
		FirstName:          req.FirstName,
		PreferredName:      req.PreferredName,
		LastName:           req.LastName,
		NativeLanguage:     req.NativeLanguage,
		PreferredLanguage:  req.PreferredLanguage,
		EmailAddress:       req.EmailAddress,
		ProfilePictureURL:  req.ProfilePictureURL,
		ProfilePicturePath: req.ProfilePicturePath,
		ThemeMode:          req.ThemeMode,
		FontStyle:          req.FontStyle,
		TimeZone:           req.TimeZone,
	}
	response := types.CreateTeacherResponse{
		TeacherID:          req.TeacherID,
		FirstName:          req.FirstName,
		PreferredName:      req.PreferredName,
		LastName:           req.LastName,
		NativeLanguage:     req.NativeLanguage,
		PreferredLanguage:  req.PreferredLanguage,
		EmailAddress:       req.EmailAddress,
		ProfilePictureURL:  req.ProfilePictureURL,
		ProfilePicturePath: req.ProfilePicturePath,
		ThemeMode:          req.ThemeMode,
		FontStyle:          req.FontStyle,
		TimeZone:           req.TimeZone,
	}

	salt, err := utils.GenerateSalt(10)
	if err != nil {
		fmt.Println("Error generating salt... using email address instead.") // Handle this better later
		salt = req.EmailAddress
	}
	fmt.Println("Salt: " + salt)
	fmt.Println("Req.password: " + req.Password)

	hashedPassword, err := utils.HashPassword(req.Password + salt)
	if err != nil {
		fmt.Println("Error hashing password and salt... returning an error") // Handle this better later
		fmt.Println("Error is: " + err.Error())
		return response, err
	}
	fmt.Println("hashedPassword: " + hashedPassword)

	newTeacher.Password = hashedPassword
	newTeacher.Salt = salt

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)

	_, err = collection.InsertOne(ctx, newTeacher)
	if err != nil {
		log.Println("Error inserting teacher into teachersCollection: " + err.Error())
		return response, err
	}

	if req.PublicKey != "" {
		teacherID := req.TeacherID
		formattedTeacherID := strings.ReplaceAll(teacherID, "-", "_")
		utils.SavePublicKey(formattedTeacherID, req.PublicKey)
	}

	return response, nil
}
