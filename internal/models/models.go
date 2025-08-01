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

// WorkoutTemplate represents a reusable workout template
// Users can create templates to specify a blueprint for future workouts.
type WorkoutTemplate struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Exercises   []Exercise `json:"exercises,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WorkoutProgram represents a pre-built workout program
// Programs can be predefined sets of workouts for specific goals.
type WorkoutProgram struct {
	ID            int               `json:"id" db:"id"`
	Name          string            `json:"name" db:"name"`
	Description   string            `json:"description" db:"description"`
	Difficulty    string            `json:"difficulty" db:"difficulty"`    // beginner, intermediate, advanced
	DurationWeeks int               `json:"duration_weeks" db:"duration_weeks"`
	Goal          string            `json:"goal" db:"goal"`          // strength, muscle_gain, fat_loss, endurance
	IsPublic      bool              `json:"is_public" db:"is_public"`
	CreatedBy     int               `json:"created_by" db:"created_by"`
	Templates     []ProgramTemplate `json:"templates,omitempty"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// CreateTemplateRequest represents the request payload for creating a workout template
type CreateTemplateRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Exercises   []Exercise `json:"exercises"`
}

// UpdateTemplateRequest represents the request payload for updating a workout template
type UpdateTemplateRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Exercises   []Exercise `json:"exercises"`
}

// CreateProgramRequest represents the request payload for creating a workout program
type CreateProgramRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Templates   []WorkoutTemplate `json:"templates"`
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
	UserID      int       `json:"user_id" db:"user_id"`
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
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Category        string    `json:"category" db:"category"`
	Description     string    `json:"description" db:"description"`
	VideoURL        string    `json:"video_url" db:"video_url"`
	Instructions    string    `json:"instructions" db:"instructions"`
	Tips            string    `json:"tips" db:"tips"`
	MuscleGroups    string    `json:"muscle_groups" db:"muscle_groups"`
	Equipment       string    `json:"equipment" db:"equipment"`
	Difficulty      string    `json:"difficulty" db:"difficulty"` // beginner, intermediate, advanced
	ImageURL        string    `json:"image_url" db:"image_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePredefinedExerciseRequest represents the request payload for creating a predefined exercise
type CreatePredefinedExerciseRequest struct {
	Name         string `json:"name" validate:"required"`
	Category     string `json:"category" validate:"required"`
	Description  string `json:"description"`
	VideoURL     string `json:"video_url"`
	Instructions string `json:"instructions"`
	Tips         string `json:"tips"`
	MuscleGroups string `json:"muscle_groups"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	ImageURL     string `json:"image_url"`
}

