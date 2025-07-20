#!/bin/bash

# Workout Tracker Docker Deployment Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="workout-tracker"
IMAGE_NAME="workout-tracker:latest"

print_help() {
    echo "Workout Tracker Deployment Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build       Build the Docker image"
    echo "  start       Start the application (development)"
    echo "  prod        Start the application (production)"
    echo "  stop        Stop the application"
    echo "  restart     Restart the application"
    echo "  logs        Show application logs"
    echo "  backup      Backup the database"
    echo "  clean       Remove containers and images"
    echo "  status      Show container status"
    echo "  shell       Open shell in container"
    echo "  help        Show this help message"
}

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        error "Docker daemon is not running"
        exit 1
    fi
}

build() {
    log "Building Docker image..."
    docker build -t $IMAGE_NAME .
    success "Image built successfully"
}

start_dev() {
    log "Starting application in development mode..."
    docker-compose up -d
    success "Application started on http://localhost:8080"
}

start_prod() {
    log "Starting application in production mode..."
    if [ ! -f .env ]; then
        warning ".env file not found. Please create one based on .env.example"
        exit 1
    fi
    docker-compose -f docker-compose.prod.yml up -d
    success "Application started in production mode"
}

stop() {
    log "Stopping application..."
    docker-compose down 2>/dev/null || docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
    success "Application stopped"
}

restart() {
    log "Restarting application..."
    stop
    sleep 2
    start_dev
}

show_logs() {
    log "Showing logs (press Ctrl+C to exit)..."
    docker logs -f $CONTAINER_NAME 2>/dev/null || docker logs -f workout-tracker-prod 2>/dev/null || error "No running container found"
}

backup() {
    log "Creating database backup..."
    mkdir -p backups
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="backups/workout_tracker_backup_${TIMESTAMP}.db"
    
    # Try to copy from running container first
    if docker ps | grep -q $CONTAINER_NAME; then
        docker cp ${CONTAINER_NAME}:/app/data/workout_tracker.db $BACKUP_FILE
        success "Database backed up to $BACKUP_FILE"
    elif docker ps | grep -q workout-tracker-prod; then
        docker cp workout-tracker-prod:/app/data/workout_tracker.db $BACKUP_FILE
        success "Database backed up to $BACKUP_FILE"
    else
        error "No running container found for backup"
        exit 1
    fi
}

clean() {
    log "Cleaning up Docker containers and images..."
    docker-compose down -v 2>/dev/null || true
    docker-compose -f docker-compose.prod.yml down -v 2>/dev/null || true
    docker rmi $IMAGE_NAME 2>/dev/null || true
    success "Cleanup completed"
}

status() {
    log "Container status:"
    echo ""
    docker ps -a | grep workout-tracker || echo "No workout-tracker containers found"
    echo ""
    log "Image status:"
    docker images | grep workout-tracker || echo "No workout-tracker images found"
}

shell() {
    if docker ps | grep -q $CONTAINER_NAME; then
        docker exec -it $CONTAINER_NAME sh
    elif docker ps | grep -q workout-tracker-prod; then
        docker exec -it workout-tracker-prod sh
    else
        error "No running container found"
        exit 1
    fi
}

# Main script logic
check_docker

case "${1:-help}" in
    "build")
        build
        ;;
    "start")
        start_dev
        ;;
    "prod")
        start_prod
        ;;
    "stop")
        stop
        ;;
    "restart")
        restart
        ;;
    "logs")
        show_logs
        ;;
    "backup")
        backup
        ;;
    "clean")
        clean
        ;;
    "status")
        status
        ;;
    "shell")
        shell
        ;;
    "help"|*)
        print_help
        ;;
esac
