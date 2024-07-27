package models

import "time"

type User struct {
	ID               int       `json:"id" readonly:"true"`
	Name             string    `json:"name" validate:"required"`
	Email            string    `json:"email" validate:"required,email" example:"string@gmail.com"`
	RegistrationDate time.Time `json:"registrationDate" readonly:"true"`
	Role             string    `json:"role" validate:"required"`
}

type Task struct {
	ID             int       `json:"id" readonly:"true"`
	Title          string    `json:"title" validate:"required"`
	Description    string    `json:"description"`
	Priority       string    `json:"priority" validate:"oneof=low medium high"`
	Status         string    `json:"status" validate:"oneof=new inprogress done"`
	RespId         int       `json:"respId" validate:"required" example:"1"`
	ProjectID      int       `json:"projectId"`
	CreationDate   time.Time `json:"creationDate" readonly:"true"`
	CompletionDate time.Time `json:"completionDate" example:"2024-09-20T15:04:05Z"`
}

type Project struct {
	ID                 int       `json:"id"  readonly:"true"`
	ProjectTitle       string    `json:"projectTitle" validate:"required"`
	ProjectDescription string    `json:"projectDescription"`
	Started            time.Time `json:"started" readonly:"true"`
	Completed          time.Time `json:"completed"  example:"2024-09-20T15:04:05Z"`
	ManagerId          int       `json:"managerId" validate:"required" example:"1"`
}
