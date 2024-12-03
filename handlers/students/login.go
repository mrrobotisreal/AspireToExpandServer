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

func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.ValidateLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("ValidateLogin request incoming...")

	result, err := validateLogin(req)

	if err != nil {
		http.Error(w, "Error validating student login.", http.StatusInternalServerError)
		return
	}

	if !result.IsValid {
		http.Error(w, "Email address or password is incorrect.", http.StatusUnauthorized)
		return
	}

	response := types.ValidateLoginResponse{
		StudentId:          result.StudentInfo.StudentId,
		FirstName:          result.StudentInfo.FirstName,
		PreferredName:      result.StudentInfo.PreferredName,
		LastName:           result.StudentInfo.LastName,
		EmailAddress:       result.StudentInfo.EmailAddress,
		NativeLanguage:     result.StudentInfo.NativeLanguage,
		PreferredLanguage:  result.StudentInfo.PreferredLanguage,
		StudentSince:       result.StudentInfo.StudentSince,
		ProfilePictureURL:  result.StudentInfo.ProfilePictureURL,
		ProfilePicturePath: result.StudentInfo.ProfilePicturePath,
		ThemeMode:          result.StudentInfo.ThemeMode,
		FontStyle:          result.StudentInfo.FontStyle,
		TimeZone:           result.StudentInfo.TimeZone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateLogin(req types.ValidateLoginRequest) (types.ValidateLoginResult, error) {
	student := types.ValidateLoginResponse{
		StudentId:          "",
		FirstName:          "",
		PreferredName:      "",
		LastName:           "",
		EmailAddress:       "",
		NativeLanguage:     "",
		PreferredLanguage:  "",
		StudentSince:       "",
		ProfilePictureURL:  "",
		ProfilePicturePath: "",
		ThemeMode:          "",
		FontStyle:          "",
		TimeZone:           "",
	}
	validateLoginResult := types.ValidateLoginResult{
		IsValid:     false,
		StudentInfo: student,
	}
	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	// Query MongoDB
	var studentResult types.Student
	err := collection.FindOne(ctx, bson.M{"emailaddress": req.EmailAddress}).Decode(&studentResult)
	if err != nil {
		fmt.Println("Error finding student account... returning error")
		fmt.Println("Error is: " + err.Error())
		return validateLoginResult, err
	}

	student.StudentId = studentResult.StudentId
	student.FirstName = studentResult.FirstName
	student.PreferredName = studentResult.PreferredName
	student.LastName = studentResult.LastName
	student.EmailAddress = studentResult.EmailAddress
	student.NativeLanguage = studentResult.NativeLanguage
	student.PreferredLanguage = studentResult.PreferredLanguage
	student.StudentSince = studentResult.StudentSince
	student.ProfilePictureURL = studentResult.ProfilePictureURL
	student.ProfilePicturePath = studentResult.ProfilePicturePath
	student.ThemeMode = studentResult.ThemeMode
	student.FontStyle = studentResult.FontStyle
	student.TimeZone = studentResult.TimeZone

	isPasswordValid := utils.CheckPasswordHash(req.Password+studentResult.Salt, studentResult.Password)
	validateLoginResult.IsValid = isPasswordValid
	validateLoginResult.StudentInfo = student

	if isPasswordValid {
		fmt.Println("Password is valid: TRUE")
	} else {
		fmt.Println("Password is valid: FALSE")
	}

	return validateLoginResult, nil
}

func ValidateGoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.ValidateGoogleLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Email: " + req.Email)
	fmt.Println("EmailVerified: " + req.EmailVerified)

	if req.EmailVerified != "true" {
		http.Error(w, "Email not verified", http.StatusBadRequest)
		fmt.Println("Email not verified...")
		return
	}

	fmt.Println("ValidateGoogleLogin request incoming...")

	response, err := validateGoogleLogin(req)
	if err != nil {
		http.Error(w, "Error validating student's Google login.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateGoogleLogin(req types.ValidateGoogleLoginRequest) (types.ValidateLoginResponse, error) {
	student := types.ValidateLoginResponse{
		StudentId:          "",
		FirstName:          "",
		PreferredName:      "",
		LastName:           "",
		EmailAddress:       "",
		NativeLanguage:     "",
		PreferredLanguage:  "",
		StudentSince:       "",
		ProfilePictureURL:  "",
		ProfilePicturePath: "",
		ThemeMode:          "",
		FontStyle:          "",
		TimeZone:           "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	// Query MongoDB
	var studentResult types.Student
	err := collection.FindOne(ctx, bson.M{"emailaddress": req.Email}).Decode(&studentResult)
	if err != nil {
		fmt.Println("Error finding student account... returning error")
		fmt.Println("Error is: " + err.Error())
		return student, err
	}

	student.StudentId = studentResult.StudentId
	student.FirstName = studentResult.FirstName
	student.PreferredName = studentResult.PreferredName
	student.LastName = studentResult.LastName
	student.EmailAddress = studentResult.EmailAddress
	student.NativeLanguage = studentResult.NativeLanguage
	student.PreferredLanguage = studentResult.PreferredLanguage
	student.StudentSince = studentResult.StudentSince
	student.ProfilePictureURL = studentResult.ProfilePictureURL
	student.ProfilePicturePath = studentResult.ProfilePicturePath
	student.ThemeMode = studentResult.ThemeMode
	student.FontStyle = studentResult.FontStyle
	student.TimeZone = studentResult.TimeZone

	return student, nil
}
