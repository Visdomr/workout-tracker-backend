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
	FullName     string    `json:"full_name" db:"full_name"`
	Bio          string    `json:"bio" db:"bio"`
	Avatar       string    `json:"avatar" db:"avatar"` // URL to profile picture
	IsActive     bool      `json:"is_active" db:"is_active"`
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

// Meal represents a nutrition entry
type Meal struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Calories    int       `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Fat         float64   `json:"fat" db:"fat"`
	Date        time.Time `json:"date" db:"date"`
	MealType    string    `json:"meal_type" db:"meal_type"` // breakfast, lunch, dinner, snack
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateMealRequest represents the request payload for creating a meal
type CreateMealRequest struct {
	Name     string  `json:"name" validate:"required"`
	Calories int     `json:"calories" validate:"required"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Date     string  `json:"date" validate:"required"`
	MealType string  `json:"meal_type" validate:"required"`
}

// BodyWeight represents a body weight entry
type BodyWeight struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Weight    float64   `json:"weight" db:"weight"`
	Unit      string    `json:"unit" db:"unit"` // lbs or kg
	Date      time.Time `json:"date" db:"date"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBodyWeightRequest represents the request payload for creating a body weight entry
type CreateBodyWeightRequest struct {
	Weight float64 `json:"weight" validate:"required"`
	Unit   string  `json:"unit" validate:"required"`
	Date   string  `json:"date" validate:"required"`
	Notes  string  `json:"notes"`
}

// BodyFat represents a body fat percentage entry
type BodyFat struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	BodyFatPct   float64   `json:"body_fat_pct" db:"body_fat_pct"`
	Date         time.Time `json:"date" db:"date"`
	Measurement  string    `json:"measurement" db:"measurement"` // method used (calipers, scale, etc.)
	Notes        string    `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBodyFatRequest represents the request payload for creating a body fat entry
type CreateBodyFatRequest struct {
	BodyFatPct  float64 `json:"body_fat_pct" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Measurement string  `json:"measurement"`
	Notes       string  `json:"notes"`
}

// BodyMeasurement represents a body measurement entry
type BodyMeasurement struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Measurement  string    `json:"measurement" db:"measurement"` // chest, waist, bicep, etc.
	Value        float64   `json:"value" db:"value"`
	Unit         string    `json:"unit" db:"unit"` // inches or cm
	Date         time.Time `json:"date" db:"date"`
	Notes        string    `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBodyMeasurementRequest represents the request payload for creating a body measurement
type CreateBodyMeasurementRequest struct {
	Measurement string  `json:"measurement" validate:"required"`
	Value       float64 `json:"value" validate:"required"`
	Unit        string  `json:"unit" validate:"required"`
	Date        string  `json:"date" validate:"required"`
	Notes       string  `json:"notes"`
}

// UserSettings represents user preferences and settings
type UserSettings struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	Theme            string    `json:"theme" db:"theme"`                         // light, dark
	Timezone         string    `json:"timezone" db:"timezone"`                   // user's timezone
	WeightUnit       string    `json:"weight_unit" db:"weight_unit"`             // kg, lbs
	DistanceUnit     string    `json:"distance_unit" db:"distance_unit"`         // km, miles
	DateFormat       string    `json:"date_format" db:"date_format"`             // YYYY-MM-DD, MM/DD/YYYY, DD/MM/YYYY
	Notifications    bool      `json:"notifications" db:"notifications"`         // email notifications
	PrivacyMode      bool      `json:"privacy_mode" db:"privacy_mode"`           // hide stats from others
	AutoLogout       int       `json:"auto_logout" db:"auto_logout"`             // minutes, 0 = never
	Language         string    `json:"language" db:"language"`                   // en, es, fr, etc.
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateProfileRequest represents a request to update user profile
type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// UpdateSettingsRequest represents a request to update user settings
type UpdateSettingsRequest struct {
	Theme            string `json:"theme"`
	Timezone         string `json:"timezone"`
	WeightUnit       string `json:"weight_unit"`
	DistanceUnit     string `json:"distance_unit"`
	DateFormat       string `json:"date_format"`
	Notifications    bool   `json:"notifications"`
	PrivacyMode      bool   `json:"privacy_mode"`
	AutoLogout       int    `json:"auto_logout"`
	Language         string `json:"language"`
}

// DeleteAccountRequest represents a request to delete account
type DeleteAccountRequest struct {
	Password        string `json:"password" validate:"required"`
	ConfirmDeletion string `json:"confirm_deletion" validate:"required"` // must be "DELETE"
}

// AnalyticsData represents comprehensive analytics data for the dashboard
type AnalyticsData struct {
	TotalWorkouts   int                    `json:"totalWorkouts"`
	AvgCalories     int                    `json:"avgCalories"`
	WeightChange    float64                `json:"weightChange"`
	WorkoutStreak   int                    `json:"workoutStreak"`
	Dates           []string               `json:"dates"`
	WorkoutCounts   []int                  `json:"workoutCounts"`
	DailyCalories   []int                  `json:"dailyCalories"`
	WeightDates     []string               `json:"weightDates"`
	Weights         []float64              `json:"weights"`
	Durations       []int                  `json:"durations"`
	ExerciseCategories *ExerciseCategoryData `json:"exerciseCategories"`
	Macros          *MacroData             `json:"macros"`
}

// ExerciseCategoryData represents exercise category distribution
type ExerciseCategoryData struct {
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

// MacroData represents macronutrient breakdown
type MacroData struct {
	Protein float64 `json:"protein"`
	Carbs   float64 `json:"carbs"`
	Fat     float64 `json:"fat"`
}
