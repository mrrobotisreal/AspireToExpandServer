package teachersHandlers

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

func ValidateTeacherLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.ValidateTeacherLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("ValidateTeacherLogin request incoming...")

	result, err := validateTeacherLogin(req)

	if err != nil {
		http.Error(w, "Error validating teacher login.", http.StatusInternalServerError)
		return
	}

	if !result.IsValid {
		http.Error(w, "Email address, TeacherID, or password is incorrect.", http.StatusUnauthorized)
		return
	}

	response := types.ValidateTeacherLoginResponse{
		TeacherID:          result.TeacherInfo.TeacherID,
		FirstName:          result.TeacherInfo.FirstName,
		PreferredName:      result.TeacherInfo.PreferredName,
		LastName:           result.TeacherInfo.LastName,
		EmailAddress:       result.TeacherInfo.EmailAddress,
		NativeLanguage:     result.TeacherInfo.NativeLanguage,
		PreferredLanguage:  result.TeacherInfo.PreferredLanguage,
		ProfilePictureURL:  result.TeacherInfo.ProfilePictureURL,
		ProfilePicturePath: result.TeacherInfo.ProfilePicturePath,
		ThemeMode:          result.TeacherInfo.ThemeMode,
		FontStyle:          result.TeacherInfo.FontStyle,
		TimeZone:           result.TeacherInfo.TimeZone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateTeacherLogin(req types.ValidateTeacherLoginRequest) (types.ValidateTeacherLoginResult, error) {
	teacher := types.ValidateTeacherLoginResponse{
		TeacherID:          "",
		FirstName:          "",
		PreferredName:      "",
		LastName:           "",
		EmailAddress:       "",
		NativeLanguage:     "",
		PreferredLanguage:  "",
		ProfilePictureURL:  "",
		ProfilePicturePath: "",
		ThemeMode:          "",
		FontStyle:          "",
		TimeZone:           "",
	}
	validateLoginResult := types.ValidateTeacherLoginResult{
		IsValid:     false,
		TeacherInfo: teacher,
	}
	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.TeachersCollection)

	// Query MongoDB
	var teacherResult types.Teacher
	err := collection.FindOne(ctx, bson.M{"emailaddress": req.EmailAddress}).Decode(&teacherResult)
	if err != nil {
		fmt.Println("Error finding teacher account... returning error")
		fmt.Println("Error is: " + err.Error())
		return validateLoginResult, err
	}

	teacher.TeacherID = teacherResult.TeacherID
	teacher.FirstName = teacherResult.FirstName
	teacher.PreferredName = teacherResult.PreferredName
	teacher.LastName = teacherResult.LastName
	teacher.EmailAddress = teacherResult.EmailAddress
	teacher.NativeLanguage = teacherResult.NativeLanguage
	teacher.PreferredLanguage = teacherResult.PreferredLanguage
	teacher.ProfilePictureURL = teacherResult.ProfilePictureURL
	teacher.ProfilePicturePath = teacherResult.ProfilePicturePath
	teacher.ThemeMode = teacherResult.ThemeMode
	teacher.FontStyle = teacherResult.FontStyle
	teacher.TimeZone = teacherResult.TimeZone

	isPasswordValid := utils.CheckPasswordHash(req.Password+teacherResult.Salt, teacherResult.Password)
	validateLoginResult.IsValid = isPasswordValid
	validateLoginResult.TeacherInfo = teacher

	if isPasswordValid {
		fmt.Println("Teacher Password is valid: TRUE")
	} else {
		fmt.Println("Teacher Password is valid: FALSE")
	}

	return validateLoginResult, nil
}
