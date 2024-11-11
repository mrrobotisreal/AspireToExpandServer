package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	FirstName         string `json:"first_name"`
	PreferredName     string `json:"preferred_name"`
	LastName          string `json:"last_name"`
	EmailAddress      string `json:"email_address"`
	Password          string `json:"password"`
	NativeLanguage    string `json:"native_language"`
	PreferredLanguage string `json:"preferred_language"`
	ThemeMode         string `json:"theme_mode"`
	FontStyle         string `json:"font_style"`
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
	StudentId         string `json:"student_id"`
	FirstName         string `json:"first_name"`
	PreferredName     string `json:"preferred_name"`
	LastName          string `json:"last_name"`
	EmailAddress      string `json:"email_address"`
	NativeLanguage    string `json:"native_language"`
	PreferredLanguage string `json:"preferred_language"`
	StudentSince      string `json:"student_since"`
	ThemeMode         string `json:"theme_mode"`
	FontStyle         string `json:"font_style"`
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
	http.HandleFunc("/validate/registration", validateRegistrationHandler)
	http.HandleFunc("/students/create", createNewStudentHandler)
	http.HandleFunc("/validate/login", validateLoginHandler)

	fmt.Println("Server running on port 8888...")
	if err := http.ListenAndServe(":8888", enableCors(http.DefaultServeMux)); err != nil {
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
	StudentId         string `json:"student_id"`
	FirstName         string `json:"first_name"`
	PreferredName     string `json:"preferred_name"`
	LastName          string `json:"last_name"`
	EmailAddress      string `json:"email_address"`
	Password          string `json:"password"`
	Salt              string `json:"salt"`
	NativeLanguage    string `json:"native_language"`
	PreferredLanguage string `json:"preferred_language"`
	StudentSince      string `json:"student_since"`
	ThemeMode         string `json:"theme_mode"`
	FontStyle         string `json:"font_style"`
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
	newStudent.ThemeMode = req.ThemeMode
	newStudent.FontStyle = req.FontStyle

	salt, err := generateSalt(10)
	if err != nil {
		fmt.Println("Error generating salt... using email address instead.") // Handle this better later
		salt = req.EmailAddress
	}
	fmt.Println("Salt: " + salt)

	hashedPassword, err := hashPassword(req.Password + salt)
	if err != nil {
		fmt.Println("Error hashing password and salt... returning an error") // Handle this better later
		fmt.Println("Error is: " + err.Error())
		return "", err
	}

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
	fmt.Println(req.EmailAddress)
	fmt.Println(req.Password)

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
		StudentId:         result.StudentInfo.StudentId,
		FirstName:         result.StudentInfo.FirstName,
		PreferredName:     result.StudentInfo.PreferredName,
		LastName:          result.StudentInfo.LastName,
		EmailAddress:      result.StudentInfo.EmailAddress,
		NativeLanguage:    result.StudentInfo.NativeLanguage,
		PreferredLanguage: result.StudentInfo.PreferredLanguage,
		StudentSince:      result.StudentInfo.StudentSince,
		ThemeMode:         result.StudentInfo.ThemeMode,
		FontStyle:         result.StudentInfo.FontStyle,
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
		StudentId:         "",
		FirstName:         "",
		PreferredName:     "",
		LastName:          "",
		EmailAddress:      "",
		NativeLanguage:    "",
		PreferredLanguage: "",
		StudentSince:      "",
		ThemeMode:         "",
		FontStyle:         "",
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
	student.ThemeMode = studentResult.ThemeMode
	student.FontStyle = studentResult.FontStyle

	isPasswordValid := checkPasswordHash(req.Password+studentResult.Salt, studentResult.Password)
	validateLoginResult.IsValid = isPasswordValid

	if isPasswordValid {
		fmt.Println("Password is valid: TRUE")
	} else {
		fmt.Println("Password is valid: FALSE")
	}

	return validateLoginResult, nil
}

// STUDENT struct
// RegistrationCode
// StudentId
// FirstName
// PreferredName
// LastName
// EmailAddress
// NativeLanguage
// PreferredLanguage
// Password
// Salt
// StudentSince (timestamp)
// ThemeMode
