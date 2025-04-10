package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io.winapps.aspirewithalina.aspirewithalinaserver/db"
	"io.winapps.aspirewithalina.aspirewithalinaserver/handlers"
	chatsHandlers "io.winapps.aspirewithalina.aspirewithalinaserver/handlers/chats"
	lessonsHandlers "io.winapps.aspirewithalina.aspirewithalinaserver/handlers/lessons"
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

	// Setup HTTPS server handlers
	// Registration handlers
	http.HandleFunc("/registration/create", handlers.CreateRegistrationHandler)
	http.HandleFunc("/validate/registration", handlers.ValidateRegistrationHandler)
	http.HandleFunc("/verifications/create", handlers.CreateVerificationHandler)

	// Login handlers
	http.HandleFunc("/validate/login", studentsHandlers.ValidateLoginHandler)
	http.HandleFunc("/validate/login/mobile", studentsHandlers.ValidateLoginMobileHandler)
	http.HandleFunc("/validate/login/google", studentsHandlers.ValidateGoogleLoginHandler)
	http.HandleFunc("/teachers/validate/login", teachersHandlers.ValidateTeacherLoginHandler)

	// Teacher CRUD handlers
	http.HandleFunc("/teacher", teachersHandlers.GetTeacherHandler)
	http.HandleFunc("/teachers", teachersHandlers.ListTeachersHandler)
	http.HandleFunc("/teachers/create", teachersHandlers.CreateTeacherHandler)
	http.HandleFunc("/teachers/delete", teachersHandlers.DeleteTeacherHandler)
	http.HandleFunc("/teachers/update", teachersHandlers.UpdateTeacherInfoHandler)

	// Student CRUD handlers
	http.HandleFunc("/students/create", studentsHandlers.CreateNewStudentHandler)
	http.HandleFunc("/students/update", studentsHandlers.UpdateStudentInfoHandler)
	http.HandleFunc("/students", studentsHandlers.ListStudentsHandler)
	http.HandleFunc("/student", studentsHandlers.GetStudentHandler)
	http.HandleFunc("/students/update/image", handlers.HandleUploadProfileImage)
	http.HandleFunc("/students/delete", studentsHandlers.HandleDeleteStudent)

	// Lessons CRUD handlers
	http.HandleFunc("/lessons/create", lessonsHandlers.CreateLessonHandler)
	http.HandleFunc("/lessons/update", lessonsHandlers.UpdateLessonHandler)
	http.HandleFunc("/lessons/delete", lessonsHandlers.DeleteLessonHandler)
	http.HandleFunc("/lessons", lessonsHandlers.ListLessonsHandler)

	// Chats/Messaging CRUD handlers
	http.HandleFunc("/chats/create", chatsHandlers.CreateChatRoomHandler)
	http.HandleFunc("/chats/delete", chatsHandlers.DeleteChatRoomHandler)
	http.HandleFunc("/chats", chatsHandlers.ListChatRoomsHandler)
	http.HandleFunc("/messages/send", chatsHandlers.SendMessageHandler)
	http.HandleFunc("/messages/delete", chatsHandlers.DeleteMessageHandler)
	http.HandleFunc("/messages/update", chatsHandlers.UpdateMessageHandler)
	http.HandleFunc("/messages", chatsHandlers.ListMessagesHandler)
	http.HandleFunc("/chatUsers/create", chatsHandlers.CreateUserHandler)
	http.HandleFunc("/chatUsers/update", chatsHandlers.UpdateUserHandler)

	// Serve profile images
	http.Handle("/uploads/profileImages/", http.StripPrefix("/uploads/profileImages", http.FileServer(http.Dir("./uploads/profileImages"))))

	certFile := "/etc/letsencrypt/live/aspirewithalina.com/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/aspirewithalina.com/privkey.pem"

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
