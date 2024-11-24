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
	"strings"
	"time"
)

func UpdateStudentInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.UpdateStudentInfoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("UpdateStudentInfo request incoming...")
	fmt.Println("req.EmailAddress: " + req.EmailAddress)
	fmt.Println("req.ProfilePictureURL: " + req.ProfilePictureURL)
	fmt.Println("req.ProfilePicturePath: " + req.ProfilePicturePath)
	fmt.Println("req.ThemeMode: " + req.ThemeMode)
	fmt.Println("req.FontStyle: " + req.FontStyle)
	fmt.Println("req.PublicKey: " + req.PublicKey)

	result, err := updateStudentInfo(req)

	if err != nil {
		http.Error(w, "Error updating student information.", http.StatusInternalServerError)
		return
	}

	response := types.UpdateStudentInfoResponse{
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
		TimeZone:           result.TimeZone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateStudentInfo(req types.UpdateStudentInfoRequest) (types.UpdateStudentInfoResponse, error) {
	student := types.UpdateStudentInfoResponse{
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

	update := bson.M{}

	if req.ThemeMode != "" {
		update["thememode"] = req.ThemeMode
	}

	if req.FontStyle != "" {
		update["fontstyle"] = req.FontStyle
	}

	if req.ProfilePictureURL != "" {
		update["profilepictureurl"] = req.ProfilePictureURL
	}

	if req.ProfilePicturePath != "" {
		update["profilepicturepath"] = req.ProfilePicturePath
	}

	if req.PreferredName != "" {
		update["preferredname"] = req.PreferredName
	}

	if req.PreferredLanguage != "" {
		update["preferredlanguage"] = req.PreferredLanguage
	}

	if req.TimeZone != "" {
		update["timezone"] = req.TimeZone
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoClient.Database(db.DbName).Collection(db.StudentsCollection)

	// Query MongoDB
	var studentResult types.Student
	err := collection.FindOneAndUpdate(ctx, bson.M{
		"$or": []bson.M{
			{"emailaddress": req.EmailAddress}, // add StudentID here later
		},
	}, bson.M{
		"$set": update,
	}).Decode(&studentResult)
	if err != nil {
		fmt.Println("Error updating student info in mongodb... returning error")
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

	if req.PublicKey != "" {
		studentID := studentResult.StudentId
		formattedStudentID := strings.ReplaceAll(studentID, "-", "_")
		utils.SavePublicKey(formattedStudentID, req.PublicKey)
	}

	return student, nil
}
