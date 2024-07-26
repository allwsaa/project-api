package repositories

import (
	"database/sql"
	"fmt"

	"github.com/allwsaa/project-api/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, registrationDate, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RegistrationDate, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, name, email, registrationDate, role FROM users WHERE id = $1", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.RegistrationDate, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with this ID %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user models.User) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO users (name, email, registrationDate, role)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Name, user.Email, user.RegistrationDate, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) UpdateUser(user models.User) error {
	_, err := r.DB.Exec(`
		UPDATE users SET name = $1, email = $2, registrationDate = $3, role = $4
		WHERE id = $5
	`, user.Name, user.Email, user.RegistrationDate, user.Role, user.ID)
	return err
}

func (r *UserRepo) DeleteUser(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *UserRepo) GetTasksByUserID(userID int) ([]models.Task, error) {
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *UserRepo) FindUsersByName(name string) ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, registrationDate, role FROM users WHERE name ILIKE '%' || $1 || '%'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RegistrationDate, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepo) FindUsersByEmail(email string) ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, registrationDate, role FROM users WHERE email ILIKE '%' || $1 || '%'", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RegistrationDate, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
