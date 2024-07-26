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

// GetTasks godoc
// @Description Get a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
func GetTasks(w http.ResponseWriter, r *http.Request) {
	repo := repositories.TaskRepo{DB: DB}
	tasks, err := repo.GetTasks()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask godoc
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task data"
// @Success 201 {object} map[string]int
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create task"
// @Router /tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task.CreationDate = time.Now()
	if task.CompletionDate.Before(task.CreationDate) && !task.CompletionDate.IsZero() {
		http.Error(w, "Invalid completion date", http.StatusBadRequest)
		return
	}
	if task.CompletionDate.IsZero() {
		task.CompletionDate = time.Now().AddDate(0, 1, 0)
	}

	repo := repositories.TaskRepo{DB: DB}
	id, err := repo.CreateTask(task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetTaskByID godoc
// @Description Get task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [get]
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
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask godoc
// @Description Update task details by its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Updated task data"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Failed to update task"
// @Router /tasks/{id} [put]
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if task.CompletionDate.Before(time.Now()) && !task.CompletionDate.IsZero() {
		http.Error(w, "Invalid completion date", http.StatusBadRequest)
		return
	}
	if task.CompletionDate.IsZero() {
		task.CompletionDate = time.Now().AddDate(0, 1, 0)
	}

	repo := repositories.TaskRepo{DB: DB}
	_, err = repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = repo.UpdateTask(task)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// DeleteTask godoc
// @Description Delete a task by its unique ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 204
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Failed to delete task"
// @Router /tasks/{id} [delete]
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
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = repo.DeleteTask(id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SearchTasksHandler godoc
// @Description Search tasks based on criteria
// @Tags tasks
// @Produce json
// @Param title query string false "Title of the task"
// @Param status query string false "Status of the task"
// @Param priority query string false "Priority of the task"
// @Param respId query string false "Assigned user ID"
// @Param projectId query string false "Project ID"
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Failed to search tasks"
// @Router /tasks/search [get]
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
