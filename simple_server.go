package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	_ "github.com/mattn/go-sqlite3"
)

type Workout struct {
	ID        int
	Name      string
	Date      time.Time
	Duration  int
	Notes     string
	Exercises []Exercise
	CreatedAt time.Time
}

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Role         string    // "admin", "member"
	Provider     string    // "local", "google", "github", "facebook"
	ProviderID   string    // OAuth provider user ID
	AvatarURL    string    // Profile picture URL
	FirstName    string
	LastName     string
	IsActive     bool
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WorkoutTemplate struct {
	ID          int
	Name        string
	Description string
	Category    string
	Exercises   []Exercise
	CreatedAt   time.Time
}

type Exercise struct {
	ID         int
	WorkoutID  int
	TemplateID int
	Name       string
	Sets       int
	Reps       string
	Weight     string
	Notes      string
	WorkoutSets []WorkoutSet
	CreatedAt  time.Time
}

type LiftDatabase struct {
	ID           int
	Name         string
	Category     string
	MuscleGroups string
	Equipment    string
	Description  string
	FormNotes    string
	CreatedAt    time.Time
}

type WorkoutSet struct {
	ID         int
	ExerciseID int
	UserID     int
	SetNumber  int
	Weight     float64
	Reps       int
	RPE        int // Rate of Perceived Exertion (1-10)
	Notes      string
	Completed  bool
	CreatedAt  time.Time
}

type PersonalRecord struct {
	ID           int
	UserID       int
	ExerciseName string
	Weight       float64
	Reps         int
	OneRepMax    float64
	WorkoutID    int
	AchievedAt   time.Time
}

type CardioActivity struct {
	ID          int
	Name        string
	Category    string // Running, Cycling, Swimming, etc.
	Description string
	Metrics     string // Distance, Time, Pace, etc.
	CreatedAt   time.Time
}

type CardioSession struct {
	ID         int
	UserID     int
	Activity   string
	Duration   int     // in seconds
	Distance   float64 // in miles/km
	Pace       string  // min/mile or min/km
	HeartRate  int     // average BPM
	Calories   int
	Notes      string
	WorkoutID  int
	CreatedAt  time.Time
}

type CardioRecord struct {
	ID           int
	UserID       int
	Activity     string
	RecordType   string  // "fastest_pace", "longest_distance", "longest_duration"
	Value        float64 // pace in seconds per mile/km, distance in miles/km, duration in seconds
	DisplayValue string  // formatted display (e.g., "6:30/mile", "10.5 miles", "2:15:30")
	WorkoutID    int
	AchievedAt   time.Time
}

