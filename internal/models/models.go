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
