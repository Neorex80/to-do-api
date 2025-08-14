# Simple Docker Testing Script for To-Do API (PowerShell)
# This script tests the Docker setup and API functionality on Windows

Write-Host "Docker Testing for To-Do API" -ForegroundColor Cyan
Write-Host "============================" -ForegroundColor Cyan

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
try {
    docker --version
    Write-Host "Docker is available" -ForegroundColor Green
} catch {
    Write-Host "Docker is not running. Please start Docker Desktop." -ForegroundColor Red
    exit 1
}

# Clean up existing containers
Write-Host "Cleaning up..." -ForegroundColor Yellow
docker stop to-do-api-test 2>$null
docker rm to-do-api-test 2>$null

# Build the image
Write-Host "Building Docker image..." -ForegroundColor Yellow
docker build -t to-do-api-test .
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "Build successful" -ForegroundColor Green

# Create data directory
New-Item -ItemType Directory -Path "test-data" -Force | Out-Null

# Run container
Write-Host "Starting container..." -ForegroundColor Yellow
docker run -d --name to-do-api-test -p 8081:8080 -v "${PWD}/test-data:/app/data" to-do-api-test
if ($LASTEXITCODE -ne 0) {
    Write-Host "Container start failed" -ForegroundColor Red
    exit 1
}
Write-Host "Container started" -ForegroundColor Green

# Wait for startup
Write-Host "Waiting for service..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

# Test health endpoint
Write-Host "Testing health endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -UseBasicParsing
    if ($response.StatusCode -eq 200) {
        Write-Host "Health check passed" -ForegroundColor Green
    }
} catch {
    Write-Host "Health check failed" -ForegroundColor Red
    docker logs to-do-api-test
    exit 1
}

# Test API endpoints
Write-Host "Testing API endpoints..." -ForegroundColor Yellow

# Create task
$taskBody = '{"title":"Docker Test","description":"Testing in Docker","status":"pending"}'
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/api/tasks" -Method POST -Headers @{"Content-Type"="application/json"} -Body $taskBody -UseBasicParsing
    Write-Host "Task creation: PASSED" -ForegroundColor Green
} catch {
    Write-Host "Task creation: FAILED" -ForegroundColor Red
}

# List tasks
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/api/tasks" -UseBasicParsing
    Write-Host "Task listing: PASSED" -ForegroundColor Green
} catch {
    Write-Host "Task listing: FAILED" -ForegroundColor Red
}

# Update task
$updateBody = '{"status":"completed"}'
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/api/tasks/1" -Method PUT -Headers @{"Content-Type"="application/json"} -Body $updateBody -UseBasicParsing
    Write-Host "Task update: PASSED" -ForegroundColor Green
} catch {
    Write-Host "Task update: FAILED" -ForegroundColor Red
}

# Show logs
Write-Host "Container logs:" -ForegroundColor Yellow
docker logs to-do-api-test --tail 5

Write-Host ""
Write-Host "============================" -ForegroundColor Cyan
Write-Host "Docker tests completed!" -ForegroundColor Green
Write-Host "API running at: http://localhost:8081" -ForegroundColor Yellow
Write-Host "To stop: docker stop to-do-api-test" -ForegroundColor Yellow
Write-Host "============================" -ForegroundColor Cyan
