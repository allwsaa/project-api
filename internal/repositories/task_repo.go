package repositories

import (
	"database/sql"
	"fmt"

	"github.com/allwsaa/project-api/internal/models"
)

type TaskRepo struct {
	DB *sql.DB
}

func (r *TaskRepo) GetTasks() ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks")
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

func (r *TaskRepo) CreateTask(task models.Task) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO tasks (title, description, priority, status, respId, projectId, creationDate, completionDate)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		task.Title, task.Description, task.Priority, task.Status, task.RespId, task.ProjectID, task.CreationDate, task.CompletionDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TaskRepo) GetTaskByID(id int) (*models.Task, error) {
	row := r.DB.QueryRow("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE id = $1", id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.RespId, &task.ProjectID, &task.CreationDate, &task.CompletionDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with this ID not found")
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) UpdateTask(task models.Task) error {
	_, err := r.DB.Exec(`
		UPDATE tasks SET title = $1, description = $2, priority = $3, status = $4, respId = $5, projectId = $6, creationDate = $7, completionDate = $8
		WHERE id = $9`,
		task.Title, task.Description, task.Priority, task.Status, task.RespId, task.ProjectID, task.CreationDate, task.CompletionDate, task.ID)
	return err
}

func (r *TaskRepo) DeleteTask(id int) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (r *TaskRepo) FindTasksByTitle(title string) ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE title ILIKE '%' || $1 || '%'", title)
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

func (r *TaskRepo) FindTasksByStatus(status string) ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE status = $1", status)
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

func (r *TaskRepo) FindTasksByPriority(priority string) ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE priority = $1", priority)
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

func (r *TaskRepo) FindTasksByAssignedUserId(userID int) ([]models.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, priority, status, respId, projectId, creationDate, completionDate FROM tasks WHERE respId = $1", userID)
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

func (r *TaskRepo) FindTasksByProject(projectID int) ([]models.Task, error) {
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
