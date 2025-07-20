package main

import (
	"log"
	"net/http"
	"os"

	"workout-tracker/internal/database"
	"workout-tracker/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.New(db)

	// Setup routes
	r := mux.NewRouter()
	
	// Authentication routes
	r.HandleFunc("/login", h.Login).Methods("GET", "POST")
	r.HandleFunc("/register", h.Register).Methods("GET", "POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")

	// Protected web routes
	r.HandleFunc("/", h.AuthMiddleware(h.Home)).Methods("GET")
	r.HandleFunc("/workouts", h.AuthMiddleware(h.GetWorkouts)).Methods("GET")
	r.HandleFunc("/workouts/new", h.AuthMiddleware(h.CreateWorkout)).Methods("GET")
	r.HandleFunc("/workouts", h.AuthMiddleware(h.CreateWorkout)).Methods("POST")
	r.HandleFunc("/workouts/{id}", h.AuthMiddleware(h.GetWorkout)).Methods("GET")
	r.HandleFunc("/workouts/{id}/edit", h.AuthMiddleware(h.UpdateWorkout)).Methods("GET")
	r.HandleFunc("/workouts/{id}/edit", h.AuthMiddleware(h.UpdateWorkout)).Methods("POST")
	r.HandleFunc("/workouts/{id}/delete", h.AuthMiddleware(h.DeleteWorkout)).Methods("POST")
	
	// API routes
	r.HandleFunc("/api/workouts", h.APIGetWorkouts).Methods("GET")
	r.HandleFunc("/api/workouts", h.APICreateWorkout).Methods("POST")
	r.HandleFunc("/api/workouts/{id}", h.APIGetWorkout).Methods("GET")
	r.HandleFunc("/api/workouts/{id}", h.UpdateWorkout).Methods("PUT")
	r.HandleFunc("/api/workouts/{id}", h.DeleteWorkout).Methods("DELETE")
	
	// Exercise API routes
	r.HandleFunc("/api/exercises", h.CreateExercise).Methods("POST")
	r.HandleFunc("/api/exercises/{id}", h.UpdateExercise).Methods("PUT")
	r.HandleFunc("/api/exercises/{id}", h.DeleteExercise).Methods("DELETE")
	
	// Set API routes
	r.HandleFunc("/api/sets", h.CreateSet).Methods("POST")
	r.HandleFunc("/api/sets/{id}", h.UpdateSet).Methods("PUT")
	r.HandleFunc("/api/sets/{id}", h.DeleteSet).Methods("DELETE")
	
	// Predefined exercise API routes
	r.HandleFunc("/api/predefined-exercises", h.GetPredefinedExercises).Methods("GET")
	r.HandleFunc("/api/predefined-exercises", h.CreatePredefinedExercise).Methods("POST")
	r.HandleFunc("/api/predefined-exercises/category/{category}", h.GetPredefinedExercisesByCategory).Methods("GET")
	
	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
