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

var DB *sql.DB

func init() {
	DB = database.GetDB()
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	repo := repositories.TaskRepo{DB: DB}
	tasks, err := repo.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	task.CreationDate = time.Now()
	if task.CompletionDate.Before(task.CreationDate) && !task.CompletionDate.IsZero() {
		http.Error(w, "no logic", http.StatusBadRequest)
		return
	}
	if task.CompletionDate.IsZero() {
		task.CompletionDate = time.Now().AddDate(0, 1, 0)
	}

	repo := repositories.TaskRepo{DB: DB}
	id, err := repo.CreateTask(task)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	repo := repositories.TaskRepo{DB: DB}
	task, err := repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if task.CompletionDate.Before(time.Now()) && !task.CompletionDate.IsZero() {
		http.Error(w, "no logic", http.StatusBadRequest)
		return
	}
	if task.CompletionDate.IsZero() {
		task.CompletionDate = time.Now().AddDate(0, 1, 0)
	}

	repo := repositories.TaskRepo{DB: DB}
	_, err = repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	err = repo.UpdateTask(task)
	if err != nil {
		http.Error(w, "fail", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	repo := repositories.TaskRepo{DB: DB}
	_, err = repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	err = repo.DeleteTask(id)
	if err != nil {
		http.Error(w, "fail", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"message": "Deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SearchTasksHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")
	respId := r.URL.Query().Get("respId")
	projectId := r.URL.Query().Get("projectId")

	var tasks []models.Task
	var err error

	repo := repositories.TaskRepo{DB: DB}

	if title != "" {
		tasks, err = repo.FindTasksByTitle(title)
	} else if status != "" {
		tasks, err = repo.FindTasksByStatus(status)
	} else if priority != "" {
		tasks, err = repo.FindTasksByPriority(priority)
	} else if respId != "" {
		assigneeID, _ := strconv.Atoi(respId)
		tasks, err = repo.FindTasksByAssignedUserId(assigneeID)
	} else if projectId != "" {
		projectID, _ := strconv.Atoi(projectId)
		tasks, err = repo.FindTasksByProject(projectID)
	} else {
		http.Error(w, "No search criteria provided", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to search tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