// UpdatePredefinedExerciseRequest represents the request payload for updating a predefined exercise
type UpdatePredefinedExerciseRequest struct {
	Name         string `json:"name" validate:"required"`
	Category     string `json:"category" validate:"required"`
	Description  string `json:"description"`
	VideoURL     string `json:"video_url"`
	Instructions string `json:"instructions"`
	Tips         string `json:"tips"`
	MuscleGroups string `json:"muscle_groups"`
	Equipment    string `json:"equipment"`
	Difficulty   string `json:"difficulty"`
	ImageURL     string `json:"image_url"`
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
	TotalWorkouts      int                     `json:"totalWorkouts"`
	AvgCalories        int                     `json:"avgCalories"`
	WeightChange       float64                 `json:"weightChange"`
	WorkoutStreak      int                     `json:"workoutStreak"`
	Dates              []string                `json:"dates"`
	WorkoutCounts      []int                   `json:"workoutCounts"`
	DailyCalories      []int                   `json:"dailyCalories"`
	WeightDates        []string                `json:"weightDates"`
	Weights            []float64               `json:"weights"`
	Durations          []int                   `json:"durations"`
	ExerciseCategories *ExerciseCategoryData   `json:"exerciseCategories"`
	Macros             *MacroData              `json:"macros"`
	
	// Enhanced progress tracking
	TotalSets          int                     `json:"totalSets"`
	TotalReps          int                     `json:"totalReps"`
	TotalVolume        float64                 `json:"totalVolume"` // Total weight * reps
	AvgVolume          float64                 `json:"avgVolume"`   // Average volume per workout
	AvgSetsPerWorkout  float64                 `json:"avgSetsPerWorkout"`
	AvgRepsPerSet      float64                 `json:"avgRepsPerSet"`
	MaxWeight          float64                 `json:"maxWeight"`   // Heaviest weight lifted
	PersonalRecords    []PersonalRecord        `json:"personalRecords"`
	StrengthProgress   []StrengthProgress      `json:"strengthProgress"`
	WorkoutIntensity   []WorkoutIntensity      `json:"workoutIntensity"`
	MostPopularExercises []ExerciseFrequency   `json:"mostPopularExercises"`
	ProgressTrends     *ProgressTrends         `json:"progressTrends"`
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

// PersonalRecord represents a personal record for an exercise
type PersonalRecord struct {
	ExerciseName string    `json:"exerciseName"`
	Weight       float64   `json:"weight"`
	Reps         int       `json:"reps"`
	Volume       float64   `json:"volume"` // weight * reps
	OneRepMax    float64   `json:"oneRepMax"`
	Date         time.Time `json:"date"`
	IsNew        bool      `json:"isNew"` // If achieved in current period
}

// StrengthProgress represents strength progression for specific exercises
type StrengthProgress struct {
	ExerciseName string                 `json:"exerciseName"`
	Category     string                 `json:"category"`
	DataPoints   []StrengthDataPoint    `json:"dataPoints"`
	Trend        string                 `json:"trend"` // "improving", "stable", "declining"
	TrendPercent float64                `json:"trendPercent"`
}

// StrengthDataPoint represents a single data point in strength progression
type StrengthDataPoint struct {
	Date      time.Time `json:"date"`
	MaxWeight float64   `json:"maxWeight"`
	MaxReps   int       `json:"maxReps"`
	Volume    float64   `json:"volume"`
	OneRepMax float64   `json:"oneRepMax"`
}

// WorkoutIntensity represents workout intensity metrics over time
type WorkoutIntensity struct {
	Date          time.Time `json:"date"`
	TotalVolume   float64   `json:"totalVolume"`
	AvgWeight     float64   `json:"avgWeight"`
	TotalSets     int       `json:"totalSets"`
	TotalReps     int       `json:"totalReps"`
	IntensityScore float64  `json:"intensityScore"` // Calculated intensity metric
}

// ExerciseFrequency represents how often exercises are performed
type ExerciseFrequency struct {
	ExerciseName string  `json:"exerciseName"`
	Category     string  `json:"category"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
	LastPerformed time.Time `json:"lastPerformed"`
}

// ProgressTrends represents overall progress trends
type ProgressTrends struct {
	Volumetrend        string  `json:"volumeTrend"`        // "up", "down", "stable"
	VolumeChangePercent float64 `json:"volumeChangePercent"`
	StrengthTrend       string  `json:"strengthTrend"`
	StrengthChangePercent float64 `json:"strengthChangePercent"`
	FrequencyTrend      string  `json:"frequencyTrend"`
	FrequencyChangePercent float64 `json:"frequencyChangePercent"`
	ConsistencyScore    float64 `json:"consistencyScore"` // 0-100
}

// Weekly/Monthly Summary Models
type WeeklySummary struct {
	WeekStart        time.Time `json:"weekStart"`
	WeekEnd          time.Time `json:"weekEnd"`
	WeekNumber       int       `json:"weekNumber"`
	Year             int       `json:"year"`
	TotalWorkouts    int       `json:"totalWorkouts"`
	TotalDuration    int       `json:"totalDuration"`
	AvgDuration      float64   `json:"avgDuration"`
	TotalVolume      float64   `json:"totalVolume"`
	TotalSets        int       `json:"totalSets"`
	TotalReps        int       `json:"totalReps"`
	UniqueExercises  int       `json:"uniqueExercises"`
	MaxWeight        float64   `json:"maxWeight"`
	AvgCalories      int       `json:"avgCalories"`
	WeightChange     float64   `json:"weightChange"`
	TopExercises     []string  `json:"topExercises"`
	PRsAchieved      int       `json:"prsAchieved"`
	ConsistencyScore float64   `json:"consistencyScore"`
	IntensityScore   float64   `json:"intensityScore"`
}

type MonthlySummary struct {
	Month               int                    `json:"month"`
	Year                int                    `json:"year"`
	MonthName           string                 `json:"monthName"`
	TotalWorkouts       int                    `json:"totalWorkouts"`
	TotalDuration       int                    `json:"totalDuration"`
	AvgDuration         float64                `json:"avgDuration"`
	TotalVolume         float64                `json:"totalVolume"`
	TotalSets           int                    `json:"totalSets"`
	TotalReps           int                    `json:"totalReps"`
	UniqueExercises     int                    `json:"uniqueExercises"`
	MaxWeight           float64                `json:"maxWeight"`
	AvgCalories         int                    `json:"avgCalories"`
	WeightChange        float64                `json:"weightChange"`
	StartWeight         float64                `json:"startWeight"`
	EndWeight           float64                `json:"endWeight"`
	TopExercises        []string               `json:"topExercises"`
	PRsAchieved         int                    `json:"prsAchieved"`
	ConsistencyScore    float64                `json:"consistencyScore"`
	IntensityScore      float64                `json:"intensityScore"`
	WeeklySummaries     []WeeklySummary        `json:"weeklySummaries"`
	CategoryBreakdown   map[string]int         `json:"categoryBreakdown"`
	ProgressHighlights  []string               `json:"progressHighlights"`
	GoalsAchieved       []string               `json:"goalsAchieved"`
	Recommendations     []string               `json:"recommendations"`
}

type YearlySummary struct {
	Year                int                    `json:"year"`
	TotalWorkouts       int                    `json:"totalWorkouts"`
	TotalDuration       int                    `json:"totalDuration"`
	AvgDuration         float64                `json:"avgDuration"`
	TotalVolume         float64                `json:"totalVolume"`
	TotalSets           int                    `json:"totalSets"`
	TotalReps           int                    `json:"totalReps"`
	UniqueExercises     int                    `json:"uniqueExercises"`
	MaxWeight           float64                `json:"maxWeight"`
	AvgCalories         int                    `json:"avgCalories"`
	TotalWeightChange   float64                `json:"totalWeightChange"`
	StartWeight         float64                `json:"startWeight"`
	EndWeight           float64                `json:"endWeight"`
	TopExercises        []string               `json:"topExercises"`
	TotalPRsAchieved    int                    `json:"totalPRsAchieved"`
	AvgConsistency      float64                `json:"avgConsistency"`
	AvgIntensity        float64                `json:"avgIntensity"`
	MonthlySummaries    []MonthlySummary       `json:"monthlySummaries"`
	QuarterlySummaries  []QuarterlySummary     `json:"quarterlySummaries"`
	Milestones          []Milestone            `json:"milestones"`
	YearHighlights      []string               `json:"yearHighlights"`
	FitnessJourney      []string               `json:"fitnessJourney"`
}

type QuarterlySummary struct {
	Quarter             int       `json:"quarter"`
	Year                int       `json:"year"`
	QuarterName         string    `json:"quarterName"`
	StartDate           time.Time `json:"startDate"`
	EndDate             time.Time `json:"endDate"`
	TotalWorkouts       int       `json:"totalWorkouts"`
	TotalVolume         float64   `json:"totalVolume"`
	AvgIntensity        float64   `json:"avgIntensity"`
	WeightChange        float64   `json:"weightChange"`
	PRsAchieved         int       `json:"prsAchieved"`
	TopAchievements     []string  `json:"topAchievements"`
	FocusAreas          []string  `json:"focusAreas"`
}

type Milestone struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AchievedAt  time.Time `json:"achievedAt"`
	Category    string    `json:"category"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
}

// Summary request/response models
type SummaryRequest struct {
	UserID    int    `json:"userId"`
	Type      string `json:"type"` // "weekly", "monthly", "yearly"
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Year      int    `json:"year,omitempty"`
	Month     int    `json:"month,omitempty"`
	Week      int    `json:"week,omitempty"`
}

// Exercise-specific progress chart models
type ExerciseProgressChart struct {
	ExerciseName    string                    `json:"exerciseName"`
	Category        string                    `json:"category"`
	TimeRange       string                    `json:"timeRange"`
	DataPoints      []ExerciseProgressPoint   `json:"dataPoints"`
	Statistics      ExerciseStatistics        `json:"statistics"`
	Milestones      []ExerciseMilestone       `json:"milestones"`
	Predictions     *ExercisePrediction       `json:"predictions,omitempty"`
	Comparisons     *ExerciseComparison       `json:"comparisons,omitempty"`
}

type ExerciseProgressPoint struct {
	Date        time.Time `json:"date"`
	MaxWeight   float64   `json:"maxWeight"`
	MaxReps     int       `json:"maxReps"`
	TotalVolume float64   `json:"totalVolume"`
	SetsCount   int       `json:"setsCount"`
	OneRepMax   float64   `json:"oneRepMax"`
	WorkoutID   int       `json:"workoutId"`
}

type ExerciseStatistics struct {
	TotalSessions    int       `json:"totalSessions"`
	TotalSets        int       `json:"totalSets"`
	TotalReps        int       `json:"totalReps"`
	TotalVolume      float64   `json:"totalVolume"`
	AverageWeight    float64   `json:"averageWeight"`
	AverageReps      float64   `json:"averageReps"`
	AverageSets      float64   `json:"averageSets"`
	MaxWeightEver    float64   `json:"maxWeightEver"`
	MaxRepsEver      int       `json:"maxRepsEver"`
	MaxVolumeDay     float64   `json:"maxVolumeDay"`
	CurrentOneRepMax float64   `json:"currentOneRepMax"`
	FirstWorkout     time.Time `json:"firstWorkout"`
	LastWorkout      time.Time `json:"lastWorkout"`
	ProgressRate     float64   `json:"progressRate"` // weight gain per week
	Consistency      float64   `json:"consistency"`  // sessions per week
}

type ExerciseMilestone struct {
	ID           int       `json:"id"`
	Type         string    `json:"type"`         // "weight", "reps", "volume", "consistency"
	Description  string    `json:"description"`  // "First 100lb bench press"
	Value        float64   `json:"value"`
	Unit         string    `json:"unit"`
	AchievedAt   time.Time `json:"achievedAt"`
	WorkoutID    int       `json:"workoutId"`
	IsPersonalRecord bool  `json:"isPersonalRecord"`
}

type ExercisePrediction struct {
	NextMilestone     *PredictedMilestone `json:"nextMilestone,omitempty"`
	WeightProgression []WeightPrediction  `json:"weightProgression"`
	TimeToGoals       []GoalPrediction    `json:"timeToGoals"`
}

type PredictedMilestone struct {
	Type            string    `json:"type"`
	Description     string    `json:"description"`
	TargetValue     float64   `json:"targetValue"`
	CurrentValue    float64   `json:"currentValue"`
	EstimatedDate   time.Time `json:"estimatedDate"`
	Confidence      float64   `json:"confidence"` // 0-100%
	WeeksEstimated  int       `json:"weeksEstimated"`
}

type WeightPrediction struct {
	Date           time.Time `json:"date"`
	PredictedWeight float64  `json:"predictedWeight"`
	ConfidenceRange struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	} `json:"confidenceRange"`
}

