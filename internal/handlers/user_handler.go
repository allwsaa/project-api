package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"time"

	"github.com/allwsaa/project-api/database"
	"github.com/allwsaa/project-api/internal/models"
	"github.com/allwsaa/project-api/internal/repositories"
	"github.com/go-chi/chi"
)

var db *sql.DB
var repository *repositories.UserRepo

func init() {
	db = database.GetDB()
	repository = &repositories.UserRepo{DB: db}
}

// GetAllUsers godoc
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAll()
	if err != nil {
		http.Error(w, "Error occured", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// CreateUser godoc
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newUser.RegistrationDate = time.Now()

	id, err := repository.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserByID godoc
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Couldn'tt find user with given ID", http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Description Update user details by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := repository.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// DeleteUser godoc
// @Description Delete a user by their unique ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteUser(id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetTasksByUserID godoc
// @Description Get a list of tasks assigned to a user by their ID
// @Tags tasks
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.Task
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tasks not found"
// @Router /users/{id}/tasks [get]
func GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	tasks, err := repository.GetTasksByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// FindUsersByName godoc
// @Description Find users by their name
// @Tags users
// @Produce json
// @Param name query string true "Name of the user"
// @Success 200 {array} models.User
// @Failure 400 {string} string "Invalid input"
// @Router /users/search [get]
func FindUsersByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	repository := repositories.UserRepo{DB: db}
	users, err := repository.FindUsersByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// FindUsersByEmail godoc
// @Description Find users by their email address
// @Tags users
// @Produce json
// @Param email query string true "Email address of the user"
// @Success 200 {array} models.User
// @Failure 400 {string} string "Invalid input"
// @Router /users/search [get]
func FindUsersByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	repository := repositories.UserRepo{DB: db}
	users, err := repository.FindUsersByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
