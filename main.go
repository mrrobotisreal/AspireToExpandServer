package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateRegistrationRequest Struct to handle incoming create registration code request
type CreateRegistrationRequest struct {
	RegistrationCode string `json:"registration_code"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	EmailAddress     string `json:"email_address"`
	IsValid          bool   `json:"is_valid"`
}

// CreateRegistrationResponse struct to handle outgoing create registration code response
type CreateRegistrationResponse struct {
	RegistrationCode string `json:"registration_code"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	EmailAddress     string `json:"email_address"`
	IsValid          bool   `json:"is_valid"`
}

// RegistrationRequest Struct to handle incoming registration code request
type RegistrationRequest struct {
	RegistrationCode string `json:"registration_code"`
}

// ValidationResponse Struct to handle outgoing validation response
type ValidationResponse struct {
	IsValid          bool   `json:"is_valid"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	EmailAddress     string `json:"email_address"`
	RegistrationCode string `json:"registration_code"`
}

// CreateNewStudentLoginRequest Struct to handle incoming create new student request
type CreateNewStudentLoginRequest struct {
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	EmailAddress       string `json:"email_address"`
	Password           string `json:"password"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
	PublicKey          string `json:"public_key"`
}

// CreateNewStudentResponse Struct to handle outgoing create new student response
type CreateNewStudentResponse struct {
	StudentId string `json:"student_id"`
}

// ValidateLoginRequest Struct to handle incoming login request
type ValidateLoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

// ValidateLoginResponse Struct to handle outgoing login response
type ValidateLoginResponse struct {
	StudentId          string `json:"student_id"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	EmailAddress       string `json:"email_address"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	StudentSince       string `json:"student_since"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

// UpdateStudentInfoRequest Struct to handle incoming updates to student info
type UpdateStudentInfoRequest struct {
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	EmailAddress       string `json:"email_address"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	StudentSince       string `json:"student_since"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
	PublicKey          string `json:"public_key"`
}

// UpdateStudentInfoResponse Struct to handle outgoing response after updating student info
type UpdateStudentInfoResponse struct {
	StudentId          string `json:"student_id"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	EmailAddress       string `json:"email_address"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	StudentSince       string `json:"student_since"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

// Assignment Struct to be stored in studentAssignmentsCollection
type Assignment struct {
	Title         string    `json:"title"`
	Subject       string    `json:"subject"`
	DocumentUrl   string    `json:"document_url"`
	DateAssigned  time.Time `json:"date_assigned"`
	DateStarted   time.Time `json:"date_started"`
	DateCompleted time.Time `json:"date_completed"`
}

// StudentAssignments Struct
type StudentAssignments struct {
	StudentId   string       `json:"student_id"`
	Assignments []Assignment `json:"assignments"`
}

// SpaceShooterGame Struct to be stored in studentGamesCollection for SpaceShooter
type SpaceShooterGame struct {
	Level         string    `json:"level"`
	Score         int       `json:"score"`
	DateStarted   time.Time `json:"date_started"`
	DateCompleted time.Time `json:"date_completed"`
}

// WordioGame Struct to be stored in studentGamesCollection for Wordio
type WordioGame struct {
	Level         string    `json:"level"`
	Score         int       `json:"score"`
	DateStarted   time.Time `json:"date_started"`
	DateCompleted time.Time `json:"date_completed"`
}

// SpellingPuddlesWord Struct word for SpellingPuddlesGame
type SpellingPuddlesWord struct {
	Word          string   `json:"word"`
	LettersChosen []string `json:"letters_chosen"`
	AudioUri      string   `json:"audio_uri"`
}

// SpellingPuddlesGame Struct to be stored in studentGamesCollection for SpellingPuddles
type SpellingPuddlesGame struct {
	Level         string    `json:"level"`
	Score         int       `json:"score"`
	DateStarted   time.Time `json:"date_started"`
	DateCompleted time.Time `json:"date_completed"`
}

// StudentGames Struct
type StudentGames struct {
	StudentId       string                `json:"student_id"`
	SpaceShooter    []SpaceShooterGame    `json:"space_shooter"`    // Well... a Space Shooter game lol
	Wordio          []WordioGame          `json:"wordio"`           // Mario-like game
	SpellingPuddles []SpellingPuddlesGame `json:"spelling_puddles"` // Rain drops containing characters fall down to spell words game
}

var mongoClient *mongo.Client
var dbName = "aspireDB"
var registrationCollection = "registrations"
var teachersCollection = "teachers"
var studentsCollection = "students"
var studentAssignmentsCollection = "assignments"
var studentGamesCollection = "games"

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Setup HTTP server
	http.HandleFunc("/registration/create", createRegistrationHandler)
	http.HandleFunc("/validate/registration", validateRegistrationHandler)
	http.HandleFunc("/students/create", createNewStudentHandler)
	http.HandleFunc("/validate/login", validateLoginHandler)
	http.HandleFunc("/teachers/create", createTeacherHandler)
	http.HandleFunc("/teachers/validate/login", validateTeacherLoginHandler)
	http.HandleFunc("/teachers/update", updateTeacherInfoHandler)
	http.HandleFunc("/students/update", updateStudentInfoHandler)
	http.HandleFunc("/students", handleFetchAllStudents)
	http.HandleFunc("/students/update/image", handleUploadProfileImage)

	// Serve profile images
	http.Handle("/uploads/profileImages/", http.StripPrefix("/uploads/profileImages", http.FileServer(http.Dir("./uploads/profileImages"))))

	fmt.Println("Server running on port 8888...")
	//if err := http.ListenAndServe(":8888", enableCors(http.DefaultServeMux)); err != nil {
	//	log.Fatal(err)
	//}
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}