type GoalPrediction struct {
	GoalType        string    `json:"goalType"`    // "weight", "reps", "volume"
	GoalValue       float64   `json:"goalValue"`
	CurrentValue    float64   `json:"currentValue"`
	EstimatedWeeks  int       `json:"estimatedWeeks"`
	EstimatedDate   time.Time `json:"estimatedDate"`
	Likelihood      string    `json:"likelihood"` // "very_likely", "likely", "possible", "unlikely"
}

type ExerciseComparison struct {
	VsPreviousPeriod *PeriodComparison `json:"vsPreviousPeriod,omitempty"`
	VsAverage        *AverageComparison `json:"vsAverage,omitempty"`
	VsGoals          *GoalsComparison   `json:"vsGoals,omitempty"`
}

type PeriodComparison struct {
	WeightChange    float64 `json:"weightChange"`
	RepsChange      float64 `json:"repsChange"`
	VolumeChange    float64 `json:"volumeChange"`
	FrequencyChange float64 `json:"frequencyChange"`
	ProgressRate    string  `json:"progressRate"` // "faster", "slower", "same"
}

type AverageComparison struct {
	WeightVsAvg    string `json:"weightVsAvg"`    // "above", "below", "average"
	RepsVsAvg      string `json:"repsVsAvg"`
	VolumeVsAvg    string `json:"volumeVsAvg"`
	FrequencyVsAvg string `json:"frequencyVsAvg"`
}

