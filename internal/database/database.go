package database

import (
	"database/sql"
	"fmt"
	"log"
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
		`CREATE TABLE IF NOT EXISTS workout_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS template_exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			template_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			order_index INTEGER DEFAULT 0,
			target_sets INTEGER DEFAULT 3,
			target_reps INTEGER DEFAULT 10,
			target_weight REAL DEFAULT 0,
			rest_time INTEGER DEFAULT 60,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (template_id) REFERENCES workout_templates(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workout_programs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			difficulty TEXT DEFAULT 'beginner',
			duration_weeks INTEGER DEFAULT 8,
			goal TEXT DEFAULT 'general',
			is_public BOOLEAN DEFAULT 1,
			created_by INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS program_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			program_id INTEGER NOT NULL,
			template_id INTEGER NOT NULL,
			day_of_week INTEGER NOT NULL,
			week_number INTEGER DEFAULT 1,
			order_index INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (program_id) REFERENCES workout_programs(id) ON DELETE CASCADE,
			FOREIGN KEY (template_id) REFERENCES workout_templates(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS template_sharing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			template_id INTEGER NOT NULL,
			owner_id INTEGER NOT NULL,
			shared_with_id INTEGER NOT NULL,
			permission TEXT DEFAULT 'view',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (template_id) REFERENCES workout_templates(id) ON DELETE CASCADE,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (shared_with_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(template_id, shared_with_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_templates_user_id ON workout_templates(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_template_exercises_template_id ON template_exercises(template_id)`,
		`CREATE INDEX IF NOT EXISTS idx_program_templates_program_id ON program_templates(program_id)`,
		`CREATE INDEX IF NOT EXISTS idx_template_sharing_owner_id ON template_sharing(owner_id)`,
		`CREATE INDEX IF NOT EXISTS idx_template_sharing_shared_with_id ON template_sharing(shared_with_id)`,
		`CREATE TABLE IF NOT EXISTS scheduled_workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			template_id INTEGER,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			scheduled_date DATETIME NOT NULL,
			scheduled_time TEXT, -- Format: HH:MM
			estimated_duration INTEGER DEFAULT 60, -- minutes
			status TEXT DEFAULT 'scheduled', -- scheduled, completed, skipped, cancelled
			workout_id INTEGER, -- Reference to actual workout when completed
			reminder_sent BOOLEAN DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (template_id) REFERENCES workout_templates(id) ON DELETE SET NULL,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS workout_reminders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			scheduled_workout_id INTEGER NOT NULL,
			reminder_type TEXT NOT NULL, -- email, push, sms
			message TEXT NOT NULL,
			scheduled_for DATETIME NOT NULL,
			status TEXT DEFAULT 'pending', -- pending, sent, failed, cancelled
			sent_at DATETIME,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (scheduled_workout_id) REFERENCES scheduled_workouts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS rest_day_recommendations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			recommended_date DATETIME NOT NULL,
			reason TEXT NOT NULL, -- high_volume, consecutive_days, muscle_group_fatigue, etc.
			intensity_score REAL DEFAULT 0, -- 0-10 scale
			volume_load REAL DEFAULT 0, -- Total volume from recent workouts
			consecutive_days INTEGER DEFAULT 0,
			muscle_groups_worked TEXT, -- JSON array of muscle groups needing rest
			status TEXT DEFAULT 'suggested', -- suggested, accepted, ignored, overridden
			user_response TEXT, -- User's response/note
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS deload_recommendations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			recommended_start_date DATETIME NOT NULL,
			recommended_end_date DATETIME NOT NULL,
			reason TEXT NOT NULL, -- fatigue_accumulation, plateau, overreaching, scheduled
			volume_reduction_percent INTEGER DEFAULT 40, -- Recommended volume reduction
			intensity_reduction_percent INTEGER DEFAULT 20, -- Recommended intensity reduction
			trigger_metrics TEXT, -- JSON object with metrics that triggered recommendation
			status TEXT DEFAULT 'suggested', -- suggested, accepted, ignored, active, completed
			user_response TEXT,
			started_at DATETIME,
			completed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workout_calendar_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			scheduled_workout_id INTEGER,
			rest_day_id INTEGER,
			deload_id INTEGER,
			event_type TEXT NOT NULL, -- workout, rest_day, deload
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			start_date DATETIME NOT NULL,
			end_date DATETIME,
			all_day BOOLEAN DEFAULT 1,
			color TEXT DEFAULT '#3788d8', -- Color for calendar display
			is_recurring BOOLEAN DEFAULT 0,
			recurrence_pattern TEXT, -- JSON object for recurring events
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (scheduled_workout_id) REFERENCES scheduled_workouts(id) ON DELETE CASCADE,
			FOREIGN KEY (rest_day_id) REFERENCES rest_day_recommendations(id) ON DELETE CASCADE,
			FOREIGN KEY (deload_id) REFERENCES deload_recommendations(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_scheduled_workouts_user_id ON scheduled_workouts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_scheduled_workouts_date ON scheduled_workouts(scheduled_date)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_reminders_user_id ON workout_reminders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_reminders_scheduled_for ON workout_reminders(scheduled_for)`,
		`CREATE INDEX IF NOT EXISTS idx_rest_day_recommendations_user_id ON rest_day_recommendations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_rest_day_recommendations_date ON rest_day_recommendations(recommended_date)`,
		`CREATE INDEX IF NOT EXISTS idx_deload_recommendations_user_id ON deload_recommendations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_calendar_events_user_id ON workout_calendar_events(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_calendar_events_start_date ON workout_calendar_events(start_date)`,
		// ========== ADVANCED ANALYTICS TABLES ==========
		`CREATE TABLE IF NOT EXISTS personal_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercise_name TEXT NOT NULL,
			weight REAL NOT NULL,
			reps INTEGER NOT NULL,
			volume REAL NOT NULL, -- weight * reps
			one_rep_max REAL NOT NULL,
			date DATETIME NOT NULL,
			workout_id INTEGER,
			set_id INTEGER,
			is_new BOOLEAN DEFAULT 1, -- If achieved in current period
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
			FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS strength_progress (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercise_name TEXT NOT NULL,
			category TEXT NOT NULL,
			date DATETIME NOT NULL,
			max_weight REAL NOT NULL,
			max_reps INTEGER NOT NULL,
			volume REAL NOT NULL,
			one_rep_max REAL NOT NULL,
			workout_id INTEGER,
			trend TEXT DEFAULT 'stable', -- improving, stable, declining
			trend_percent REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workout_intensity (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			workout_id INTEGER NOT NULL,
			date DATETIME NOT NULL,
			total_volume REAL NOT NULL,
			avg_weight REAL NOT NULL,
			total_sets INTEGER NOT NULL,
			total_reps INTEGER NOT NULL,
			intensity_score REAL NOT NULL, -- Calculated intensity metric (0-10)
			duration_minutes INTEGER DEFAULT 0,
			avg_rest_time INTEGER DEFAULT 0, -- seconds
			rpe_score REAL DEFAULT 0, -- Rate of Perceived Exertion (1-10)
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exercise_frequency (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercise_name TEXT NOT NULL,
			category TEXT NOT NULL,
			count INTEGER DEFAULT 1,
			percentage REAL DEFAULT 0,
			last_performed DATETIME NOT NULL,
			first_performed DATETIME NOT NULL,
			total_volume REAL DEFAULT 0,
			total_sets INTEGER DEFAULT 0,
			total_reps INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS progress_trends (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			trend_type TEXT NOT NULL, -- volume, strength, frequency
			trend_direction TEXT NOT NULL, -- up, down, stable
			change_percent REAL NOT NULL,
			confidence_score REAL DEFAULT 0, -- 0-100
			time_period TEXT NOT NULL, -- week, month, quarter, year
			start_date DATETIME NOT NULL,
			end_date DATETIME NOT NULL,
			data_points TEXT DEFAULT '[]', -- JSON array of data points
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS weekly_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			week_start DATETIME NOT NULL,
			week_end DATETIME NOT NULL,
			week_number INTEGER NOT NULL,
			year INTEGER NOT NULL,
			total_workouts INTEGER DEFAULT 0,
			total_duration INTEGER DEFAULT 0, -- minutes
			avg_duration REAL DEFAULT 0,
			total_volume REAL DEFAULT 0,
			total_sets INTEGER DEFAULT 0,
			total_reps INTEGER DEFAULT 0,
			unique_exercises INTEGER DEFAULT 0,
			max_weight REAL DEFAULT 0,
			avg_calories INTEGER DEFAULT 0,
			weight_change REAL DEFAULT 0,
			top_exercises TEXT DEFAULT '[]', -- JSON array
			prs_achieved INTEGER DEFAULT 0,
			consistency_score REAL DEFAULT 0, -- 0-100
			intensity_score REAL DEFAULT 0, -- 0-10
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS monthly_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			month INTEGER NOT NULL,
			year INTEGER NOT NULL,
			month_name TEXT NOT NULL,
			total_workouts INTEGER DEFAULT 0,
			total_duration INTEGER DEFAULT 0,
			avg_duration REAL DEFAULT 0,
			total_volume REAL DEFAULT 0,
			total_sets INTEGER DEFAULT 0,
			total_reps INTEGER DEFAULT 0,
			unique_exercises INTEGER DEFAULT 0,
			max_weight REAL DEFAULT 0,
			avg_calories INTEGER DEFAULT 0,
			weight_change REAL DEFAULT 0,
			start_weight REAL DEFAULT 0,
			end_weight REAL DEFAULT 0,
			top_exercises TEXT DEFAULT '[]',
			prs_achieved INTEGER DEFAULT 0,
			consistency_score REAL DEFAULT 0,
			intensity_score REAL DEFAULT 0,
			category_breakdown TEXT DEFAULT '{}', -- JSON object
			progress_highlights TEXT DEFAULT '[]', -- JSON array
			goals_achieved TEXT DEFAULT '[]',
			recommendations TEXT DEFAULT '[]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS yearly_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			year INTEGER NOT NULL,
			total_workouts INTEGER DEFAULT 0,
			total_duration INTEGER DEFAULT 0,
			avg_duration REAL DEFAULT 0,
			total_volume REAL DEFAULT 0,
			total_sets INTEGER DEFAULT 0,
			total_reps INTEGER DEFAULT 0,
			unique_exercises INTEGER DEFAULT 0,
			max_weight REAL DEFAULT 0,
			avg_calories INTEGER DEFAULT 0,
			total_weight_change REAL DEFAULT 0,
			start_weight REAL DEFAULT 0,
			end_weight REAL DEFAULT 0,
			top_exercises TEXT DEFAULT '[]',
			total_prs_achieved INTEGER DEFAULT 0,
			avg_consistency REAL DEFAULT 0,
			avg_intensity REAL DEFAULT 0,
			year_highlights TEXT DEFAULT '[]',
			fitness_journey TEXT DEFAULT '[]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS quarterly_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			quarter INTEGER NOT NULL,
			year INTEGER NOT NULL,
			quarter_name TEXT NOT NULL,
			start_date DATETIME NOT NULL,
			end_date DATETIME NOT NULL,
			total_workouts INTEGER DEFAULT 0,
			total_volume REAL DEFAULT 0,
			avg_intensity REAL DEFAULT 0,
			weight_change REAL DEFAULT 0,
			prs_achieved INTEGER DEFAULT 0,
			top_achievements TEXT DEFAULT '[]',
			focus_areas TEXT DEFAULT '[]',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS milestones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			category TEXT NOT NULL, -- strength, endurance, consistency, weight_loss, etc.
			value REAL NOT NULL,
			unit TEXT NOT NULL, -- lbs, kg, reps, days, etc.
			achieved_at DATETIME NOT NULL,
			workout_id INTEGER,
			exercise_name TEXT,
			is_personal_record BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exercise_progress_charts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercise_name TEXT NOT NULL,
			category TEXT NOT NULL,
			time_range TEXT NOT NULL, -- last_30_days, last_3_months, last_6_months, last_year
			data_points TEXT DEFAULT '[]', -- JSON array of ExerciseProgressPoint
			statistics TEXT DEFAULT '{}', -- JSON object with ExerciseStatistics
			milestones TEXT DEFAULT '[]', -- JSON array of ExerciseMilestone
			predictions TEXT DEFAULT '{}', -- JSON object with ExercisePrediction
			comparisons TEXT DEFAULT '{}', -- JSON object with ExerciseComparison
			last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exercise_comparisons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercises TEXT NOT NULL, -- JSON array of exercise names
			time_range TEXT NOT NULL,
			metric TEXT NOT NULL, -- weight, volume, frequency
			data_series TEXT DEFAULT '[]', -- JSON array of ComparisonDataSeries
			rankings TEXT DEFAULT '[]', -- JSON array of ExerciseRanking
			insights TEXT DEFAULT '[]', -- JSON array of insights
			last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS workout_analytics_cache (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			cache_key TEXT NOT NULL, -- unique identifier for cached data
			cache_type TEXT NOT NULL, -- analytics_data, summary, progress_chart, etc.
			data TEXT NOT NULL, -- JSON data
			parameters TEXT DEFAULT '{}', -- JSON object with query parameters
			expires_at DATETIME NOT NULL,
			last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(user_id, cache_key)
		)`,
		`CREATE TABLE IF NOT EXISTS template_usage (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			template_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			workout_id INTEGER NOT NULL,
			used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (template_id) REFERENCES workout_templates(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
		)`,
		// Indexes for advanced analytics tables
		`CREATE INDEX IF NOT EXISTS idx_personal_records_user_id ON personal_records(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_personal_records_exercise_name ON personal_records(exercise_name)`,
		`CREATE INDEX IF NOT EXISTS idx_personal_records_date ON personal_records(date)`,
		`CREATE INDEX IF NOT EXISTS idx_strength_progress_user_id ON strength_progress(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_strength_progress_exercise_name ON strength_progress(exercise_name)`,
		`CREATE INDEX IF NOT EXISTS idx_strength_progress_date ON strength_progress(date)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_intensity_user_id ON workout_intensity(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_intensity_workout_id ON workout_intensity(workout_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_intensity_date ON workout_intensity(date)`,
		`CREATE INDEX IF NOT EXISTS idx_exercise_frequency_user_id ON exercise_frequency(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_exercise_frequency_exercise_name ON exercise_frequency(exercise_name)`,
		`CREATE INDEX IF NOT EXISTS idx_progress_trends_user_id ON progress_trends(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_progress_trends_type_period ON progress_trends(trend_type, time_period)`,
		`CREATE INDEX IF NOT EXISTS idx_weekly_summaries_user_id ON weekly_summaries(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_weekly_summaries_week ON weekly_summaries(year, week_number)`,
		`CREATE INDEX IF NOT EXISTS idx_monthly_summaries_user_id ON monthly_summaries(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_monthly_summaries_month ON monthly_summaries(year, month)`,
		`CREATE INDEX IF NOT EXISTS idx_yearly_summaries_user_id ON yearly_summaries(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_yearly_summaries_year ON yearly_summaries(year)`,
		`CREATE INDEX IF NOT EXISTS idx_quarterly_summaries_user_id ON quarterly_summaries(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_quarterly_summaries_quarter ON quarterly_summaries(year, quarter)`,
		`CREATE INDEX IF NOT EXISTS idx_milestones_user_id ON milestones(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_milestones_category ON milestones(category)`,
		`CREATE INDEX IF NOT EXISTS idx_milestones_achieved_at ON milestones(achieved_at)`,
		`CREATE INDEX IF NOT EXISTS idx_exercise_progress_charts_user_id ON exercise_progress_charts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_exercise_progress_charts_exercise ON exercise_progress_charts(exercise_name)`,
		`CREATE INDEX IF NOT EXISTS idx_exercise_comparisons_user_id ON exercise_comparisons(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_analytics_cache_user_id ON workout_analytics_cache(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_analytics_cache_key ON workout_analytics_cache(cache_key)`,
		`CREATE INDEX IF NOT EXISTS idx_workout_analytics_cache_expires ON workout_analytics_cache(expires_at)`,
		`CREATE INDEX IF NOT EXISTS idx_template_usage_template_id ON template_usage(template_id)`,
		`CREATE INDEX IF NOT EXISTS idx_template_usage_user_id ON template_usage(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_template_usage_used_at ON template_usage(used_at)`,
		`CREATE TABLE IF NOT EXISTS achievements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL, -- strength, consistency, milestone, etc.
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			badge_url TEXT DEFAULT '', -- URL to badge image
			achieved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS badges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			badge_url TEXT DEFAULT '',
			criteria TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS streaks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL, -- workout, nutrition
			start_date DATETIME NOT NULL,
			end_date DATETIME,
			count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS points (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL, -- workout, nutrition, achievement
			value INTEGER DEFAULT 0,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			start_date DATETIME NOT NULL,
			end_date DATETIME NOT NULL,
			goal_type TEXT NOT NULL, -- steps, streak, workout, etc.
			goal_value INTEGER NOT NULL,
			reward_points INTEGER DEFAULT 0,
			badge_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS user_challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			challenge_id INTEGER NOT NULL,
			progress INTEGER DEFAULT 0,
			status TEXT DEFAULT 'ongoing', -- ongoing, completed, failed
			completed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (challenge_id) REFERENCES challenges(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS export_jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			export_type TEXT NOT NULL, -- csv, json, backup
			data_types TEXT NOT NULL, -- JSON array of data types to export
			status TEXT DEFAULT 'pending', -- pending, processing, completed, failed
			file_path TEXT DEFAULT '',
			file_size INTEGER DEFAULT 0,
			download_count INTEGER DEFAULT 0,
			export_options TEXT DEFAULT '{}', -- JSON object with export options
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			expires_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS import_jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			import_type TEXT NOT NULL, -- csv, json, myfitnesspal, strava, etc.
			data_types TEXT NOT NULL, -- JSON array of data types to import
			status TEXT DEFAULT 'pending', -- pending, processing, completed, failed, validation_error
			file_path TEXT NOT NULL,
			file_size INTEGER DEFAULT 0,
			total_records INTEGER DEFAULT 0,
			processed_records INTEGER DEFAULT 0,
			successful_records INTEGER DEFAULT 0,
			failed_records INTEGER DEFAULT 0,
			import_options TEXT DEFAULT '{}', -- JSON object with import options
			validation_errors TEXT DEFAULT '[]', -- JSON array of validation errors
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS data_sync_configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			provider TEXT NOT NULL, -- strava, myfitnesspal, garmin, fitbit, etc.
			access_token TEXT NOT NULL,
			refresh_token TEXT,
			token_expires_at DATETIME,
			sync_enabled BOOLEAN DEFAULT 1,
			last_sync_at DATETIME,
			sync_frequency TEXT DEFAULT 'daily', -- manual, hourly, daily, weekly
			data_types TEXT DEFAULT '[]', -- JSON array of data types to sync
			sync_options TEXT DEFAULT '{}', -- JSON object with sync options
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sync_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			sync_config_id INTEGER NOT NULL,
			sync_type TEXT NOT NULL, -- import, export
			status TEXT NOT NULL, -- success, partial, failed
			records_processed INTEGER DEFAULT 0,
			records_successful INTEGER DEFAULT 0,
			records_failed INTEGER DEFAULT 0,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			error_details TEXT,
			sync_summary TEXT, -- JSON object with sync summary
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (sync_config_id) REFERENCES data_sync_configs(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS backup_configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			backup_frequency TEXT DEFAULT 'weekly', -- manual, daily, weekly, monthly
			include_workouts BOOLEAN DEFAULT 1,
			include_nutrition BOOLEAN DEFAULT 1,
			include_body_metrics BOOLEAN DEFAULT 1,
			include_templates BOOLEAN DEFAULT 1,
			include_settings BOOLEAN DEFAULT 1,
			include_media BOOLEAN DEFAULT 0,
			compression_enabled BOOLEAN DEFAULT 1,
			encryption_enabled BOOLEAN DEFAULT 0,
			retention_days INTEGER DEFAULT 90,
			last_backup_at DATETIME,
			next_backup_at DATETIME,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS file_uploads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			filename TEXT NOT NULL,
			original_filename TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			mime_type TEXT NOT NULL,
			file_hash TEXT NOT NULL,
			upload_type TEXT NOT NULL, -- import, profile_picture, etc.
			status TEXT DEFAULT 'uploaded', -- uploaded, processing, processed, deleted
			processed_at DATETIME,
			deleted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_export_jobs_user_id ON export_jobs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_export_jobs_status ON export_jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_import_jobs_user_id ON import_jobs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_import_jobs_status ON import_jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_data_sync_configs_user_id ON data_sync_configs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sync_logs_user_id ON sync_logs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sync_logs_config_id ON sync_logs(sync_config_id)`,
		`CREATE INDEX IF NOT EXISTS idx_backup_configs_user_id ON backup_configs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_file_uploads_user_id ON file_uploads(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_file_uploads_hash ON file_uploads(file_hash)`,
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

	// Check if user_id column exists in workouts table
	var userIDColumnExists int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('workouts') WHERE name='user_id'`).Scan(&userIDColumnExists)
	if err != nil {
		return fmt.Errorf("failed to check workouts user_id column existence: %v", err)
	}

	// If user_id column doesn't exist in workouts table, add it
	if userIDColumnExists == 0 {
		workoutMigrations := []string{
			`ALTER TABLE workouts ADD COLUMN user_id INTEGER DEFAULT 1`,
			`CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id)`,
		}

		for _, migration := range workoutMigrations {
			if _, err := db.Exec(migration); err != nil {
				return fmt.Errorf("failed to run workouts migration: %s, error: %v", migration, err)
			}
		}
		
		log.Println("Added user_id column to workouts table with default value 1 for existing workouts")
	}

	return nil
}
