package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"context"

	"workout-tracker/internal/database"
	"workout-tracker/internal/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Handler holds the database connection and templates
type Handler struct {
	db        *database.DB
	templates *template.Template
	store     *sessions.CookieStore
}

// New creates a new handler instance
func New(db *database.DB) *Handler {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"hasPrefix": strings.HasPrefix,
		"contains": strings.Contains,
		"eq":       func(a, b interface{}) bool { return a == b },
		"ne":       func(a, b interface{}) bool { return a != b },
	}

	// Load templates with proper parsing for inheritance
	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/*.html"))
	
	// Get session secret from environment variable
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Println("WARNING: SESSION_SECRET not set, using default (insecure for production)")
		sessionSecret = "default-insecure-secret-change-in-production"
	}
	
	store := sessions.NewCookieStore([]byte(sessionSecret))
	// Configure session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	return &Handler{
		db:        db,
		templates: templates,
		store:     store,
	}
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := h.getUserByUsername(username)
		if err != nil || !checkPasswordHash(password, user.PasswordHash) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

	// Get or create session
	session, _ := h.store.Get(r, "session-name")
	log.Printf("LOGIN - Before setting session values: %+v", session.Values)
	
	// Set session values
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Values["created_at"] = time.Now().Unix()
	session.Values["login_time"] = time.Now().Format(time.RFC3339)
	
	// Save session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	
	log.Printf("LOGIN - Session saved for user %d (%s) with values: %+v", user.ID, user.Username, session.Values)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Login",
	}

	err := h.templates.ExecuteTemplate(w, "login_simple.html", data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		passwordHash, err := hashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err = h.createUser(user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Register",
	}

	err := h.templates.ExecuteTemplate(w, "register_simple.html", data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// Home renders the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get recent workouts for the dashboard
	workouts, err := h.getRecentWorkoutsByUser(userID, 5)
	if err != nil {
		http.Error(w, "Failed to load workouts", http.StatusInternalServerError)
		return
	}

	// Get workout statistics
	stats, err := h.getWorkoutStats()
	if err != nil {
		http.Error(w, "Failed to load statistics", http.StatusInternalServerError)
		return
	}

	data := struct {
		Workouts     []models.Workout
		Stats        map[string]interface{}
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Workouts:     workouts,
		Stats:        stats,
		Title:        "Workout Tracker",
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "index_dashboard.html", data); err != nil {
		log.Printf("Template error: %v", err)
	http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GetWorkouts returns all workouts for the current user
func (h *Handler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	workouts, err := h.getAllWorkoutsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to load workouts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Workouts     []models.Workout
		Title        string
		UserSettings *models.UserSettings
	}{
		Workouts:     workouts,
		Title:        "All Workouts",
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "workouts.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// CreateWorkout handles workout creation
func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("CreateWorkout POST request received")
		// Get current user ID from session
		userID, err := h.getCurrentUserID(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		dateStr := r.FormValue("date")
		notes := r.FormValue("notes")
		log.Printf("Form values - Name: %s, Date: %s, Notes: %s", name, dateStr, notes)

		// Parse date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		// Create workout
		workout := models.Workout{
			Name:      name,
			Date:      date,
			Notes:     notes,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		id, err := h.createWorkoutWithUser(workout, userID)
		if err != nil {
			http.Error(w, "Failed to create workout", http.StatusInternalServerError)
			return
		}

		// Redirect to the new workout
		http.Redirect(w, r, "/workouts/"+strconv.Itoa(id), http.StatusSeeOther)
		return
	}

	// GET request - show create form
	data := struct {
		Title        string
		UserSettings *models.UserSettings
	}{
		Title:        "Create Workout",
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "create_workout.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GetWorkout returns a specific workout
func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GetWorkout - Invalid ID: %v", err)
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Get current user ID from session
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	log.Printf("GetWorkout - Fetching workout with ID: %d for user: %d", id, userID)
	workout, err := h.getWorkoutByIDWithUser(id, userID)
	if err != nil {
		log.Printf("GetWorkout - Database error: %v", err)
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	log.Printf("GetWorkout - Retrieved workout: %+v", workout)
	data := struct {
		Workout      models.Workout
		Title        string
		UserSettings *models.UserSettings
	}{
		Workout:      workout,
		Title:        "Workout: " + workout.Name,
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	log.Printf("GetWorkout - About to render template with data: %+v", data)
	// Execute the workout_detail.html template
	err = h.templates.ExecuteTemplate(w, "workout_detail.html", data)
	if err != nil {
		log.Printf("GetWorkout - Template error: %v", err)
		// Don't call http.Error after template execution might have started
		return
	}
	log.Printf("GetWorkout - Template rendered successfully")
}

// UpdateWorkout handles workout updates
func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Get current user ID from session
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" || r.Method == "PUT" {
		var name, dateStr, durationStr, notes string
		
		if r.Method == "PUT" {
			// Handle JSON for API requests
			var req models.UpdateWorkoutRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			name = req.Name
			dateStr = req.Date
			durationStr = strconv.Itoa(req.Duration)
			notes = req.Notes
		} else {
			// Handle form data for web requests
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Check for method override
			if r.FormValue("_method") != "PUT" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			
			name = r.FormValue("name")
			dateStr = r.FormValue("date")
			durationStr = r.FormValue("duration")
			notes = r.FormValue("notes")
		}

		// Parse date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		// Parse duration
		var duration int
		if durationStr != "" {
			duration, err = strconv.Atoi(durationStr)
			if err != nil {
				http.Error(w, "Invalid duration format", http.StatusBadRequest)
				return
			}
		}

		// Create updated workout
		workout := models.Workout{
			ID:        id,
			Name:      name,
			Date:      date,
			Duration:  duration,
			Notes:     notes,
			UpdatedAt: time.Now(),
		}

		err = h.updateWorkoutWithUser(workout, userID)
		if err != nil {
			http.Error(w, "Failed to update workout", http.StatusInternalServerError)
			return
		}

		if r.Method == "PUT" {
			// API request - return JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(workout)
		} else {
			// Web form request - redirect
			http.Redirect(w, r, "/workouts/"+strconv.Itoa(id), http.StatusSeeOther)
		}
		return
	}

	// GET request - show edit form
	workout, err := h.getWorkoutByID(id)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	data := struct {
		Workout      models.Workout
		Title        string
		UserSettings *models.UserSettings
	}{
		Workout:      workout,
		Title:        "Edit Workout: " + workout.Name,
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "edit_workout.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DeleteWorkout handles workout deletion
func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	// Check if this is a POST request with method override
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		if r.FormValue("_method") != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Get current user ID from session
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err = h.deleteWorkoutWithUser(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/workouts", http.StatusSeeOther)
}

// CreateExercise handles exercise creation
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exercise := models.Exercise{
		WorkoutID: req.WorkoutID,
		Name:      req.Name,
		Category:  req.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := h.createExercise(exercise)
	if err != nil {
		http.Error(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	exercise.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exercise)
}

// CreateSet handles set creation
func (h *Handler) CreateSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	set := models.Set{
		ExerciseID: req.ExerciseID,
		SetNumber:  req.SetNumber,
		Reps:       req.Reps,
		Weight:     req.Weight,
		Distance:   req.Distance,
		Duration:   req.Duration,
		RestTime:   req.RestTime,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	id, err := h.createSet(set)
	if err != nil {
		http.Error(w, "Failed to create set", http.StatusInternalServerError)
		return
	}

	set.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(set)
}

// UpdateExercise handles exercise updates
func (h *Handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	if r.Method != "PUT" && (r.Method != "POST" || r.FormValue("_method") != "PUT") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UpdateExerciseRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Handle form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		req.Name = r.FormValue("name")
		req.Category = r.FormValue("category")
	}

	exercise := models.Exercise{
		ID:        id,
		Name:      req.Name,
		Category:  req.Category,
		UpdatedAt: time.Now(),
	}

	err = h.updateExercise(exercise)
	if err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
	} else {
		// Redirect back to workout page
		workoutID := r.FormValue("workout_id")
		http.Redirect(w, r, "/workouts/"+workoutID, http.StatusSeeOther)
	}
}

// DeleteExercise handles exercise deletion
func (h *Handler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	if r.Method != "DELETE" && (r.Method != "POST" || r.FormValue("_method") != "DELETE") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = h.deleteExercise(id)
	if err != nil {
		http.Error(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		// Redirect back to workout page
		workoutID := r.FormValue("workout_id")
		http.Redirect(w, r, "/workouts/"+workoutID, http.StatusSeeOther)
	}
}

// UpdateSet handles set updates
func (h *Handler) UpdateSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	if r.Method != "PUT" && (r.Method != "POST" || r.FormValue("_method") != "PUT") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UpdateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	set := models.Set{
		ID:         id,
		SetNumber:  req.SetNumber,
		Reps:       req.Reps,
		Weight:     req.Weight,
		Distance:   req.Distance,
		Duration:   req.Duration,
		RestTime:   req.RestTime,
		UpdatedAt:  time.Now(),
	}

	err = h.updateSet(set)
	if err != nil {
		http.Error(w, "Failed to update set", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(set)
}

// DeleteSet handles set deletion
func (h *Handler) DeleteSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid set ID", http.StatusBadRequest)
		return
	}

	if r.Method != "DELETE" && (r.Method != "POST" || r.FormValue("_method") != "DELETE") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = h.deleteSet(id)
	if err != nil {
		http.Error(w, "Failed to delete set", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout request - Method: %s, URL: %s", r.Method, r.URL.Path)
	session, err := h.store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}
	log.Printf("Before logout - Session values: %+v", session.Values)
	session.Options.MaxAge = -1 // This deletes the session
	session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
	}
	log.Printf("After logout - Session values: %+v", session.Values)
	log.Printf("Session cleared, redirecting to login")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// ClearSession clears all sessions for debugging
func (h *Handler) ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	session.Options.MaxAge = -1 // This will delete the session
	session.Save(r, w)
	w.Write([]byte("Session cleared. <a href='/'>Go to home</a> to test login."))
}

// Context key for user ID
type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware checks if user is authenticated and adds user ID to context
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug: print all cookies received
		log.Printf("Auth check - URL: %s, Cookies received: %v", r.URL.Path, r.Cookies())
		session, err := h.store.Get(r, "session-name")
		if err != nil {
			log.Printf("Error getting session: %v", err)
		}
		auth, ok := session.Values["authenticated"]
		userID, userIDOk := session.Values["user_id"].(int)
		log.Printf("Auth check - URL: %s, Session exists: %v, Auth value: %v, Type: %T, UserID: %v, Session Values: %+v", r.URL.Path, ok, auth, auth, userID, session.Values)
		if !ok || auth != true || !userIDOk {
			log.Printf("Redirecting to login from %s", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)
		log.Printf("Authentication successful for user %d, proceeding to %s", userID, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Helper function to get user ID from context (preferred for authenticated requests)
func (h *Handler) getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user not authenticated")
	}
	return userID, nil
}

// Helper function to get current user ID from session (fallback)
func (h *Handler) getCurrentUserID(r *http.Request) (int, error) {
	// First try to get from context (set by AuthMiddleware)
	if userID, err := h.getUserIDFromContext(r); err == nil {
		return userID, nil
	}
	
	// Fallback to session
	session, err := h.store.Get(r, "session-name")
	if err != nil {
		return 0, err
	}
	
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0, fmt.Errorf("user not authenticated")
	}
	
	return userID, nil
}

// Helper function to get user settings for templates
func (h *Handler) getUserSettingsForTemplate(r *http.Request) *models.UserSettings {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return nil
	}
	
	settings, err := h.getUserSettings(userID)
	if err != nil {
		// Return default settings if error
		return &models.UserSettings{
			Theme: "light",
			Language: "en",
		}
	}
	return &settings
}

// Helper functions for password hashing
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// APIHandler returns a handler for API routes
func (h *Handler) APIHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/workouts", h.APIGetWorkouts).Methods("GET")
	r.HandleFunc("/api/workouts", h.APICreateWorkout).Methods("POST")
	r.HandleFunc("/api/workouts/{id}", h.APIGetWorkout).Methods("GET")
	return r
}

// APIGetWorkouts returns workouts as JSON
func (h *Handler) APIGetWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.getAllWorkouts()
	if err != nil {
		http.Error(w, "Failed to load workouts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// APICreateWorkout creates a workout via API
func (h *Handler) APICreateWorkout(w http.ResponseWriter, r *http.Request) {
	var req models.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	workout := models.Workout{
		Name:      req.Name,
		Date:      date,
		Duration:  req.Duration,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := h.createWorkout(workout)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	workout.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// APIGetWorkout returns a specific workout as JSON
func (h *Handler) APIGetWorkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.getWorkoutByID(id)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// GetPredefinedExercises returns all predefined exercises
func (h *Handler) GetPredefinedExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.getAllPredefinedExercises()
	if err != nil {
		http.Error(w, "Failed to load predefined exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

// GetPredefinedExercisesByCategory returns predefined exercises by category
func (h *Handler) GetPredefinedExercisesByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	
	exercises, err := h.getPredefinedExercisesByCategory(category)
	if err != nil {
		http.Error(w, "Failed to load predefined exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

// CreatePredefinedExercise creates a new predefined exercise
func (h *Handler) CreatePredefinedExercise(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePredefinedExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exercise := models.PredefinedExercise{
		Name:         req.Name,
		Category:     req.Category,
		Description:  req.Description,
		VideoURL:     req.VideoURL,
		Instructions: req.Instructions,
		Tips:         req.Tips,
		MuscleGroups: req.MuscleGroups,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		ImageURL:     req.ImageURL,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := h.createPredefinedExercise(exercise)
	if err != nil {
		http.Error(w, "Failed to create predefined exercise", http.StatusInternalServerError)
		return
	}

	exercise.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exercise)
}

// CreateMeal handles meal creation
func (h *Handler) CreateMeal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateMealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	meal := models.Meal{
		UserID:    userID,
		Name:      req.Name,
		Calories:  req.Calories,
		Protein:   req.Protein,
		Carbs:     req.Carbs,
		Fat:       req.Fat,
		Date:      date,
		MealType:  req.MealType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := h.createMeal(meal)
	if err != nil {
		http.Error(w, "Failed to create meal", http.StatusInternalServerError)
		return
	}

	meal.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(meal)
}

// CreateBodyWeight handles body weight creation
func (h *Handler) CreateBodyWeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateBodyWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	bodyWeight := models.BodyWeight{
		UserID:    userID,
		Weight:    req.Weight,
		Unit:      req.Unit,
		Date:      date,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := h.createBodyWeight(bodyWeight)
	if err != nil {
		http.Error(w, "Failed to create body weight entry", http.StatusInternalServerError)
		return
	}

	bodyWeight.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bodyWeight)
}

// CreateBodyFat handles body fat creation
func (h *Handler) CreateBodyFat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateBodyFatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	bodyFat := models.BodyFat{
		UserID:      userID,
		BodyFatPct:  req.BodyFatPct,
		Date:        date,
		Measurement: req.Measurement,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.createBodyFat(bodyFat)
	if err != nil {
		http.Error(w, "Failed to create body fat entry", http.StatusInternalServerError)
		return
	}

	bodyFat.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bodyFat)
}

// AnalyticsPage renders the analytics dashboard
type AnalyticsPage struct {
	Title      string
	Analytics  models.AnalyticsData
	ErrMsg     string
	UserSettings *models.UserSettings
}

func (h *Handler) Analytics(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get analytics data for default 14 days
	analyticsData, err := h.getAnalyticsData(userID, 14)
	if err != nil {
		log.Printf("Failed to get analytics data: %v", err)
		http.Error(w, "Failed to load analytics data", http.StatusInternalServerError)
		return
	}

	data := AnalyticsPage{
		Title:        "Analytics Dashboard",
		Analytics:    *analyticsData,
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "analytics.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// AnalyticsAPI returns analytics data as JSON
func (h *Handler) AnalyticsAPI(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	days := 14
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	analyticsData, err := h.getAnalyticsData(userID, days)
	if err != nil {
		log.Printf("Failed to get analytics data: %v", err)
		http.Error(w, "Failed to load analytics data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analyticsData)
}

// WeeklySummaryAPI returns weekly summary data as JSON
func (h *Handler) WeeklySummaryAPI(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get year and week from query parameters
	yearParam := r.URL.Query().Get("year")
	weekParam := r.URL.Query().Get("week")

	// If no parameters provided, use current week
	var year, week int
	var err error

	if yearParam == "" || weekParam == "" {
		now := time.Now()
		year, week = now.ISOWeek()
	} else {
		year, err = strconv.Atoi(yearParam)
		if err != nil || year < 2020 || year > 2030 {
			http.Error(w, "Invalid year parameter", http.StatusBadRequest)
			return
		}

		week, err = strconv.Atoi(weekParam)
		if err != nil || week < 1 || week > 53 {
			http.Error(w, "Invalid week parameter", http.StatusBadRequest)
			return
		}
	}

	summary, err := h.getWeeklySummary(userID, year, week)
	if err != nil {
		log.Printf("Error getting weekly summary: %v", err)
		http.Error(w, "Failed to get weekly summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// MonthlySummaryAPI returns monthly summary data as JSON
func (h *Handler) MonthlySummaryAPI(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get year and month from query parameters
	yearParam := r.URL.Query().Get("year")
	monthParam := r.URL.Query().Get("month")

	// If no parameters provided, use current month
	var year, month int
	var err error

	if yearParam == "" || monthParam == "" {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	} else {
		year, err = strconv.Atoi(yearParam)
		if err != nil || year < 2020 || year > 2030 {
			http.Error(w, "Invalid year parameter", http.StatusBadRequest)
			return
		}

		month, err = strconv.Atoi(monthParam)
		if err != nil || month < 1 || month > 12 {
			http.Error(w, "Invalid month parameter", http.StatusBadRequest)
			return
		}
	}

	summary, err := h.getMonthlySummary(userID, year, month)
	if err != nil {
		log.Printf("Error getting monthly summary: %v", err)
		http.Error(w, "Failed to get monthly summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// Profile renders the user profile page
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get user information
	user, err := h.getUserByID(userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Failed to load user profile", http.StatusInternalServerError)
		return
	}

	// Get recent workouts for the profile
	workouts, err := h.getRecentWorkouts(5)
	if err != nil {
		log.Printf("Failed to get recent workouts: %v", err)
		// Don't error out, just continue without workouts
	}

	// Get workout statistics
	stats, err := h.getWorkoutStats()
	if err != nil {
		log.Printf("Failed to get workout stats: %v", err)
		// Provide default stats if there's an error
		stats = map[string]interface{}{
			"total_workouts":    0,
			"this_week_workouts": 0,
			"avg_duration":      0,
		}
	}

	data := struct {
		User         *models.User
		Workouts     []models.Workout
		Stats        map[string]interface{}
		Title        string
		UserSettings *models.UserSettings
	}{
		User:         &user,
		Workouts:     workouts,
		Stats:        stats,
		Title:        "Profile",
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "profile.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// CreateBodyMeasurement handles body measurement creation
func (h *Handler) CreateBodyMeasurement(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateBodyMeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	bodyMeasurement := models.BodyMeasurement{
		UserID:      userID,
		Measurement: req.Measurement,
		Value:       req.Value,
		Unit:        req.Unit,
		Date:        date,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.createBodyMeasurement(bodyMeasurement)
	if err != nil {
		http.Error(w, "Failed to create body measurement entry", http.StatusInternalServerError)
		return
	}

	bodyMeasurement.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bodyMeasurement)
}

// ExerciseLibrary displays the exercise library page
func (h *Handler) ExerciseLibrary(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.getAllPredefinedExercises()
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	data := struct {
		Exercises []models.PredefinedExercise
		Title     string
	}{
		Exercises: exercises,
		Title:     "Exercise Library",
	}

	if err := h.templates.ExecuteTemplate(w, "exercise_library.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// ProgressStats displays progress statistics
func (h *Handler) ProgressStats(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get workout stats
	workoutStats, err := h.getWorkoutStats()
	if err != nil {
		http.Error(w, "Failed to load workout statistics", http.StatusInternalServerError)
		return
	}

	// Get body weight data
	bodyWeights, err := h.getBodyWeightsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to load body weight data", http.StatusInternalServerError)
		return
	}

	data := struct {
		WorkoutStats map[string]interface{}
		BodyWeights  []models.BodyWeight
		Title        string
	}{
		WorkoutStats: workoutStats,
		BodyWeights:  bodyWeights,
		Title:        "Progress Stats",
	}

	if err := h.templates.ExecuteTemplate(w, "progress_stats.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// LogMeals displays the meal logging page
func (h *Handler) LogMeals(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get today's meals
	today := time.Now()
	meals, err := h.getMealsByUserAndDate(userID, today)
	if err != nil {
		http.Error(w, "Failed to load meals", http.StatusInternalServerError)
		return
	}

	data := struct {
		Meals []models.Meal
		Date  string
		Title string
	}{
		Meals: meals,
		Date:  today.Format("2006-01-02"),
		Title: "Log Meals",
	}

	if err := h.templates.ExecuteTemplate(w, "log_meals.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// BodyWeight displays the body weight tracking page
func (h *Handler) BodyWeight(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	bodyWeights, err := h.getBodyWeightsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to load body weights", http.StatusInternalServerError)
		return
	}

	data := struct {
		BodyWeights []models.BodyWeight
		Title       string
	}{
		BodyWeights: bodyWeights,
		Title:       "Body Weight Tracking",
	}

	if err := h.templates.ExecuteTemplate(w, "body_weight.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// BodyFatPage displays the body fat tracking page
func (h *Handler) BodyFatPage(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	bodyFats, err := h.getBodyFatsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to load body fat data", http.StatusInternalServerError)
		return
	}

	data := struct {
		BodyFats []models.BodyFat
		Title    string
	}{
		BodyFats: bodyFats,
		Title:    "Body Fat Tracking",
	}

	if err := h.templates.ExecuteTemplate(w, "body_fat.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// BodyMeasurements displays the body measurements tracking page
func (h *Handler) BodyMeasurements(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	bodyMeasurements, err := h.getBodyMeasurementsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to load body measurements", http.StatusInternalServerError)
		return
	}

	data := struct {
		BodyMeasurements []models.BodyMeasurement
		Title            string
	}{
		BodyMeasurements: bodyMeasurements,
		Title:            "Body Measurements",
	}

	if err := h.templates.ExecuteTemplate(w, "body_measurements.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// TemplatesList displays the templates list page
func (h *Handler) TemplatesList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Title:        "Workout Templates",
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "templates_list.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// TemplateDetails displays a specific template
func (h *Handler) TemplateDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get template with exercises
	template, err := h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	data := struct {
		Template     *models.WorkoutTemplate
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Template:     &template,
		Title:        template.Name,
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "template_details.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// TemplateEdit displays the template edit page
func (h *Handler) TemplateEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get template with exercises
	template, err := h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	data := struct {
		Template     *models.WorkoutTemplate
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Template:     &template,
		Title:        "Edit " + template.Name,
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "template_edit.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// ProgramsList displays the programs list page
func (h *Handler) ProgramsList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Title:        "Workout Programs",
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "programs_list.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// ProgramDetails displays a specific program with weekly schedule
func (h *Handler) ProgramDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	programID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	// Authentication check for protected content
	session, _ := h.store.Get(r, "session-name")
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get program with templates
	program, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	data := struct {
		Program      *models.WorkoutProgram
		Title        string
		UserSettings *models.UserSettings
		CurrentPath  string
	}{
		Program:      &program,
		Title:        program.Name,
		UserSettings: h.getUserSettingsForTemplate(r),
		CurrentPath:  r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "program_details.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// ProgramEdit displays the program edit page
func (h *Handler) ProgramEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	programID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get program with templates
	program, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	// Get available templates for dropdown
	availableTemplates, err := h.getWorkoutTemplatesByUserID(userID)
	if err != nil {
		log.Printf("Failed to load available templates: %v", err)
		availableTemplates = []models.WorkoutTemplate{} // Continue with empty list
	}

	// Convert available templates to JSON for JavaScript
	availableTemplatesJSON, _ := json.Marshal(availableTemplates)

	data := struct {
		Program                *models.WorkoutProgram
		AvailableTemplates     []models.WorkoutTemplate
		AvailableTemplatesJSON string
		Title                  string
		UserSettings           *models.UserSettings
		CurrentPath            string
	}{
		Program:                &program,
		AvailableTemplates:     availableTemplates,
		AvailableTemplatesJSON: string(availableTemplatesJSON),
		Title:                  "Edit " + program.Name,
		UserSettings:           h.getUserSettingsForTemplate(r),
		CurrentPath:            r.URL.Path,
	}

	if err := h.templates.ExecuteTemplate(w, "program_edit.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// AccountSettings displays the account settings page
func (h *Handler) AccountSettings(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	log.Printf("ACCOUNT_SETTINGS - All session values: %+v", session.Values)
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Printf("ACCOUNT_SETTINGS - User ID not found in session or wrong type")
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}
	log.Printf("ACCOUNT_SETTINGS - Retrieved user ID: %d", userID)

	// Get user profile
	user, err := h.getUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to load user profile", http.StatusInternalServerError)
		return
	}

	// Get user settings - this will automatically create default settings if none exist
	settings, err := h.getUserSettings(userID)
	if err != nil {
		log.Printf("Failed to get or create user settings: %v", err)
		http.Error(w, "Failed to load user settings", http.StatusInternalServerError)
		return
	}

data := struct {
		User     models.User
		Settings models.UserSettings
		Title    string
	}{
		User:     user,
		Settings: settings,
		Title:    "Account Settings",
	}

	log.Printf("AccountSettings Data: %+v", data)

	if err := h.templates.ExecuteTemplate(w, "account_settings.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// UpdateProfile handles profile updates
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		req.Username = r.FormValue("username")
		req.Email = r.FormValue("email")
		req.FullName = r.FormValue("full_name")
		req.Bio = r.FormValue("bio")
	}

	err := h.updateUserProfile(userID, req)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
	} else {
		http.Redirect(w, r, "/account-settings?updated=profile", http.StatusSeeOther)
	}
}

// UpdateSettings handles settings updates
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.UpdateSettingsRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		req.Theme = r.FormValue("theme")
		req.Timezone = r.FormValue("timezone")
		req.WeightUnit = r.FormValue("weight_unit")
		req.DistanceUnit = r.FormValue("distance_unit")
		req.DateFormat = r.FormValue("date_format")
		req.Notifications = r.FormValue("notifications") == "true"
		req.PrivacyMode = r.FormValue("privacy_mode") == "true"
		req.Language = r.FormValue("language")
		
		// Parse auto logout
		autoLogoutStr := r.FormValue("auto_logout")
		if autoLogoutStr != "" {
			autoLogout, err := strconv.Atoi(autoLogoutStr)
			if err != nil {
				req.AutoLogout = 0
			} else {
				req.AutoLogout = autoLogout
			}
		}
	}

	err := h.updateUserSettings(userID, req)
	if err != nil {
		log.Printf("Failed to update settings: %v", err)
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
	} else {
		http.Redirect(w, r, "/account-settings?updated=settings", http.StatusSeeOther)
	}
}

// ChangePassword handles password changes
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		req.CurrentPassword = r.FormValue("current_password")
		req.NewPassword = r.FormValue("new_password")
		req.ConfirmPassword = r.FormValue("confirm_password")
	}

	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Get current user to verify password
	user, err := h.getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !checkPasswordHash(req.CurrentPassword, user.PasswordHash) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update password
	err = h.updateUserPassword(userID, newPasswordHash)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
	} else {
		http.Redirect(w, r, "/account-settings?updated=password", http.StatusSeeOther)
	}
}

// DeleteAccount handles account deletion
func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.DeleteAccountRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		req.Password = r.FormValue("password")
		req.ConfirmDeletion = r.FormValue("confirm_deletion")
	}

	// Validate confirmation
	if req.ConfirmDeletion != "DELETE" {
		http.Error(w, "Must type DELETE to confirm account deletion", http.StatusBadRequest)
		return
	}

	// Get current user to verify password
	user, err := h.getUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify password
	if !checkPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Password is incorrect", http.StatusUnauthorized)
		return
	}

	// Delete account
	err = h.deleteUserAccount(userID)
	if err != nil {
		log.Printf("Failed to delete account: %v", err)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	// Clear session
	session.Values["authenticated"] = false
	delete(session.Values, "user_id")
	session.Save(r, w)

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
	} else {
		http.Redirect(w, r, "/login?deleted=true", http.StatusSeeOther)
	}
}

// GetExerciseProgressChart returns exercise-specific progress chart data as JSON
func (h *Handler) GetExerciseProgressChart(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get exercise name from URL path or query parameters
	vars := mux.Vars(r)
	exerciseName := vars["exercise"]
	if exerciseName == "" {
		exerciseName = r.URL.Query().Get("exercise")
	}
	if exerciseName == "" {
		http.Error(w, "Exercise name is required", http.StatusBadRequest)
		return
	}

	// Get date range from query parameters
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Default to last 90 days if no date range provided
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, 0, -90).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// Validate date formats
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		http.Error(w, "Invalid start_date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		http.Error(w, "Invalid end_date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Get exercise progress chart data
	chartData, err := h.getExerciseProgressChart(userID, exerciseName, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get exercise progress chart for %s: %v", exerciseName, err)
		http.Error(w, "Failed to load exercise progress data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chartData)
}

// GetExerciseList returns a list of exercises performed by the user
func (h *Handler) GetExerciseList(w http.ResponseWriter, r *http.Request) {
	// Authentication check for protected content
	session, _ := h.store.Get(r, "session-name")
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Get date range from query parameters for filtering
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Default to last 365 days if no date range provided
	if startDate == "" || endDate == "" {
		now := time.Now()
		startDate = now.AddDate(-1, 0, 0).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// Query for unique exercise names within the date range
	// Note: For now, we don't filter by user since the workouts table doesn't have user_id
	// This should be updated when user-specific workouts are implemented
	query := `
		SELECT DISTINCT e.name, e.category, COUNT(*) as frequency,
		       MAX(w.date) as last_performed
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.name, e.category
		ORDER BY frequency DESC, e.name
	`

	// TODO: Add user_id to workouts table and filter by userID
	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		log.Printf("Failed to query exercise list: %v", err)
		http.Error(w, "Failed to load exercise list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ExerciseInfo struct {
		Name          string    `json:"name"`
		Category      string    `json:"category"`
		Frequency     int       `json:"frequency"`
		LastPerformed time.Time `json:"last_performed"`
	}

	var exercises []ExerciseInfo
	for rows.Next() {
		var ex ExerciseInfo
		err := rows.Scan(&ex.Name, &ex.Category, &ex.Frequency, &ex.LastPerformed)
		if err != nil {
			log.Printf("Error scanning exercise row: %v", err)
			continue
		}
		exercises = append(exercises, ex)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating exercise rows: %v", err)
		http.Error(w, "Failed to process exercise list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}

// ========== WORKOUT TEMPLATE HANDLERS ==========

// GetWorkoutTemplates returns all workout templates for the current user
func (h *Handler) GetWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	templates, err := h.getWorkoutTemplatesByUserID(userID)
	if err != nil {
		log.Printf("Failed to get workout templates: %v", err)
		http.Error(w, "Failed to load workout templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// GetWorkoutTemplate returns a specific workout template by ID
func (h *Handler) GetWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	template, err := h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		log.Printf("Failed to get workout template: %v", err)
		http.Error(w, "Workout template not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// CreateWorkoutTemplate creates a new workout template
func (h *Handler) CreateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateWorkoutTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	template := models.WorkoutTemplate{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.createWorkoutTemplate(template)
	if err != nil {
		log.Printf("Failed to create workout template: %v", err)
		http.Error(w, "Failed to create workout template", http.StatusInternalServerError)
		return
	}

	// Create template exercises
	for index, exerciseReq := range req.Exercises {
		exercise := models.TemplateExercise{
			TemplateID:   id,
			Name:         exerciseReq.Name,
			Category:     exerciseReq.Category,
			OrderIndex:   index,
			TargetSets:   exerciseReq.TargetSets,
			TargetReps:   exerciseReq.TargetReps,
			TargetWeight: exerciseReq.TargetWeight,
			RestTime:     exerciseReq.RestTime,
			Notes:        exerciseReq.Notes,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := h.createTemplateExercise(exercise)
		if err != nil {
			log.Printf("Failed to create template exercise: %v", err)
			// Continue with other exercises even if one fails
		}
	}

	// Return the created template with exercises
	createdTemplate, err := h.getWorkoutTemplateByID(id, userID)
	if err != nil {
		log.Printf("Failed to get created template: %v", err)
		http.Error(w, "Template created but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTemplate)
}

// UpdateWorkoutTemplate updates an existing workout template
func (h *Handler) UpdateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && (r.Method != "POST" || r.FormValue("_method") != "PUT") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateWorkoutTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Verify template belongs to user
	_, err = h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	template := models.WorkoutTemplate{
		ID:          templateID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	err = h.updateWorkoutTemplate(template)
	if err != nil {
		log.Printf("Failed to update workout template: %v", err)
		http.Error(w, "Failed to update workout template", http.StatusInternalServerError)
		return
	}

	// Update exercises if provided
	if req.Exercises != nil {
		// Delete existing exercises
		err = h.deleteTemplateExercisesByTemplateID(templateID)
		if err != nil {
			log.Printf("Failed to delete existing template exercises: %v", err)
		}

		// Create new exercises
		for i, exerciseReq := range req.Exercises {
			exercise := models.TemplateExercise{
				TemplateID:   templateID,
				Name:         exerciseReq.Name,
				Category:     exerciseReq.Category,
				OrderIndex:   i,
				TargetSets:   exerciseReq.TargetSets,
				TargetReps:   exerciseReq.TargetReps,
				TargetWeight: exerciseReq.TargetWeight,
				RestTime:     exerciseReq.RestTime,
				Notes:        exerciseReq.Notes,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			_, err := h.createTemplateExercise(exercise)
			if err != nil {
				log.Printf("Failed to create template exercise: %v", err)
			}
		}
	}

	// Return updated template
	updatedTemplate, err := h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		log.Printf("Failed to get updated template: %v", err)
		http.Error(w, "Template updated but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTemplate)
}

// DeleteWorkoutTemplate deletes a workout template
func (h *Handler) DeleteWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" && (r.Method != "POST" || r.FormValue("_method") != "DELETE") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	// Verify template belongs to user
	_, err = h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = h.deleteWorkoutTemplate(templateID)
	if err != nil {
		log.Printf("Failed to delete workout template: %v", err)
		http.Error(w, "Failed to delete workout template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========== WORKOUT PROGRAM HANDLERS ==========

// GetWorkoutPrograms returns all workout programs
func (h *Handler) GetWorkoutPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := h.getAllWorkoutPrograms()
	if err != nil {
		log.Printf("Failed to get workout programs: %v", err)
		http.Error(w, "Failed to load workout programs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programs)
}

// GetWorkoutProgram returns a specific workout program by ID
func (h *Handler) GetWorkoutProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	programID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	program, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		log.Printf("Failed to get workout program: %v", err)
		http.Error(w, "Workout program not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

// CreateWorkoutProgram creates a new workout program
func (h *Handler) CreateWorkoutProgram(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	var req models.CreateWorkoutProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	program := models.WorkoutProgram{
		Name:          req.Name,
		Description:   req.Description,
		Difficulty:    req.Difficulty,
		DurationWeeks: req.DurationWeeks,
		Goal:          req.Goal,
		IsPublic:      req.IsPublic,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	id, err := h.createWorkoutProgram(program)
	if err != nil {
		log.Printf("Failed to create workout program: %v", err)
		http.Error(w, "Failed to create workout program", http.StatusInternalServerError)
		return
	}

	// Create program templates
	for _, templateReq := range req.Templates {
		programTemplate := models.ProgramTemplate{
			ProgramID:  id,
			TemplateID: templateReq.TemplateID,
			DayOfWeek:  templateReq.DayOfWeek,
			WeekNumber: templateReq.WeekNumber,
			OrderIndex: templateReq.OrderIndex,
			CreatedAt:  time.Now(),
		}
		_, err := h.createProgramTemplate(programTemplate)
		if err != nil {
			log.Printf("Failed to create program template: %v", err)
			// Continue with other templates even if one fails
		}
	}

	// Return the created program with templates
	createdProgram, err := h.getWorkoutProgramByID(id)
	if err != nil {
		log.Printf("Failed to get created program: %v", err)
		http.Error(w, "Program created but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProgram)
}

// UpdateWorkoutProgram updates an existing workout program
func (h *Handler) UpdateWorkoutProgram(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && (r.Method != "POST" || r.FormValue("_method") != "PUT") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	programID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateWorkoutProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Verify program exists and user has permission to edit
	existingProgram, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	// Check if user created the program or if it's public
	if existingProgram.CreatedBy != userID && !existingProgram.IsPublic {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	program := models.WorkoutProgram{
		ID:            programID,
		Name:          req.Name,
		Description:   req.Description,
		Difficulty:    req.Difficulty,
		DurationWeeks: req.DurationWeeks,
		Goal:          req.Goal,
		IsPublic:      req.IsPublic,
		CreatedBy:     existingProgram.CreatedBy,
		UpdatedAt:     time.Now(),
	}

	err = h.updateWorkoutProgram(program)
	if err != nil {
		log.Printf("Failed to update workout program: %v", err)
		http.Error(w, "Failed to update workout program", http.StatusInternalServerError)
		return
	}

	// Update program templates if provided
	if req.Templates != nil {
		// Delete existing program templates
		err = h.deleteProgramTemplatesByProgramID(programID)
		if err != nil {
			log.Printf("Failed to delete existing program templates: %v", err)
		}

		// Create new program templates
		for _, templateReq := range req.Templates {
			programTemplate := models.ProgramTemplate{
				ProgramID:  programID,
				TemplateID: templateReq.TemplateID,
				DayOfWeek:  templateReq.DayOfWeek,
				WeekNumber: templateReq.WeekNumber,
				OrderIndex: templateReq.OrderIndex,
				CreatedAt:  time.Now(),
			}
			_, err := h.createProgramTemplate(programTemplate)
			if err != nil {
				log.Printf("Failed to create program template: %v", err)
			}
		}
	}

	// Return updated program
	updatedProgram, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		log.Printf("Failed to get updated program: %v", err)
		http.Error(w, "Program updated but failed to retrieve", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProgram)
}

// DeleteWorkoutProgram deletes a workout program
func (h *Handler) DeleteWorkoutProgram(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" && (r.Method != "POST" || r.FormValue("_method") != "DELETE") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	programID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	// Verify program exists and user has permission to delete
	existingProgram, err := h.getWorkoutProgramByID(programID)
	if err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	// Only the creator can delete a program
	if existingProgram.CreatedBy != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	err = h.deleteWorkoutProgram(programID)
	if err != nil {
		log.Printf("Failed to delete workout program: %v", err)
		http.Error(w, "Failed to delete workout program", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ShareWorkoutTemplate handles template sharing
func (h *Handler) ShareWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	var req models.ShareTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Verify template belongs to user
	_, err = h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	sharing := models.TemplateSharing{
		TemplateID:   templateID,
		OwnerID:      userID,
		SharedWithID: req.SharedWithID,
		Permission:   req.Permission,
		CreatedAt:    time.Now(),
	}

	_, err = h.createTemplateSharing(sharing)
	if err != nil {
		log.Printf("Failed to share template: %v", err)
		http.Error(w, "Failed to share template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Template shared successfully"})
}

// GetSharedTemplates returns templates shared with the current user
func (h *Handler) GetSharedTemplates(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	templates, err := h.getSharedTemplates(userID)
	if err != nil {
		log.Printf("Failed to get shared templates: %v", err)
		http.Error(w, "Failed to load shared templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// CreateWorkoutFromTemplate creates a workout from a template with optional customizations
func (h *Handler) CreateWorkoutFromTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	templateID, err := strconv.Atoi(vars["template_id"])
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	var req models.CreateWorkoutFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set template ID from URL parameter
	req.TemplateID = templateID

	// Get the template (check if user owns it or has access to shared template)
	template, err := h.getWorkoutTemplateByID(templateID, userID)
	if err != nil {
		http.Error(w, "Template not found or access denied", http.StatusNotFound)
		return
	}

	// Parse workout date
	workoutDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format (YYYY-MM-DD required)", http.StatusBadRequest)
		return
	}

	// Create workout name (use custom name or template name)
	workoutName := req.Name
	if workoutName == "" {
		workoutName = template.Name
	}

	// Create the workout
	workout := models.Workout{
		Name:      workoutName,
		Date:      workoutDate,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	workoutID, err := h.createWorkoutWithUser(workout, userID)
	if err != nil {
		log.Printf("Failed to create workout from template: %v", err)
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	// Create customization map for quick lookup
	customizationMap := make(map[string]models.ExerciseCustomization)
	for _, custom := range req.Customizations {
		customizationMap[custom.ExerciseName] = custom
	}

	// Create exercises from template
	for _, templateExercise := range template.Exercises {
		// Check if this exercise should be skipped
		if customization, exists := customizationMap[templateExercise.Name]; exists && customization.Skip {
			continue
		}

		// Create exercise
		exercise := models.Exercise{
			WorkoutID: workoutID,
			Name:      templateExercise.Name,
			Category:  templateExercise.Category,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		exerciseID, err := h.createExerciseForWorkout(exercise)
		if err != nil {
			log.Printf("Failed to create exercise from template: %v", err)
			continue
		}

		// For now, we'll set default values since Exercise doesn't have target fields
		// This should be updated when template exercise support is fully implemented
		targetSets := 3    // default
		targetReps := 10   // default
		targetWeight := 0.0 // default
		restTime := 60     // default 60 seconds

		// Apply customizations if provided
		if customization, exists := customizationMap[templateExercise.Name]; exists {
			if customization.TargetSets > 0 {
				targetSets = customization.TargetSets
			}
			if customization.TargetReps > 0 {
				targetReps = customization.TargetReps
			}
			if customization.TargetWeight > 0 {
				targetWeight = customization.TargetWeight
			}
		}

		// Create sets based on target sets
		for setNum := 1; setNum <= targetSets; setNum++ {
			set := models.Set{
				ExerciseID: exerciseID,
				SetNumber:  setNum,
				Reps:       targetReps,
				Weight:     targetWeight,
				RestTime:   restTime,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			_, err := h.createSet(set)
			if err != nil {
				log.Printf("Failed to create set from template: %v", err)
			}
		}

		// Add exercise to workout's exercise list
		exercise.ID = exerciseID
		workout.Exercises = append(workout.Exercises, exercise)
	}

	// Record template usage for analytics
	usage := models.TemplateUsage{
		TemplateID: templateID,
		UserID:     userID,
		WorkoutID:  workoutID,
		UsedAt:     time.Now(),
		CreatedAt:  time.Now(),
	}

	_, err = h.createTemplateUsage(usage)
	if err != nil {
		log.Printf("Failed to record template usage: %v", err)
		// Don't fail the request for usage tracking failure
	}

	// Get the complete workout with exercises
	completeWorkout, err := h.getWorkoutByIDWithUser(workoutID, userID)
	if err != nil {
		log.Printf("Failed to get complete workout: %v", err)
		// Return basic response if we can't get the complete workout
		response := models.WorkoutFromTemplateResponse{
			Workout:      models.Workout{ID: workoutID, Name: workoutName, Date: workoutDate},
			TemplateUsed: models.WorkoutTemplate{ID: templateID, Name: template.Name},
			Message:      "Workout created successfully from template",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return complete response
	response := models.WorkoutFromTemplateResponse{
		Workout:      completeWorkout,
		TemplateUsed: models.WorkoutTemplate{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			UserID:      template.UserID,
			CreatedAt:   template.CreatedAt,
			UpdatedAt:   template.UpdatedAt,
		},
		Message: "Workout created successfully from template",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
