package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"to-do-api/database"
	"to-do-api/handlers"
	"to-do-api/middleware"
	"to-do-api/models"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Initialize repository and handlers
	taskRepo := models.NewSQLiteTaskRepository(db)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// Create router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)
	router.Use(middleware.Gzip)

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

	// Static file serving
	staticFS := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(middleware.WithCacheControl(http.StripPrefix("/static/", staticFS), "public, max-age=604800, immutable"))
	
	// Root route serves the frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./static/index.html")
	}).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Health check: http://localhost:%s/health", port)
		log.Printf("UI: http://localhost:%s/", port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
