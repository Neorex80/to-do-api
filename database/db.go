package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database connection and creates tables
func InitDB() (*sql.DB, error) {
	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./tasks.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Apply performance-oriented PRAGMAs and connection pool tuning
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA temp_store=MEMORY;"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Hour)

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	createTasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		due_date DATETIME,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	// Create index on status for better query performance
	createStatusIndex := `
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	`

	// Create index on created_at for better sorting performance
	createCreatedAtIndex := `
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`

	// Execute table creation
	if _, err := db.Exec(createTasksTable); err != nil {
		return err
	}

	// Execute index creation
	if _, err := db.Exec(createStatusIndex); err != nil {
		return err
	}

	if _, err := db.Exec(createCreatedAtIndex); err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

// CloseDB closes the database connection gracefully
func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}