type GoalsComparison struct {
	WeightGoal   *GoalStatus `json:"weightGoal,omitempty"`
	RepsGoal     *GoalStatus `json:"repsGoal,omitempty"`
	VolumeGoal   *GoalStatus `json:"volumeGoal,omitempty"`
	FrequencyGoal *GoalStatus `json:"frequencyGoal,omitempty"`
}

type GoalStatus struct {
	Target     float64 `json:"target"`
	Current    float64 `json:"current"`
	Progress   float64 `json:"progress"` // percentage 0-100
	Status     string  `json:"status"`   // "achieved", "on_track", "behind", "at_risk"
	Timeframe  string  `json:"timeframe"`
}

// Exercise comparison models
type ExerciseComparisonChart struct {
	Exercises   []string                  `json:"exercises"`
	TimeRange   string                    `json:"timeRange"`
	Metric      string                    `json:"metric"` // "weight", "volume", "frequency"
	DataSeries  []ComparisonDataSeries    `json:"dataSeries"`
	Rankings    []ExerciseRanking         `json:"rankings"`
	Insights    []string                  `json:"insights"`
}

type ComparisonDataSeries struct {
	ExerciseName string                    `json:"exerciseName"`
	Color        string                    `json:"color"`
	DataPoints   []ComparisonDataPoint     `json:"dataPoints"`
}

type ComparisonDataPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type ExerciseRanking struct {
	Rank         int     `json:"rank"`
	ExerciseName string  `json:"exerciseName"`
	Value        float64 `json:"value"`
	Change       float64 `json:"change"`
	Trend        string  `json:"trend"`
}

// ========== ENHANCED TEMPLATE AND PROGRAM MODELS ==========

// TemplateExercise represents an exercise within a workout template
type TemplateExercise struct {
	ID           int       `json:"id" db:"id"`
	TemplateID   int       `json:"template_id" db:"template_id"`
	Name         string    `json:"name" db:"name"`
	Category     string    `json:"category" db:"category"`
	OrderIndex   int       `json:"order_index" db:"order_index"`
	TargetSets   int       `json:"target_sets" db:"target_sets"`
	TargetReps   int       `json:"target_reps" db:"target_reps"`
	TargetWeight float64   `json:"target_weight" db:"target_weight"`
	RestTime     int       `json:"rest_time" db:"rest_time"` // in seconds
	Notes        string    `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ProgramTemplate represents the relationship between programs and templates
type ProgramTemplate struct {
	ID              int              `json:"id" db:"id"`
	ProgramID       int              `json:"program_id" db:"program_id"`
	TemplateID      int              `json:"template_id" db:"template_id"`
	DayOfWeek       int              `json:"day_of_week" db:"day_of_week"` // 0-6 (Sunday-Saturday)
	WeekNumber      int              `json:"week_number" db:"week_number"`
	OrderIndex      int              `json:"order_index" db:"order_index"`
	WorkoutTemplate *WorkoutTemplate `json:"workout_template,omitempty"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
}

