package types

//===============//
// STUDENT TYPES //
//===============//

// Student struct that determines how students will be stored in the database
type Student struct {
	StudentId          string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

// StudentInfo struct that determines the info that can be retrieved about a student securely (i.e. no passwords, salts, etc.)
type StudentInfo struct {
	StudentId          string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

// ListStudentsRequest struct to handle incoming request to list all students
type ListStudentsRequest struct{}

// ListStudentsResponse struct to handle outgoing response to list all students
type ListStudentsResponse struct {
	Students []Student `json:"students"`
	Page     int64     `json:"page"`
}

// GetStudentRequest struct to handle incoming request for a student
type GetStudentRequest struct{}

// GetStudentResponse struct to handle outgoing response for a student
type GetStudentResponse struct {
	Student StudentInfo `json:"student"`
}

// VerifyUserEmailResponse struct to handle outgoing user verification response
type VerifyUserEmailResponse struct{}

// UserVerification struct
type UserVerification struct {
	VerificationCode string `bson:"verificationCode" json:"verificationCode"`
}

// CreateVerificationRequest struct to handle incoming request to create a verification object
type CreateVerificationRequest struct {
	Token            string `bson:"token" json:"token"`
	Email            string `bson:"email" json:"email"`
	IsVerified       bool   `bson:"isVerified" json:"isVerified"`
	RegistrationCode string `bson:"registrationCode" json:"registrationCode"`
	IsRegistered     bool   `bson:"isRegistered" json:"isRegistered"`
}

// CreateVerificationResponse struct to handle outgoing response to create a verification object
type CreateVerificationResponse struct {
	IsCreated bool `bson:"isCreated" json:"isCreated"`
}

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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

// CreateNewStudentResponse Struct to handle outgoing create new student response
type CreateNewStudentResponse struct {
	StudentId string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
}

// ValidateLoginRequest Struct to handle incoming login request
type ValidateLoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

// ValidateLoginResponse Struct to handle outgoing login response
type ValidateLoginResponse struct {
	StudentId          string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

type ValidateLoginResult struct {
	IsValid     bool                  `json:"is_valid"`
	StudentInfo ValidateLoginResponse `json:"student_info"`
}

type ValidateGoogleLoginRequest struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	// Later add picture and name and JWT here
}

// UpdateStudentInfoRequest Struct to handle incoming updates to student info
type UpdateStudentInfoRequest struct {
	StudentId          string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

// UpdateStudentInfoResponse Struct to handle outgoing response after updating student info
type UpdateStudentInfoResponse struct {
	StudentId          string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
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
	LessonsRemaining   int64  `json:"lessons_remaining"`
	LessonsCompleted   int64  `json:"lessons_completed"`
}

// DeleteStudentRequest struct to handle incoming request to delete a student
type DeleteStudentRequest struct {
	StudentId string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
}

// DeleteStudentResponse struct to handle outgoing response to delete a student
type DeleteStudentResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

//===============//
// TEACHER TYPES //
//===============//

// Teacher struct that determines how a teacher is stored in the database
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
	LessonsTaught      int64  `json:"lessons_taught"`
}

// TeacherInfo struct that determines the info that can be retrieved about a teacher securely (i.e. no passwords, salts, etc.)
type TeacherInfo struct {
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
	LessonsTaught      int64  `json:"lessons_taught"`
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
	LessonsTaught      int64  `json:"lessons_taught"`
}

type GetTeacherRequest struct {
	TeacherID string `json:"teacherID"`
}

type GetTeacherResponse struct {
	Teacher TeacherInfo `json:"teacher"`
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
	LessonsTaught      int64  `json:"lessons_taught"`
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
	LessonsTaught      int64  `json:"lessons_taught"`
}

// ListTeachersRequest struct to handle incoming request for listing all teachers
type ListTeachersRequest struct{}

// ListTeachersResponse struct to handle outgoing response for listing all teachers
type ListTeachersResponse struct {
	Teachers []TeacherInfo `json:"teachers"`
}

// DeleteTeacherRequest struct to handle incoming request to delete a teacher
type DeleteTeacherRequest struct{}

// DeleteTeacherResponse struct to handle outgoing response to delete a teacher
type DeleteTeacherResponse struct {
	IsDeleted bool `json:"is_deleted"`
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
	LessonsTaught      int64  `json:"lessons_taught"`
}

type ValidateTeacherLoginResult struct {
	TeacherInfo ValidateTeacherLoginResponse
	IsValid     bool
}

//==============//
// LESSON TYPES //
//==============//

// Lesson struct to be stored in lessonsCollection
type Lesson struct {
	LessonID          string `json:"lessonID"`
	TeacherID         string `json:"teacherID"`
	StudentId         string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
	Subject           string `json:"subject"`
	ScheduledDateTime int64  `json:"scheduled_date_time"`
	Room              int64  `json:"room"`
	IsCanceled        bool   `json:"is_canceled"`
	IsCompleted       bool   `json:"is_completed"`
	TimesRescheduled  int64  `json:"times_rescheduled"`
	IsStudentLate     bool   `json:"is_student_late"`    // Only true when 5 minutes or more late
	IsTeacherLate     bool   `json:"is_teacher_late"`    // Only true when 5 minutes or more late
	IsConnectionLost  bool   `json:"is_connection_lost"` // Need to really think about this implementation
}

// CreateLessonRequest struct to handle incoming request for creating a new lesson
type CreateLessonRequest struct {
	TeacherID         string `json:"teacherID"`
	StudentId         string `json:"student_id"` // TODO: Update to be like TeacherID; needs done in Electron apps too
	Subject           string `json:"subject"`
	ScheduledDateTime int64  `json:"scheduled_date_time"`
	Room              int64  `json:"room"`
}

// CreateLessonResponse struct to handle outgoing response for creating a new lesson
type CreateLessonResponse struct {
	Lesson Lesson `json:"lesson"`
}

// UpdateLessonRequest struct to handle incoming request for updating an existing lesson
type UpdateLessonRequest struct {
	LessonID          string `json:"lessonID"`
	Subject           string `json:"subject"`
	ScheduledDateTime int64  `json:"scheduled_date_time"`
	Room              int64  `json:"room"`
	IsCanceled        bool   `json:"is_canceled"`
	IsCompleted       bool   `json:"is_completed"`
	TimesRescheduled  int64  `json:"times_rescheduled"`
	IsStudentLate     bool   `json:"is_student_late"`    // Only true when 5 minutes or more late
	IsTeacherLate     bool   `json:"is_teacher_late"`    // Only true when 5 minutes or more late
	IsConnectionLost  bool   `json:"is_connection_lost"` // Need to really think about this implementation
}

// UpdateLessonResponse struct to handle outgoing response for updating an existing lesson
type UpdateLessonResponse struct {
	Lesson Lesson `json:"lesson"`
}

// DeleteLessonRequest struct to handle incoming request for deleting an existing lesson
type DeleteLessonRequest struct {
	LessonID string `json:"lessonID"`
}

// DeleteLessonResponse struct to handle outgoing response for deleting an existing lesson
type DeleteLessonResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

// ListLessonsRequest struct to handle incoming request for listing lessons
type ListLessonsRequest struct {
	UserID      string `json:"ID"` // This is just "UserID" because teacherID and studentID will be used interchangeably
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
	IsCanceled  bool   `json:"is_canceled"`  // IsCanceled=true & IsCompleted=false returns only canceled classes
	IsCompleted bool   `json:"is_completed"` // IsCanceled=false & IsCompleted=true returns only completed classes,
	// IsCanceled=false & IsCompleted=false returns only upcoming classes
}

// ListLessonsResponse struct to handle outgoing response for listing lessons
type ListLessonsResponse struct {
	Lessons []Lesson `json:"lessons"`
	Page    int64    `json:"page"`
}

//====================//
// CHAT/MESSAGE TYPES //
//====================//

// CreateUserRequest struct to handle incoming request to create a new user for chat
type CreateUserRequest struct {
	UserId            string `bson:"userId" json:"userId"`
	UserType          string `bson:"userType" json:"userType"`
	PreferredName     string `bson:"preferredName" json:"preferredName"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	ProfilePictureURL string `bson:"profilePictureUrl" json:"profilePictureUrl"`
}

// CreateUserResponse struct to handle outgoing response to create a new user for chat
type CreateUserResponse struct {
	IsCreated bool `bson:"isCreated" json:"isCreated"`
}

// UpdateUserRequest struct to handle incoming request to update a chat user
type UpdateUserRequest struct {
	UserId            string `bson:"userId" json:"userId"`
	UserType          string `bson:"userType" json:"userType"`
	PreferredName     string `bson:"preferredName" json:"preferredName"`
	ProfilePictureURL string `bson:"profilePictureUrl" json:"profilePictureUrl"`
}

// UpdateUserResponse struct to handle outgoing response to update a chat user
type UpdateUserResponse struct {
	IsUpdated bool `bson:"isUpdated" json:"isUpdated"`
}

// CreateChatRoomRequest struct to handle incoming request to create a new chat room
type CreateChatRoomRequest struct{}

// CreateChatRoomResponse struct to handle outgoing response to create a new chat room
type CreateChatRoomResponse struct{}

// DeleteChatRoomRequest struct to handle incoming request to delete a chat room
type DeleteChatRoomRequest struct{}

// DeleteChatRoomResponse struct to handle outgoing response to delete a chat room
type DeleteChatRoomResponse struct{}

// ListChatRoomsRequest struct to handle incoming request to list all chat rooms
type ListChatRoomsRequest struct{}

// ListChatRoomsResponse struct to handle outgoing response to list all chat rooms
type ListChatRoomsResponse struct{}

// SendMessageRequest struct to handle incoming request to send a chat message
type SendMessageRequest struct{}

// SendMessageResponse struct to handle outgoing response to send a chat message
type SendMessageResponse struct{}

// UpdateMessageRequest struct to handle incoming request to update a chat message
type UpdateMessageRequest struct{}

// UpdateMessageResponse struct to handle outgoing response to update a chat message
type UpdateMessageResponse struct{}

// DeleteMessageRequest struct to handle incoming request to delete a chat message
type DeleteMessageRequest struct{}

// DeleteMessageResponse struct to handle outgoing response to delete a chat message
type DeleteMessageResponse struct{}

// ListMessagesRequest struct to handle incoming request for all chat messages in a chat room
type ListMessagesRequest struct{}

// ListMessagesResponse struct to handle outgoing response for all chat messages in a chat room
type ListMessagesResponse struct{}

//==================//
// ASSIGNMENT TYPES //
//==================//

// Assignment Struct to be stored in studentAssignmentsCollection
type Assignment struct {
	AssignmentID  string `json:"assignmentID"`
	Title         string `json:"title"`
	Subject       string `json:"subject"`
	DocumentUrl   string `json:"document_url"`
	DateAssigned  int64  `json:"date_assigned"`
	DateStarted   int64  `json:"date_started"`
	DateCompleted int64  `json:"date_completed"`
}

// StudentAssignments Struct
type StudentAssignments struct {
	StudentId   string       `json:"student_id"`  // TODO: Update to be like TeacherID; needs done in Electron apps too
	Assignments []Assignment `json:"assignments"` // TODO: Probably remove this and change it to array of string IDs so it's not too large
}

//============//
// GAME TYPES //
//============//

// SpaceShooterGame Struct to be stored in studentGamesCollection for SpaceShooter
type SpaceShooterGame struct {
	Level         string `json:"level"`
	Score         int    `json:"score"`
	DateStarted   int64  `json:"date_started"`
	DateCompleted int64  `json:"date_completed"`
}

// WordioGame Struct to be stored in studentGamesCollection for Wordio
type WordioGame struct {
	Level         string `json:"level"`
	Score         int    `json:"score"`
	DateStarted   int64  `json:"date_started"`
	DateCompleted int64  `json:"date_completed"`
}

// SpellingPuddlesWord Struct word for SpellingPuddlesGame
type SpellingPuddlesWord struct {
	Word          string   `json:"word"`
	LettersChosen []string `json:"letters_chosen"`
	AudioUri      string   `json:"audio_uri"`
}

// SpellingPuddlesGame Struct to be stored in studentGamesCollection for SpellingPuddles
type SpellingPuddlesGame struct {
	Level         string `json:"level"`
	Score         int    `json:"score"`
	DateStarted   int64  `json:"date_started"`
	DateCompleted int64  `json:"date_completed"`
}

// StudentGames Struct
type StudentGames struct {
	StudentId       string                `json:"student_id"`       // TODO: Update to be like TeacherID; needs done in Electron apps too
	SpaceShooter    []SpaceShooterGame    `json:"space_shooter"`    // Well... a Space Shooter game lol
	Wordio          []WordioGame          `json:"wordio"`           // Mario-like game
	SpellingPuddles []SpellingPuddlesGame `json:"spelling_puddles"` // Rain drops containing characters fall down to spell words game
}
