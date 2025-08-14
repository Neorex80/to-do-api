# Docker Testing Guide

This guide provides comprehensive instructions for testing the To-Do API using Docker.

## Prerequisites

1. **Install Docker Desktop**
   - Windows: Download from [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)
   - macOS: Download from [Docker Desktop for Mac](https://docs.docker.com/desktop/install/mac-install/)
   - Linux: Follow [Docker Engine installation](https://docs.docker.com/engine/install/)

2. **Start Docker Desktop**
   - Ensure Docker Desktop is running before proceeding
   - Check with: `docker --version`

## Quick Start

### Option 1: Using Docker Compose (Recommended)

```bash
# Build and start the service
docker-compose up --build

# Or run in background
docker-compose up --build -d

# View logs
docker-compose logs -f to-do-api

# Stop the service
docker-compose down
```

### Option 2: Using Docker Commands

```bash
# Build the image
docker build -t to-do-api .

# Run the container
docker run -d \
  --name to-do-api \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  to-do-api

# View logs
docker logs -f to-do-api

# Stop and remove
docker stop to-do-api
docker rm to-do-api
```

## Testing the API

Once the container is running, test the endpoints:

### 1. Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "to-do-api"
}
```

### 2. Create a Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Docker Test Task",
    "description": "Testing the API in Docker",
    "status": "pending"
  }'
```

### 3. List All Tasks
```bash
curl http://localhost:8080/api/tasks
```

### 4. Get Specific Task
```bash
curl http://localhost:8080/api/tasks/1
```

### 5. Update a Task
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

### 6. Delete a Task
```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```

### 7. Filter by Status
```bash
curl "http://localhost:8080/api/tasks?status=pending"
```

## PowerShell Testing (Windows)

If using PowerShell instead of curl:

```powershell
# Health check
Invoke-WebRequest -Uri "http://localhost:8080/health"

# Create task
Invoke-WebRequest -Uri "http://localhost:8080/api/tasks" `
  -Method POST `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"Docker Test","description":"Testing in Docker","status":"pending"}'

# List tasks
Invoke-WebRequest -Uri "http://localhost:8080/api/tasks"
```

## Docker Container Management

### View Running Containers
```bash
docker ps
```

### View Container Logs
```bash
docker logs to-do-api

# Follow logs in real-time
docker logs -f to-do-api
```

### Execute Commands in Container
```bash
# Access container shell
docker exec -it to-do-api sh

# Check database file
docker exec to-do-api ls -la /app/data/

# View SQLite database
docker exec -it to-do-api sqlite3 /app/data/tasks.db ".tables"
```

### Container Health Check
```bash
# Check health status
docker inspect --format='{{.State.Health.Status}}' to-do-api

# View health check logs
docker inspect --format='{{range .State.Health.Log}}{{.Output}}{{end}}' to-do-api
```

## Data Persistence

### Volume Mounting
The Docker setup includes volume mounting for data persistence:

```bash
# Create data directory
mkdir -p ./data

# Run with volume mount
docker run -d \
  --name to-do-api \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  to-do-api
```

### Database File Location
- **Host**: `./data/tasks.db`
- **Container**: `/app/data/tasks.db`

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Find process using port 8080
   netstat -ano | findstr :8080
   
   # Kill the process (Windows)
   taskkill /PID <PID> /F
   
   # Or use different port
   docker run -p 8081:8080 to-do-api
   ```

2. **Docker Desktop Not Running**
   ```
   Error: error during connect: Head "http://%2F%2F.%2Fpipe%2FdockerDesktopLinuxEngine/_ping"
   ```
   **Solution**: Start Docker Desktop application

3. **Build Failures**
   ```bash
   # Clean build (remove cache)
   docker build --no-cache -t to-do-api .
   
   # Check build logs
   docker build -t to-do-api . --progress=plain
   ```

4. **Container Won't Start**
   ```bash
   # Check container logs
   docker logs to-do-api
   
   # Run interactively for debugging
   docker run -it --rm to-do-api sh
   ```

### Performance Testing

```bash
# Install Apache Bench (optional)
# Windows: Download from Apache website
# macOS: brew install httpie
# Linux: apt-get install apache2-utils

# Test API performance
ab -n 1000 -c 10 http://localhost:8080/health

# Or use curl for simple load test
for i in {1..100}; do
  curl -s http://localhost:8080/health > /dev/null
  echo "Request $i completed"
done
```

## Multi-Platform Testing

### Build for Different Architectures
```bash
# Build for ARM64 (Apple Silicon)
docker buildx build --platform linux/arm64 -t to-do-api:arm64 .

# Build for AMD64 (Intel/AMD)
docker buildx build --platform linux/amd64 -t to-do-api:amd64 .

# Build multi-platform
docker buildx build --platform linux/amd64,linux/arm64 -t to-do-api:latest .
```

## Production Simulation

### Environment Variables
```bash
# Test with production-like settings
docker run -d \
  --name to-do-api-prod \
  -p 8080:8080 \
  -e PORT=8080 \
  -e DB_PATH=/app/data/tasks.db \
  -e GIN_MODE=release \
  -v $(pwd)/data:/app/data \
  --restart unless-stopped \
  to-do-api
```

### Resource Limits
```bash
# Run with resource constraints
docker run -d \
  --name to-do-api-limited \
  -p 8080:8080 \
  --memory=128m \
  --cpus=0.5 \
  -v $(pwd)/data:/app/data \
  to-do-api
```

## Cleanup

### Remove Everything
```bash
# Stop and remove container
docker stop to-do-api
docker rm to-do-api

# Remove image
docker rmi to-do-api

# Remove volumes (careful - this deletes data!)
docker volume prune

# Using docker-compose
docker-compose down --volumes --rmi all
```

## Integration with CI/CD

### GitHub Actions Example
```yaml
name: Docker Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        run: docker build -t to-do-api .
      - name: Run container
        run: docker run -d --name test-api -p 8080:8080 to-do-api
      - name: Wait for startup
        run: sleep 10
      - name: Test health endpoint
        run: curl -f http://localhost:8080/health
      - name: Test API endpoints
        run: |
          curl -X POST http://localhost:8080/api/tasks \
            -H "Content-Type: application/json" \
            -d '{"title":"Test","status":"pending"}'
```

This comprehensive testing approach ensures your To-Do API works correctly in containerized environments and is ready for production deployment.
