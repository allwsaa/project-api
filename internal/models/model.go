package models

import "time"

type User struct {
	ID               int    `json:"id"`
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email"`
	RegistrationDate string `json:"registrationDate"`
	Role             string `json:"role"`
}

type Task struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Priority       string    `json:"priority" validate:"oneof=low medium high"`
	Status         string    `json:"status" validate:"oneof=new inprogress done"`
	RespId         int       `json:"respId"`
	ProjectID      int       `json:"projectId"`
	CreationDate   time.Time `json:"creationDate"`
	CompletionDate time.Time `json:"completionDate"`
}

type Project struct {
	ID                 int       `json:"id"`
	ProjectTitle       string    `json:"projectTitle"`
	ProjectDescription string    `json:"projectDescription"`
	Started            time.Time `json:"started"`
	Completed          time.Time `json:"completed"`
	ManagerId          int       `json:"managerId"`
}
