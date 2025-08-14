# To-Do API

A RESTful API built with Go for managing tasks (to-do items). This API provides full CRUD operations for tasks with SQLite storage and is designed for easy deployment on platforms like Render or Railway.

## Features

- ✅ Full CRUD operations (Create, Read, Update, Delete)
- ✅ SQLite database for data persistence
- ✅ JSON request/response format
- ✅ Input validation and error handling
- ✅ CORS support for web clients
- ✅ Health check endpoint
- ✅ Graceful shutdown
- ✅ Docker support
- ✅ Environment variable configuration

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Gorilla Mux
- **Database**: SQLite
- **Deployment**: Docker, Render, Railway

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite (for local development)

### Local Development

1. Clone the repository:
```bash
git clone <your-repo-url>
cd to-do-api
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Using Docker

1. **Quick Start with Docker Compose (Recommended):**
```bash
docker-compose up --build
```

2. **Manual Docker Commands:**
```bash
# Build the image
docker build -t to-do-api .

# Run the container with data persistence
docker run -d --name to-do-api -p 8080:8080 -v $(pwd)/data:/app/data to-do-api
```

3. **For detailed Docker testing instructions, see [DOCKER_TESTING.md](DOCKER_TESTING.md)**

## API Documentation

### Base URL
- Local: `http://localhost:8080`
- Production: `https://your-app.onrender.com` (or your deployment URL)

### Endpoints

#### Health Check
```
GET /health
```
Returns the health status of the API.

**Response:**
```json
{
  "status": "healthy",
  "service": "to-do-api"
}
```

#### Get API Information
```
GET /
```
Returns basic API information and available endpoints.

#### Create Task
```
POST /api/tasks
```

**Request Body:**
```json
{
  "title": "Learn Go",
  "description": "Complete Go tutorial and build an API",
  "due_date": "2024-12-31T23:59:59Z",
  "status": "pending"
}
```

**Response:**
```json
{
  "message": "Task created successfully",
  "data": {
    "id": 1,
    "title": "Learn Go",
    "description": "Complete Go tutorial and build an API",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Get All Tasks
```
GET /api/tasks
```

**Query Parameters:**
- `status` (optional): Filter tasks by status (`pending`, `in_progress`, `completed`)

**Example:**
```
GET /api/tasks?status=pending
```

**Response:**
```json
{
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Learn Go",
      "description": "Complete Go tutorial and build an API",
      "due_date": "2024-12-31T23:59:59Z",
      "status": "pending",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### Get Single Task
```
GET /api/tasks/{id}
```

**Response:**
```json
{
  "message": "Task retrieved successfully",
  "data": {
    "id": 1,
    "title": "Learn Go",
    "description": "Complete Go tutorial and build an API",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Update Task
```
PUT /api/tasks/{id}
```

**Request Body (partial updates allowed):**
```json
{
  "status": "completed"
}
```

**Response:**
```json
{
  "message": "Task updated successfully",
  "data": {
    "id": 1,
    "title": "Learn Go",
    "description": "Complete Go tutorial and build an API",
    "due_date": "2024-12-31T23:59:59Z",
    "status": "completed",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T12:45:00Z"
  }
}
```

#### Delete Task
```
DELETE /api/tasks/{id}
```

**Response:**
```json
{
  "message": "Task deleted successfully"
}
```

### Task Model

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | integer | auto | Unique identifier |
| title | string | yes | Task title |
| description | string | no | Task description |
| due_date | datetime | no | Due date in ISO 8601 format |
| status | string | no | Task status (default: "pending") |
| created_at | datetime | auto | Creation timestamp |
| updated_at | datetime | auto | Last update timestamp |

### Status Values
- `pending` - Task is not started
- `in_progress` - Task is being worked on
- `completed` - Task is finished

## Error Handling

The API returns standardized error responses:

```json
{
  "error": "Validation failed",
  "message": "title is required"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |
| DB_PATH | ./tasks.db | SQLite database file path |

## Deployment

### Deploy to Render

1. Create a new Web Service on [Render](https://render.com)
2. Connect your GitHub repository
3. Use the following settings:
   - **Build Command**: `go build -o main .`
   - **Start Command**: `./main`
   - **Environment**: Go

### Deploy to Railway

1. Create a new project on [Railway](https://railway.app)
2. Connect your GitHub repository
3. Railway will automatically detect the Go application
4. The service will be deployed using the Dockerfile

### Environment Variables for Production

Set these environment variables in your deployment platform:
- `PORT`: Usually set automatically by the platform
- `DB_PATH`: `/app/data/tasks.db` (or your preferred path)

## Testing the API

### Using curl

Create a task:
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Task",
    "description": "This is a test task",
    "status": "pending"
  }'
```

Get all tasks:
```bash
curl http://localhost:8080/api/tasks
```

Update a task:
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

Delete a task:
```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```

### Using a REST Client

You can also test the API using tools like:
- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- [Thunder Client](https://www.thunderclient.com/) (VS Code extension)

## Project Structure

```
to-do-api/
├── main.go                 # Entry point and server setup
├── models/
│   └── task.go            # Task model and repository
├── handlers/
│   └── task_handlers.go   # HTTP handlers
├── database/
│   └── db.go              # Database connection
├── middleware/
│   └── cors.go            # CORS middleware
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies
├── Dockerfile             # Docker configuration
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
