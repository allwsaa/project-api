package main

import (
	"log"
	"net/http"

	"github.com/allwsaa/project-api/database"
	"github.com/allwsaa/project-api/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Project API
// @version 1.0
// @BasePath /
func main() {
	database.SetupDB()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/users", handlers.GetAllUsers)
	r.Post("/users", handlers.CreateUser)
	r.Get("/users/{id}", handlers.GetUserByID)
	r.Put("/users/{id}", handlers.UpdateUser)
	r.Delete("/users/{id}", handlers.DeleteUser)
	r.Get("/users/{id}/tasks", handlers.GetTasksByUserID)
	r.Get("/users/search", handlers.FindUsersByName)
	r.Get("/users/search", handlers.FindUsersByEmail)

	r.Get("/tasks", handlers.GetTasks)
	r.Post("/tasks", handlers.CreateTask)
	r.Get("/tasks/{id}", handlers.GetTaskByID)
	r.Put("/tasks/{id}", handlers.UpdateTask)
	r.Delete("/tasks/{id}", handlers.DeleteTask)
	r.Get("/tasks/search", handlers.SearchTasksHandler)

	r.Get("/projects", handlers.GetProjects)
	r.Post("/projects", handlers.CreateProject)
	r.Get("/projects/{id}", handlers.GetProjectByID)
	r.Put("/projects/{id}", handlers.UpdateProject)
	r.Delete("/projects/{id}", handlers.DeleteProject)
	r.Get("/projects/{id}/tasks", handlers.GetTasksByProjectID)
	r.Get("/projects/search/title", handlers.SearchProjectsByTitle)
	r.Get("/projects/search/manager", handlers.SearchProjectsByManager)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
