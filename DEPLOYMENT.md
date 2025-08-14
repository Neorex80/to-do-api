# Deployment Guide

This guide covers deploying the To-Do API to popular cloud platforms.

## Deploy to Railway

Railway is recommended for its simplicity and automatic Docker deployment.

### Steps:

1. **Create a Railway Account**
   - Go to [Railway](https://railway.app)
   - Sign up with GitHub

2. **Deploy from GitHub**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway will automatically detect the Dockerfile and deploy

3. **Environment Variables**
   - Railway will automatically set `PORT`
   - Optionally set `DB_PATH` if you want a custom database location

4. **Domain**
   - Railway provides a free domain: `your-app.up.railway.app`
   - You can add a custom domain in the settings

### Railway Configuration

The `railway.toml` file is included for optimal configuration:
- Health check on `/health`
- Automatic restarts on failure
- Dockerfile-based builds

## Deploy to Render

Render is another excellent option with a generous free tier.

### Steps:

1. **Create a Render Account**
   - Go to [Render](https://render.com)
   - Sign up with GitHub

2. **Create a Web Service**
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Choose the repository

3. **Configuration**
   - **Name**: `to-do-api` (or your preferred name)
   - **Environment**: `Docker`
   - **Build Command**: Leave empty (uses Dockerfile)
   - **Start Command**: Leave empty (uses Dockerfile CMD)

4. **Environment Variables**
   - `PORT`: Automatically set by Render
   - `DB_PATH`: `/opt/render/project/src/data/tasks.db` (optional)

5. **Deploy**
   - Click "Create Web Service"
   - Render will build and deploy automatically

### Render Features
- Free tier with 750 hours/month
- Automatic SSL certificates
- Custom domains
- Automatic deploys on git push

## Deploy to Heroku

### Steps:

1. **Install Heroku CLI**
   ```bash
   # Install Heroku CLI from https://devcenter.heroku.com/articles/heroku-cli
   ```

2. **Login and Create App**
   ```bash
   heroku login
   heroku create your-app-name
   ```

3. **Set Buildpack**
   ```bash
   heroku buildpacks:set heroku/go
   ```

4. **Deploy**
   ```bash
   git push heroku main
   ```

5. **Environment Variables**
   ```bash
   heroku config:set DB_PATH=/app/data/tasks.db
   ```

## Deploy to Google Cloud Run

### Steps:

1. **Build and Push to Container Registry**
   ```bash
   # Build the image
   docker build -t gcr.io/YOUR_PROJECT_ID/to-do-api .
   
   # Push to registry
   docker push gcr.io/YOUR_PROJECT_ID/to-do-api
   ```

2. **Deploy to Cloud Run**
   ```bash
   gcloud run deploy to-do-api \
     --image gcr.io/YOUR_PROJECT_ID/to-do-api \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

## Deploy to AWS (using Docker)

### Using AWS App Runner:

1. **Push to ECR**
   ```bash
   # Create ECR repository
   aws ecr create-repository --repository-name to-do-api
   
   # Build and push
   docker build -t to-do-api .
   docker tag to-do-api:latest YOUR_ACCOUNT.dkr.ecr.REGION.amazonaws.com/to-do-api:latest
   docker push YOUR_ACCOUNT.dkr.ecr.REGION.amazonaws.com/to-do-api:latest
   ```

2. **Create App Runner Service**
   - Go to AWS App Runner console
   - Create service from container registry
   - Configure auto-scaling and health checks

## Environment Variables

All platforms support these environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port (usually set by platform) |
| `DB_PATH` | ./tasks.db | SQLite database file path |

## Health Checks

All platforms can use the health check endpoint:
- **URL**: `/health`
- **Method**: GET
- **Expected Response**: 200 OK with JSON status

## Database Persistence

**Important**: SQLite files are stored on the container filesystem. For production use:

1. **Railway/Render**: Files persist between deployments
2. **Heroku**: Files are ephemeral (lost on restart)
3. **Cloud platforms**: Consider using managed databases for production

For production, consider migrating to:
- PostgreSQL (recommended)
- MySQL
- MongoDB

## Monitoring

After deployment, monitor your API:

1. **Health Check**: `GET /health`
2. **Logs**: Check platform-specific logging
3. **Metrics**: Monitor response times and error rates

## Custom Domains

Most platforms support custom domains:
- **Railway**: Settings → Domains
- **Render**: Settings → Custom Domains
- **Heroku**: Settings → Domains

## SSL/HTTPS

All mentioned platforms provide automatic SSL certificates for both their domains and custom domains.

## Troubleshooting

### Common Issues:

1. **Build Failures**
   - Check Dockerfile syntax
   - Ensure all dependencies are in go.mod

2. **Runtime Errors**
   - Check logs for CGO/SQLite issues
   - Verify environment variables

3. **Database Issues**
   - Ensure DB_PATH directory exists
   - Check file permissions

### Getting Help:

- Check platform-specific documentation
- Review application logs
- Test locally with Docker first

## Cost Considerations

- **Railway**: Free tier, then usage-based pricing
- **Render**: Free tier with limitations, paid plans available
- **Heroku**: Free tier discontinued, paid plans only
- **Cloud providers**: Pay-per-use, can be cost-effective for low traffic

Choose the platform that best fits your needs and budget!
