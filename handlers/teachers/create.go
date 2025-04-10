package teachersHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
		LessonsTaught:      result.LessonsTaught,
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
		LessonsTaught:      0,
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
		LessonsTaught:      0,
	}

	fmt.Println("TeacherID:", newTeacher.TeacherID)
	fmt.Println("FirstName:", newTeacher.FirstName)
	fmt.Println("PreferredName:", newTeacher.PreferredName)
	fmt.Println("LastName:", newTeacher.LastName)
	fmt.Println("NativeLanguage:", newTeacher.NativeLanguage)
	fmt.Println("PreferredLanguage:", newTeacher.PreferredLanguage)
	fmt.Println("EmailAddress:", newTeacher.EmailAddress)
	fmt.Println("ProfilePictureURL:", newTeacher.ProfilePictureURL)
	fmt.Println("ProfilePicturePath:", newTeacher.ProfilePicturePath)
	fmt.Println("ThemeMode:", newTeacher.ThemeMode)
	fmt.Println("FontStyle:", newTeacher.FontStyle)
	fmt.Println("TimeZone:", newTeacher.TimeZone)
	fmt.Println("LessonsTaught:", newTeacher.LessonsTaught)

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

	usersCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := db.MongoClient.Database(db.DbName).Collection(db.UsersCollection)
	_, err = usersCollection.InsertOne(usersCtx, bson.M{
		"userId":            newTeacher.TeacherID,
		"userType":          "teacher",
		"preferredName":     newTeacher.PreferredName,
		"firstName":         newTeacher.FirstName,
		"lastName":          newTeacher.LastName,
		"profilePictureUrl": newTeacher.ProfilePictureURL,
	})
	if err != nil {
		log.Println("Error inserting user into usersCollection: " + err.Error())
		return response, err
	}

	return response, nil
}
