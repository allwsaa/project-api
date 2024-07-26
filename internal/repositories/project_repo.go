package repositories

import (
	"database/sql"
	"fmt"

	"github.com/allwsaa/project-api/internal/models"
)

type ProjectRepo struct {
	DB *sql.DB
}

func (r *ProjectRepo) GetAllProjects() ([]models.Project, error) {
	rows, err := r.DB.Query("SELECT id, projectTitle, projectDescription, started FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.ProjectTitle, &project.ProjectDescription, &project.Started); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepo) CreateProject(project models.Project) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO projects (projectTitle, projectDescription, started, completed, managerId)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		project.ProjectTitle, project.ProjectDescription, project.Started, project.Completed, project.ManagerId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProjectRepo) GetProjectByID(id int) (*models.Project, error) {
	row := r.DB.QueryRow("SELECT id, projectTitle, projectDescription, started, completed, managerId FROM projects WHERE id = $1", id)
	var project models.Project
	err := row.Scan(&project.ID, &project.ProjectTitle, &project.ProjectDescription, &project.Started, &project.Completed, &project.ManagerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with this ID not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) UpdateProject(project models.Project) error {
	_, err := r.DB.Exec(`
		UPDATE projects SET projectTitle = $1, projectDescription = $2, started = $3, completed = $4, managerId = $5
		WHERE id = $6`,
		project.ProjectTitle, project.ProjectDescription, project.Started, project.Completed, project.ManagerId, project.ID)
	return err
}

func (r *ProjectRepo) DeleteProject(id int) error {
	_, err := r.DB.Exec("DELETE FROM projects WHERE id = $1", id)
	return err
}

func (r *ProjectRepo) SearchProjectsByTitle(title string) ([]models.Project, error) {
	rows, err := r.DB.Query("SELECT id, projectTitle, projectDescription, started, completed, managerId FROM projects WHERE projectTitle ILIKE $1", "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.ProjectTitle, &project.ProjectDescription, &project.Started, &project.Completed, &project.ManagerId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepo) SearchProjectsByManager(managerId int) ([]models.Project, error) {
	rows, err := r.DB.Query("SELECT id, projectTitle, projectDescription, started, completed, managerId FROM projects WHERE managerId = $1", managerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.ProjectTitle, &project.ProjectDescription, &project.Started, &project.Completed, &project.ManagerId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepo) GetTasksByProjectID(projectID int) ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE projectId = $1", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.RespId, &task.ProjectID, &task.CreationDate, &task.CompletionDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
