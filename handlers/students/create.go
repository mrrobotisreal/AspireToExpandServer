package studentsHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"io.winapps.aspirewithalina.aspirewithalinaserver/utils"
	"net/http"
	"strings"
	"time"
)

func CreateNewStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateNewStudentLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("CreateNewStudent request incoming...")

	newStudentId, err := createNewStudent(req)

	if err != nil {
		http.Error(w, "Error saving new student. Data was not saved.", http.StatusInternalServerError)
		return
	}

	response := types.CreateNewStudentResponse{StudentId: newStudentId}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createNewStudent(req types.CreateNewStudentLoginRequest) (string, error) {
	var newStudent types.Student
	newStudent.StudentId = uuid.New().String()
	newStudent.FirstName = req.FirstName
	newStudent.PreferredName = req.PreferredName
	newStudent.LastName = req.LastName
	newStudent.EmailAddress = req.EmailAddress
	newStudent.NativeLanguage = req.NativeLanguage
	newStudent.PreferredLanguage = req.PreferredLanguage
	newStudent.ProfilePictureURL = req.ProfilePictureURL
	newStudent.ProfilePicturePath = req.ProfilePicturePath
	newStudent.ThemeMode = req.ThemeMode
	newStudent.FontStyle = req.FontStyle
	newStudent.TimeZone = req.TimeZone
	newStudent.LessonsRemaining = 0
	newStudent.LessonsCompleted = 0

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
		return "", err
	}
	fmt.Println("hashedPassword: " + hashedPassword)

	newStudent.Password = hashedPassword
	newStudent.Salt = salt
	newStudent.StudentSince = time.Now().UTC().String()

	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	result, err := collection.InsertOne(ctx, newStudent)

	if err != nil {
		return "", err
	} else {
		fmt.Println(result)
	}

	if req.PublicKey != "" {
		studentID := newStudent.StudentId
		formattedStudentID := strings.ReplaceAll(studentID, "-", "_")
		utils.SavePublicKey(formattedStudentID, req.PublicKey)
	}

	usersCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := db.MongoClient.Database(db.DbName).Collection(db.UsersCollection)
	usersResult, err := usersCollection.InsertOne(usersCtx, bson.M{
		"userId":            newStudent.StudentId,
		"userType":          "student",
		"preferredName":     newStudent.PreferredName,
		"firstName":         newStudent.FirstName,
		"lastName":          newStudent.LastName,
		"profilePictureUrl": newStudent.ProfilePictureURL,
	})
	if err != nil {
		return "", err
	} else {
		fmt.Println("Users result:")
		fmt.Println(usersResult)
	}

	return newStudent.StudentId, err
}