// CORS middleware
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Generate a new salt
func generateSalt(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to Base64
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// Hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)

	return string(bytes), err
}

// Check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func createRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("CreateRegistration request incoming...")
	fmt.Println(req.RegistrationCode)
	fmt.Println(req.FirstName)
	fmt.Println(req.LastName)
	fmt.Println(req.EmailAddress)
	fmt.Println(req.IsValid)

	// Validate registration code
	result, err := createRegistration(req)
	if err != nil {
		fmt.Println("Error Creating registration result...")
		http.Error(w, "Error Creating registration result...", http.StatusInternalServerError)
		return
	}

	response := CreateRegistrationResponse{
		RegistrationCode: result.RegistrationCode,
		FirstName:        result.FirstName,
		LastName:         result.LastName,
		EmailAddress:     result.EmailAddress,
		IsValid:          result.IsValid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createRegistration(req CreateRegistrationRequest) (CreateRegistrationResponse, error) {
	response := CreateRegistrationResponse{
		RegistrationCode: "abcd-1234",
		FirstName:        "Bob",
		LastName:         "Smith",
		EmailAddress:     "bob@gmail.com",
		IsValid:          true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(dbName).Collection(registrationCollection)

	_, err := collection.InsertOne(ctx, bson.M{
		"registrationcode": req.RegistrationCode,
		"firstname":        req.FirstName,
		"lastname":         req.LastName,
		"emailaddress":     req.EmailAddress,
	})
	if err != nil {
		fmt.Println("Error inserting registration code...")
		fmt.Println("Error is: " + err.Error())
		return response, err
	}

	fmt.Println("Successfully inserted new registration...")
	response.RegistrationCode = req.RegistrationCode
	response.FirstName = req.FirstName
	response.LastName = req.LastName
	response.EmailAddress = req.EmailAddress
	response.IsValid = req.IsValid
	return response, nil
}

// HTTP handler for validating registration code
func validateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON body
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Registration request incoming...")
	fmt.Println(req.RegistrationCode)

	// Validate registration code
	result, err := validateRegistrationCode(req.RegistrationCode)
	if err != nil {
		fmt.Println("Error validating registration result...")
		http.Error(w, "Error validating registration result...", http.StatusInternalServerError)
		return
	}

	response := ValidationResponse{
		IsValid:          result.IsValid,
		FirstName:        result.Result.FirstName,
		LastName:         result.Result.LastName,
		EmailAddress:     result.Result.EmailAddress,
		RegistrationCode: result.Result.RegistrationCode,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Registration struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	EmailAddress     string `json:"email_address"`
	RegistrationCode string `json:"registration_code"`
}
type RegistrationValidationResult struct {
	IsValid bool         `json:"is_valid"`
	Result  Registration `json:"result"`
}

// Function to validate the registration code
func validateRegistrationCode(code string) (RegistrationValidationResult, error) {
	// For testing purposes, auto-validate the test code
	if code == "test-registration-code" {
		return RegistrationValidationResult{
			IsValid: false,
			Result: Registration{
				FirstName:        "Student",
				LastName:         "Pupil",
				EmailAddress:     "student@alina.edu",
				RegistrationCode: "Student-12345678",
			},
		}, nil
	}

	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(dbName).Collection(registrationCollection)

	// Query MongoDB
	var result Registration
	err := collection.FindOne(ctx, bson.M{"registrationcode": code}).Decode(&result)
	if err != nil {
		fmt.Println("Error finding registration... returning error.")
		fmt.Println("Error: " + err.Error())
		return RegistrationValidationResult{
			IsValid: false,
			Result: Registration{
				FirstName:        "",
				LastName:         "",
				EmailAddress:     "",
				RegistrationCode: "",
			},
		}, err
	}

	// If no error, code exists, so it's valid
	return RegistrationValidationResult{
		IsValid: true,
		Result:  result,
	}, nil
}

func createNewStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateNewStudentLoginRequest
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

	response := CreateNewStudentResponse{StudentId: newStudentId}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Student struct {
	StudentId          string `json:"student_id"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	EmailAddress       string `json:"email_address"`
	Password           string `json:"password"`
	Salt               string `json:"salt"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	StudentSince       string `json:"student_since"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

func createNewStudent(req CreateNewStudentLoginRequest) (string, error) {
	var newStudent Student
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

	salt, err := generateSalt(10)
	if err != nil {
		fmt.Println("Error generating salt... using email address instead.") // Handle this better later
		salt = req.EmailAddress
	}
	fmt.Println("Salt: " + salt)
	fmt.Println("Req.password: " + req.Password)

	hashedPassword, err := hashPassword(req.Password + salt)
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

	collection := mongoClient.Database(dbName).Collection(studentsCollection)

	result, err := collection.InsertOne(ctx, newStudent)

	if err != nil {
		return "", err
	} else {
		fmt.Println(result)
	}

	if req.PublicKey != "" {
		studentID := newStudent.StudentId
		formattedStudentID := strings.ReplaceAll(studentID, "-", "_")
		savePublicKey(formattedStudentID, req.PublicKey)
	}

	return newStudent.StudentId, err
}

func validateLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateLoginRequest
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

	response := ValidateLoginResponse{
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

type ValidateLoginResult struct {
	IsValid     bool                  `json:"is_valid"`
	StudentInfo ValidateLoginResponse `json:"student_info"`
}

func validateLogin(req ValidateLoginRequest) (ValidateLoginResult, error) {
	student := ValidateLoginResponse{
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
	validateLoginResult := ValidateLoginResult{
		IsValid:     false,
		StudentInfo: student,
	}
	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(dbName).Collection(studentsCollection)

	// Query MongoDB
	var studentResult Student
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

	isPasswordValid := checkPasswordHash(req.Password+studentResult.Salt, studentResult.Password)
	validateLoginResult.IsValid = isPasswordValid
	validateLoginResult.StudentInfo = student

	if isPasswordValid {
		fmt.Println("Password is valid: TRUE")
	} else {
		fmt.Println("Password is valid: FALSE")
	}

	return validateLoginResult, nil
}

type Teacher struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	Password           string `json:"password"`
	Salt               string `json:"salt"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

type CreateTeacherRequest struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	Password           string `json:"password"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
	PublicKey          string `json:"public_key"`
}

type CreateTeacherResponse struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

func createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTeacherRequest
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

	response := CreateTeacherResponse{
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

func createTeacher(req CreateTeacherRequest) (CreateTeacherResponse, error) {
	newTeacher := Teacher{
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
	response := CreateTeacherResponse{
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

	salt, err := generateSalt(10)
	if err != nil {
		fmt.Println("Error generating salt... using email address instead.") // Handle this better later
		salt = req.EmailAddress
	}
	fmt.Println("Salt: " + salt)
	fmt.Println("Req.password: " + req.Password)

	hashedPassword, err := hashPassword(req.Password + salt)
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

	collection := mongoClient.Database(dbName).Collection(teachersCollection)

	_, err = collection.InsertOne(ctx, newTeacher)
	if err != nil {
		log.Println("Error inserting teacher into teachersCollection: " + err.Error())
		return response, err
	}

	if req.PublicKey != "" {
		teacherID := req.TeacherID
		formattedTeacherID := strings.ReplaceAll(teacherID, "-", "_")
		savePublicKey(formattedTeacherID, req.PublicKey)
	}

	return response, nil
}

type UpdateTeacherInfoRequest struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
	PublicKey          string `json:"public_key"`
}

type UpdateTeacherInfoResponse struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

func updateTeacherInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateTeacherInfoRequest
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

func updateTeacherInfo(req UpdateTeacherInfoRequest) (UpdateTeacherInfoResponse, error) {
	updateTeacherInfoResponse := UpdateTeacherInfoResponse{
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

	collection := mongoClient.Database(dbName).Collection(studentsCollection)

	var teacherResult Teacher
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
		savePublicKey(formattedTeacherID, req.PublicKey)
	}

	return updateTeacherInfoResponse, nil
}

type ValidateTeacherLoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	TeacherID    string `json:"teacherID"`
}

type ValidateTeacherLoginResponse struct {
	TeacherID          string `json:"teacherID"`
	FirstName          string `json:"first_name"`
	PreferredName      string `json:"preferred_name"`
	LastName           string `json:"last_name"`
	NativeLanguage     string `json:"native_language"`
	PreferredLanguage  string `json:"preferred_language"`
	EmailAddress       string `json:"email_address"`
	ProfilePictureURL  string `json:"profile_picture_url"`
	ProfilePicturePath string `json:"profile_picture_path"`
	ThemeMode          string `json:"theme_mode"`
	FontStyle          string `json:"font_style"`
	TimeZone           string `json:"time_zone"`
}

func validateTeacherLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateTeacherLoginRequest
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

	response := ValidateTeacherLoginResponse{
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

type ValidateTeacherLoginResult struct {
	TeacherInfo ValidateTeacherLoginResponse
	IsValid     bool
}

func validateTeacherLogin(req ValidateTeacherLoginRequest) (ValidateTeacherLoginResult, error) {
	teacher := ValidateTeacherLoginResponse{
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
	validateLoginResult := ValidateTeacherLoginResult{
		IsValid:     false,
		TeacherInfo: teacher,
	}
	// Check MongoDB for the registration code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(dbName).Collection(teachersCollection)

	// Query MongoDB
	var teacherResult Teacher
	err := collection.FindOne(ctx, bson.M{"emailaddress": req.EmailAddress}).Decode(&teacherResult)
	if err != nil {
		fmt.Println("Error finding student account... returning error")
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

	isPasswordValid := checkPasswordHash(req.Password+teacherResult.Salt, teacherResult.Password)
	validateLoginResult.IsValid = isPasswordValid
	validateLoginResult.TeacherInfo = teacher

	if isPasswordValid {
		fmt.Println("Teacher Password is valid: TRUE")
	} else {
		fmt.Println("Teacher Password is valid: FALSE")
	}

	return validateLoginResult, nil
}

func updateStudentInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateStudentInfoRequest
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

	response := UpdateStudentInfoResponse{
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

func updateStudentInfo(req UpdateStudentInfoRequest) (UpdateStudentInfoResponse, error) {
	student := UpdateStudentInfoResponse{
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

	collection := mongoClient.Database(dbName).Collection(studentsCollection)

	// Query MongoDB
	var studentResult Student
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
		savePublicKey(formattedStudentID, req.PublicKey)
	}

	return student, nil
}

func handleFetchAllStudents(w http.ResponseWriter, r *http.Request) {
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

	collection := mongoClient.Database(dbName).Collection(studentsCollection)

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

func handleUploadProfileImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // This limits images to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join("uploads", "profileImages", handler.Filename)
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	// defer dst.Close() // <- close doesn't exist?

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to copy/save file", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("http://%s:8888/uploads/profileImages/%s", os.Getenv("IP_ADDRESS"), handler.Filename)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"imageURL": "%s"}`, imageURL)
}

func savePublicKey(userID, publicKey string) error {
	dir := "./keys"
	os.MkdirAll(dir, os.ModePerm)
	filePath := filepath.Join(dir, fmt.Sprintf("%s.pem", userID))
	return os.WriteFile(filePath, []byte(publicKey), 0644)
}

func loadPublicKey(userID string) (*rsa.PublicKey, error) {
	filePath := fmt.Sprintf("./keys/%s.pem", userID)
	publicKeyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing public key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func generateAndEncryptSymmetricKey(userID string) (string, error) {
	publicKey, err := loadPublicKey(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to load public key: %v", err)
	}

	symmetricKey := make([]byte, 32)
	if _, err := rand.Read(symmetricKey); err != nil {
		return "", fmt.Errorf("Failed to generate symmetric key: %v", err)
	}

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to encrypt symmetric key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}
