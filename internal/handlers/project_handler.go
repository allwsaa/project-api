package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/allwsaa/project-api/database"
	"github.com/allwsaa/project-api/internal/models"
	"github.com/allwsaa/project-api/internal/repositories"
	"github.com/go-chi/chi"
)

var Datab *sql.DB

func init() {
	Datab = database.GetDB()
}

// GetProjects godoc
// @Description Get a list of all projects
// @Tags projects
// @Produce json
// @Success 200 {array} models.Project
// @Failure 500 {string} string "Internal server error"
// @Router /projects [get]
func GetProjects(w http.ResponseWriter, r *http.Request) {
	repo := repositories.ProjectRepo{DB: DB}
	projects, err := repo.GetAllProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

// CreateProject godoc
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project data"
// @Success 201 {object} models.Project
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /projects [post]
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	id, err := repo.CreateProject(project)
	if err != nil {
		http.Error(w, "failed to create project", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetProjectByID godoc
// @Description Get project by ID
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Project not found"
// @Router /projects/{id} [get]
func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	project, err := repo.GetProjectByID(id)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// UpdateProject godoc
// @Summary Update project
// @Description Update project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project data"
// @Success 200 {object} models.Project
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Project not found"
// @Failure 500 {string} string "Internal server error"
// @Router /projects/{id} [put]
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	project.ID = id

	repo := repositories.ProjectRepo{DB: DB}
	if err := repo.UpdateProject(project); err != nil {
		http.Error(w, "failed to update project", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// DeleteProject godoc
// @Summary Delete project
// @Description Delete project
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Project not found"
// @Router /projects/{id} [delete]
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	if err := repo.DeleteProject(id); err != nil {
		http.Error(w, "failed to delete project", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"message": "Deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetTasksByProjectID godoc
// @Summary Get tasks by project ID
// @Description Get a list of tasks associated with a project by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} models.Task
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tasks not found"
// @Router /projects/{id}/tasks [get]
func GetTasksByProjectID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	tasks, err := repo.GetTasksByProjectID(id)
	if err != nil {
		http.Error(w, "failed to get tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// SearchProjectsByTitle godoc
// @Summary Search projects by title
// @Description Search projects based on title
// @Tags projects
// @Produce json
// @Param title query string true "Title of the project"
// @Success 200 {array} models.Project
// @Failure 400 {string} string "Invalid input"
// @Router /projects/search/title [get]
func SearchProjectsByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	projects, err := repo.SearchProjectsByTitle(title)
	if err != nil {
		http.Error(w, "failed to search projects", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

// SearchProjectsByManager godoc
// @Summary Search projects by manager
// @Description Search projects based on manager's ID
// @Tags projects
// @Produce json
// @Param managerId query int true "Manager's ID"
// @Success 200 {array} models.Project
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /projects/search/manager [get]
func SearchProjectsByManager(w http.ResponseWriter, r *http.Request) {
	managerIDStr := r.URL.Query().Get("managerId")
	if managerIDStr == "" {
		http.Error(w, "manager id is required", http.StatusBadRequest)
		return
	}

	managerID, err := strconv.Atoi(managerIDStr)
	if err != nil {
		http.Error(w, "Invalid manager id", http.StatusBadRequest)
		return
	}

	repo := repositories.ProjectRepo{DB: DB}
	projects, err := repo.SearchProjectsByManager(managerID)
	if err != nil {
		http.Error(w, "failed to search projects", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}
