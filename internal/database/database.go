package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// Initialize creates and returns a new database connection
func Initialize() (*DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "workout_tracker.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	dbWrapper := &DB{db}
	
	// Create tables if they don't exist
	if err := dbWrapper.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	// Run migrations to update existing tables
	if err := dbWrapper.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return dbWrapper, nil
}

// createTables creates the necessary tables if they don't exist
func (db *DB) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT DEFAULT '',
			bio TEXT DEFAULT '',
			avatar TEXT DEFAULT '',
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			date DATETIME NOT NULL,
			duration INTEGER DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			exercise_id INTEGER NOT NULL,
			set_number INTEGER NOT NULL,
			reps INTEGER DEFAULT 0,
			weight REAL DEFAULT 0.0,
			distance REAL DEFAULT 0.0,
			duration INTEGER DEFAULT 0,
			rest_time INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_workouts_date ON workouts(date)`,
		`CREATE INDEX IF NOT EXISTS idx_exercises_workout_id ON exercises(workout_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id)`,
		`CREATE TABLE IF NOT EXISTS predefined_exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			description TEXT DEFAULT '',
			video_url TEXT DEFAULT '',
			instructions TEXT DEFAULT '',
			tips TEXT DEFAULT '',
			muscle_groups TEXT DEFAULT '',
			equipment TEXT DEFAULT '',
			difficulty TEXT DEFAULT 'beginner',
			image_url TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS meals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			calories INTEGER NOT NULL,
			protein REAL DEFAULT 0.0,
			carbs REAL DEFAULT 0.0,
			fat REAL DEFAULT 0.0,
			date DATETIME NOT NULL,
			meal_type TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS body_weights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			weight REAL NOT NULL,
			unit TEXT NOT NULL,
			date DATETIME NOT NULL,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS body_fats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			body_fat_pct REAL NOT NULL,
			date DATETIME NOT NULL,
			measurement TEXT DEFAULT '',
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS body_measurements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			measurement TEXT NOT NULL,
			value REAL NOT NULL,
			unit TEXT NOT NULL,
			date DATETIME NOT NULL,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS user_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			theme TEXT DEFAULT 'light',
			timezone TEXT DEFAULT 'UTC',
			weight_unit TEXT DEFAULT 'lbs',
			distance_unit TEXT DEFAULT 'miles',
			date_format TEXT DEFAULT 'MM/DD/YYYY',
			notifications BOOLEAN DEFAULT 1,
			privacy_mode BOOLEAN DEFAULT 0,
			auto_logout INTEGER DEFAULT 0,
			language TEXT DEFAULT 'en',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %v", query, err)
		}
	}

	return nil
}

// runMigrations updates existing tables to add new columns
func (db *DB) runMigrations() error {
	// Check if full_name column exists in users table
	var columnExists int
	err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='full_name'`).Scan(&columnExists)
	if err != nil {
		return fmt.Errorf("failed to check column existence: %v", err)
	}

	// If columns don't exist, add them
	if columnExists == 0 {
		migrations := []string{
			`ALTER TABLE users ADD COLUMN full_name TEXT DEFAULT ''`,
			`ALTER TABLE users ADD COLUMN bio TEXT DEFAULT ''`,
			`ALTER TABLE users ADD COLUMN avatar TEXT DEFAULT ''`,
			`ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT 1`,
		}

		for _, migration := range migrations {
			if _, err := db.Exec(migration); err != nil {
				return fmt.Errorf("failed to run migration: %s, error: %v", migration, err)
			}
		}
	}

	// Check if video_url column exists in predefined_exercises table
	var videoColumnExists int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('predefined_exercises') WHERE name='video_url'`).Scan(&videoColumnExists)
	if err != nil {
		return fmt.Errorf("failed to check predefined_exercises column existence: %v", err)
	}

	// If new predefined exercise columns don't exist, add them
	if videoColumnExists == 0 {
		exerciseMigrations := []string{
			`ALTER TABLE predefined_exercises ADD COLUMN video_url TEXT DEFAULT ''`,
			`ALTER TABLE predefined_exercises ADD COLUMN instructions TEXT DEFAULT ''`,
			`ALTER TABLE predefined_exercises ADD COLUMN tips TEXT DEFAULT ''`,
			`ALTER TABLE predefined_exercises ADD COLUMN muscle_groups TEXT DEFAULT ''`,
			`ALTER TABLE predefined_exercises ADD COLUMN equipment TEXT DEFAULT ''`,
			`ALTER TABLE predefined_exercises ADD COLUMN difficulty TEXT DEFAULT 'beginner'`,
			`ALTER TABLE predefined_exercises ADD COLUMN image_url TEXT DEFAULT ''`,
		}

		for _, migration := range exerciseMigrations {
			if _, err := db.Exec(migration); err != nil {
				return fmt.Errorf("failed to run predefined_exercises migration: %s, error: %v", migration, err)
			}
		}
	}

	return nil
}
