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

		session, _ := h.store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Values["user_id"] = user.ID
		session.Save(r, w)

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
		Workouts []models.Workout
		Stats    map[string]interface{}
		Title    string
	}{
		Workouts: workouts,
		Stats:    stats,
		Title:    "Workout Tracker",
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
		Workouts []models.Workout
		Title    string
	}{
		Workouts: workouts,
		Title:    "All Workouts",
	}

	if err := h.templates.ExecuteTemplate(w, "workouts.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// CreateWorkout handles workout creation
func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		dateStr := r.FormValue("date")
		notes := r.FormValue("notes")

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
		Title string
	}{
		Title: "Create Workout",
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
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.getWorkoutByID(id)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	data := struct {
		Workout models.Workout
		Title   string
	}{
		Workout: workout,
		Title:   "Workout: " + workout.Name,
	}

	if err := h.templates.ExecuteTemplate(w, "workout_detail.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
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
		Workout models.Workout
		Title   string
	}{
		Workout: workout,
		Title:   "Edit Workout: " + workout.Name,
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
	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = false
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// AuthMiddleware checks if user is authenticated
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.store.Get(r, "session-name")
		auth, ok := session.Values["authenticated"]
		if !ok || auth != true {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
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