// Enhanced WorkoutTemplate with TemplateExercise
type WorkoutTemplateWithExercises struct {
	ID          int                 `json:"id"`
	UserID      int                 `json:"user_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Exercises   []TemplateExercise  `json:"exercises"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// Enhanced WorkoutProgram with comprehensive data
type WorkoutProgramWithDetails struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Difficulty    string            `json:"difficulty"`    // beginner, intermediate, advanced
	DurationWeeks int               `json:"duration_weeks"`
	Goal          string            `json:"goal"`          // strength, muscle_gain, fat_loss, endurance
	IsPublic      bool              `json:"is_public"`
	CreatedBy     int               `json:"created_by"`
	Templates     []ProgramTemplate `json:"templates"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// Template sharing between users
type TemplateSharing struct {
	ID           int       `json:"id" db:"id"`
	TemplateID   int       `json:"template_id" db:"template_id"`
	OwnerID      int       `json:"owner_id" db:"owner_id"`
	SharedWithID int       `json:"shared_with_id" db:"shared_with_id"`
	Permission   string    `json:"permission" db:"permission"` // "view" or "edit"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ========== IMPORT/EXPORT MODELS ==========

// ExportJob represents an export job for user data
type ExportJob struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	ExportType    string    `json:"export_type" db:"export_type"` // csv, json, backup
	DataTypes     []string  `json:"data_types" db:"data_types"` // Types of data to export
	Status        string    `json:"status" db:"status"` // pending, processing, completed, failed
	FilePath      string    `json:"file_path" db:"file_path"`
	FileSize      int       `json:"file_size" db:"file_size"`
	DownloadCount int       `json:"download_count" db:"download_count"`
	StartedAt     *time.Time `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage  *string   `json:"error_message" db:"error_message"`
	ExpiresAt     *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ImportJob represents an import job for user data
type ImportJob struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	ImportType        string    `json:"import_type" db:"import_type"` // csv, json, myfitnesspal, strava, etc.
	DataTypes         []string  `json:"data_types" db:"data_types"` // JSON array of data types to import
	Status            string    `json:"status" db:"status"` // pending, processing, completed, failed, validation_error
	FilePath          string    `json:"file_path" db:"file_path"`
	FileSize          int       `json:"file_size" db:"file_size"`
	TotalRecords      int       `json:"total_records" db:"total_records"`
	ProcessedRecords  int       `json:"processed_records" db:"processed_records"`
	SuccessfulRecords int       `json:"successful_records" db:"successful_records"`
	FailedRecords     int       `json:"failed_records" db:"failed_records"`
	StartedAt         *time.Time `json:"started_at" db:"started_at"`
	CompletedAt       *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage      *string   `json:"error_message" db:"error_message"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// DataSyncConfig represents configuration for data synchronization with external providers
type DataSyncConfig struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Provider      string    `json:"provider" db:"provider"` // strava, myfitnesspal, garmin, fitbit, etc.
	AccessToken   string    `json:"access_token" db:"access_token"`
	RefreshToken  *string   `json:"refresh_token" db:"refresh_token"`
	TokenExpiresAt *time.Time `json:"token_expires_at" db:"token_expires_at"`
	SyncEnabled   bool      `json:"sync_enabled" db:"sync_enabled"`
	LastSyncAt    *time.Time `json:"last_sync_at" db:"last_sync_at"`
	SyncFrequency string    `json:"sync_frequency" db:"sync_frequency"` // manual, hourly, daily, weekly
	DataTypes     []string  `json:"data_types" db:"data_types"`
	SyncOptions   JSONValue `json:"sync_options" db:"sync_options"` // JSON object with sync options
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// SyncLog represents a log of a data sync operation
type SyncLog struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	SyncConfigID     int       `json:"sync_config_id" db:"sync_config_id"`
	SyncType         string    `json:"sync_type" db:"sync_type"` // import, export
	Status           string    `json:"status" db:"status"` // success, partial, failed
	RecordsProcessed int       `json:"records_processed" db:"records_processed"`
	RecordsSuccessful int      `json:"records_successful" db:"records_successful"`
	RecordsFailed    int       `json:"records_failed" db:"records_failed"`
	StartTime        time.Time `json:"start_time" db:"start_time"`
	EndTime          *time.Time `json:"end_time" db:"end_time"`
	ErrorDetails     *string   `json:"error_details" db:"error_details"`
	SyncSummary      JSONValue `json:"sync_summary" db:"sync_summary"` // JSON object with sync summary
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// BackupConfig represents configuration for data backup
type BackupConfig struct {
	ID                  int       `json:"id" db:"id"`
	UserID              int       `json:"user_id" db:"user_id"`
	BackupFrequency     string    `json:"backup_frequency" db:"backup_frequency"` // manual, daily, weekly, monthly
	IncludeWorkouts     bool      `json:"include_workouts" db:"include_workouts"`
	IncludeNutrition    bool      `json:"include_nutrition" db:"include_nutrition"`
	IncludeBodyMetrics  bool      `json:"include_body_metrics" db:"include_body_metrics"`
	IncludeTemplates    bool      `json:"include_templates" db:"include_templates"`
	IncludeSettings     bool      `json:"include_settings" db:"include_settings"`
	IncludeMedia        bool      `json:"include_media" db:"include_media"`
	CompressionEnabled  bool      `json:"compression_enabled" db:"compression_enabled"`
	EncryptionEnabled   bool      `json:"encryption_enabled" db:"encryption_enabled"`
	RetentionDays       int       `json:"retention_days" db:"retention_days"`
	LastBackupAt        *time.Time `json:"last_backup_at" db:"last_backup_at"`
	NextBackupAt        *time.Time `json:"next_backup_at" db:"next_backup_at"`
	IsActive            bool      `json:"is_active" db:"is_active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// FileUpload represents a file uploaded by a user
type FileUpload struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	Filename         string    `json:"filename" db:"filename"`
	OriginalFilename string    `json:"original_filename" db:"original_filename"`
	FilePath         string    `json:"file_path" db:"file_path"`
	FileSize         int       `json:"file_size" db:"file_size"`
	MimeType         string    `json:"mime_type" db:"mime_type"`
	FileHash         string    `json:"file_hash" db:"file_hash"`
	UploadType       string    `json:"upload_type" db:"upload_type"` // import, profile_picture, etc.
	Status           string    `json:"status" db:"status"` // uploaded, processing, processed, deleted
	ProcessedAt      *time.Time `json:"processed_at" db:"processed_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Request models for import/export operations
type ExportDataRequest struct {
	DataTypes    []string `json:"data_types"`
	ExportType   string   `json:"export_type"` // csv, json, backup
	ExportOptions JSONValue `json:"export_options"`
}

type ImportDataRequest struct {
	FilePath    string   `json:"file_path"`
	ImportType  string   `json:"import_type"`
	DataTypes   []string `json:"data_types"`
	ImportOptions JSONValue `json:"import_options"`
}

type APIIntegrationConfig struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"access_token"`
	SyncOptions JSONValue `json:"sync_options"`
}

