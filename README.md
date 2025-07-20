<<<<<<< HEAD
# workout-tracker-backend
Backend of a workout tracker app written in Go.
=======
# Workout Tracker

A web-based workout tracking application built with Go and SQLite.

## Features

- **Dashboard**: Overview of recent workouts and statistics
- **Workout Management**: Create, view, and manage workouts
- **Exercise Tracking**: Add exercises to workouts with detailed sets information
- **Responsive Design**: Works on desktop and mobile devices
- **SQLite Database**: Lightweight, file-based database storage
- **REST API**: JSON API endpoints for programmatic access

## Project Structure

```
workout-tracker/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── database.go          # Database connection and setup
│   ├── handlers/
│   │   ├── handlers.go          # HTTP handlers
│   │   └── database.go          # Database operations
│   └── models/
│       └── models.go            # Data models
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css        # Stylesheet
│   │   └── js/
│   │       └── app.js           # JavaScript functionality
│   └── templates/
│       ├── base.html            # Base template
│       ├── index.html           # Dashboard
│       ├── workouts.html        # Workout listing
│       └── workout_detail.html  # Workout details
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.19 or higher
- SQLite3 (for database)

## Installation

1. Clone or download the project
2. Navigate to the project directory:
   ```bash
   cd workout-tracker
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

5. Open your browser and visit: `http://localhost:8080`

## Usage

### Web Interface

1. **Dashboard**: View recent workouts and statistics
2. **Create Workout**: Click "New Workout" to create a new workout session
3. **View Workouts**: Browse all your workouts on the workouts page
4. **Workout Details**: Click on a workout to view detailed information

### API Endpoints

- `GET /api/workouts` - List all workouts
- `POST /api/workouts` - Create a new workout
- `GET /api/workouts/{id}` - Get a specific workout

### Example API Usage

Create a workout:
```bash
curl -X POST http://localhost:8080/api/workouts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Morning Workout",
    "date": "2024-01-15",
    "duration": 45,
    "notes": "Great session!"
  }'
```

## Configuration

The application uses the following environment variables:

- `PORT`: Server port (default: 8080)
- `DATABASE_PATH`: SQLite database file path (default: workout_tracker.db)

## Database Schema

The application uses three main tables:

### Workouts
- `id`: Primary key
- `name`: Workout name
- `date`: Workout date
- `duration`: Duration in minutes
- `notes`: Optional notes
- `created_at`, `updated_at`: Timestamps

### Exercises
- `id`: Primary key
- `workout_id`: Foreign key to workouts
- `name`: Exercise name
- `category`: Exercise category (strength, cardio, flexibility, sports)
- `created_at`, `updated_at`: Timestamps

### Sets
- `id`: Primary key
- `exercise_id`: Foreign key to exercises
- `set_number`: Set number within exercise
- `reps`: Number of repetitions
- `weight`: Weight used (lbs)
- `distance`: Distance covered (miles)
- `duration`: Duration in seconds
- `rest_time`: Rest time in seconds
- `created_at`, `updated_at`: Timestamps

## Development

### Adding New Features

1. **Models**: Define data structures in `internal/models/models.go`
2. **Database**: Add database operations in `internal/handlers/database.go`
3. **Handlers**: Create HTTP handlers in `internal/handlers/handlers.go`
4. **Routes**: Register routes in `cmd/server/main.go`
5. **Templates**: Create HTML templates in `web/templates/`
6. **Styles**: Add CSS to `web/static/css/style.css`
7. **JavaScript**: Add interactivity in `web/static/js/app.js`

### Build for Production

```bash
go build -o workout-tracker cmd/server/main.go
```

### Running Tests

```bash
go test ./...
```

## Future Enhancements

- [ ] User authentication and authorization
- [ ] Exercise library with predefined exercises
- [ ] Workout templates and programs
- [ ] Progress tracking and analytics
- [ ] Data import/export functionality
- [ ] Mobile app (React Native or Flutter)
- [ ] Social features (sharing workouts)
- [ ] Integration with fitness trackers

## License

This project is open source and available under the MIT License.
>>>>>>> ae94d92 (Initial commit: Go workout tracker backend with Docker support)
