<div align="center">

# 🚀 To-Do API

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
![Railway](https://img.shields.io/badge/Railway-131415?style=for-the-badge&logo=railway&logoColor=white)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Neorex80/to-do-api?style=for-the-badge)](https://goreportcard.com/report/github.com/Neorex80/to-do-api)
[![Docker Build](https://img.shields.io/badge/Docker-Build%20Ready-brightgreen?style=for-the-badge&logo=docker)](https://hub.docker.com/)

**🎯 A lightning-fast RESTful API built with Go for managing tasks**

[🚀 Quick Start](#-quick-start) • [📚 API Docs](#-api-endpoints) • [🐳 Docker](#-docker-deployment) • [☁️ Deploy](#️-deployment)

</div>

---

## ✨ Features

🔥 **Full CRUD Operations** - Create, Read, Update, Delete tasks  
💾 **SQLite Database** - Lightweight & persistent storage  
🐳 **Docker Ready** - One-click containerized deployment  
🌐 **CORS Enabled** - Ready for web frontend integration  
⚡ **Health Monitoring** - Built-in health check endpoint  
🛡️ **Error Handling** - Proper HTTP status codes & validation  

## 🚀 Quick Start

### 🐳 Docker (Recommended)
```bash
git clone https://github.com/Neorex80/to-do-api.git
cd to-do-api
docker-compose up --build
```

### 🔧 Local Development
```bash
go mod download
go run main.go
```

**🌟 Server runs on:** `http://localhost:8080`

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | 💚 Health check |
| `GET` | `/api/tasks` | 📋 Get all tasks |
| `POST` | `/api/tasks` | ➕ Create task |
| `GET` | `/api/tasks/{id}` | 🔍 Get specific task |
| `PUT` | `/api/tasks/{id}` | ✏️ Update task |
| `DELETE` | `/api/tasks/{id}` | 🗑️ Delete task |

### 🧪 Quick Test
```bash
# Health check
curl http://localhost:8080/health

# Create a task
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"My Task","description":"Test task","status":"pending"}'

# Get all tasks
curl http://localhost:8080/api/tasks
```

## 🐳 Docker Deployment

```bash
# Build & Run
docker build -t to-do-api .
docker run -p 8080:8080 to-do-api

# With data persistence
docker run -p 8080:8080 -v $(pwd)/data:/app/data to-do-api
```

## ☁️ Deployment

### 🚄 Railway (One-Click Deploy)
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template?template=https://github.com/Neorex80/to-do-api)

### 🎨 Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

### 📦 Manual Steps
1. **Push to GitHub** ✅
2. **Connect to Railway/Render** 🔗
3. **Auto-deploy with Docker** 🚀
4. **Get live URL** 🌐

## 🛠️ Tech Stack

- **Backend:** Go 1.21+ with Gorilla Mux
- **Database:** SQLite
- **Containerization:** Docker
- **Deployment:** Railway, Render, Heroku

## 📄 Task Model

```json
{
  "id": 1,
  "title": "Learn Go",
  "description": "Build awesome APIs",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Options:** `pending` | `in_progress` | `completed`

## 🤝 Contributing

1. 🍴 Fork the repo
2. 🌿 Create feature branch
3. 💻 Make changes
4. 🚀 Submit PR

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

---

<div align="center">

**⭐ Star this repo if you found it helpful!**

Made with ❤️ and Go

</div>

# To-Do API

A lightweight Go + SQLite to-do API with a minimal frontend.

## Run locally

```bash
# build and run
GO111MODULE=on go run ./...
# open http://localhost:8080
```

## Docker

```bash
docker build -t todo-api .
docker run --rm -p 8080:8080 -e DB_PATH=/app/data/tasks.db -v $(pwd)/.data:/app/data todo-api
```

## Endpoints
- GET `/health`
- GET `/api/tasks?status=&limit=&offset=&sort_by=&sort_order=`
- GET `/api/tasks/{id}`
- POST `/api/tasks`
- PUT `/api/tasks/{id}`
- DELETE `/api/tasks/{id}`

## Frontend
- Served at `/` with static assets under `/static/`
- Links to GitHub and a Star button for quick access

## Performance optimizations
- SQLite PRAGMAs: WAL, synchronous=NORMAL, temp_store=MEMORY, busy_timeout
- Connection pool tuned (max open/idle, conn lifetime)
- Pagination and server-side filtering for task list
- Gzip compression and cache-control for static assets
- Docker image slimmed via `-trimpath`, `-s -w` and minimal runtime