type JSONValue map[string]interface{}

// Request models for template and program management

// Achievement represents a user's achievement
type Achievement struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Type         string    `json:"type" db:"type"` // strength, consistency, milestone, etc.
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	BadgeURL     string    `json:"badge_url" db:"badge_url"`
	AchievedAt   time.Time `json:"achieved_at" db:"achieved_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Badge represents a badge that can be awarded
type Badge struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	BadgeURL    string    `json:"badge_url" db:"badge_url"`
	Criteria    string    `json:"criteria" db:"criteria"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Streak represents a user's streak in workouts or nutrition
type Streak struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Type       string    `json:"type" db:"type"` // workout, nutrition
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    *time.Time `json:"end_date" db:"end_date"`
	Count      int       `json:"count" db:"count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Points represent points a user earns
type Points struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"` // workout, nutrition, achievement
	Value       int       `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Challenge represents a fitness challenge
type Challenge struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	GoalType     string    `json:"goal_type" db:"goal_type"` // steps, streak, workout, etc.
	GoalValue    int       `json:"goal_value" db:"goal_value"`
	RewardPoints int       `json:"reward_points" db:"reward_points"`
	BadgeID      *int      `json:"badge_id" db:"badge_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserChallenge represents a user's progress in a challenge
type UserChallenge struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	ChallengeID  int       `json:"challenge_id" db:"challenge_id"`
	Progress     int       `json:"progress" db:"progress"`
	Status       string    `json:"status" db:"status"` // ongoing, completed, failed
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Request models for challenges and achievements
type CompleteChallengeRequest struct {
	ChallengeID int    `json:"challenge_id" validate:"required"`
	Status      string `json:"status" validate:"required"` // completed, failed
}

