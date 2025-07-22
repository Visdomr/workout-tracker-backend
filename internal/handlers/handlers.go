package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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
	// Load templates with proper parsing for inheritance
	templates := template.Must(template.ParseGlob("web/templates/*.html"))
	
	store := sessions.NewCookieStore([]byte("something-very-secret"))
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
	// Get recent workouts for the dashboard
	workouts, err := h.getRecentWorkouts(5)
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
	}{
		Workouts:     workouts,
		Stats:        stats,
		Title:        "Workout Tracker",
		UserSettings: h.getUserSettingsForTemplate(r),
	}

	if err := h.templates.ExecuteTemplate(w, "index_dashboard.html", data); err != nil {
		log.Printf("Template error: %v", err)
	http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GetWorkouts returns all workouts
func (h *Handler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.getAllWorkouts()
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

		id, err := h.createWorkout(workout)
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

	log.Printf("GetWorkout - Fetching workout with ID: %d", id)
	workout, err := h.getWorkoutByID(id)
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

	if r.Method == "POST" {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Check for method override
		if r.FormValue("_method") != "PUT" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		dateStr := r.FormValue("date")
		durationStr := r.FormValue("duration")
		notes := r.FormValue("notes")

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

		err = h.updateWorkout(workout)
		if err != nil {
			http.Error(w, "Failed to update workout", http.StatusInternalServerError)
			return
		}

		// Redirect to the updated workout
		http.Redirect(w, r, "/workouts/"+strconv.Itoa(id), http.StatusSeeOther)
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

	err = h.deleteWorkout(id)
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

// AuthMiddleware checks if user is authenticated
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug: print all cookies received
		log.Printf("Auth check - URL: %s, Cookies received: %v", r.URL.Path, r.Cookies())
		session, err := h.store.Get(r, "session-name")
		if err != nil {
			log.Printf("Error getting session: %v", err)
		}
		auth, ok := session.Values["authenticated"]
		log.Printf("Auth check - URL: %s, Session exists: %v, Auth value: %v, Type: %T, Session Values: %+v", r.URL.Path, ok, auth, auth, session.Values)
		if !ok || auth != true {
			log.Printf("Redirecting to login from %s", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		log.Printf("Authentication successful, proceeding to %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
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
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
