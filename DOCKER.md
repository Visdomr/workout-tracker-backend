# Docker Deployment Guide

This document describes how to run the Workout Tracker application using Docker containers.

## Quick Start

### Development Mode
```bash
# Build and start the application
./deploy.sh build
./deploy.sh start

# Access the application at http://localhost:8080
```

### Production Mode
```bash
# Create environment file
cp .env.example .env
# Edit .env with your production settings

# Start in production mode
./deploy.sh prod
```

## Available Commands

The `deploy.sh` script provides several management commands:

- `./deploy.sh build` - Build the Docker image
- `./deploy.sh start` - Start in development mode
- `./deploy.sh prod` - Start in production mode (requires .env file)
- `./deploy.sh stop` - Stop the application
- `./deploy.sh restart` - Restart the application
- `./deploy.sh logs` - Show application logs
- `./deploy.sh backup` - Backup the database
- `./deploy.sh clean` - Remove containers and images
- `./deploy.sh status` - Show container status
- `./deploy.sh shell` - Open shell in running container

## Manual Docker Commands

If you prefer to use Docker commands directly:

### Development
```bash
# Build the image
docker build -t workout-tracker:latest .

# Run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Production
```bash
# Create .env file first
cp .env.example .env

# Run production compose
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop
docker-compose -f docker-compose.prod.yml down
```

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Server Configuration
PORT=8080
DATABASE_PATH=/app/data/workout_tracker.db

# Security (IMPORTANT: Change in production!)
SESSION_SECRET=your-secure-session-secret-here

# OAuth Configuration (optional)
GOOGLE_CLIENT_ID=your-google-client-id.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### Data Persistence

The application uses Docker volumes to persist data:
- `workout_data` - Contains the SQLite database
- Database backups can be created with `./deploy.sh backup`

## Deployment Options

### 1. Development (Local)
- Uses `docker-compose.yml`
- Suitable for local development and testing
- Minimal resource limits

### 2. Production (Local/Server)
- Uses `docker-compose.prod.yml`
- Requires `.env` file
- Includes resource limits and health checks
- Optional nginx reverse proxy (commented out)

### 3. Raspberry Pi
The application runs well on Raspberry Pi 4/5:
```bash
# On Raspberry Pi (ARM64)
./deploy.sh build
./deploy.sh start
```

### 4. Cloud Deployment
For cloud platforms (AWS, GCP, Azure):
1. Build and push image to registry
2. Deploy using cloud-specific orchestration
3. Configure external database if needed

## Proxmox Integration

To deploy in your Proxmox environment:

1. **Container approach** (recommended):
   ```bash
   # In your Ubuntu Docker VM
   git clone <your-repo>
   cd workout-tracker
   ./deploy.sh prod
   ```

2. **VM approach**:
   - Create new Ubuntu VM in Proxmox
   - Install Docker
   - Deploy application

## Backup and Restore

### Backup
```bash
# Create backup
./deploy.sh backup

# Backups are stored in ./backups/ directory
```

### Restore
```bash
# Stop application
./deploy.sh stop

# Copy backup to volume (manual process)
# Start application
./deploy.sh start
```

## Monitoring

### Health Checks
The container includes health checks that verify the application is responding on port 8080.

### Logs
```bash
# View logs
./deploy.sh logs

# Or with Docker Compose
docker-compose logs -f workout-tracker
```

### Resource Usage
```bash
# Check resource usage
docker stats workout-tracker
```

## Troubleshooting

### Container won't start
```bash
# Check container status
./deploy.sh status

# View detailed logs
./deploy.sh logs

# Rebuild if needed
./deploy.sh clean
./deploy.sh build
./deploy.sh start
```

### Database issues
```bash
# Access container shell
./deploy.sh shell

# Check database file
ls -la /app/data/
```

### Port conflicts
If port 8080 is already in use, modify the port mapping in `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Changed from 8080:8080
```

## Security Considerations

1. **Change SESSION_SECRET** in production
2. **Use proper firewall rules** for your Proxmox setup
3. **Regular backups** of the database
4. **Update base images** periodically for security patches
5. **Consider using HTTPS** with nginx reverse proxy

## Next Steps

1. Test locally with `./deploy.sh start`
2. Configure production environment variables
3. Deploy to your preferred platform
4. Set up regular backups
5. Configure monitoring/alerting if needed