// ScheduledWorkout represents a workout scheduled for a specific date/time
type ScheduledWorkout struct {
	ID                int              `json:"id" db:"id"`
	UserID            int              `json:"user_id" db:"user_id"`
	TemplateID        *int             `json:"template_id" db:"template_id"`
	Title             string           `json:"title" db:"title"`
	Description       string           `json:"description" db:"description"`
	ScheduledDate     time.Time        `json:"scheduled_date" db:"scheduled_date"`
	ScheduledTime     *string          `json:"scheduled_time" db:"scheduled_time"` // Format: HH:MM
	EstimatedDuration int              `json:"estimated_duration" db:"estimated_duration"` // minutes
	Status            string           `json:"status" db:"status"` // scheduled, completed, skipped, cancelled
	WorkoutID         *int             `json:"workout_id" db:"workout_id"`
	ReminderSent      bool             `json:"reminder_sent" db:"reminder_sent"`
	Notes             string           `json:"notes" db:"notes"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
	// Related objects
	Template          *WorkoutTemplate `json:"template,omitempty"`
	Workout           *Workout         `json:"workout,omitempty"`
	Reminders         []WorkoutReminder `json:"reminders,omitempty"`
}

// WorkoutReminder represents a reminder for a scheduled workout
type WorkoutReminder struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ScheduledWorkoutID int       `json:"scheduled_workout_id" db:"scheduled_workout_id"`
	ReminderType       string    `json:"reminder_type" db:"reminder_type"` // email, push, sms
	Message            string    `json:"message" db:"message"`
	ScheduledFor       time.Time `json:"scheduled_for" db:"scheduled_for"`
	Status             string    `json:"status" db:"status"` // pending, sent, failed, cancelled
	SentAt             *time.Time `json:"sent_at" db:"sent_at"`
	ErrorMessage       *string   `json:"error_message" db:"error_message"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// RestDayRecommendation represents an AI recommendation for rest days
type RestDayRecommendation struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	RecommendedDate    time.Time `json:"recommended_date" db:"recommended_date"`
	Reason             string    `json:"reason" db:"reason"` // high_volume, consecutive_days, muscle_group_fatigue
	IntensityScore     float64   `json:"intensity_score" db:"intensity_score"` // 0-10 scale
	VolumeLoad         float64   `json:"volume_load" db:"volume_load"` // Total volume from recent workouts
	ConsecutiveDays    int       `json:"consecutive_days" db:"consecutive_days"`
	MuscleGroupsWorked string    `json:"muscle_groups_worked" db:"muscle_groups_worked"` // JSON array
	Status             string    `json:"status" db:"status"` // suggested, accepted, ignored, overridden
	UserResponse       *string   `json:"user_response" db:"user_response"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// DeloadRecommendation represents an AI recommendation for deload weeks
type DeloadRecommendation struct {
	ID                        int        `json:"id" db:"id"`
	UserID                    int        `json:"user_id" db:"user_id"`
	RecommendedStartDate      time.Time  `json:"recommended_start_date" db:"recommended_start_date"`
	RecommendedEndDate        time.Time  `json:"recommended_end_date" db:"recommended_end_date"`
	Reason                    string     `json:"reason" db:"reason"` // fatigue_accumulation, plateau, overreaching, scheduled
	VolumeReductionPercent    int        `json:"volume_reduction_percent" db:"volume_reduction_percent"`
	IntensityReductionPercent int        `json:"intensity_reduction_percent" db:"intensity_reduction_percent"`
	TriggerMetrics            string     `json:"trigger_metrics" db:"trigger_metrics"` // JSON object
	Status                    string     `json:"status" db:"status"` // suggested, accepted, ignored, active, completed
	UserResponse              *string    `json:"user_response" db:"user_response"`
	StartedAt                 *time.Time `json:"started_at" db:"started_at"`
	CompletedAt               *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
}

// WorkoutCalendarEvent represents an event in the workout calendar
type WorkoutCalendarEvent struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ScheduledWorkoutID *int      `json:"scheduled_workout_id" db:"scheduled_workout_id"`
	RestDayID          *int      `json:"rest_day_id" db:"rest_day_id"`
	DeloadID           *int      `json:"deload_id" db:"deload_id"`
	EventType          string    `json:"event_type" db:"event_type"` // workout, rest_day, deload
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	StartDate          time.Time `json:"start_date" db:"start_date"`
	EndDate            *time.Time `json:"end_date" db:"end_date"`
	AllDay             bool      `json:"all_day" db:"all_day"`
	Color              string    `json:"color" db:"color"`
	IsRecurring        bool      `json:"is_recurring" db:"is_recurring"`
	RecurrencePattern  *string   `json:"recurrence_pattern" db:"recurrence_pattern"` // JSON object
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Request models for scheduling operations
type CreateScheduledWorkoutRequest struct {
	TemplateID        *int   `json:"template_id"`
	Title             string `json:"title" validate:"required"`
	Description       string `json:"description"`
	ScheduledDate     string `json:"scheduled_date" validate:"required"` // YYYY-MM-DD
	ScheduledTime     string `json:"scheduled_time"` // HH:MM
	EstimatedDuration int    `json:"estimated_duration"`
	Notes             string `json:"notes"`
}

type UpdateScheduledWorkoutRequest struct {
	Title             string `json:"title" validate:"required"`
	Description       string `json:"description"`
	ScheduledDate     string `json:"scheduled_date" validate:"required"`
	ScheduledTime     string `json:"scheduled_time"`
	EstimatedDuration int    `json:"estimated_duration"`
	Status            string `json:"status"`
	Notes             string `json:"notes"`
}

type CreateReminderRequest struct {
	ReminderType    string `json:"reminder_type" validate:"required"` // email, push, sms
	Message         string `json:"message" validate:"required"`
	MinutesBefore   int    `json:"minutes_before" validate:"required"` // Minutes before workout
}

type RestDayResponseRequest struct {
	Status       string `json:"status" validate:"required"` // accepted, ignored, overridden
	UserResponse string `json:"user_response"`
}

type DeloadResponseRequest struct {
	Status       string `json:"status" validate:"required"` // accepted, ignored, overridden
	UserResponse string `json:"user_response"`
}

// Calendar view models
type CalendarView struct {
	StartDate time.Time               `json:"start_date"`
	EndDate   time.Time               `json:"end_date"`
	Events    []WorkoutCalendarEvent  `json:"events"`
	Workouts  []ScheduledWorkout      `json:"workouts"`
	RestDays  []RestDayRecommendation `json:"rest_days"`
	Deloads   []DeloadRecommendation  `json:"deloads"`
}

type MonthlyCalendarView struct {
	Year   int                     `json:"year"`
	Month  int                     `json:"month"`
	Days   []CalendarDay           `json:"days"`
	Events []WorkoutCalendarEvent `json:"events"`
}

type CalendarDay struct {
	Date            time.Time             `json:"date"`
	IsCurrentMonth  bool                  `json:"is_current_month"`
	IsToday         bool                  `json:"is_today"`
	Workouts        []ScheduledWorkout    `json:"workouts"`
	RestDay         *RestDayRecommendation `json:"rest_day,omitempty"`
	Deload          *DeloadRecommendation  `json:"deload,omitempty"`
	WorkoutCount    int                   `json:"workout_count"`
}

// Analytics for planning
type PlanningAnalytics struct {
	TotalScheduled       int     `json:"total_scheduled"`
	CompletedWorkouts    int     `json:"completed_workouts"`
	SkippedWorkouts      int     `json:"skipped_workouts"`
	CompletionRate       float64 `json:"completion_rate"`
	AverageRestDays      float64 `json:"average_rest_days"`
	ConsistencyScore     float64 `json:"consistency_score"`
	UpcomingWorkouts     int     `json:"upcoming_workouts"`
	OverdueWorkouts      int     `json:"overdue_workouts"`
	ActiveDeloads        int     `json:"active_deloads"`
	RestDayCompliance    float64 `json:"rest_day_compliance"`
	DeloadCompliance     float64 `json:"deload_compliance"`
}

// Recommendation engine data structures
type RecommendationMetrics struct {
	UserID              int                  `json:"user_id"`
	RecentWorkouts      []Workout           `json:"recent_workouts"`
	WeeklyVolume        []WeeklyVolumeData  `json:"weekly_volume"`
	ConsecutiveDays     int                 `json:"consecutive_days"`
	LastRestDay         *time.Time          `json:"last_rest_day"`
	MuscleGroupFatigue  map[string]float64  `json:"muscle_group_fatigue"`
	IntensityTrend      string              `json:"intensity_trend"`
	PerformanceMetrics  PerformanceMetrics  `json:"performance_metrics"`
}

type WeeklyVolumeData struct {
	WeekStart   time.Time `json:"week_start"`
	TotalVolume float64   `json:"total_volume"`
	Workouts    int       `json:"workouts"`
	AvgRPE      float64   `json:"avg_rpe"`
}

type PerformanceMetrics struct {
	VolumeDecline     float64 `json:"volume_decline"`
	StrengthPlateau   bool    `json:"strength_plateau"`
	FatigueIndicators int     `json:"fatigue_indicators"`
	RecoveryScore     float64 `json:"recovery_score"`
}

// Request models for template and program management
type CreateWorkoutTemplateRequest struct {
	Name        string                         `json:"name" validate:"required"`
	Description string                         `json:"description"`
	Exercises   []CreateTemplateExerciseRequest `json:"exercises"`
}

type CreateTemplateExerciseRequest struct {
	Name         string  `json:"name" validate:"required"`
	Category     string  `json:"category" validate:"required"`
	TargetSets   int     `json:"target_sets"`
	TargetReps   int     `json:"target_reps"`
	TargetWeight float64 `json:"target_weight"`
	RestTime     int     `json:"rest_time"`
	Notes        string  `json:"notes"`
}

type UpdateWorkoutTemplateRequest struct {
	Name        string                         `json:"name" validate:"required"`
	Description string                         `json:"description"`
	Exercises   []CreateTemplateExerciseRequest `json:"exercises"`
}

type CreateWorkoutProgramRequest struct {
	Name          string                         `json:"name" validate:"required"`
	Description   string                         `json:"description"`
	Difficulty    string                         `json:"difficulty"`
	DurationWeeks int                            `json:"duration_weeks"`
	Goal          string                         `json:"goal"`
	IsPublic      bool                           `json:"is_public"`
	Templates     []CreateProgramTemplateRequest `json:"templates"`
}

type CreateProgramTemplateRequest struct {
	TemplateID int `json:"template_id" validate:"required"`
	DayOfWeek  int `json:"day_of_week"`
	WeekNumber int `json:"week_number"`
	OrderIndex int `json:"order_index"`
}

type UpdateWorkoutProgramRequest struct {
	Name          string                         `json:"name" validate:"required"`
	Description   string                         `json:"description"`
	Difficulty    string                         `json:"difficulty"`
	DurationWeeks int                            `json:"duration_weeks"`
	Goal          string                         `json:"goal"`
	IsPublic      bool                           `json:"is_public"`
	Templates     []CreateProgramTemplateRequest `json:"templates"`
}

// Template sharing request
type ShareTemplateRequest struct {
	SharedWithID int    `json:"shared_with_id"`
	Permission   string `json:"permission"` // "view" or "edit"
}

// ========== TEMPLATE-BASED WORKOUT CREATION MODELS ==========

// Create workout from template request
type CreateWorkoutFromTemplateRequest struct {
	TemplateID     int                      `json:"template_id"`
	Name           string                   `json:"name"`           // Optional: override template name
	Date           string                   `json:"date"`           // Workout date (YYYY-MM-DD)
	Notes          string                   `json:"notes"`          // Optional workout notes
	Customizations []ExerciseCustomization `json:"customizations"` // Optional exercise customizations
}

// Exercise customization for template-based workouts
type ExerciseCustomization struct {
	ExerciseName string  `json:"exercise_name"`
	TargetSets   int     `json:"target_sets"`
	TargetReps   int     `json:"target_reps"`
	TargetWeight float64 `json:"target_weight"`
	Skip         bool    `json:"skip"` // Skip this exercise
}

// Template usage tracking
type TemplateUsage struct {
	ID         int       `json:"id"`
	TemplateID int       `json:"template_id"`
	UserID     int       `json:"user_id"`
	WorkoutID  int       `json:"workout_id"`
	UsedAt     time.Time `json:"used_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// Template-based workout creation response
type WorkoutFromTemplateResponse struct {
	Workout      Workout       `json:"workout"`
	TemplateUsed WorkoutTemplate `json:"template_used"`
	Message      string        `json:"message"`
}
