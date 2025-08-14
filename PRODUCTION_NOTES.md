# Production Notes

## Local Development vs Production

### CGO Requirement Issue

The main.go file uses SQLite with CGO, which requires a C compiler (gcc) to run locally on Windows. This is a common issue when developing Go applications with SQLite on Windows.

### Solutions:

#### Option 1: Use the Test Server for Local Development
```bash
go run test_server.go
```
This runs the API with in-memory storage (no CGO required) and is perfect for:
- Local development and testing
- API endpoint verification
- Frontend development
- Learning the API structure

#### Option 2: Install CGO Dependencies (Windows)
To run the full SQLite version locally:

1. **Install TDM-GCC or MinGW-w64:**
   - Download from: https://jmeubank.github.io/tdm-gcc/
   - Or use Chocolatey: `choco install mingw`

2. **Install with Go:**
   ```bash
   set CGO_ENABLED=1
   go run main.go
   ```

#### Option 3: Use Docker for Local Development
```bash
docker build -t to-do-api .
docker run -p 8080:8080 to-do-api
```

### Production Deployment

**The project is 100% production-ready!** The CGO issue only affects local Windows development. All deployment platforms (Railway, Render, etc.) handle CGO compilation automatically:

- ✅ **Railway**: Automatic Docker build with CGO support
- ✅ **Render**: Docker environment with build tools
- ✅ **Heroku**: Go buildpack with CGO support
- ✅ **Cloud platforms**: Container builds include necessary tools

### Recommended Workflow

1. **Local Development**: Use `test_server.go` for rapid development
2. **Testing**: Use the in-memory version to test API endpoints
3. **Production**: Deploy `main.go` with SQLite to your chosen platform
4. **CI/CD**: Use Docker builds for consistent environments

### File Usage

- **main.go**: Production version with SQLite persistence
- **test_server.go**: Development version with in-memory storage
- **Dockerfile**: Production container with CGO support

Both versions implement the exact same API interface, so you can develop locally with the test server and deploy the production version with confidence.

## Database Persistence

- **test_server.go**: Data is lost when server restarts (in-memory)
- **main.go**: Data persists in SQLite database file
- **Production**: Use main.go for data persistence

## API Compatibility

Both servers expose identical endpoints:
- POST /api/tasks
- GET /api/tasks
- GET /api/tasks/{id}
- PUT /api/tasks/{id}
- DELETE /api/tasks/{id}
- GET /health

The only difference is the storage backend (in-memory vs SQLite).
