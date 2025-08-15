package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"to-do-api/models"

	"github.com/gorilla/mux"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	repo models.TaskRepository
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(repo models.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateTask handles POST /api/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskReq models.TaskRequest
	
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}
	
	if err := taskReq.Validate(); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}
	
	task, err := h.repo.Create(&taskReq)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create task", "")
		return
	}
	
	h.sendSuccessResponse(w, http.StatusCreated, "Task created successfully", task)
}

// GetTasks handles GET /api/tasks
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Query params: status, limit, offset, sort_by, sort_order
	q := r.URL.Query()
	status := q.Get("status")
	limit := 50
	offset := 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			if n < 1 {
				limit = 1
			} else if n > 100 {
				limit = 100
			} else {
				limit = n
			}
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	sortBy := q.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := q.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	var filterStatusPtr *string
	if status != "" {
		// Validate status
		if !isValidStatus(status) {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid status", "Status must be one of: pending, in_progress, completed")
			return
		}
		filterStatusPtr = &status
	}

	tasks, err := h.repo.GetAllPaginated(filterStatusPtr, limit, offset, sortBy, sortOrder)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch tasks", "")
		return
	}
	
	// Return empty array instead of null if no tasks
	if tasks == nil {
		tasks = []models.Task{}
	}
	
	h.sendSuccessResponse(w, http.StatusOK, "Tasks retrieved successfully", tasks)
}

// GetTask handles GET /api/tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid task ID", "Task ID must be a number")
		return
	}
	
	task, err := h.repo.GetByID(id)
	if err != nil {
		log.Printf("Error fetching task: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch task", "")
		return
	}
	
	if task == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Task not found", "")
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, "Task retrieved successfully", task)
}

// UpdateTask handles PUT /api/tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid task ID", "Task ID must be a number")
		return
	}
	
	var taskReq models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}
	
	// For updates, we allow partial updates, so we don't require title
	if taskReq.Status != "" && !isValidStatus(taskReq.Status) {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid status", "Status must be one of: pending, in_progress, completed")
		return
	}
	
	task, err := h.repo.Update(id, &taskReq)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to update task", "")
		return
	}
	
	if task == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Task not found", "")
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, "Task updated successfully", task)
}

// DeleteTask handles DELETE /api/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid task ID", "Task ID must be a number")
		return
	}
	
	err = h.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendErrorResponse(w, http.StatusNotFound, "Task not found", "")
			return
		}
		log.Printf("Error deleting task: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete task", "")
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, "Task deleted successfully", nil)
}

// HealthCheck handles GET /health
func (h *TaskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "to-do-api",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends a standardized error response
func (h *TaskHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, error string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// sendSuccessResponse sends a standardized success response
func (h *TaskHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := SuccessResponse{
		Message: message,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
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
