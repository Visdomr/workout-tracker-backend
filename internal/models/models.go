package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Workout represents a workout session
type Workout struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Date        time.Time `json:"date" db:"date"`
	Duration    int       `json:"duration" db:"duration"` // in minutes
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Exercises   []Exercise `json:"exercises,omitempty"`
}

// Exercise represents an exercise within a workout
type Exercise struct {
	ID          int    `json:"id" db:"id"`
	WorkoutID   int    `json:"workout_id" db:"workout_id"`
	Name        string `json:"name" db:"name"`
	Category    string `json:"category" db:"category"` // e.g., "strength", "cardio", "flexibility"
	Sets        []Set  `json:"sets,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Set represents a set within an exercise
type Set struct {
	ID          int     `json:"id" db:"id"`
	ExerciseID  int     `json:"exercise_id" db:"exercise_id"`
	SetNumber   int     `json:"set_number" db:"set_number"`
	Reps        int     `json:"reps" db:"reps"`
	Weight      float64 `json:"weight" db:"weight"` // in pounds or kg
	Distance    float64 `json:"distance" db:"distance"` // for cardio exercises
	Duration    int     `json:"duration" db:"duration"` // in seconds
	RestTime    int     `json:"rest_time" db:"rest_time"` // in seconds
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateWorkoutRequest represents the request payload for creating a workout
type CreateWorkoutRequest struct {
	Name      string `json:"name" validate:"required"`
	Date      string `json:"date" validate:"required"`
	Duration  int    `json:"duration"`
	Notes     string `json:"notes"`
}

// CreateExerciseRequest represents the request payload for creating an exercise
type CreateExerciseRequest struct {
	WorkoutID int    `json:"workout_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Category  string `json:"category" validate:"required"`
}

// CreateSetRequest represents the request payload for creating a set
type CreateSetRequest struct {
	ExerciseID int     `json:"exercise_id" validate:"required"`
	SetNumber  int     `json:"set_number" validate:"required"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Distance   float64 `json:"distance"`
	Duration   int     `json:"duration"`
	RestTime   int     `json:"rest_time"`
}

// UpdateWorkoutRequest represents the request payload for updating a workout
type UpdateWorkoutRequest struct {
	Name      string `json:"name" validate:"required"`
	Date      string `json:"date" validate:"required"`
	Duration  int    `json:"duration"`
	Notes     string `json:"notes"`
}

// UpdateExerciseRequest represents the request payload for updating an exercise
type UpdateExerciseRequest struct {
	Name      string `json:"name" validate:"required"`
	Category  string `json:"category" validate:"required"`
}

// UpdateSetRequest represents the request payload for updating a set
type UpdateSetRequest struct {
	SetNumber  int     `json:"set_number" validate:"required"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Distance   float64 `json:"distance"`
	Duration   int     `json:"duration"`
	RestTime   int     `json:"rest_time"`
}

// PredefinedExercise represents a predefined exercise in the library
type PredefinedExercise struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePredefinedExerciseRequest represents the request payload for creating a predefined exercise
type CreatePredefinedExerciseRequest struct {
	Name        string `json:"name" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Description string `json:"description"`
}

// UpdatePredefinedExerciseRequest represents the request payload for updating a predefined exercise
type UpdatePredefinedExerciseRequest struct {
	Name        string `json:"name" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Description string `json:"description"`
}