type ScheduledWorkout struct {
	ID          int
	UserID      int
	TemplateID  int
	Title       string
	Description string
	ScheduledAt time.Time
	Status      string // "scheduled", "completed", "skipped"
	WorkoutID   int    // Reference to actual workout when completed
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Nutrition tracking structures
type NutritionEntry struct {
	ID          int
	UserID      int
	Date        time.Time
	MealType    string  // "breakfast", "lunch", "dinner", "snack"
	FoodName    string
	Quantity    float64
	Unit        string  // "grams", "cups", "pieces", etc.
	Calories    int
	Protein     float64 // in grams
	Carbs       float64 // in grams
	Fat         float64 // in grams
	Fiber       float64 // in grams
	Sugar       float64 // in grams
	Sodium      int     // in mg
	Notes       string
	CreatedAt   time.Time
}

type DailyNutrition struct {
	ID             int
	UserID         int
	Date           time.Time
	TotalCalories  int
	TotalProtein   float64
	TotalCarbs     float64
	TotalFat       float64
	TotalFiber     float64
	TotalSugar     float64
	TotalSodium    int
	WaterIntake    float64 // in liters
	CalorieGoal    int
	ProteinGoal    float64
	CarbsGoal      float64
	FatGoal        float64
	WeightKg       float64
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type FoodDatabase struct {
	ID              int
	Name            string
	Brand           string
	Category        string
	ServingSize     float64
	ServingUnit     string
	CaloriesPer100g int
	ProteinPer100g  float64
	CarbsPer100g    float64
	FatPer100g      float64
	FiberPer100g    float64
	SugarPer100g    float64
	SodiumPer100g   int
	Barcode         string
	CreatedAt       time.Time
}

// Advanced analytics structures
type BodyComposition struct {
	ID               int
	UserID           int
	Date             time.Time
	WeightKg         float64
	BodyFatPercent   float64
	MuscleMassKg     float64
	BodyWaterPercent float64
	BoneMassKg       float64
	BMI              float64
	Notes            string
	CreatedAt        time.Time
}

type BodyMeasurements struct {
	ID           int
	UserID       int
	Date         time.Time
	NeckCm       float64
	ChestCm      float64
	WaistCm      float64
	HipsCm       float64
	BicepCm      float64
	ForearmCm    float64
	ThighCm      float64
	CalfCm       float64
	ShouldersCm  float64
	Notes        string
	CreatedAt    time.Time
}

type WorkoutVolume struct {
	ID              int
	UserID          int
	Date            time.Time
	TotalVolume     float64 // Total weight × reps
	TotalSets       int
	TotalReps       int
	TotalWeight     float64
	AverageRPE      float64
	WorkoutDuration int // in minutes
	MuscleGroups    string // JSON array of muscle groups worked
	ExerciseCount   int
	CreatedAt       time.Time
}

type ProgressPhoto struct {
	ID          int
	UserID      int
	Date        time.Time
	PhotoType   string // "front", "side", "back", "pose"
	FilePath    string
	WeightKg    float64
	BodyFat     float64
	Description string
	IsPublic    bool
	CreatedAt   time.Time
}

type WorkoutAnalytics struct {
	ID                    int
	UserID                int
	WeekStart             time.Time
	TotalWorkouts         int
	TotalVolume           float64
	TotalDuration         int
	AverageRPE            float64
	MostWorkedMuscleGroup string
	VolumeIncrease        float64 // Percentage increase from previous week
	ConsistencyScore      float64 // 0-100 based on planned vs actual workouts
	RecoveryScore         float64 // Based on RPE trends and workout frequency
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// Goal tracking structures
type Goal struct {
	ID            int
	UserID        int
	GoalType      string    // "weight", "body_fat", "muscle_mass", "measurement"
	GoalCategory  string    // For measurements: "chest", "waist", "bicep", etc.
	Title         string    // User-friendly goal name
	Description   string    // Detailed description
	CurrentValue  float64   // Starting value
	TargetValue   float64   // Target value
	Unit          string    // "kg", "lbs", "%", "cm", "inches"
	TargetDate    time.Time // When user wants to achieve this
	IsActive      bool      // Whether this goal is currently active
	Priority      int       // 1 (high) to 3 (low)
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time // When goal was achieved (if completed)
}

type GoalProgress struct {
	ID           int
	GoalID       int
	CurrentValue float64
	ProgressDate time.Time
	Note         string
	CreatedAt    time.Time
}

type Server struct {
	db           *sql.DB
	templates    *template.Template
	store        *sessions.CookieStore
	googleOAuth  *oauth2.Config
}

// OAuth provider user info
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewServer() *Server {
	db, err := sql.Open("sqlite3", "simple_tracker.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT,
			role TEXT DEFAULT 'member',
			provider TEXT DEFAULT 'local',
			provider_id TEXT,
			avatar_url TEXT,
			first_name TEXT,
			last_name TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			last_login DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			date DATETIME NOT NULL,
			duration INTEGER DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS workout_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			category TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_id INTEGER,
			template_id INTEGER,
			name TEXT NOT NULL,
			sets INTEGER DEFAULT 1,
			reps TEXT DEFAULT '',
			weight TEXT DEFAULT '',
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (workout_id) REFERENCES workouts(id),
			FOREIGN KEY (template_id) REFERENCES workout_templates(id)
		)`,
		`CREATE TABLE IF NOT EXISTS lift_database (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			category TEXT NOT NULL,
			muscle_groups TEXT NOT NULL,
			equipment TEXT DEFAULT '',
			description TEXT DEFAULT '',
			form_notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workout_sets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			exercise_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			set_number INTEGER NOT NULL,
			weight REAL DEFAULT 0,
			reps INTEGER NOT NULL,
			rpe INTEGER DEFAULT NULL,
			notes TEXT DEFAULT '',
			completed BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (exercise_id) REFERENCES exercises(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS personal_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exercise_name TEXT NOT NULL,
			weight REAL NOT NULL,
			reps INTEGER NOT NULL,
			one_rep_max REAL NOT NULL,
			workout_id INTEGER,
			achieved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (workout_id) REFERENCES workouts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cardio_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			category TEXT NOT NULL,
			description TEXT DEFAULT '',
			metrics TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cardio_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			activity TEXT NOT NULL,
			duration INTEGER DEFAULT 0,
			distance REAL DEFAULT 0,
			pace TEXT DEFAULT '',
			heart_rate INTEGER DEFAULT 0,
			calories INTEGER DEFAULT 0,
			notes TEXT DEFAULT '',
			workout_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (workout_id) REFERENCES workouts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cardio_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			activity TEXT NOT NULL,
			record_type TEXT NOT NULL,
			value REAL NOT NULL,
			display_value TEXT NOT NULL,
			workout_id INTEGER,
			achieved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (workout_id) REFERENCES workouts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS scheduled_workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			template_id INTEGER,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			scheduled_at DATETIME NOT NULL,
			status TEXT DEFAULT 'scheduled',
			workout_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (template_id) REFERENCES workout_templates(id),
			FOREIGN KEY (workout_id) REFERENCES workouts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS nutrition_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			meal_type TEXT NOT NULL,
			food_name TEXT NOT NULL,
			quantity REAL NOT NULL,
			unit TEXT NOT NULL,
			calories INTEGER NOT NULL,
			protein REAL DEFAULT 0,
			carbs REAL DEFAULT 0,
			fat REAL DEFAULT 0,
			fiber REAL DEFAULT 0,
			sugar REAL DEFAULT 0,
			sodium INTEGER DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS daily_nutrition (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL UNIQUE,
			total_calories INTEGER DEFAULT 0,
			total_protein REAL DEFAULT 0,
			total_carbs REAL DEFAULT 0,
			total_fat REAL DEFAULT 0,
			total_fiber REAL DEFAULT 0,
			total_sugar REAL DEFAULT 0,
			total_sodium INTEGER DEFAULT 0,
			water_intake REAL DEFAULT 0,
			calorie_goal INTEGER DEFAULT 2000,
			protein_goal REAL DEFAULT 150,
			carbs_goal REAL DEFAULT 250,
			fat_goal REAL DEFAULT 67,
			weight_kg REAL,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS food_database (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			brand TEXT DEFAULT '',
			category TEXT NOT NULL,
			serving_size REAL DEFAULT 100,
			serving_unit TEXT DEFAULT 'g',
			calories_per_100g INTEGER NOT NULL,
			protein_per_100g REAL DEFAULT 0,
			carbs_per_100g REAL DEFAULT 0,
			fat_per_100g REAL DEFAULT 0,
			fiber_per_100g REAL DEFAULT 0,
			sugar_per_100g REAL DEFAULT 0,
			sodium_per_100g INTEGER DEFAULT 0,
			barcode TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS body_composition (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			weight_kg REAL DEFAULT 0,
			body_fat_percent REAL DEFAULT 0,
			muscle_mass_kg REAL DEFAULT 0,
			body_water_percent REAL DEFAULT 0,
			bone_mass_kg REAL DEFAULT 0,
			bmi REAL DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS body_measurements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			neck_cm REAL DEFAULT 0,
			chest_cm REAL DEFAULT 0,
			waist_cm REAL DEFAULT 0,
			hips_cm REAL DEFAULT 0,
			bicep_cm REAL DEFAULT 0,
			forearm_cm REAL DEFAULT 0,
			thigh_cm REAL DEFAULT 0,
			calf_cm REAL DEFAULT 0,
			shoulders_cm REAL DEFAULT 0,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS workout_volume (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			total_volume REAL DEFAULT 0,
			total_sets INTEGER DEFAULT 0,
			total_reps INTEGER DEFAULT 0,
			total_weight REAL DEFAULT 0,
			average_rpe REAL DEFAULT 0,
			workout_duration INTEGER DEFAULT 0,
			muscle_groups TEXT DEFAULT '',
			exercise_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS progress_photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			photo_type TEXT NOT NULL,
			file_path TEXT NOT NULL,
			weight_kg REAL DEFAULT 0,
			body_fat REAL DEFAULT 0,
			description TEXT DEFAULT '',
			is_public BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS workout_analytics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			week_start DATE NOT NULL,
			total_workouts INTEGER DEFAULT 0,
			total_volume REAL DEFAULT 0,
			total_duration INTEGER DEFAULT 0,
			average_rpe REAL DEFAULT 0,
			most_worked_muscle_group TEXT DEFAULT '',
			volume_increase REAL DEFAULT 0,
			consistency_score REAL DEFAULT 0,
			recovery_score REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS goals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			goal_type TEXT NOT NULL,
			goal_category TEXT DEFAULT '',
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			current_value REAL NOT NULL,
			target_value REAL NOT NULL,
			unit TEXT NOT NULL,
			target_date DATE,
			is_active BOOLEAN DEFAULT TRUE,
			priority INTEGER DEFAULT 2,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS goal_progress (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			goal_id INTEGER NOT NULL,
			current_value REAL NOT NULL,
			progress_date DATE NOT NULL,
			note TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (goal_id) REFERENCES goals(id)
		)`,
	}

	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Seed workout templates, lift database, cardio activities, and food database
	s := &Server{db: db}
	s.seedWorkoutTemplates()
	s.seedLiftDatabase()
	s.seedCardioActivities()
	s.seedFoodDatabase()
	s.createDefaultAdmin()

	// Parse templates with custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a interface{}, b interface{}) interface{} {
			switch v1 := a.(type) {
			case int:
				switch v2 := b.(type) {
				case int:
					return v1 * v2
				case float64:
					return float64(v1) * v2
				}
			case float64:
				switch v2 := b.(type) {
				case int:
					return v1 * float64(v2)
				case float64:
					return v1 * v2
				}
			}
			return 0
		},
		"div": func(a interface{}, b interface{}) interface{} {
			switch v1 := a.(type) {
			case int:
				switch v2 := b.(type) {
				case int:
					if v2 == 0 {
						return 0
					}
					return v1 / v2
				case float64:
					if v2 == 0 {
						return 0
					}
					return float64(v1) / v2
				}
			case float64:
				switch v2 := b.(type) {
				case int:
					if v2 == 0 {
						return 0
					}
					return v1 / float64(v2)
				case float64:
					if v2 == 0 {
						return 0
					}
					return v1 / v2
				}
			}
			return 0
		},
		"gt": func(a interface{}, b interface{}) bool {
			switch v1 := a.(type) {
			case int:
				switch v2 := b.(type) {
				case int:
					return v1 > v2
				case float64:
					return float64(v1) > v2
				}
			case float64:
				switch v2 := b.(type) {
				case int:
					return v1 > float64(v2)
				case float64:
					return v1 > v2
				}
			}
			return false
		},
		"toJson": func(data interface{}) string {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return "[]"
			}
			return string(jsonData)
		},
	}
	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("simple_templates/*.html"))

	// Initialize session store
	store := sessions.NewCookieStore([]byte("workout-secret-key-change-in-production"))

	// Initialize Google OAuth configuration
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	
	var googleOAuth *oauth2.Config
	if clientID != "" && clientSecret != "" {
		googleOAuth = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		}
		log.Println("✅ Google OAuth configured successfully")
	} else {
		log.Println("⚠️  Google OAuth not configured - set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables")
		log.Println("   OAuth login will be disabled until credentials are provided")
	}

	return &Server{
		db:          db,
		templates:   templates,
		store:       store,
		googleOAuth: googleOAuth,
	}
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	workouts, err := s.getWorkoutsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get quick stats for dashboard
	workoutStats, err := s.getWorkoutStats(userID)
	if err != nil {
		log.Println("Error fetching workout stats for home:", err)
	}

	// Get recent PRs for home dashboard
	recentPRs, err := s.getRecentPRs(userID, 3)
	if err != nil {
		log.Println("Error fetching recent PRs for home:", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title        string
		Workouts     []Workout
		Username     string
		UserRole     string
		WorkoutStats WorkoutStats
		RecentPRs    []PersonalRecord
	}{
		Title:        "Workout Tracker",
		Workouts:     workouts,
		Username:     username,
		UserRole:     userRole,
		WorkoutStats: workoutStats,
		RecentPRs:    recentPRs,
	}

	err = s.templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		dateStr := r.FormValue("date")
		notes := r.FormValue("notes")
		templateIDStr := r.FormValue("template_id")

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}

		userID := s.getCurrentUserID(r)
		result, err := s.db.Exec("INSERT INTO workouts (user_id, name, date, notes, created_at) VALUES (?, ?, ?, ?, ?)",
			userID, name, date, notes, time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If a template was selected, copy its exercises to the new workout
		if templateIDStr != "" {
			workoutID, err := result.LastInsertId()
			if err == nil {
				s.copyTemplateExercises(int(workoutID), templateIDStr)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET request - show create form with templates
	templates, err := s.getWorkoutTemplates()
	if err != nil {
		log.Println("Error fetching templates:", err)
	}

	data := struct {
		Title     string
		Templates []WorkoutTemplate
	}{
		Title:     "Create Workout",
		Templates: templates,
	}

	err = s.templates.ExecuteTemplate(w, "create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) EditWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if r.Method == "POST" {
		name := r.FormValue("name")
		dateStr := r.FormValue("date")
		notes := r.FormValue("notes")

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}

		userID := s.getCurrentUserID(r)
		_, err = s.db.Exec("UPDATE workouts SET name = ?, date = ?, notes = ? WHERE id = ? AND user_id = ?",
			name, date, notes, id, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET request - show edit form
	userID := s.getCurrentUserID(r)
	workout, err := s.getWorkoutByID(id, userID)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title   string
		Workout Workout
	}{
		Title:   "Edit Workout",
		Workout: workout,
	}

	err = s.templates.ExecuteTemplate(w, "edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	userID := s.getCurrentUserID(r)
	_, err := s.db.Exec("DELETE FROM workouts WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) getWorkoutByID(id string, userID int) (Workout, error) {
	var w Workout
	err := s.db.QueryRow("SELECT id, name, date, duration, notes, created_at FROM workouts WHERE id = ? AND user_id = ?", id, userID).Scan(
		&w.ID, &w.Name, &w.Date, &w.Duration, &w.Notes, &w.CreatedAt)
	return w, err
}

func (s *Server) getWorkoutsByUser(userID int) ([]Workout, error) {
	rows, err := s.db.Query("SELECT id, name, date, duration, notes, created_at FROM workouts WHERE user_id = ? ORDER BY date DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []Workout
	for rows.Next() {
		var w Workout
		err := rows.Scan(&w.ID, &w.Name, &w.Date, &w.Duration, &w.Notes, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

// Authentication handlers
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := s.getUserByUsername(username)
		if err != nil || !s.checkPasswordHash(password, user.PasswordHash) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		session, _ := s.store.Get(r, "workout-session")
		session.Values["authenticated"] = true
		session.Values["user_id"] = user.ID
		session.Values["username"] = user.Username
		session.Values["role"] = user.Role
		session.Save(r, w)
		
		// Update last login time
		now := time.Now()
		_, err = s.db.Exec("UPDATE users SET last_login = ?, updated_at = ? WHERE id = ?", now, now, user.ID)
		if err != nil {
			log.Printf("Failed to update last login: %v", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := struct {
		Title        string
		OAuthEnabled bool
	}{
		Title:        "Login - Workout Tracker",
		OAuthEnabled: s.googleOAuth != nil,
	}

	err := s.templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate password requirements
		if err := validatePassword(password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		passwordHash, err := s.hashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = s.db.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
			username, email, passwordHash, time.Now())
		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Register - Workout Tracker",
	}

	err := s.templates.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "workout-session")
	session.Values["authenticated"] = false
	delete(session.Values, "user_id")
	delete(session.Values, "username")
	delete(session.Values, "role")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Helper functions
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters long")
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*' || char == '(' || char == ')' || char == '-' || char == '_' || char == '=' || char == '+' || char == '[' || char == ']' || char == '{' || char == '}' || char == '|' || char == ';' || char == ':' || char == '"' || char == '\'' || char == '<' || char == '>' || char == ',' || char == '.' || char == '?' || char == '/':
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

func (s *Server) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Server) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Server) getUserByUsername(username string) (User, error) {
	var user User
	err := s.db.QueryRow("SELECT id, username, email, password_hash, role, created_at FROM users WHERE username = ? AND is_active = true", username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	return user, err
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "workout-session")
		auth, ok := session.Values["authenticated"]
		if !ok || auth != true {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// Role-based middleware
func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "workout-session")
		role, ok := session.Values["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

func (s *Server) getCurrentUserID(r *http.Request) int {
	session, _ := s.store.Get(r, "workout-session")
	userID, _ := session.Values["user_id"].(int)
	return userID
}

// Progress tracking handlers
func (s *Server) ProgressTracker(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	
	// Get lift database
	lifts, err := s.getLiftDatabase()
	if err != nil {
		log.Println("Error fetching lifts:", err)
	}
	
	// Get cardio activities
	cardioActivities, err := s.getCardioActivities()
	if err != nil {
		log.Println("Error fetching cardio activities:", err)
	}
	
	// Get user's detailed PRs for trend analysis
	prs, err := s.getDetailedPRs(userID)
	if err != nil {
		log.Println("Error fetching detailed PRs:", err)
	}

	// Get user's detailed cardio records for trend analysis
	cardioRecords, err := s.getDetailedCardioRecords(userID)
	if err != nil {
		log.Println("Error fetching detailed cardio records:", err)
	}

	// Get workout statistics for analytics dashboard
	workoutStats, err := s.getWorkoutStats(userID)
	if err != nil {
		log.Println("Error fetching workout stats:", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)

	data := struct {
		Title            string
		Username         string
		Lifts            []LiftDatabase
		CardioActivities []CardioActivity
		PRs              []PersonalRecord
		CardioRecords    []CardioRecord
		WorkoutStats     WorkoutStats
	}{
		Title:            "Progress Tracker",
		Username:         username,
		Lifts:            lifts,
		CardioActivities: cardioActivities,
		PRs:              prs,
		CardioRecords:    cardioRecords,
		WorkoutStats:     workoutStats,
	}

	err = s.templates.ExecuteTemplate(w, "progress.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ExportData(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}
	
	switch format {
	case "json":
		s.exportJSON(w, userID)
	case "xml":
		s.exportXML(w, userID)
	default:
		s.exportCSV(w, userID)
	}
}

func (s *Server) exportCSV(w http.ResponseWriter, userID int) {
	fileName := fmt.Sprintf("workout_export_%d.csv", userID)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/csv")
	
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()
	
	// Write CSV Headers
	csvWriter.Write([]string{"Workout ID", "Date", "Workout Name", "Exercise Name", "Set Number", "Reps", "Weight", "RPE", "Notes", "Duration", "Workout Notes"})
	
	// Write workout data
	rows, err := s.db.Query(`
		SELECT w.id, w.date, w.name AS workout_name, e.name AS exercise_name, 
		       ws.set_number, ws.reps, ws.weight, ws.rpe, ws.notes, w.duration, w.notes AS workout_notes
		FROM workouts w
		JOIN exercises e ON e.workout_id = w.id
		JOIN workout_sets ws ON ws.exercise_id = e.id
		WHERE w.user_id = ?
		ORDER BY w.date, w.id, e.id, ws.set_number`, userID)
	if err != nil {
		log.Println("Error exporting CSV data:", err)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var workoutID int
		var date time.Time
		var workoutName, exerciseName, notes, workoutNotes string
		var setNumber, reps, rpe, duration int
		var weight float64
		err := rows.Scan(&workoutID, &date, &workoutName, &exerciseName, &setNumber, &reps, &weight, &rpe, &notes, &duration, &workoutNotes)
		if err != nil {
			continue
		}
		
		dateFormatted := date.Format("2006-01-02")
		weightStr := fmt.Sprintf("%.1f", weight)
		rpeStr := ""
		if rpe > 0 {
			rpeStr = strconv.Itoa(rpe)
		}
		
		csvWriter.Write([]string{
			strconv.Itoa(workoutID),
			dateFormatted,
			workoutName,
			exerciseName,
			strconv.Itoa(setNumber),
			strconv.Itoa(reps),
			weightStr,
			rpeStr,
			notes,
			strconv.Itoa(duration),
			workoutNotes,
		})
	}
}

// Export data structures
type ExportWorkout struct {
	ID        int                `json:"id" xml:"id"`
	Name      string             `json:"name" xml:"name"`
	Date      string             `json:"date" xml:"date"`
	Duration  int                `json:"duration" xml:"duration"`
	Notes     string             `json:"notes" xml:"notes"`
	Exercises []ExportExercise   `json:"exercises" xml:"exercises>exercise"`
}

type ExportExercise struct {
	Name string         `json:"name" xml:"name"`
	Sets []ExportSet    `json:"sets" xml:"sets>set"`
}

type ExportSet struct {
	SetNumber int     `json:"set_number" xml:"set_number"`
	Reps      int     `json:"reps" xml:"reps"`
	Weight    float64 `json:"weight" xml:"weight"`
	RPE       int     `json:"rpe" xml:"rpe"`
	Notes     string  `json:"notes" xml:"notes"`
}

type ExportData struct {
	ExportedAt string          `json:"exported_at" xml:"exported_at"`
	UserID     int             `json:"user_id" xml:"user_id"`
	Workouts   []ExportWorkout `json:"workouts" xml:"workouts>workout"`
}

func (s *Server) exportJSON(w http.ResponseWriter, userID int) {
	data, err := s.getExportData(userID)
	if err != nil {
		log.Println("Error getting export data:", err)
		http.Error(w, "Error preparing export data", http.StatusInternalServerError)
		return
	}
	
	fileName := fmt.Sprintf("workout_export_%d.json", userID)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/json")
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

func (s *Server) exportXML(w http.ResponseWriter, userID int) {
	data, err := s.getExportData(userID)
	if err != nil {
		log.Println("Error getting export data:", err)
		http.Error(w, "Error preparing export data", http.StatusInternalServerError)
		return
	}
	
	fileName := fmt.Sprintf("workout_export_%d.xml", userID)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/xml")
	
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	xmlData, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Error encoding XML:", err)
		return
	}
	w.Write(xmlData)
}

func (s *Server) getExportData(userID int) (*ExportData, error) {
	workouts, err := s.getWorkoutsByUser(userID)
	if err != nil {
		return nil, err
	}
	
	exportWorkouts := make([]ExportWorkout, 0)
	for _, workout := range workouts {
		exportWorkout := ExportWorkout{
			ID:       workout.ID,
			Name:     workout.Name,
			Date:     workout.Date.Format("2006-01-02"),
			Duration: workout.Duration,
			Notes:    workout.Notes,
		}
		
		// Get exercises for this workout
		workoutWithExercises, err := s.getWorkoutWithExercises(strconv.Itoa(workout.ID), userID)
		if err != nil {
			continue
		}
		
		for _, exercise := range workoutWithExercises.Exercises {
			exportExercise := ExportExercise{
				Name: exercise.Name,
				Sets: make([]ExportSet, 0),
			}
			
			for _, set := range exercise.WorkoutSets {
				exportSet := ExportSet{
					SetNumber: set.SetNumber,
					Reps:      set.Reps,
					Weight:    set.Weight,
					RPE:       set.RPE,
					Notes:     set.Notes,
				}
				exportExercise.Sets = append(exportExercise.Sets, exportSet)
			}
			
			exportWorkout.Exercises = append(exportWorkout.Exercises, exportExercise)
		}
		
		exportWorkouts = append(exportWorkouts, exportWorkout)
	}
	
	return &ExportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		UserID:     userID,
		Workouts:   exportWorkouts,
	}, nil
}

// Data Import Functionality
func (s *Server) ImportData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID := s.getCurrentUserID(r)
		
		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		
		file, handler, err := r.FormFile("import_file")
		if err != nil {
			http.Error(w, "Error getting file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		
		// Determine file type
		fileName := handler.Filename
		var importErr error
		
		if strings.HasSuffix(strings.ToLower(fileName), ".json") {
			importErr = s.importFromJSON(file, userID)
		} else if strings.HasSuffix(strings.ToLower(fileName), ".xml") {
			importErr = s.importFromXML(file, userID)
		} else if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
			importErr = s.importFromCSV(file, userID)
		} else {
			http.Error(w, "Unsupported file format. Please use JSON, XML, or CSV.", http.StatusBadRequest)
			return
		}
		
		if importErr != nil {
			log.Printf("Import error: %v", importErr)
			http.Error(w, fmt.Sprintf("Import failed: %v", importErr), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(w, r, "/import?success=1", http.StatusSeeOther)
		return
	}
	
	// GET request - show import form
	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	success := r.URL.Query().Get("success") == "1"
	
	data := struct {
		Title    string
		Username string
		Success  bool
	}{
		Title:    "Import Data",
		Username: username,
		Success:  success,
	}
	
	err := s.templates.ExecuteTemplate(w, "import.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) importFromJSON(file multipart.File, userID int) error {
	var importData ExportData
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&importData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}
	
	return s.importWorkouts(importData.Workouts, userID)
}

func (s *Server) importFromXML(file multipart.File, userID int) error {
	var importData ExportData
	decoder := xml.NewDecoder(file)
	err := decoder.Decode(&importData)
	if err != nil {
		return fmt.Errorf("error parsing XML: %v", err)
	}
	
	return s.importWorkouts(importData.Workouts, userID)
}

func (s *Server) importFromCSV(file multipart.File, userID int) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}
	
	if len(records) < 2 {
		return fmt.Errorf("CSV file appears to be empty or invalid")
	}
	
	// Parse CSV into workout data
	workoutMap := make(map[int]*ExportWorkout)
	
	for i, record := range records[1:] { // Skip header
		if len(record) < 8 {
			continue
		}
		
		workoutID, _ := strconv.Atoi(record[0])
		if workoutID == 0 {
			workoutID = i + 1 // Generate ID if not present
		}
		
		// Create workout if not exists
		if _, exists := workoutMap[workoutID]; !exists {
			workoutMap[workoutID] = &ExportWorkout{
				ID:   workoutID,
				Date: record[1],
				Name: record[2],
			}
			if len(record) > 9 {
				duration, _ := strconv.Atoi(record[9])
				workoutMap[workoutID].Duration = duration
			}
			if len(record) > 10 {
				workoutMap[workoutID].Notes = record[10]
			}
		}
		
		// Find or create exercise
		workout := workoutMap[workoutID]
		exerciseName := record[3]
		var exercise *ExportExercise
		
		for i := range workout.Exercises {
			if workout.Exercises[i].Name == exerciseName {
				exercise = &workout.Exercises[i]
				break
			}
		}
		
		if exercise == nil {
			workout.Exercises = append(workout.Exercises, ExportExercise{
				Name: exerciseName,
				Sets: make([]ExportSet, 0),
			})
			exercise = &workout.Exercises[len(workout.Exercises)-1]
		}
		
		// Add set
		setNumber, _ := strconv.Atoi(record[4])
		reps, _ := strconv.Atoi(record[5])
		weight, _ := strconv.ParseFloat(record[6], 64)
		rpe, _ := strconv.Atoi(record[7])
		notes := ""
		if len(record) > 8 {
			notes = record[8]
		}
		
		exercise.Sets = append(exercise.Sets, ExportSet{
			SetNumber: setNumber,
			Reps:      reps,
			Weight:    weight,
			RPE:       rpe,
			Notes:     notes,
		})
	}
	
	// Convert map to slice
	workouts := make([]ExportWorkout, 0)
	for _, workout := range workoutMap {
		workouts = append(workouts, *workout)
	}
	
	return s.importWorkouts(workouts, userID)
}

func (s *Server) importWorkouts(workouts []ExportWorkout, userID int) error {
	for _, workout := range workouts {
		// Parse date
		date, err := time.Parse("2006-01-02", workout.Date)
		if err != nil {
			date = time.Now() // Fallback to current date
		}
		
		// Create workout
		result, err := s.db.Exec(`
			INSERT INTO workouts (user_id, name, date, duration, notes, created_at) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			userID, workout.Name, date, workout.Duration, workout.Notes, time.Now())
		if err != nil {
			continue // Skip this workout if it fails
		}
		
		workoutID, err := result.LastInsertId()
		if err != nil {
			continue
		}
		
		// Create exercises and sets
		for _, exercise := range workout.Exercises {
			exerciseResult, err := s.db.Exec(`
				INSERT INTO exercises (workout_id, name, sets, reps, weight, notes, created_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				workoutID, exercise.Name, len(exercise.Sets), "", "", "", time.Now())
			if err != nil {
				continue
			}
			
			exerciseID, err := exerciseResult.LastInsertId()
			if err != nil {
				continue
			}
			
			// Create sets
			for _, set := range exercise.Sets {
				_, err = s.db.Exec(`
					INSERT INTO workout_sets (exercise_id, user_id, set_number, weight, reps, rpe, notes, created_at) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
					exerciseID, userID, set.SetNumber, set.Weight, set.Reps, set.RPE, set.Notes, time.Now())
				// Continue even if set creation fails
			}
		}
	}
	
	return nil
}

// OAuth Handlers
func (s *Server) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if s.googleOAuth == nil {
		http.Error(w, "Google OAuth is not configured. Please contact the administrator.", http.StatusServiceUnavailable)
		return
	}
	
	state := "random-state-string" // In production, use a random string and store in session
	url := s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if s.googleOAuth == nil {
		http.Error(w, "Google OAuth is not configured. Please contact the administrator.", http.StatusServiceUnavailable)
		return
	}
	
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	
	// Verify state parameter for CSRF protection
	if state != "random-state-string" {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	
	// Exchange code for token
	token, err := s.googleOAuth.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("OAuth exchange error: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	
	// Get user info from Google
	client := s.googleOAuth.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		http.Error(w, "Failed to read user info", http.StatusInternalServerError)
		return
	}
	
	var googleUser GoogleUser
	err = json.Unmarshal(data, &googleUser)
	if err != nil {
		log.Printf("Failed to unmarshal user info: %v", err)
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}
	
	// Create or update user in database
	user, err := s.createOrUpdateOAuthUser(googleUser)
	if err != nil {
		log.Printf("Failed to create/update user: %v", err)
		http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		return
	}
	
	// Create session
	session, _ := s.store.Get(r, "workout-session")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session: %v", err)
	}
	
	// Update last login time
	now := time.Now()
	_, err = s.db.Exec("UPDATE users SET last_login = ?, updated_at = ? WHERE id = ?", now, now, user.ID)
	if err != nil {
		log.Printf("Failed to update last login: %v", err)
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) createOrUpdateOAuthUser(googleUser GoogleUser) (User, error) {
	var user User
	
	// Check if user exists by email or provider ID
	err := s.db.QueryRow(`
		SELECT id, username, email, role, provider, provider_id, avatar_url, first_name, last_name, is_active, created_at
		FROM users 
		WHERE email = ? OR (provider = 'google' AND provider_id = ?)`,
		googleUser.Email, googleUser.ID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Provider, &user.ProviderID, 
		&user.AvatarURL, &user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		// Create new user
		username := generateUsernameFromEmail(googleUser.Email)
		
		// Make sure username is unique
		username = s.ensureUniqueUsername(username)
		
		result, err := s.db.Exec(`
			INSERT INTO users (username, email, role, provider, provider_id, avatar_url, first_name, last_name, is_active, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			username, googleUser.Email, "member", "google", googleUser.ID, googleUser.Picture, 
			googleUser.GivenName, googleUser.FamilyName, true, time.Now(), time.Now())
		if err != nil {
			return user, err
		}
		
		userID, err := result.LastInsertId()
		if err != nil {
			return user, err
		}
		
		user.ID = int(userID)
		user.Username = username
		user.Email = googleUser.Email
		user.Role = "member"
		user.Provider = "google"
		user.ProviderID = googleUser.ID
		user.AvatarURL = googleUser.Picture
		user.FirstName = googleUser.GivenName
		user.LastName = googleUser.FamilyName
		user.IsActive = true
		user.CreatedAt = time.Now()
	} else if err != nil {
		return user, err
	} else {
		// Update existing user with latest info from Google
		_, err = s.db.Exec(`
			UPDATE users SET avatar_url = ?, first_name = ?, last_name = ?, updated_at = ? 
			WHERE id = ?`,
			googleUser.Picture, googleUser.GivenName, googleUser.FamilyName, time.Now(), user.ID)
		if err != nil {
			return user, err
		}
		
		// Update the user struct with latest info
		user.AvatarURL = googleUser.Picture
		user.FirstName = googleUser.GivenName
		user.LastName = googleUser.FamilyName
	}
	
	return user, nil
}

func generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "user"
}

func (s *Server) ensureUniqueUsername(baseUsername string) string {
	username := baseUsername
	counter := 1
	
	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
		if err != nil || count == 0 {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}
	
	return username
}

// Admin Panel Handlers
func (s *Server) AdminPanel(w http.ResponseWriter, r *http.Request) {
	// Get all users
	users, err := s.getAllUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	
	// Get user statistics
	userStats := make(map[int]map[string]interface{})
	for _, user := range users {
		stats, err := s.getWorkoutStats(user.ID)
		if err == nil {
			userStats[user.ID] = map[string]interface{}{
				"TotalWorkouts": stats.TotalWorkouts,
				"WorkoutsThisWeek": stats.WorkoutsThisWeek,
				"TotalSets": stats.TotalSets,
			}
		}
	}
	
	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	
	data := struct {
		Title     string
		Username  string
		Users     []User
		UserStats map[int]map[string]interface{}
	}{
		Title:     "Admin Panel",
		Username:  username,
		Users:     users,
		UserStats: userStats,
	}
	
	err = s.templates.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) AdminUserManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		userIDStr := r.FormValue("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		
		switch action {
		case "promote":
			err = s.updateUserRole(userID, "admin")
		case "demote":
			err = s.updateUserRole(userID, "member")
		case "activate":
			err = s.updateUserStatus(userID, true)
		case "deactivate":
			err = s.updateUserStatus(userID, false)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	}
	
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) getAllUsers() ([]User, error) {
	rows, err := s.db.Query(`
		SELECT id, username, email, role, provider, avatar_url, first_name, last_name, is_active, last_login, created_at
		FROM users 
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		var lastLogin sql.NullTime
		var avatarURL, firstName, lastName sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Provider,
			&avatarURL, &firstName, &lastName, &user.IsActive, &lastLogin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}
		if firstName.Valid {
			user.FirstName = firstName.String
		}
		if lastName.Valid {
			user.LastName = lastName.String
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Server) updateUserRole(userID int, role string) error {
	_, err := s.db.Exec("UPDATE users SET role = ?, updated_at = ? WHERE id = ?", role, time.Now(), userID)
	return err
}

func (s *Server) updateUserStatus(userID int, isActive bool) error {
	_, err := s.db.Exec("UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?", isActive, time.Now(), userID)
	return err
}

func (s *Server) Schedule(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	
	if r.Method == "POST" {
		// Handle scheduling a new workout
		action := r.FormValue("action")
		
		if action == "schedule" {
			title := r.FormValue("title")
			description := r.FormValue("description")
			scheduledDateStr := r.FormValue("scheduled_date")
			scheduledTimeStr := r.FormValue("scheduled_time")
			templateIDStr := r.FormValue("template_id")
			
			// Parse date and time
			dateTimeStr := scheduledDateStr + " " + scheduledTimeStr + ":00"
			scheduledAt, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
			if err != nil {
				http.Error(w, "Invalid date/time format", http.StatusBadRequest)
				return
			}
			
			var templateID *int
			if templateIDStr != "" {
				tid, err := strconv.Atoi(templateIDStr)
				if err == nil {
					templateID = &tid
				}
			}
			
			// Insert scheduled workout
			_, err = s.db.Exec("INSERT INTO scheduled_workouts (user_id, template_id, title, description, scheduled_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
				userID, templateID, title, description, scheduledAt, time.Now(), time.Now())
			if err != nil {
				http.Error(w, "Error scheduling workout: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else if action == "complete" {
			// Mark scheduled workout as completed
			scheduledID := r.FormValue("scheduled_id")
			_, err := s.db.Exec("UPDATE scheduled_workouts SET status = 'completed', updated_at = ? WHERE id = ? AND user_id = ?",
				time.Now(), scheduledID, userID)
			if err != nil {
				http.Error(w, "Error updating workout status", http.StatusInternalServerError)
				return
			}
		} else if action == "skip" {
			// Mark scheduled workout as skipped
			scheduledID := r.FormValue("scheduled_id")
			_, err := s.db.Exec("UPDATE scheduled_workouts SET status = 'skipped', updated_at = ? WHERE id = ? AND user_id = ?",
				time.Now(), scheduledID, userID)
			if err != nil {
				http.Error(w, "Error updating workout status", http.StatusInternalServerError)
				return
			}
		}
		
		http.Redirect(w, r, "/schedule", http.StatusSeeOther)
		return
	}
	
	// GET - Show schedule calendar
	// Get current month's scheduled workouts
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	
	scheduledWorkouts, err := s.getScheduledWorkouts(userID, startOfMonth, endOfMonth)
	if err != nil {
		log.Println("Error fetching scheduled workouts:", err)
	}
	
	// Get workout templates for scheduling
	templates, err := s.getWorkoutTemplates()
	if err != nil {
		log.Println("Error fetching templates:", err)
	}
	
	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	
	data := struct {
		Title             string
		Username          string
		ScheduledWorkouts []ScheduledWorkout
		Templates         []WorkoutTemplate
		CurrentMonth      time.Time
	}{
		Title:             "Workout Schedule",
		Username:          username,
		ScheduledWorkouts: scheduledWorkouts,
		Templates:         templates,
		CurrentMonth:      startOfMonth,
	}
	
	err = s.templates.ExecuteTemplate(w, "schedule.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) LogWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workoutID := vars["id"]
	userID := s.getCurrentUserID(r)

	if r.Method == "POST" {
		// Handle set logging
		exerciseID := r.FormValue("exercise_id")
		setNumber := r.FormValue("set_number")
		weight := r.FormValue("weight")
		reps := r.FormValue("reps")
		rpe := r.FormValue("rpe")
		notes := r.FormValue("notes")
		exerciseName := r.FormValue("exercise_name")

		// Convert string values
		weightFloat, _ := strconv.ParseFloat(weight, 64)
		repsInt, _ := strconv.Atoi(reps)
		rpeInt, _ := strconv.Atoi(rpe)
		setInt, _ := strconv.Atoi(setNumber)
		exerciseIDInt, _ := strconv.Atoi(exerciseID)
		workoutIDInt, _ := strconv.Atoi(workoutID)

		// Save the set
		_, err := s.db.Exec("INSERT INTO workout_sets (exercise_id, user_id, set_number, weight, reps, rpe, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			exerciseIDInt, userID, setInt, weightFloat, repsInt, rpeInt, notes, time.Now())
		if err != nil {
			http.Error(w, "Error saving set: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for PR
		if weightFloat > 0 && repsInt > 0 {
			s.updatePersonalRecord(userID, exerciseName, weightFloat, repsInt, workoutIDInt)
		}

		http.Redirect(w, r, "/workout/"+workoutID+"/log", http.StatusSeeOther)
		return
	}

	// GET - Show workout logging interface
	workout, err := s.getWorkoutWithExercises(workoutID, userID)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)

	data := struct {
		Title    string
		Username string
		Workout  Workout
	}{
		Title:    "Log Workout",
		Username: username,
		Workout:  workout,
	}

	err = s.templates.ExecuteTemplate(w, "log_workout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func (s *Server) seedWorkoutTemplates() {
	// Check if templates already exist to avoid duplicates
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM workout_templates").Scan(&count)
	if err != nil || count > 0 {
		return // Templates already exist or error occurred
	}

	templates := []WorkoutTemplate{
		{Name: "Full Body Workout", Description: "A balanced workout for all major muscle groups.", Category: "Strength", Exercises: []Exercise{
			{Name: "Push-ups", Sets: 3, Reps: "10-15"},
			{Name: "Pull-ups", Sets: 3, Reps: "5-10"},
			{Name: "Squats", Sets: 3, Reps: "15-20"},
			{Name: "Planks", Sets: 3, Reps: "30s-1min"},
		}},
		{Name: "Upper Body Strength", Description: "Focus on chest, back, shoulders and arms.", Category: "Strength", Exercises: []Exercise{
			{Name: "Bench Press", Sets: 4, Reps: "8-12"},
			{Name: "Rows", Sets: 4, Reps: "8-12"},
			{Name: "Overhead Press", Sets: 3, Reps: "8-10"},
			{Name: "Bicep Curls", Sets: 3, Reps: "12-15"},
			{Name: "Tricep Dips", Sets: 3, Reps: "10-15"},
		}},
		{Name: "Lower Body Power", Description: "Build strong legs and glutes.", Category: "Strength", Exercises: []Exercise{
			{Name: "Squats", Sets: 4, Reps: "12-15"},
			{Name: "Deadlifts", Sets: 4, Reps: "8-10"},
			{Name: "Lunges", Sets: 3, Reps: "12 each leg"},
			{Name: "Calf Raises", Sets: 3, Reps: "15-20"},
		}},
		{Name: "Cardio Blast", Description: "High intensity cardio exercises.", Category: "Cardio", Exercises: []Exercise{
			{Name: "Running", Sets: 1, Reps: "30 min"},
			{Name: "Jump Rope", Sets: 3, Reps: "2 min"},
			{Name: "Burpees", Sets: 3, Reps: "10-20"},
			{Name: "Mountain Climbers", Sets: 3, Reps: "30s"},
		}},
		{Name: "HIIT Training", Description: "Quick high-intensity interval training.", Category: "Cardio", Exercises: []Exercise{
			{Name: "High Knees", Sets: 4, Reps: "30s on, 15s rest"},
			{Name: "Jumping Jacks", Sets: 4, Reps: "30s on, 15s rest"},
			{Name: "Sprint Intervals", Sets: 8, Reps: "20s on, 40s rest"},
		}},
		{Name: "Yoga Flow", Description: "Gentle stretching and flexibility.", Category: "Flexibility", Exercises: []Exercise{
			{Name: "Sun Salutation", Sets: 5, Reps: "1 flow"},
			{Name: "Downward Dog", Sets: 1, Reps: "1 min hold"},
			{Name: "Warrior Poses", Sets: 2, Reps: "30s each side"},
			{Name: "Child's Pose", Sets: 1, Reps: "2 min hold"},
		}},
	}

	for _, template := range templates {
		result, err := s.db.Exec("INSERT INTO workout_templates (name, description, category, created_at) VALUES (?, ?, ?, ?)", template.Name, template.Description, template.Category, time.Now())
		if err != nil {
			log.Println("Failed to seed template:", template.Name, err)
			continue
		}
		templateID, err := result.LastInsertId()
		if err != nil {
			log.Println("Failed to get template ID:", err)
			continue
		}

		for _, exercise := range template.Exercises {
			_, err := s.db.Exec("INSERT INTO exercises (template_id, name, sets, reps, created_at) VALUES (?, ?, ?, ?, ?)", templateID, exercise.Name, exercise.Sets, exercise.Reps, time.Now())
			if err != nil {
				log.Println("Failed to seed exercise:", exercise.Name, err)
			}
		}
	}
}

func (s *Server) getWorkoutTemplates() ([]WorkoutTemplate, error) {
	rows, err := s.db.Query("SELECT id, name, description, category FROM workout_templates ORDER BY category, name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []WorkoutTemplate
	for rows.Next() {
		var t WorkoutTemplate
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Category)
		if err != nil {
			return nil, err
		}

		// Get exercises for this template
		exerciseRows, err := s.db.Query("SELECT name, sets, reps FROM exercises WHERE template_id = ? ORDER BY id", t.ID)
		if err != nil {
			return nil, err
		}
		defer exerciseRows.Close()

		for exerciseRows.Next() {
			var e Exercise
			err := exerciseRows.Scan(&e.Name, &e.Sets, &e.Reps)
			if err != nil {
				return nil, err
			}
			t.Exercises = append(t.Exercises, e)
		}
		templates = append(templates, t)
	}

	return templates, nil
}

func (s *Server) copyTemplateExercises(workoutID int, templateIDStr string) {
	// Get exercises from the template first, then close the query
	rows, err := s.db.Query("SELECT name, sets, reps FROM exercises WHERE template_id = ? ORDER BY id", templateIDStr)
	if err != nil {
		log.Println("Error fetching template exercises:", err)
		return
	}

	// Read all exercises into memory first
	var exercises []Exercise
	for rows.Next() {
		var e Exercise
		err := rows.Scan(&e.Name, &e.Sets, &e.Reps)
		if err != nil {
			log.Println("Error scanning exercise:", err)
			continue
		}
		exercises = append(exercises, e)
	}
	rows.Close() // Close the query before doing inserts

	// Copy exercises to the new workout
	for _, exercise := range exercises {
		_, err = s.db.Exec("INSERT INTO exercises (workout_id, name, sets, reps, created_at) VALUES (?, ?, ?, ?, ?)",
			workoutID, exercise.Name, exercise.Sets, exercise.Reps, time.Now())
		if err != nil {
			log.Println("Error copying exercise:", exercise.Name, err)
		}
	}
}

// Helper functions for strength tracking
func calculateOneRepMax(weight float64, reps int) float64 {
	// Using Brzycki formula: Weight * (36 / (37 - reps))
	if reps == 1 {
		return weight
	}
	if reps > 15 {
		reps = 15 // Cap at 15 reps for accuracy
	}
	return weight * (36.0 / (37.0 - float64(reps)))
}

func (s *Server) updatePersonalRecord(userID int, exerciseName string, weight float64, reps int, workoutID int) {
	oneRepMax := calculateOneRepMax(weight, reps)
	
	// Check if this is a new PR for this rep range (1RM, 3RM, 5RM, 8RM, 10RM)
	repRanges := []int{1, 3, 5, 8, 10}
	for _, repRange := range repRanges {
		if reps <= repRange {
			// Check existing PR for this rep range
			var existingWeight float64
			err := s.db.QueryRow("SELECT weight FROM personal_records WHERE user_id = ? AND exercise_name = ? AND reps <= ? ORDER BY weight DESC LIMIT 1", 
				userID, exerciseName, repRange).Scan(&existingWeight)
			
			if err != nil || weight > existingWeight {
				// New PR! Insert record
				_, err = s.db.Exec("INSERT INTO personal_records (user_id, exercise_name, weight, reps, one_rep_max, workout_id, achieved_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
					userID, exerciseName, weight, reps, oneRepMax, workoutID, time.Now())
				if err != nil {
					log.Println("Error recording PR:", err)
				}
				break // Only record one PR per lift
			}
		}
	}
}

func (s *Server) getLiftDatabase() ([]LiftDatabase, error) {
	rows, err := s.db.Query("SELECT id, name, category, muscle_groups, equipment, description, form_notes FROM lift_database ORDER BY category, name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lifts []LiftDatabase
	for rows.Next() {
		var lift LiftDatabase
		err := rows.Scan(&lift.ID, &lift.Name, &lift.Category, &lift.MuscleGroups, &lift.Equipment, &lift.Description, &lift.FormNotes)
		if err != nil {
			return nil, err
		}
		lifts = append(lifts, lift)
	}
	return lifts, nil
}

func (s *Server) getRecentPRs(userID int, limit int) ([]PersonalRecord, error) {
	rows, err := s.db.Query("SELECT id, exercise_name, weight, reps, one_rep_max, achieved_at FROM personal_records WHERE user_id = ? ORDER BY achieved_at DESC LIMIT ?", userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []PersonalRecord
	for rows.Next() {
		var pr PersonalRecord
		err := rows.Scan(&pr.ID, &pr.ExerciseName, &pr.Weight, &pr.Reps, &pr.OneRepMax, &pr.AchievedAt)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (s *Server) getAllPRs(userID int) (map[string][]PersonalRecord, error) {
	rows, err := s.db.Query("SELECT id, exercise_name, weight, reps, one_rep_max, achieved_at FROM personal_records WHERE user_id = ? ORDER BY exercise_name, achieved_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prsMap := make(map[string][]PersonalRecord)
	for rows.Next() {
		var pr PersonalRecord
		err := rows.Scan(&pr.ID, &pr.ExerciseName, &pr.Weight, &pr.Reps, &pr.OneRepMax, &pr.AchievedAt)
		if err != nil {
			return nil, err
		}
		prsMap[pr.ExerciseName] = append(prsMap[pr.ExerciseName], pr)
	}
	return prsMap, nil
}

func (s *Server) getCardioActivities() ([]CardioActivity, error) {
	rows, err := s.db.Query("SELECT id, name, category, description, metrics FROM cardio_activities ORDER BY category, name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []CardioActivity
	for rows.Next() {
		var activity CardioActivity
		err := rows.Scan(&activity.ID, &activity.Name, &activity.Category, &activity.Description, &activity.Metrics)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (s *Server) getRecentCardioRecords(userID int, limit int) ([]CardioRecord, error) {
	rows, err := s.db.Query("SELECT id, activity, record_type, display_value, achieved_at FROM cardio_records WHERE user_id = ? ORDER BY achieved_at DESC LIMIT ?", userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []CardioRecord
	for rows.Next() {
		var record CardioRecord
		err := rows.Scan(&record.ID, &record.Activity, &record.RecordType, &record.DisplayValue, &record.AchievedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// New detailed functions for enhanced progress tracking
func (s *Server) getDetailedPRs(userID int) ([]PersonalRecord, error) {
	rows, err := s.db.Query(`
		SELECT pr.id, pr.exercise_name, pr.weight, pr.reps, pr.one_rep_max, pr.achieved_at, 
		       w.name as workout_name
		FROM personal_records pr
		LEFT JOIN workouts w ON pr.workout_id = w.id
		WHERE pr.user_id = ? 
		ORDER BY pr.achieved_at DESC 
		LIMIT 20`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []PersonalRecord
	for rows.Next() {
		var pr PersonalRecord
		var workoutName sql.NullString
		err := rows.Scan(&pr.ID, &pr.ExerciseName, &pr.Weight, &pr.Reps, &pr.OneRepMax, &pr.AchievedAt, &workoutName)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (s *Server) getDetailedCardioRecords(userID int) ([]CardioRecord, error) {
	rows, err := s.db.Query(`
		SELECT cr.id, cr.activity, cr.record_type, cr.display_value, cr.achieved_at,
		       w.name as workout_name
		FROM cardio_records cr
		LEFT JOIN workouts w ON cr.workout_id = w.id
		WHERE cr.user_id = ? 
		ORDER BY cr.achieved_at DESC 
		LIMIT 20`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []CardioRecord
	for rows.Next() {
		var record CardioRecord
		var workoutName sql.NullString
		err := rows.Scan(&record.ID, &record.Activity, &record.RecordType, &record.DisplayValue, &record.AchievedAt, &workoutName)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// WorkoutStats holds workout statistics for analytics
type WorkoutStats struct {
	TotalWorkouts    int
	WorkoutsThisWeek int
	WorkoutsThisMonth int
	TotalSets        int
	TotalReps        int
	AverageWorkoutsPerWeek float64
	MostActiveDay    string
	TotalWeightLifted float64
}

func (s *Server) getWorkoutStats(userID int) (WorkoutStats, error) {
	var stats WorkoutStats
	
	// Get total workouts
	err := s.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE user_id = ?", userID).Scan(&stats.TotalWorkouts)
	if err != nil {
		return stats, err
	}
	
	// Get workouts this week
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	err = s.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE user_id = ? AND created_at >= ?", userID, startOfWeek).Scan(&stats.WorkoutsThisWeek)
	if err != nil {
		return stats, err
	}
	
	// Get workouts this month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err = s.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE user_id = ? AND created_at >= ?", userID, startOfMonth).Scan(&stats.WorkoutsThisMonth)
	if err != nil {
		return stats, err
	}
	
	// Get total sets and reps
	err = s.db.QueryRow(`
		SELECT COUNT(ws.id) as total_sets, COALESCE(SUM(ws.reps), 0) as total_reps
		FROM workout_sets ws 
		JOIN exercises e ON ws.exercise_id = e.id 
		JOIN workouts w ON e.workout_id = w.id 
		WHERE ws.user_id = ?`, userID).Scan(&stats.TotalSets, &stats.TotalReps)
	if err != nil {
		return stats, err
	}
	
	// Get total weight lifted
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(ws.weight * ws.reps), 0) as total_weight
		FROM workout_sets ws 
		JOIN exercises e ON ws.exercise_id = e.id 
		JOIN workouts w ON e.workout_id = w.id 
		WHERE ws.user_id = ?`, userID).Scan(&stats.TotalWeightLifted)
	if err != nil {
		return stats, err
	}
	
	// Calculate average workouts per week (last 12 weeks)
	last12Weeks := now.AddDate(0, 0, -84) // 12 weeks ago
	var workoutsLast12Weeks int
	err = s.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE user_id = ? AND created_at >= ?", userID, last12Weeks).Scan(&workoutsLast12Weeks)
	if err == nil {
		weeksPassed := float64(now.Sub(last12Weeks).Hours()) / (24 * 7)
		if weeksPassed > 0 {
			stats.AverageWorkoutsPerWeek = float64(workoutsLast12Weeks) / weeksPassed
		}
	}
	
	// Get most active day of week
	var dayName string
	err = s.db.QueryRow(`
		SELECT 
			CASE strftime('%w', date) 
				WHEN '0' THEN 'Sunday'
				WHEN '1' THEN 'Monday'
				WHEN '2' THEN 'Tuesday'
				WHEN '3' THEN 'Wednesday'
				WHEN '4' THEN 'Thursday'
				WHEN '5' THEN 'Friday'
				WHEN '6' THEN 'Saturday'
			END as day_name
		FROM workouts 
		WHERE user_id = ?
		GROUP BY strftime('%w', date)
		ORDER BY COUNT(*) DESC
		LIMIT 1`, userID).Scan(&dayName)
	if err == nil {
		stats.MostActiveDay = dayName
	} else {
		stats.MostActiveDay = "No data"
	}
	
	return stats, nil
}

func (s *Server) getScheduledWorkouts(userID int, startDate, endDate time.Time) ([]ScheduledWorkout, error) {
	rows, err := s.db.Query("SELECT id, user_id, template_id, title, description, scheduled_at, status, workout_id, created_at, updated_at FROM scheduled_workouts WHERE user_id = ? AND scheduled_at BETWEEN ? AND ? ORDER BY scheduled_at", userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []ScheduledWorkout
	for rows.Next() {
		var workout ScheduledWorkout
		var templateID sql.NullInt64
		var workoutID sql.NullInt64
		err := rows.Scan(&workout.ID, &workout.UserID, &templateID, &workout.Title, &workout.Description, &workout.ScheduledAt, &workout.Status, &workoutID, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if templateID.Valid {
			workout.TemplateID = int(templateID.Int64)
		}
		if workoutID.Valid {
			workout.WorkoutID = int(workoutID.Int64)
		}
		workouts = append(workouts, workout)
	}
	return workouts, nil
}

func (s *Server) getWorkoutWithExercises(workoutID string, userID int) (Workout, error) {
	workout, err := s.getWorkoutByID(workoutID, userID)
	if err != nil {
		return workout, err
	}

	// Get exercises for this workout
	rows, err := s.db.Query("SELECT id, name, sets, reps, weight, notes FROM exercises WHERE workout_id = ? ORDER BY id", workoutID)
	if err != nil {
		return workout, err
	}
	defer rows.Close()

	for rows.Next() {
		var exercise Exercise
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Sets, &exercise.Reps, &exercise.Weight, &exercise.Notes)
		if err != nil {
			return workout, err
		}

		// Get sets for this exercise
		setRows, err := s.db.Query("SELECT id, set_number, weight, reps, rpe, notes FROM workout_sets WHERE exercise_id = ? ORDER BY set_number", exercise.ID)
		if err == nil {
			for setRows.Next() {
				var workoutSet WorkoutSet
				err := setRows.Scan(&workoutSet.ID, &workoutSet.SetNumber, &workoutSet.Weight, &workoutSet.Reps, &workoutSet.RPE, &workoutSet.Notes)
				if err == nil {
					exercise.WorkoutSets = append(exercise.WorkoutSets, workoutSet)
				}
			}
			setRows.Close()
		}

		workout.Exercises = append(workout.Exercises, exercise)
	}

	return workout, nil
}

func (s *Server) seedLiftDatabase() {
	// Check if lifts already exist to avoid duplicates
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM lift_database").Scan(&count)
	if err != nil || count > 0 {
		return // Lifts already exist or error occurred
	}

	lifts := []LiftDatabase{
		// Chest
		{Name: "Bench Press", Category: "Compound", MuscleGroups: "Chest, Triceps, Shoulders", Equipment: "Barbell", Description: "The king of upper body exercises", FormNotes: "Keep shoulder blades retracted, feet flat on floor, controlled descent"},
		{Name: "Incline Bench Press", Category: "Compound", MuscleGroups: "Upper Chest, Triceps, Shoulders", Equipment: "Barbell", Description: "Targets upper chest development", FormNotes: "30-45 degree incline, same form as flat bench"},
		{Name: "Dumbbell Press", Category: "Compound", MuscleGroups: "Chest, Triceps, Shoulders", Equipment: "Dumbbells", Description: "Greater range of motion than barbell", FormNotes: "Control the weight, don't let dumbbells drift too wide"},
		{Name: "Pec Deck", Category: "Isolation", MuscleGroups: "Chest", Equipment: "Machine", Description: "Chest isolation movement", FormNotes: "Squeeze at the peak, control on the way back"},

		// Back
		{Name: "Deadlift", Category: "Compound", MuscleGroups: "Back, Glutes, Hamstrings, Traps", Equipment: "Barbell", Description: "The ultimate full-body strength exercise", FormNotes: "Keep bar close to body, chest up, neutral spine, drive through heels"},
		{Name: "Pull-ups", Category: "Compound", MuscleGroups: "Lats, Rhomboids, Biceps", Equipment: "Pull-up Bar", Description: "Bodyweight back builder", FormNotes: "Full range of motion, control the negative, avoid swinging"},
		{Name: "Barbell Rows", Category: "Compound", MuscleGroups: "Lats, Rhomboids, Mid-traps", Equipment: "Barbell", Description: "Builds thick, strong back", FormNotes: "Hinge at hips, pull to lower chest/upper abdomen"},
		{Name: "Chin-ups", Category: "Compound", MuscleGroups: "Biceps, Lats", Equipment: "Bar", Description: "Underhand grip variation for biceps", FormNotes: "Full range of motion, control the negative, avoid swinging"},

		// Legs
		{Name: "Squat", Category: "Compound", MuscleGroups: "Quadriceps, Glutes, Hamstrings", Equipment: "Barbell", Description: "The king of lower body exercises", FormNotes: "Feet shoulder-width apart, knees track over toes, go to at least parallel"},
		{Name: "Front Squat", Category: "Compound", MuscleGroups: "Quadriceps, Glutes, Core", Equipment: "Barbell", Description: "More quad-focused squat variation", FormNotes: "Keep elbows high, upright torso, bar rests on front delts"},
		{Name: "Romanian Deadlift", Category: "Compound", MuscleGroups: "Hamstrings, Glutes, Lower Back", Equipment: "Barbell", Description: "Hip-hinge movement for posterior chain", FormNotes: "Push hips back, slight knee bend, feel stretch in hamstrings"},
		{Name: "Leg Press", Category: "Compound", MuscleGroups: "Quadriceps, Glutes", Equipment: "Machine", Description: "Lower body development on a sled", FormNotes: "Knees track over toes, don't lock out knees at top"},

		// Shoulders
		{Name: "Overhead Press", Category: "Compound", MuscleGroups: "Shoulders, Triceps, Core", Equipment: "Barbell", Description: "Standing military press", FormNotes: "Core tight, press straight up, lockout overhead"},
		{Name: "Dumbbell Shoulder Press", Category: "Compound", MuscleGroups: "Shoulders, Triceps", Equipment: "Dumbbells", Description: "Seated or standing shoulder development", FormNotes: "Control the weight, don't arch back excessively"},
		{Name: "Lateral Raise", Category: "Isolation", MuscleGroups: "Shoulders", Equipment: "Dumbbells", Description: "Deltoid side isolation exercise", FormNotes: "Raise arms to the side, don't swing"},

		// Arms
		{Name: "Barbell Curls", Category: "Isolation", MuscleGroups: "Biceps", Equipment: "Barbell", Description: "Classic bicep builder", FormNotes: "Keep elbows stable, control the negative, full range of motion"},
		{Name: "Close-Grip Bench Press", Category: "Compound", MuscleGroups: "Triceps, Chest", Equipment: "Barbell", Description: "Tricep-focused bench variation", FormNotes: "Hands about shoulder-width apart, keep elbows close to body"},
		{Name: "Dips", Category: "Compound", MuscleGroups: "Triceps, Chest, Shoulders", Equipment: "Dip Bars", Description: "Bodyweight tricep and chest builder", FormNotes: "Lean forward slightly for chest emphasis, stay upright for triceps"},
		{Name: "Tricep Pushdowns", Category: "Isolation", MuscleGroups: "Triceps", Equipment: "Cable", Description: "Tricep cable exercise", FormNotes: "Elbows at side, full extension at bottom"},
	}

	for _, lift := range lifts {
		_, err := s.db.Exec("INSERT INTO lift_database (name, category, muscle_groups, equipment, description, form_notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			lift.Name, lift.Category, lift.MuscleGroups, lift.Equipment, lift.Description, lift.FormNotes, time.Now())
		if err != nil {
			log.Println("Failed to seed lift:", lift.Name, err)
		}
	}
}

func (s *Server) seedCardioActivities() {
	// Check if cardio activities already exist to avoid duplicates
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM cardio_activities").Scan(&count)
	if err != nil || count > 0 {
		return // Activities already exist or error occurred
	}

	activities := []CardioActivity{
		// Running
		{Name: "Running", Category: "Running", Description: "Outdoor or treadmill running", Metrics: "Distance, Time, Pace, Heart Rate"},
		{Name: "Jogging", Category: "Running", Description: "Light jogging for endurance", Metrics: "Distance, Time, Pace, Heart Rate"},
		{Name: "Trail Running", Category: "Running", Description: "Off-road running on trails", Metrics: "Distance, Time, Elevation, Heart Rate"},
		{Name: "Sprints", Category: "Running", Description: "High-intensity sprint intervals", Metrics: "Distance, Time, Rest Periods"},
		
		// Cycling
		{Name: "Road Cycling", Category: "Cycling", Description: "Road bike cycling", Metrics: "Distance, Time, Speed, Heart Rate"},
		{Name: "Mountain Biking", Category: "Cycling", Description: "Off-road cycling", Metrics: "Distance, Time, Elevation, Heart Rate"},
		{Name: "Stationary Bike", Category: "Cycling", Description: "Indoor cycling/spin bike", Metrics: "Time, Resistance, Heart Rate, Calories"},
		
		// Swimming
		{Name: "Swimming", Category: "Swimming", Description: "Pool or open water swimming", Metrics: "Distance, Time, Stroke Rate"},
		{Name: "Water Jogging", Category: "Swimming", Description: "Aqua jogging for low impact cardio", Metrics: "Time, Intensity, Heart Rate"},
		
		// Other
		{Name: "Rowing", Category: "Rowing", Description: "Rowing machine or water rowing", Metrics: "Distance, Time, Stroke Rate, Heart Rate"},
		{Name: "Elliptical", Category: "Machine", Description: "Elliptical trainer workout", Metrics: "Time, Resistance, Heart Rate, Calories"},
		{Name: "Stairmaster", Category: "Machine", Description: "Stair climbing machine", Metrics: "Time, Level, Heart Rate, Calories"},
		{Name: "Jump Rope", Category: "Bodyweight", Description: "Skipping rope cardio", Metrics: "Time, Jumps per minute, Rest periods"},
		{Name: "HIIT", Category: "Interval", Description: "High-intensity interval training", Metrics: "Work time, Rest time, Rounds, Heart Rate"},
	}

	for _, activity := range activities {
		_, err := s.db.Exec("INSERT INTO cardio_activities (name, category, description, metrics, created_at) VALUES (?, ?, ?, ?, ?)",
			activity.Name, activity.Category, activity.Description, activity.Metrics, time.Now())
		if err != nil {
			log.Println("Failed to seed cardio activity:", activity.Name, err)
		}
	}
}

func (s *Server) seedFoodDatabase() {
	// Check if foods already exist to avoid duplicates
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM food_database").Scan(&count)
	if err != nil || count > 0 {
		return // Foods already exist or error occurred
	}

	foods := []FoodDatabase{
		// Common proteins
		{Name: "Chicken Breast", Brand: "Generic", Category: "Meat", ServingSize: 100, ServingUnit: "g", CaloriesPer100g: 165, ProteinPer100g: 31, CarbsPer100g: 0, FatPer100g: 3.6, FiberPer100g: 0, SugarPer100g: 0, SodiumPer100g: 74},
		{Name: "Salmon", Brand: "Generic", Category: "Fish", ServingSize: 100, ServingUnit: "g", CaloriesPer100g: 208, ProteinPer100g: 22, CarbsPer100g: 0, FatPer100g: 12, FiberPer100g: 0, SugarPer100g: 0, SodiumPer100g: 59},
		{Name: "Eggs", Brand: "Generic", Category: "Dairy", ServingSize: 50, ServingUnit: "g", CaloriesPer100g: 155, ProteinPer100g: 13, CarbsPer100g: 1.1, FatPer100g: 11, FiberPer100g: 0, SugarPer100g: 1.1, SodiumPer100g: 124},
		{Name: "Greek Yogurt", Brand: "Generic", Category: "Dairy", ServingSize: 170, ServingUnit: "g", CaloriesPer100g: 59, ProteinPer100g: 10, CarbsPer100g: 3.6, FatPer100g: 0.4, FiberPer100g: 0, SugarPer100g: 3.2, SodiumPer100g: 36},

		// Carbohydrates
		{Name: "Brown Rice", Brand: "Generic", Category: "Grains", ServingSize: 100, ServingUnit: "g", CaloriesPer100g: 111, ProteinPer100g: 2.6, CarbsPer100g: 23, FatPer100g: 0.9, FiberPer100g: 1.8, SugarPer100g: 0.4, SodiumPer100g: 5},
		{Name: "Oats", Brand: "Generic", Category: "Grains", ServingSize: 40, ServingUnit: "g", CaloriesPer100g: 389, ProteinPer100g: 16.9, CarbsPer100g: 66.3, FatPer100g: 6.9, FiberPer100g: 10.6, SugarPer100g: 0.4, SodiumPer100g: 2},
		{Name: "Sweet Potato", Brand: "Generic", Category: "Vegetables", ServingSize: 128, ServingUnit: "g", CaloriesPer100g: 86, ProteinPer100g: 1.6, CarbsPer100g: 20.1, FatPer100g: 0.1, FiberPer100g: 3, SugarPer100g: 4.2, SodiumPer100g: 54},
		{Name: "Banana", Brand: "Generic", Category: "Fruits", ServingSize: 118, ServingUnit: "g", CaloriesPer100g: 89, ProteinPer100g: 1.1, CarbsPer100g: 22.8, FatPer100g: 0.3, FiberPer100g: 2.6, SugarPer100g: 12.2, SodiumPer100g: 1},

		// Fats
		{Name: "Avocado", Brand: "Generic", Category: "Fruits", ServingSize: 150, ServingUnit: "g", CaloriesPer100g: 160, ProteinPer100g: 2, CarbsPer100g: 8.5, FatPer100g: 14.7, FiberPer100g: 6.7, SugarPer100g: 0.7, SodiumPer100g: 7},
		{Name: "Almonds", Brand: "Generic", Category: "Nuts", ServingSize: 28, ServingUnit: "g", CaloriesPer100g: 579, ProteinPer100g: 21.2, CarbsPer100g: 21.6, FatPer100g: 49.9, FiberPer100g: 12.5, SugarPer100g: 4.4, SodiumPer100g: 1},
		{Name: "Olive Oil", Brand: "Generic", Category: "Oils", ServingSize: 14, ServingUnit: "g", CaloriesPer100g: 884, ProteinPer100g: 0, CarbsPer100g: 0, FatPer100g: 100, FiberPer100g: 0, SugarPer100g: 0, SodiumPer100g: 2},

		// Vegetables
		{Name: "Broccoli", Brand: "Generic", Category: "Vegetables", ServingSize: 91, ServingUnit: "g", CaloriesPer100g: 34, ProteinPer100g: 2.8, CarbsPer100g: 6.6, FatPer100g: 0.4, FiberPer100g: 2.6, SugarPer100g: 1.5, SodiumPer100g: 33},
		{Name: "Spinach", Brand: "Generic", Category: "Vegetables", ServingSize: 30, ServingUnit: "g", CaloriesPer100g: 23, ProteinPer100g: 2.9, CarbsPer100g: 3.6, FatPer100g: 0.4, FiberPer100g: 2.2, SugarPer100g: 0.4, SodiumPer100g: 79},
	}

	for _, food := range foods {
		_, err := s.db.Exec("INSERT INTO food_database (name, brand, category, serving_size, serving_unit, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, sugar_per_100g, sodium_per_100g, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			food.Name, food.Brand, food.Category, food.ServingSize, food.ServingUnit, food.CaloriesPer100g, food.ProteinPer100g, food.CarbsPer100g, food.FatPer100g, food.FiberPer100g, food.SugarPer100g, food.SodiumPer100g, time.Now())
		if err != nil {
			log.Println("Failed to seed food:", food.Name, err)
		}
	}
}

func (s *Server) createDefaultAdmin() {
	// Check if admin user already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil || count > 0 {
		return // Admin already exists or error occurred
	}
	
	// Create default admin user
	passwordHash, err := s.hashPassword("Admin123!")
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}
	
	_, err = s.db.Exec(`
		INSERT INTO users (username, email, password_hash, role, provider, first_name, last_name, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"admin", "admin@workouttracker.com", passwordHash, "admin", "local", "Admin", "User", true, time.Now(), time.Now())
	if err != nil {
		log.Printf("Failed to create default admin user: %v", err)
		return
	}
	
	log.Println("✅ Default admin user created:")
	log.Println("   Username: admin")
	log.Println("   Password: Admin123!")
	log.Println("   Please change this password after first login!")
}

// Nutrition tracking handlers
func (s *Server) NutritionTracker(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	today := time.Now().Format("2006-01-02")

	if r.Method == "POST" {
		// Handle form submission for adding nutrition entry
		r.ParseForm()
		
		// Parse form data
		date := r.FormValue("date")
		if date == "" {
			date = today
		}
		
		mealType := r.FormValue("meal_type")
		foodName := r.FormValue("food_name")
		quantityStr := r.FormValue("quantity")
		unit := r.FormValue("unit")
		caloriesStr := r.FormValue("calories")
		proteinStr := r.FormValue("protein")
		carbsStr := r.FormValue("carbs")
		fatStr := r.FormValue("fat")
		notes := r.FormValue("notes")

		// Convert strings to appropriate types
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		calories, _ := strconv.Atoi(caloriesStr)
		protein, _ := strconv.ParseFloat(proteinStr, 64)
		carbs, _ := strconv.ParseFloat(carbsStr, 64)
		fat, _ := strconv.ParseFloat(fatStr, 64)

		// Insert nutrition entry
		_, err := s.db.Exec(`
			INSERT INTO nutrition_entries (user_id, date, meal_type, food_name, quantity, unit, calories, protein, carbs, fat, notes, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, date, mealType, foodName, quantity, unit, calories, protein, carbs, fat, notes, time.Now())
		
		if err != nil {
			log.Printf("Error adding nutrition entry: %v", err)
			http.Error(w, "Error adding nutrition entry", http.StatusInternalServerError)
			return
		}

		// Update daily nutrition totals
		err = s.updateDailyNutritionTotals(userID, date)
		if err != nil {
			log.Printf("Error updating daily totals: %v", err)
		}

		// Redirect to prevent re-submission
		http.Redirect(w, r, "/nutrition", http.StatusSeeOther)
		return
	}

	// Get nutrition entries for today
	entries, err := s.getNutritionEntries(userID, today)
	if err != nil {
		log.Printf("Error fetching nutrition entries: %v", err)
	}

	// Get daily nutrition summary
	dailySummary, err := s.getDailyNutrition(userID, today)
	if err != nil {
		log.Printf("Error fetching daily nutrition: %v", err)
	}

	// Get food database items for autocomplete
	foodItems, err := s.getFoodDatabaseItems()
	if err != nil {
		log.Printf("Error fetching food items: %v", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title        string
		Username     string
		UserRole     string
		Entries      []NutritionEntry
		DailySummary DailyNutrition
		FoodItems    []FoodDatabase
		Today        string
	}{
		Title:        "Nutrition Tracker",
		Username:     username,
		UserRole:     userRole,
		Entries:      entries,
		DailySummary: dailySummary,
		FoodItems:    foodItems,
		Today:        today,
	}

	err = s.templates.ExecuteTemplate(w, "nutrition.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) DeleteNutritionEntry(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	vars := mux.Vars(r)
	entryID := vars["id"]

	// Get the entry date before deleting for updating daily totals
	var date string
	err := s.db.QueryRow("SELECT date FROM nutrition_entries WHERE id = ? AND user_id = ?", entryID, userID).Scan(&date)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Delete the entry
	_, err = s.db.Exec("DELETE FROM nutrition_entries WHERE id = ? AND user_id = ?", entryID, userID)
	if err != nil {
		log.Printf("Error deleting nutrition entry: %v", err)
		http.Error(w, "Error deleting entry", http.StatusInternalServerError)
		return
	}

	// Update daily nutrition totals
	err = s.updateDailyNutritionTotals(userID, date)
	if err != nil {
		log.Printf("Error updating daily totals: %v", err)
	}

	http.Redirect(w, r, "/nutrition", http.StatusSeeOther)
}

// Helper functions for nutrition tracking
func (s *Server) getNutritionEntries(userID int, date string) ([]NutritionEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, date, meal_type, food_name, quantity, unit, calories, protein, carbs, fat, notes, created_at 
		FROM nutrition_entries 
		WHERE user_id = ? AND date = ? 
		ORDER BY created_at DESC`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []NutritionEntry
	for rows.Next() {
		var entry NutritionEntry
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.MealType, &entry.FoodName, 
			&entry.Quantity, &entry.Unit, &entry.Calories, &entry.Protein, &entry.Carbs, 
			&entry.Fat, &entry.Notes, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Server) getDailyNutrition(userID int, date string) (DailyNutrition, error) {
	var daily DailyNutrition
	err := s.db.QueryRow(`
		SELECT id, user_id, date, total_calories, total_protein, total_carbs, total_fat, 
		       calorie_goal, protein_goal, carbs_goal, fat_goal, water_intake, weight_kg, notes, created_at, updated_at
		FROM daily_nutrition 
		WHERE user_id = ? AND date = ?`, userID, date).Scan(
		&daily.ID, &daily.UserID, &daily.Date, &daily.TotalCalories, &daily.TotalProtein, 
		&daily.TotalCarbs, &daily.TotalFat, &daily.CalorieGoal, &daily.ProteinGoal, 
		&daily.CarbsGoal, &daily.FatGoal, &daily.WaterIntake, &daily.WeightKg, &daily.Notes, 
		&daily.CreatedAt, &daily.UpdatedAt)
	
	if err == sql.ErrNoRows {
		// Create a new daily nutrition entry with default goals
		daily = DailyNutrition{
			UserID:      userID,
			Date:        parseDate(date),
			CalorieGoal: 2000,
			ProteinGoal: 150,
			CarbsGoal:   250,
			FatGoal:     67,
		}
		return daily, nil
	} else if err != nil {
		return daily, err
	}
	
	return daily, nil
}

func (s *Server) getFoodDatabaseItems() ([]FoodDatabase, error) {
	rows, err := s.db.Query(`
		SELECT id, name, brand, category, serving_size, serving_unit, calories_per_100g, 
		       protein_per_100g, carbs_per_100g, fat_per_100g, fiber_per_100g, sugar_per_100g, 
		       sodium_per_100g, created_at 
		FROM food_database 
		ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []FoodDatabase
	for rows.Next() {
		var food FoodDatabase
		err := rows.Scan(&food.ID, &food.Name, &food.Brand, &food.Category, 
			&food.ServingSize, &food.ServingUnit, &food.CaloriesPer100g, 
			&food.ProteinPer100g, &food.CarbsPer100g, &food.FatPer100g, 
			&food.FiberPer100g, &food.SugarPer100g, &food.SodiumPer100g, &food.CreatedAt)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (s *Server) updateDailyNutritionTotals(userID int, date string) error {
	// Calculate totals from nutrition entries for the day
	var totalCalories int
	var totalProtein, totalCarbs, totalFat float64
	
	err := s.db.QueryRow(`
		SELECT COALESCE(SUM(calories), 0), COALESCE(SUM(protein), 0), 
		       COALESCE(SUM(carbs), 0), COALESCE(SUM(fat), 0)
		FROM nutrition_entries 
		WHERE user_id = ? AND date = ?`, userID, date).Scan(
		&totalCalories, &totalProtein, &totalCarbs, &totalFat)
	
	if err != nil {
		return err
	}

	// Insert or update daily nutrition record
	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO daily_nutrition 
		(user_id, date, total_calories, total_protein, total_carbs, total_fat, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, date, totalCalories, totalProtein, totalCarbs, totalFat, time.Now())

	return err
}

func parseDate(dateStr string) time.Time {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return parsed
}

// Analytics Dashboard
func (s *Server) AnalyticsDashboard(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)

	// Get comprehensive analytics data
	volumeData, err := s.getVolumeAnalytics(userID)
	if err != nil {
		log.Printf("Error fetching volume analytics: %v", err)
	}

	strengthProgress, err := s.getStrengthProgress(userID)
	if err != nil {
		log.Printf("Error fetching strength progress: %v", err)
	}

	bodyCompositionData, err := s.getBodyCompositionTrend(userID)
	if err != nil {
		log.Printf("Error fetching body composition data: %v", err)
	}

	measurementData, err := s.getBodyMeasurementTrend(userID)
	if err != nil {
		log.Printf("Error fetching measurement data: %v", err)
	}

	workoutFrequency, err := s.getWorkoutFrequencyAnalysis(userID)
	if err != nil {
		log.Printf("Error fetching workout frequency: %v", err)
	}

	muscleGroupBalance, err := s.getMuscleGroupBalance(userID)
	if err != nil {
		log.Printf("Error fetching muscle group balance: %v", err)
	}

	performanceTrends, err := s.getPerformanceTrends(userID)
	if err != nil {
		log.Printf("Error fetching performance trends: %v", err)
	}

	insights, err := s.generateInsights(userID)
	if err != nil {
		log.Printf("Error generating insights: %v", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title               string
		Username            string
		UserRole            string
		VolumeData          interface{}
		StrengthProgress    interface{}
		BodyCompositionData interface{}
		MeasurementData     interface{}
		WorkoutFrequency    interface{}
		MuscleGroupBalance  interface{}
		PerformanceTrends   interface{}
		Insights            interface{}
	}{
		Title:               "Advanced Analytics",
		Username:            username,
		UserRole:            userRole,
		VolumeData:          volumeData,
		StrengthProgress:    strengthProgress,
		BodyCompositionData: bodyCompositionData,
		MeasurementData:     measurementData,
		WorkoutFrequency:    workoutFrequency,
		MuscleGroupBalance:  muscleGroupBalance,
		PerformanceTrends:   performanceTrends,
		Insights:            insights,
	}

	err = s.templates.ExecuteTemplate(w, "analytics.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// Analytics helper functions
func (s *Server) getVolumeAnalytics(userID int) (interface{}, error) {
	// Get volume data for the last 12 weeks
	rows, err := s.db.Query(`
		SELECT 
			DATE(date, 'weekday 0', '-6 days') as week_start,
			SUM(ws.weight * ws.reps) as weekly_volume,
			COUNT(DISTINCT w.id) as workout_count
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ? 
		AND w.date >= date('now', '-12 weeks')
		GROUP BY week_start
		ORDER BY week_start DESC`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type WeeklyVolume struct {
		WeekStart     string  `json:"week_start"`
		Volume        float64 `json:"volume"`
		WorkoutCount  int     `json:"workout_count"`
	}

	var data []WeeklyVolume
	for rows.Next() {
		var item WeeklyVolume
		err := rows.Scan(&item.WeekStart, &item.Volume, &item.WorkoutCount)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (s *Server) getStrengthProgress(userID int) (interface{}, error) {
	// Get strength progress for major lifts
	rows, err := s.db.Query(`
		SELECT 
			e.name,
			MAX(ws.weight * ws.reps * 1.0278 - ws.reps * 0.0278) as estimated_1rm,
			MAX(ws.weight) as max_weight,
			w.date
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ?
		AND e.name IN ('Bench Press', 'Squat', 'Deadlift', 'Overhead Press')
		AND w.date >= date('now', '-6 months')
		GROUP BY e.name, DATE(w.date, 'start of week')
		ORDER BY e.name, w.date`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type StrengthData struct {
		Exercise     string  `json:"exercise"`
		Estimated1RM float64 `json:"estimated_1rm"`
		MaxWeight    float64 `json:"max_weight"`
		Date         string  `json:"date"`
	}

	var data []StrengthData
	for rows.Next() {
		var item StrengthData
		err := rows.Scan(&item.Exercise, &item.Estimated1RM, &item.MaxWeight, &item.Date)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (s *Server) getBodyCompositionTrend(userID int) (interface{}, error) {
	rows, err := s.db.Query(`
		SELECT date, weight_kg, body_fat_percent, muscle_mass_kg, bmi
		FROM body_composition 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT 52`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type BodyData struct {
		Date           string  `json:"date"`
		Weight         float64 `json:"weight"`
		BodyFat        float64 `json:"body_fat"`
		MuscleMass     float64 `json:"muscle_mass"`
		BMI            float64 `json:"bmi"`
	}

	var data []BodyData
	for rows.Next() {
		var item BodyData
		err := rows.Scan(&item.Date, &item.Weight, &item.BodyFat, &item.MuscleMass, &item.BMI)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (s *Server) getBodyMeasurementTrend(userID int) (interface{}, error) {
	rows, err := s.db.Query(`
		SELECT date, chest_cm, waist_cm, bicep_cm, thigh_cm
		FROM body_measurements 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT 26`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type MeasurementData struct {
		Date   string  `json:"date"`
		Chest  float64 `json:"chest"`
		Waist  float64 `json:"waist"`
		Bicep  float64 `json:"bicep"`
		Thigh  float64 `json:"thigh"`
	}

	var data []MeasurementData
	for rows.Next() {
		var item MeasurementData
		err := rows.Scan(&item.Date, &item.Chest, &item.Waist, &item.Bicep, &item.Thigh)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (s *Server) getWorkoutFrequencyAnalysis(userID int) (interface{}, error) {
	rows, err := s.db.Query(`
		SELECT 
			strftime('%w', date) as day_of_week,
			COUNT(*) as workout_count
		FROM workouts 
		WHERE user_id = ? 
		AND date >= date('now', '-3 months')
		GROUP BY day_of_week
		ORDER BY day_of_week`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type FrequencyData struct {
		DayOfWeek    string `json:"day"`
		WorkoutCount int    `json:"count"`
	}

	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	var data []FrequencyData
	for rows.Next() {
		var dayNum int
		var count int
		err := rows.Scan(&dayNum, &count)
		if err != nil {
			return nil, err
		}
		data = append(data, FrequencyData{
			DayOfWeek:    dayNames[dayNum],
			WorkoutCount: count,
		})
	}

	return data, nil
}

func (s *Server) getMuscleGroupBalance(userID int) (interface{}, error) {
	// Analyze muscle group distribution over recent workouts
	rows, err := s.db.Query(`
		SELECT 
			CASE 
				WHEN e.name LIKE '%bench%' OR e.name LIKE '%press%' OR e.name LIKE '%fly%' OR e.name LIKE '%dip%' THEN 'Chest'
				WHEN e.name LIKE '%row%' OR e.name LIKE '%pull%' OR e.name LIKE '%lat%' OR e.name LIKE '%deadlift%' THEN 'Back'
				WHEN e.name LIKE '%squat%' OR e.name LIKE '%leg press%' THEN 'Legs'
				WHEN e.name LIKE '%curl%' THEN 'Biceps'
				WHEN e.name LIKE '%extension%' OR e.name LIKE '%pushdown%' THEN 'Triceps'
				WHEN e.name LIKE '%shoulder%' OR e.name LIKE '%lateral%' OR e.name LIKE '%overhead%' THEN 'Shoulders'
				ELSE 'Other'
			END as muscle_group,
			SUM(ws.weight * ws.reps) as total_volume,
			COUNT(*) as set_count
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ?
		AND w.date >= date('now', '-1 month')
		GROUP BY muscle_group
		ORDER BY total_volume DESC`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type MuscleGroupData struct {
		MuscleGroup  string  `json:"muscle_group"`
		TotalVolume  float64 `json:"total_volume"`
		SetCount     int     `json:"set_count"`
	}

	var data []MuscleGroupData
	for rows.Next() {
		var item MuscleGroupData
		err := rows.Scan(&item.MuscleGroup, &item.TotalVolume, &item.SetCount)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (s *Server) getPerformanceTrends(userID int) (interface{}, error) {
	// Calculate performance trends and metrics
	row := s.db.QueryRow(`
		SELECT 
			COUNT(DISTINCT w.id) as total_workouts,
			AVG(ws.rpe) as avg_rpe,
			SUM(ws.weight * ws.reps) as total_volume,
			COUNT(ws.id) as total_sets
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ?
		AND w.date >= date('now', '-1 month')`, userID)

	type PerformanceMetrics struct {
		TotalWorkouts int     `json:"total_workouts"`
		AvgRPE        float64 `json:"avg_rpe"`
		TotalVolume   float64 `json:"total_volume"`
		TotalSets     int     `json:"total_sets"`
		VolumePerWorkout float64 `json:"volume_per_workout"`
		SetsPerWorkout   float64 `json:"sets_per_workout"`
	}

	var metrics PerformanceMetrics
	err := row.Scan(&metrics.TotalWorkouts, &metrics.AvgRPE, &metrics.TotalVolume, &metrics.TotalSets)
	if err != nil {
		return nil, err
	}

	if metrics.TotalWorkouts > 0 {
		metrics.VolumePerWorkout = metrics.TotalVolume / float64(metrics.TotalWorkouts)
		metrics.SetsPerWorkout = float64(metrics.TotalSets) / float64(metrics.TotalWorkouts)
	}

	return metrics, nil
}

func (s *Server) generateInsights(userID int) (interface{}, error) {
	// Generate AI-like insights based on user data
	type Insight struct {
		Type        string `json:"type"` // "success", "warning", "info", "tip"
		Title       string `json:"title"`
		Description string `json:"description"`
		Action      string `json:"action"`
	}

	var insights []Insight

	// Check workout consistency
	var workoutsThisWeek int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM workouts 
		WHERE user_id = ? 
		AND date >= date('now', 'weekday 0', '-6 days')`, userID).Scan(&workoutsThisWeek)
	
	if err == nil {
		if workoutsThisWeek >= 4 {
			insights = append(insights, Insight{
				Type:        "success",
				Title:       "Great Consistency!",
				Description: fmt.Sprintf("You've completed %d workouts this week. Keep up the excellent work!", workoutsThisWeek),
				Action:      "Continue your current routine",
			})
		} else if workoutsThisWeek == 0 {
			insights = append(insights, Insight{
				Type:        "warning",
				Title:       "Time to Get Moving",
				Description: "You haven't logged any workouts this week. Consistency is key to reaching your fitness goals.",
				Action:      "Schedule a workout for today",
			})
		} else {
			insights = append(insights, Insight{
				Type:        "info",
				Title:       "Room for Improvement",
				Description: fmt.Sprintf("You've completed %d workouts this week. Try to aim for 3-4 sessions for optimal results.", workoutsThisWeek),
				Action:      "Plan your remaining workouts",
			})
		}
	}

	// Check for volume progression
	var thisWeekVolume, lastWeekVolume float64
	s.db.QueryRow(`
		SELECT COALESCE(SUM(ws.weight * ws.reps), 0)
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ?
		AND w.date >= date('now', 'weekday 0', '-6 days')`, userID).Scan(&thisWeekVolume)
	
	s.db.QueryRow(`
		SELECT COALESCE(SUM(ws.weight * ws.reps), 0)
		FROM workouts w
		JOIN exercises e ON w.id = e.workout_id
		JOIN workout_sets ws ON e.id = ws.exercise_id
		WHERE w.user_id = ?
		AND w.date >= date('now', 'weekday 0', '-13 days')
		AND w.date < date('now', 'weekday 0', '-6 days')`, userID).Scan(&lastWeekVolume)

	if thisWeekVolume > 0 && lastWeekVolume > 0 {
		volumeIncrease := ((thisWeekVolume - lastWeekVolume) / lastWeekVolume) * 100
		if volumeIncrease > 10 {
			insights = append(insights, Insight{
				Type:        "success",
				Title:       "Volume Progression Detected!",
				Description: fmt.Sprintf("Your training volume increased by %.1f%% this week. Excellent progressive overload!", volumeIncrease),
				Action:      "Keep pushing those limits safely",
			})
		} else if volumeIncrease < -10 {
			insights = append(insights, Insight{
				Type:        "info",
				Title:       "Deload Week?",
				Description: fmt.Sprintf("Your volume decreased by %.1f%% this week. This could be a planned deload or recovery period.", -volumeIncrease),
				Action:      "Consider if this aligns with your program",
			})
		}
	}

	// Add a fitness tip
	insights = append(insights, Insight{
		Type:        "tip",
		Title:       "Pro Tip: Recovery Matters",
		Description: "Don't forget that muscles grow during rest, not just during workouts. Ensure you're getting adequate sleep and nutrition.",
		Action:      "Track your sleep and nutrition",
	})

	return insights, nil
}

// Body Composition Tracking Handlers
func (s *Server) BodyCompositionTracker(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)

	if r.Method == "POST" {
		// Handle form submission for adding body composition entry
		r.ParseForm()
		
		date := r.FormValue("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		
		weightStr := r.FormValue("weight_kg")
		bodyFatStr := r.FormValue("body_fat_percent")
		muscleMassStr := r.FormValue("muscle_mass_kg")
		bodyWaterStr := r.FormValue("body_water_percent")
		boneMassStr := r.FormValue("bone_mass_kg")
		notes := r.FormValue("notes")

		// Convert strings to appropriate types
		weight, _ := strconv.ParseFloat(weightStr, 64)
		bodyFat, _ := strconv.ParseFloat(bodyFatStr, 64)
		muscleMass, _ := strconv.ParseFloat(muscleMassStr, 64)
		bodyWater, _ := strconv.ParseFloat(bodyWaterStr, 64)
		boneMass, _ := strconv.ParseFloat(boneMassStr, 64)

		// Calculate BMI if height is available (for now, use a default or get from user profile)
		bmi := 0.0
		if weight > 0 {
			// Assuming average height for BMI calculation - this should come from user profile
			heightM := 1.75 // Default height in meters
			bmi = weight / (heightM * heightM)
		}

		// Insert body composition entry
		_, err := s.db.Exec(`
			INSERT INTO body_composition (user_id, date, weight_kg, body_fat_percent, muscle_mass_kg, body_water_percent, bone_mass_kg, bmi, notes, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, date, weight, bodyFat, muscleMass, bodyWater, boneMass, bmi, notes, time.Now())
		
		if err != nil {
			log.Printf("Error adding body composition entry: %v", err)
			http.Error(w, "Error adding body composition entry", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/body-composition", http.StatusSeeOther)
		return
	}

	// Get recent body composition entries
	entries, err := s.getBodyCompositionEntries(userID, 30) // Last 30 entries
	if err != nil {
		log.Printf("Error fetching body composition entries: %v", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title    string
		Username string
		UserRole string
		Entries  []BodyComposition
		Today    string
	}{
		Title:    "Body Composition Tracking",
		Username: username,
		UserRole: userRole,
		Entries:  entries,
		Today:    time.Now().Format("2006-01-02"),
	}

	err = s.templates.ExecuteTemplate(w, "body-composition.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) BodyMeasurementTracker(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)

	if r.Method == "POST" {
		// Handle form submission for adding body measurements
		r.ParseForm()
		
		date := r.FormValue("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		
		neckStr := r.FormValue("neck_cm")
		chestStr := r.FormValue("chest_cm")
		waistStr := r.FormValue("waist_cm")
		hipsStr := r.FormValue("hips_cm")
		bicepStr := r.FormValue("bicep_cm")
		forearmStr := r.FormValue("forearm_cm")
		thighStr := r.FormValue("thigh_cm")
		calfStr := r.FormValue("calf_cm")
		shoulderStr := r.FormValue("shoulders_cm")
		notes := r.FormValue("notes")

		// Convert strings to appropriate types
		neck, _ := strconv.ParseFloat(neckStr, 64)
		chest, _ := strconv.ParseFloat(chestStr, 64)
		waist, _ := strconv.ParseFloat(waistStr, 64)
		hips, _ := strconv.ParseFloat(hipsStr, 64)
		bicep, _ := strconv.ParseFloat(bicepStr, 64)
		forearm, _ := strconv.ParseFloat(forearmStr, 64)
		thigh, _ := strconv.ParseFloat(thighStr, 64)
		calf, _ := strconv.ParseFloat(calfStr, 64)
		shoulder, _ := strconv.ParseFloat(shoulderStr, 64)

		// Insert body measurements entry
		_, err := s.db.Exec(`
			INSERT INTO body_measurements (user_id, date, neck_cm, chest_cm, waist_cm, hips_cm, bicep_cm, forearm_cm, thigh_cm, calf_cm, shoulders_cm, notes, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, date, neck, chest, waist, hips, bicep, forearm, thigh, calf, shoulder, notes, time.Now())
		
		if err != nil {
			log.Printf("Error adding body measurements entry: %v", err)
			http.Error(w, "Error adding body measurements entry", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/body-measurements", http.StatusSeeOther)
		return
	}

	// Get recent body measurement entries
	entries, err := s.getBodyMeasurementEntries(userID, 30) // Last 30 entries
	if err != nil {
		log.Printf("Error fetching body measurement entries: %v", err)
	}

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title    string
		Username string
		UserRole string
		Entries  []BodyMeasurements
		Today    string
	}{
		Title:    "Body Measurements Tracking",
		Username: username,
		UserRole: userRole,
		Entries:  entries,
		Today:    time.Now().Format("2006-01-02"),
	}

	err = s.templates.ExecuteTemplate(w, "body-measurements.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) DeleteBodyComposition(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	vars := mux.Vars(r)
	entryID := vars["id"]

	_, err := s.db.Exec("DELETE FROM body_composition WHERE id = ? AND user_id = ?", entryID, userID)
	if err != nil {
		log.Printf("Error deleting body composition entry: %v", err)
		http.Error(w, "Error deleting entry", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/body-composition", http.StatusSeeOther)
}

func (s *Server) DeleteBodyMeasurement(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	vars := mux.Vars(r)
	entryID := vars["id"]

	_, err := s.db.Exec("DELETE FROM body_measurements WHERE id = ? AND user_id = ?", entryID, userID)
	if err != nil {
		log.Printf("Error deleting body measurement entry: %v", err)
			http.Error(w, "Error deleting entry", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/body-measurements", http.StatusSeeOther)
}

// Helper functions for body composition and measurements
func (s *Server) getBodyCompositionEntries(userID int, limit int) ([]BodyComposition, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, date, weight_kg, body_fat_percent, muscle_mass_kg, body_water_percent, bone_mass_kg, bmi, notes, created_at 
		FROM body_composition 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []BodyComposition
	for rows.Next() {
		var entry BodyComposition
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.WeightKg, 
			&entry.BodyFatPercent, &entry.MuscleMassKg, &entry.BodyWaterPercent, 
			&entry.BoneMassKg, &entry.BMI, &entry.Notes, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Server) getBodyMeasurementEntries(userID int, limit int) ([]BodyMeasurements, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, date, neck_cm, chest_cm, waist_cm, hips_cm, bicep_cm, forearm_cm, thigh_cm, calf_cm, shoulders_cm, notes, created_at 
		FROM body_measurements 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []BodyMeasurements
	for rows.Next() {
		var entry BodyMeasurements
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.NeckCm, 
			&entry.ChestCm, &entry.WaistCm, &entry.HipsCm, &entry.BicepCm, 
			&entry.ForearmCm, &entry.ThighCm, &entry.CalfCm, &entry.ShouldersCm, 
			&entry.Notes, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Goal Management Handlers
func (s *Server) GoalManager(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)

	if r.Method == "POST" {
		// Handle goal creation
		r.ParseForm()
		
		goalType := r.FormValue("goal_type")
		goalCategory := r.FormValue("goal_category")
		title := r.FormValue("title")
		description := r.FormValue("description")
		currentValueStr := r.FormValue("current_value")
		targetValueStr := r.FormValue("target_value")
		unit := r.FormValue("unit")
		targetDateStr := r.FormValue("target_date")
		priorityStr := r.FormValue("priority")

		// Parse values
		currentValue, _ := strconv.ParseFloat(currentValueStr, 64)
		targetValue, _ := strconv.ParseFloat(targetValueStr, 64)
		priority, _ := strconv.Atoi(priorityStr)
		if priority == 0 {
			priority = 2 // Default to medium priority
		}
		
		var targetDate time.Time
		if targetDateStr != "" {
			targetDate, _ = time.Parse("2006-01-02", targetDateStr)
		} else {
			// Default to 3 months from now if no date specified
			targetDate = time.Now().AddDate(0, 3, 0)
		}

		// Insert goal
		_, err := s.db.Exec(`
			INSERT INTO goals (user_id, goal_type, goal_category, title, description, current_value, target_value, unit, target_date, priority, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, goalType, goalCategory, title, description, currentValue, targetValue, unit, targetDate, priority, time.Now(), time.Now())
		
		if err != nil {
			log.Printf("Error adding goal: %v", err)
			http.Error(w, "Error adding goal", http.StatusInternalServerError)
			return
		}

		// Add initial progress entry
		goalID := s.getLastInsertID()
		_, err = s.db.Exec(`
			INSERT INTO goal_progress (goal_id, current_value, progress_date, note, created_at) 
			VALUES (?, ?, ?, ?, ?)`,
			goalID, currentValue, time.Now().Format("2006-01-02"), "Initial value", time.Now())

		http.Redirect(w, r, "/goals", http.StatusSeeOther)
		return
	}

	// Get user's goals
	goals, err := s.getUserGoals(userID)
	if err != nil {
		log.Printf("Error fetching goals: %v", err)
	}

	// Get latest body composition data for smart goal suggestions
	latestBodyComp, _ := s.getLatestBodyComposition(userID)
	latestMeasurements, _ := s.getLatestBodyMeasurements(userID)

	session, _ := s.store.Get(r, "workout-session")
	username, _ := session.Values["username"].(string)
	userRole, _ := session.Values["role"].(string)

	data := struct {
		Title               string
		Username            string
		UserRole            string
		Goals               []Goal
		LatestBodyComp      *BodyComposition
		LatestMeasurements  *BodyMeasurements
		Today               string
	}{
		Title:               "Goal Setting",
		Username:            username,
		UserRole:            userRole,
		Goals:               goals,
		LatestBodyComp:      latestBodyComp,
		LatestMeasurements:  latestMeasurements,
		Today:               time.Now().Format("2006-01-02"),
	}

	err = s.templates.ExecuteTemplate(w, "goals.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Delete goal progress first
	_, err := s.db.Exec("DELETE FROM goal_progress WHERE goal_id = ?", goalID)
	if err != nil {
		log.Printf("Error deleting goal progress: %v", err)
	}

	// Delete the goal
	_, err = s.db.Exec("DELETE FROM goals WHERE id = ? AND user_id = ?", goalID, userID)
	if err != nil {
		log.Printf("Error deleting goal: %v", err)
		http.Error(w, "Error deleting goal", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/goals", http.StatusSeeOther)
}

func (s *Server) CompleteGoal(w http.ResponseWriter, r *http.Request) {
	userID := s.getCurrentUserID(r)
	vars := mux.Vars(r)
	goalID := vars["id"]

	// Mark goal as completed
	_, err := s.db.Exec("UPDATE goals SET is_active = FALSE, completed_at = ? WHERE id = ? AND user_id = ?", time.Now(), goalID, userID)
	if err != nil {
		log.Printf("Error completing goal: %v", err)
		http.Error(w, "Error completing goal", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/goals", http.StatusSeeOther)
}

// Goal helper functions
func (s *Server) getUserGoals(userID int) ([]Goal, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, goal_type, goal_category, title, description, current_value, target_value, unit, target_date, is_active, priority, created_at, updated_at, completed_at
		FROM goals 
		WHERE user_id = ? 
		ORDER BY is_active DESC, priority ASC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		var completedAt sql.NullTime
		err := rows.Scan(&goal.ID, &goal.UserID, &goal.GoalType, &goal.GoalCategory, &goal.Title, &goal.Description,
			&goal.CurrentValue, &goal.TargetValue, &goal.Unit, &goal.TargetDate, &goal.IsActive, &goal.Priority,
			&goal.CreatedAt, &goal.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			goal.CompletedAt = &completedAt.Time
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func (s *Server) getLatestBodyComposition(userID int) (*BodyComposition, error) {
	var comp BodyComposition
	err := s.db.QueryRow(`
		SELECT id, user_id, date, weight_kg, body_fat_percent, muscle_mass_kg, body_water_percent, bone_mass_kg, bmi, notes, created_at
		FROM body_composition 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT 1`, userID).Scan(
		&comp.ID, &comp.UserID, &comp.Date, &comp.WeightKg, &comp.BodyFatPercent, &comp.MuscleMassKg,
		&comp.BodyWaterPercent, &comp.BoneMassKg, &comp.BMI, &comp.Notes, &comp.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	
	return &comp, nil
}

func (s *Server) getLatestBodyMeasurements(userID int) (*BodyMeasurements, error) {
	var meas BodyMeasurements
	err := s.db.QueryRow(`
		SELECT id, user_id, date, neck_cm, chest_cm, waist_cm, hips_cm, bicep_cm, forearm_cm, thigh_cm, calf_cm, shoulders_cm, notes, created_at
		FROM body_measurements 
		WHERE user_id = ? 
		ORDER BY date DESC 
		LIMIT 1`, userID).Scan(
		&meas.ID, &meas.UserID, &meas.Date, &meas.NeckCm, &meas.ChestCm, &meas.WaistCm, &meas.HipsCm,
		&meas.BicepCm, &meas.ForearmCm, &meas.ThighCm, &meas.CalfCm, &meas.ShouldersCm, &meas.Notes, &meas.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	
	return &meas, nil
}

func (s *Server) getLastInsertID() int {
	var id int
	err := s.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0
	}
	return id
}

func main() {
	server := NewServer()
	defer server.db.Close()

	r := mux.NewRouter()
	
	// Public routes
	r.HandleFunc("/login", server.Login).Methods("GET", "POST")
	r.HandleFunc("/register", server.Register).Methods("GET", "POST")
	
	// OAuth routes
	r.HandleFunc("/auth/google", server.GoogleLogin).Methods("GET")
	r.HandleFunc("/auth/google/callback", server.GoogleCallback).Methods("GET")
	
	// Protected routes
	r.HandleFunc("/", server.requireAuth(server.Home)).Methods("GET")
	r.HandleFunc("/create", server.requireAuth(server.CreateWorkout)).Methods("GET", "POST")
	r.HandleFunc("/edit/{id}", server.requireAuth(server.EditWorkout)).Methods("GET", "POST")
	r.HandleFunc("/delete/{id}", server.requireAuth(server.DeleteWorkout)).Methods("POST")
	r.HandleFunc("/progress", server.requireAuth(server.ProgressTracker)).Methods("GET")
	r.HandleFunc("/schedule", server.requireAuth(server.Schedule)).Methods("GET", "POST")
	r.HandleFunc("/workout/{id}/log", server.requireAuth(server.LogWorkout)).Methods("GET", "POST")
	r.HandleFunc("/export", server.requireAuth(server.ExportData)).Methods("GET")
	r.HandleFunc("/import", server.requireAuth(server.ImportData)).Methods("GET", "POST")
	r.HandleFunc("/nutrition", server.requireAuth(server.NutritionTracker)).Methods("GET", "POST")
	r.HandleFunc("/nutrition/entry/{id}/delete", server.requireAuth(server.DeleteNutritionEntry)).Methods("POST")
	r.HandleFunc("/analytics", server.requireAuth(server.AnalyticsDashboard)).Methods("GET")
	r.HandleFunc("/body-composition", server.requireAuth(server.BodyCompositionTracker)).Methods("GET", "POST")
	r.HandleFunc("/body-measurements", server.requireAuth(server.BodyMeasurementTracker)).Methods("GET", "POST")
	r.HandleFunc("/body-composition/{id}/delete", server.requireAuth(server.DeleteBodyComposition)).Methods("POST")
	r.HandleFunc("/body-measurements/{id}/delete", server.requireAuth(server.DeleteBodyMeasurement)).Methods("POST")
	r.HandleFunc("/goals", server.requireAuth(server.GoalManager)).Methods("GET", "POST")
	r.HandleFunc("/goals/{id}/delete", server.requireAuth(server.DeleteGoal)).Methods("POST")
	r.HandleFunc("/goals/{id}/complete", server.requireAuth(server.CompleteGoal)).Methods("POST")
	r.HandleFunc("/logout", server.requireAuth(server.Logout)).Methods("POST")
	
	// Admin routes
	r.HandleFunc("/admin", server.requireAdmin(server.AdminPanel)).Methods("GET")
	r.HandleFunc("/admin/users", server.requireAdmin(server.AdminUserManagement)).Methods("POST")
	
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	log.Println("Simple server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
