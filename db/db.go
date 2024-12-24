package db

import "go.mongodb.org/mongo-driver/mongo"

var MongoClient *mongo.Client
var DbName = "aspireDB"
var RegistrationCollection = "registrations"
var VerificationsCollection = "verifications"
var TeachersCollection = "teachers"
var StudentsCollection = "students"
var UsersCollection = "users"
var LessonsCollection = "lessons"
var StudentAssignmentsCollection = "assignments"
var StudentGamesCollection = "games"
