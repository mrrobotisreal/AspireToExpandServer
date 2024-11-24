package teachersHandlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/types"
	"io.winapps.aspirewithalina.aspirewithalinaserver/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

func UpdateTeacherInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req types.UpdateTeacherInfoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := updateTeacherInfo(req)
	if err != nil {
		http.Error(w, "Error updating teacher info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func updateTeacherInfo(req types.UpdateTeacherInfoRequest) (types.UpdateTeacherInfoResponse, error) {
	updateTeacherInfoResponse := types.UpdateTeacherInfoResponse{
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

	var teacherResult types.Teacher
	err := collection.FindOneAndUpdate(ctx, bson.M{
		"$or": []bson.M{
			{"teacherID": req.TeacherID},
			{"emailaddress": req.EmailAddress},
		},
	}, bson.M{
		"$set": update,
	}).Decode(&teacherResult)
	if err != nil {
		log.Println("Error finding and updating teacher info: ", err.Error())
		return updateTeacherInfoResponse, nil
	}

	if req.PublicKey != "" {
		teacherID := teacherResult.TeacherID
		formattedTeacherID := strings.ReplaceAll(teacherID, "-", "_")
		utils.SavePublicKey(formattedTeacherID, req.PublicKey)
	}

	return updateTeacherInfoResponse, nil
}
