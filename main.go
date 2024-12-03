package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/handlers"
	studentsHandlers "io.winapps.aspirewithalina.aspirewithalinaserver/handlers/students"
	teachersHandlers "io.winapps.aspirewithalina.aspirewithalinaserver/handlers/teachers"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	db.MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.MongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Setup HTTP server
	http.HandleFunc("/registration/create", handlers.CreateRegistrationHandler)
	http.HandleFunc("/validate/registration", handlers.ValidateRegistrationHandler)
	http.HandleFunc("/students/create", studentsHandlers.CreateNewStudentHandler)
	http.HandleFunc("/validate/login", studentsHandlers.ValidateLoginHandler)
	http.HandleFunc("/validate/login/google", studentsHandlers.ValidateGoogleLoginHandler)
	http.HandleFunc("/teachers/create", teachersHandlers.CreateTeacherHandler)
	http.HandleFunc("/teachers/validate/login", teachersHandlers.ValidateTeacherLoginHandler)
	http.HandleFunc("/teachers/update", teachersHandlers.UpdateTeacherInfoHandler)
	http.HandleFunc("/students/update", studentsHandlers.UpdateStudentInfoHandler)
	http.HandleFunc("/students", studentsHandlers.HandleFetchAllStudents)
	http.HandleFunc("/students/update/image", handlers.HandleUploadProfileImage)

	// Serve profile images
	http.Handle("/uploads/profileImages/", http.StripPrefix("/uploads/profileImages", http.FileServer(http.Dir("./uploads/profileImages"))))

	certFile := "/etc/letsencrypt/live/aspirewithalina.com/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/aspirewithalina.com/prikey.pem"

	fmt.Println("Server running on port 8888...")
	//if err := http.ListenAndServe(":8888", enableCors(http.DefaultServeMux)); err != nil {
	//	log.Fatal(err)
	//}

	//if err := http.ListenAndServe(":8888", nil); err != nil {
	//	log.Fatal(err)
	//}

	if err := http.ListenAndServeTLS(":8888", certFile, keyFile, nil); err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
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
