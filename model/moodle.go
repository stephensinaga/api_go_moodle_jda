package model

type MoodleUser struct {
	Username  string
	Password  string
	Firstname string
	Lastname  string
	Email     string
}

type Enrolment struct {
	UserID   int `json:"userid"`
	CourseID int `json:"courseid"`
	RoleID   int `json:"roleid"`
}
