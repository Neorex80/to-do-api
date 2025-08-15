package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"to-do-api/handlers"
	"to-do-api/middleware"
	"to-do-api/models"

	"github.com/gorilla/mux"
)

// InMemoryTaskRepository implements TaskRepository using in-memory storage
// This is used for testing purposes to avoid database dependencies
type InMemoryTaskRepository struct {
	tasks  map[int]*models.Task
	nextID int
	mutex  sync.RWMutex
}

// NewInMemoryTaskRepository creates a new in-memory task repository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

// Create creates a new task
func (r *InMemoryTaskRepository) Create(taskReq *models.TaskRequest) (*models.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	status := taskReq.Status
	if status == "" {
		status = "pending"
	}

	now := time.Now()
	task := &models.Task{
		ID:          r.nextID,
		Title:       taskReq.Title,
		Description: taskReq.Description,
		DueDate:     taskReq.DueDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.tasks[r.nextID] = task
	r.nextID++

	return task, nil
}

// GetAll retrieves all tasks
func (r *InMemoryTaskRepository) GetAll() ([]models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// GetByID retrieves a task by ID
func (r *InMemoryTaskRepository) GetByID(id int) (*models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}

	return task, nil
}

// Update updates a task
func (r *InMemoryTaskRepository) Update(id int, taskReq *models.TaskRequest) (*models.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}

	// Update fields if provided
	if taskReq.Title != "" {
		task.Title = taskReq.Title
	}
	if taskReq.Description != "" {
		task.Description = taskReq.Description
	}
	if taskReq.DueDate != nil {
		task.DueDate = taskReq.DueDate
	}
	if taskReq.Status != "" {
		task.Status = taskReq.Status
	}

	task.UpdatedAt = time.Now()
	r.tasks[id] = task

	return task, nil
}

// Delete deletes a task
func (r *InMemoryTaskRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return nil // Return nil for not found to match SQL behavior
	}

	delete(r.tasks, id)
	return nil
}

// GetByStatus retrieves tasks by status
func (r *InMemoryTaskRepository) GetByStatus(status string) ([]models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tasks []models.Task
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

// GetAllPaginated retrieves tasks with optional filtering, sorting, and pagination
func (r *InMemoryTaskRepository) GetAllPaginated(filterStatus *string, limit int, offset int, sortBy string, sortOrder string) ([]models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// For simplicity in test mode, we'll just return all tasks with basic filtering
	// In a real implementation, we would implement proper pagination and sorting
	
	var tasks []models.Task
	for _, task := range r.tasks {
		// Apply status filter if provided
		if filterStatus != nil && *filterStatus != "" && task.Status != *filterStatus {
			continue
		}
		
		tasks = append(tasks, *task)
	}

	// Apply basic sorting (by ID for simplicity)
	// In a real implementation, we would sort by the specified field and order
	
	// Apply pagination
	if offset < len(tasks) {
		end := offset + limit
		if end > len(tasks) {
			end = len(tasks)
		}
		tasks = tasks[offset:end]
	} else {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func main() {
	log.Println("Starting To-Do API with in-memory storage...")

	// Initialize in-memory repository
	taskRepo := NewInMemoryTaskRepository()
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// Create some sample data
	sampleTasks := []*models.TaskRequest{
		{
			Title:       "Learn Go",
			Description: "Complete Go tutorial and build an API",
			Status:      "pending",
		},
		{
			Title:       "Build REST API",
			Description: "Create a full-featured REST API with CRUD operations",
			Status:      "in_progress",
		},
		{
			Title:       "Deploy to Production",
			Description: "Deploy the API to Render or Railway",
			Status:      "pending",
		},
	}

	for _, taskReq := range sampleTasks {
		taskRepo.Create(taskReq)
	}

	// Create router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Task routes
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")

	// Health check route
	router.HandleFunc("/health", taskHandler.HealthCheck).Methods("GET")

	// Root route for basic info
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service": "To-Do API (Test Mode)",
			"version": "1.0.0",
			"storage": "in-memory",
			"endpoints": {
				"health": "GET /health",
				"tasks": {
					"create": "POST /api/tasks",
					"list": "GET /api/tasks",
					"get": "GET /api/tasks/{id}",
					"update": "PUT /api/tasks/{id}",
					"delete": "DELETE /api/tasks/{id}"
				}
			},
			"note": "This is running with in-memory storage for testing. Use main.go with SQLite for production."
		}`))
	}).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("API documentation: http://localhost:%s/", port)
	log.Printf("Sample tasks have been created for testing")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
