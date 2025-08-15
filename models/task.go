package models

import (
	"database/sql"
	"strings"
	"time"
)

// Task represents a task in the to-do list
type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TaskRequest represents the request payload for creating/updating tasks
type TaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      string     `json:"status"`
}

// Validate validates the task request
func (tr *TaskRequest) Validate() error {
	if tr.Title == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	
	if tr.Status != "" && !isValidStatus(tr.Status) {
		return &ValidationError{Field: "status", Message: "status must be one of: pending, in_progress, completed"}
	}
	
	return nil
}

// isValidStatus checks if the status is valid
func isValidStatus(status string) bool {
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// TaskRepository defines the interface for task database operations
type TaskRepository interface {
	Create(task *TaskRequest) (*Task, error)
	GetAll() ([]Task, error)
	GetByID(id int) (*Task, error)
	Update(id int, task *TaskRequest) (*Task, error)
	Delete(id int) error
	GetByStatus(status string) ([]Task, error)
	GetAllPaginated(filterStatus *string, limit int, offset int, sortBy string, sortOrder string) ([]Task, error)
}

// SQLiteTaskRepository implements TaskRepository for SQLite
type SQLiteTaskRepository struct {
	db *sql.DB
}

// NewSQLiteTaskRepository creates a new SQLite task repository
func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

// Create creates a new task
func (r *SQLiteTaskRepository) Create(taskReq *TaskRequest) (*Task, error) {
	// Set default status if not provided
	status := taskReq.Status
	if status == "" {
		status = "pending"
	}
	
	query := `
		INSERT INTO tasks (title, description, due_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	result, err := r.db.Exec(query, taskReq.Title, taskReq.Description, taskReq.DueDate, status, now, now)
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(int(id))
}

// GetAll retrieves all tasks
func (r *SQLiteTaskRepository) GetAll() ([]Task, error) {
	query := `
		SELECT id, title, description, due_date, status, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetAllPaginated retrieves tasks with optional filtering, sorting, and pagination
func (r *SQLiteTaskRepository) GetAllPaginated(filterStatus *string, limit int, offset int, sortBy string, sortOrder string) ([]Task, error) {
	allowedSort := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"due_date":   true,
		"id":         true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	base := `
		SELECT id, title, description, due_date, status, created_at, updated_at
		FROM tasks
	`
	args := make([]interface{}, 0, 3)
	if filterStatus != nil && *filterStatus != "" {
		base += " WHERE status = ?"
		args = append(args, *filterStatus)
	}
	base += " ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetByID retrieves a task by ID
func (r *SQLiteTaskRepository) GetByID(id int) (*Task, error) {
	query := `
		SELECT id, title, description, due_date, status, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	
	var task Task
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &task, nil
}

// Update updates a task
func (r *SQLiteTaskRepository) Update(id int, taskReq *TaskRequest) (*Task, error) {
	// First check if task exists
	existingTask, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingTask == nil {
		return nil, nil
	}
	
	// Update only provided fields
	title := taskReq.Title
	if title == "" {
		title = existingTask.Title
	}
	
	description := taskReq.Description
	status := taskReq.Status
	if status == "" {
		status = existingTask.Status
	}
	
	dueDate := taskReq.DueDate
	if dueDate == nil {
		dueDate = existingTask.DueDate
	}
	
	query := `
		UPDATE tasks
		SET title = ?, description = ?, due_date = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	_, err = r.db.Exec(query, title, description, dueDate, status, now, id)
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(id)
}

// Delete deletes a task
func (r *SQLiteTaskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

// GetByStatus retrieves tasks by status
func (r *SQLiteTaskRepository) GetByStatus(status string) ([]Task, error) {
	query := `
		SELECT id, title, description, due_date, status, created_at, updated_at
		FROM tasks
		WHERE status = ?
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}
