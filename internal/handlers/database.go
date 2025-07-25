package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"workout-tracker/internal/models"
)

// getRecentWorkouts returns the most recent workouts
func (h *Handler) getRecentWorkouts(limit int) ([]models.Workout, error) {
	query := `
		SELECT id, name, date, duration, notes, created_at, updated_at
		FROM workouts
		ORDER BY date DESC, created_at DESC
		LIMIT ?
	`
	
	rows, err := h.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var w models.Workout
		err := rows.Scan(&w.ID, &w.Name, &w.Date, &w.Duration, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

// getAllWorkouts returns all workouts
func (h *Handler) getAllWorkouts() ([]models.Workout, error) {
	query := `
		SELECT id, name, date, duration, notes, created_at, updated_at
		FROM workouts
		ORDER BY date DESC, created_at DESC
	`
	
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var w models.Workout
		err := rows.Scan(&w.ID, &w.Name, &w.Date, &w.Duration, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

// getWorkoutByID returns a workout by ID with its exercises and sets
func (h *Handler) getWorkoutByID(id int) (models.Workout, error) {
	var w models.Workout
	
	// Get workout
	query := `
		SELECT id, name, date, duration, notes, created_at, updated_at
		FROM workouts
		WHERE id = ?
	`
	
	err := h.db.QueryRow(query, id).Scan(&w.ID, &w.Name, &w.Date, &w.Duration, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return w, fmt.Errorf("workout not found")
		}
		return w, err
	}

	// Get exercises for this workout
	exercises, err := h.getExercisesByWorkoutID(id)
	if err != nil {
		return w, err
	}
	w.Exercises = exercises

	return w, nil
}

// getExercisesByWorkoutID returns exercises for a workout
func (h *Handler) getExercisesByWorkoutID(workoutID int) ([]models.Exercise, error) {
	query := `
		SELECT id, workout_id, name, category, created_at, updated_at
		FROM exercises
		WHERE workout_id = ?
		ORDER BY created_at ASC
	`
	
	rows, err := h.db.Query(query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		err := rows.Scan(&e.ID, &e.WorkoutID, &e.Name, &e.Category, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		
		// Get sets for this exercise
		sets, err := h.getSetsByExerciseID(e.ID)
		if err != nil {
			return nil, err
		}
		e.Sets = sets
		
		exercises = append(exercises, e)
	}

	return exercises, nil
}

// getSetsByExerciseID returns sets for an exercise
func (h *Handler) getSetsByExerciseID(exerciseID int) ([]models.Set, error) {
	query := `
		SELECT id, exercise_id, set_number, reps, weight, distance, duration, rest_time, created_at, updated_at
		FROM sets
		WHERE exercise_id = ?
		ORDER BY set_number ASC
	`
	
	rows, err := h.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []models.Set
	for rows.Next() {
		var s models.Set
		err := rows.Scan(&s.ID, &s.ExerciseID, &s.SetNumber, &s.Reps, &s.Weight, &s.Distance, &s.Duration, &s.RestTime, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sets = append(sets, s)
	}

	return sets, nil
}

// createWorkout creates a new workout and returns its ID
func (h *Handler) createWorkout(workout models.Workout) (int, error) {
	query := `
		INSERT INTO workouts (name, date, duration, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, workout.Name, workout.Date, workout.Duration, workout.Notes, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// createExercise creates a new exercise and returns its ID
func (h *Handler) createExercise(exercise models.Exercise) (int, error) {
	query := `
		INSERT INTO exercises (workout_id, name, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, exercise.WorkoutID, exercise.Name, exercise.Category, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getAllPredefinedExercises returns all predefined exercises
func (h *Handler) getAllPredefinedExercises() ([]models.PredefinedExercise, error) {
	query := `
		SELECT id, name, category, description, video_url, instructions, tips, muscle_groups, equipment, difficulty, image_url, created_at, updated_at
		FROM predefined_exercises
		ORDER BY category, name
	`
	
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.PredefinedExercise
	for rows.Next() {
		var e models.PredefinedExercise
		err := rows.Scan(&e.ID, &e.Name, &e.Category, &e.Description, &e.VideoURL, &e.Instructions, &e.Tips, &e.MuscleGroups, &e.Equipment, &e.Difficulty, &e.ImageURL, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}

	return exercises, nil
}

// getPredefinedExercisesByCategory returns predefined exercises by category
func (h *Handler) getPredefinedExercisesByCategory(category string) ([]models.PredefinedExercise, error) {
	query := `
		SELECT id, name, category, description, video_url, instructions, tips, muscle_groups, equipment, difficulty, image_url, created_at, updated_at
		FROM predefined_exercises
		WHERE category = ?
		ORDER BY name
	`
	
	rows, err := h.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.PredefinedExercise
	for rows.Next() {
		var e models.PredefinedExercise
		err := rows.Scan(&e.ID, &e.Name, &e.Category, &e.Description, &e.VideoURL, &e.Instructions, &e.Tips, &e.MuscleGroups, &e.Equipment, &e.Difficulty, &e.ImageURL, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}

	return exercises, nil
}

// createPredefinedExercise creates a new predefined exercise and returns its ID
func (h *Handler) createPredefinedExercise(exercise models.PredefinedExercise) (int, error) {
	query := `
		INSERT INTO predefined_exercises (name, category, description, video_url, instructions, tips, muscle_groups, equipment, difficulty, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, exercise.Name, exercise.Category, exercise.Description, exercise.VideoURL, exercise.Instructions, exercise.Tips, exercise.MuscleGroups, exercise.Equipment, exercise.Difficulty, exercise.ImageURL, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// createSet creates a new set and returns its ID
func (h *Handler) createSet(set models.Set) (int, error) {
	query := `
		INSERT INTO sets (exercise_id, set_number, reps, weight, distance, duration, rest_time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, set.ExerciseID, set.SetNumber, set.Reps, set.Weight, set.Distance, set.Duration, set.RestTime, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// deleteWorkout deletes a workout and all associated exercises and sets
func (h *Handler) deleteWorkout(id int) error {
	query := `DELETE FROM workouts WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

// updateWorkout updates an existing workout
func (h *Handler) updateWorkout(workout models.Workout) error {
	query := `
		UPDATE workouts 
		SET name = ?, date = ?, duration = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := h.db.Exec(query, workout.Name, workout.Date, workout.Duration, workout.Notes, time.Now(), workout.ID)
	return err
}

// deleteExercise deletes an exercise and all associated sets
func (h *Handler) deleteExercise(id int) error {
	query := `DELETE FROM exercises WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

// updateExercise updates an existing exercise
func (h *Handler) updateExercise(exercise models.Exercise) error {
	query := `
		UPDATE exercises 
		SET name = ?, category = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := h.db.Exec(query, exercise.Name, exercise.Category, time.Now(), exercise.ID)
	return err
}

// updateSet updates an existing set
func (h *Handler) updateSet(set models.Set) error {
	query := `
		UPDATE sets 
		SET set_number = ?, reps = ?, weight = ?, distance = ?, duration = ?, rest_time = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := h.db.Exec(query, set.SetNumber, set.Reps, set.Weight, set.Distance, set.Duration, set.RestTime, time.Now(), set.ID)
	return err
}

// deleteSet deletes a specific set
func (h *Handler) deleteSet(id int) error {
	query := `DELETE FROM sets WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

// getWorkoutStats returns basic statistics about workouts
func (h *Handler) getWorkoutStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total workouts
	var totalWorkouts int
	err := h.db.QueryRow("SELECT COUNT(*) FROM workouts").Scan(&totalWorkouts)
	if err != nil {
		return nil, err
	}
	stats["total_workouts"] = totalWorkouts
	
	// Workouts this week
	var thisWeekWorkouts int
	err = h.db.QueryRow(`
		SELECT COUNT(*) FROM workouts 
		WHERE date >= date('now', '-7 days')
	`).Scan(&thisWeekWorkouts)
	if err != nil {
		return nil, err
	}
	stats["this_week_workouts"] = thisWeekWorkouts
	
	// Average workout duration
	var avgDuration sql.NullFloat64
	err = h.db.QueryRow(`
		SELECT AVG(duration) FROM workouts 
		WHERE duration > 0
	`).Scan(&avgDuration)
	if err != nil {
		return nil, err
	}
	if avgDuration.Valid {
		stats["avg_duration"] = int(avgDuration.Float64)
	} else {
		stats["avg_duration"] = 0
	}
	
	return stats, nil
}

// getAnalyticsData returns comprehensive analytics data for a user over specified days
func (h *Handler) getAnalyticsData(userID int, days int) (*models.AnalyticsData, error) {
	analytics := &models.AnalyticsData{}

	// Calculate date range
	startDate := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	// Generate date labels
	dates := make([]string, 0, days)
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+1+i)
		dates = append(dates, date.Format("01/02"))
	}
	analytics.Dates = dates

	// Get basic stats
	var totalWorkouts int
	h.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE date >= ? AND date <= ?", startDate, endDate).Scan(&totalWorkouts)
	analytics.TotalWorkouts = totalWorkouts

	// Get workout frequency by day
	workoutCounts := make([]int, days)
	rows, err := h.db.Query(`
		SELECT DATE(date) as workout_date, COUNT(*) as count 
		FROM workouts 
		WHERE date >= ? AND date <= ?
		GROUP BY DATE(date)
		ORDER BY workout_date
	`, startDate, endDate)
	if err == nil {
		defer rows.Close()
		workoutMap := make(map[string]int)
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			workoutMap[date] = count
		}
		// Fill in the workout counts for each day
		for i := 0; i < days; i++ {
			date := time.Now().AddDate(0, 0, -days+1+i).Format("2006-01-02")
			workoutCounts[i] = workoutMap[date]
		}
	}
	analytics.WorkoutCounts = workoutCounts

	// Get workout durations by day
	durations := make([]int, days)
	rows, err = h.db.Query(`
		SELECT DATE(date) as workout_date, AVG(duration) as avg_duration 
		FROM workouts 
		WHERE date >= ? AND date <= ? AND duration > 0
		GROUP BY DATE(date)
		ORDER BY workout_date
	`, startDate, endDate)
	if err == nil {
		defer rows.Close()
		durationMap := make(map[string]int)
		for rows.Next() {
			var date string
			var avgDuration sql.NullFloat64
			rows.Scan(&date, &avgDuration)
			if avgDuration.Valid {
				durationMap[date] = int(avgDuration.Float64)
			}
		}
		// Fill in the durations for each day
		for i := 0; i < days; i++ {
			date := time.Now().AddDate(0, 0, -days+1+i).Format("2006-01-02")
			durations[i] = durationMap[date]
		}
	}
	analytics.Durations = durations

	// Get daily calories (if user has meal data)
	dailyCalories := make([]int, days)
	rows, err = h.db.Query(`
		SELECT DATE(date) as meal_date, SUM(calories) as total_calories 
		FROM meals 
		WHERE user_id = ? AND date >= ? AND date <= ?
		GROUP BY DATE(date)
		ORDER BY meal_date
	`, userID, startDate, endDate)
	if err == nil {
		defer rows.Close()
		caloriesMap := make(map[string]int)
		totalCaloriesSum := 0
		daysWithData := 0
		for rows.Next() {
			var date string
			var calories int
			rows.Scan(&date, &calories)
			caloriesMap[date] = calories
			totalCaloriesSum += calories
			daysWithData++
		}
		// Fill in the calories for each day
		for i := 0; i < days; i++ {
			date := time.Now().AddDate(0, 0, -days+1+i).Format("2006-01-02")
			dailyCalories[i] = caloriesMap[date]
		}
		// Calculate average calories
		if daysWithData > 0 {
			analytics.AvgCalories = totalCaloriesSum / daysWithData
		}
	}
	analytics.DailyCalories = dailyCalories

	// Get weight progress data
	weightDates := []string{}
	weights := []float64{}
	rows, err = h.db.Query(`
		SELECT date, weight 
		FROM body_weights 
		WHERE user_id = ? AND date >= ? AND date <= ?
		ORDER BY date
	`, userID, startDate, endDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date time.Time
			var weight float64
			rows.Scan(&date, &weight)
			weightDates = append(weightDates, date.Format("01/02"))
			weights = append(weights, weight)
		}
	}
	analytics.WeightDates = weightDates
	analytics.Weights = weights

	// Calculate weight change
	if len(weights) >= 2 {
		firstWeight := weights[0]
		lastWeight := weights[len(weights)-1]
		analytics.WeightChange = lastWeight - firstWeight
	}

	// Get exercise categories distribution
	exerciseCategories := &models.ExerciseCategoryData{}
	rows, err = h.db.Query(`
		SELECT e.category, COUNT(*) as count
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.category
		ORDER BY count DESC
	`, startDate, endDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var category string
			var count int
			rows.Scan(&category, &count)
			exerciseCategories.Labels = append(exerciseCategories.Labels, category)
			exerciseCategories.Values = append(exerciseCategories.Values, count)
		}
	}
	analytics.ExerciseCategories = exerciseCategories

	// Get macronutrient breakdown
	macros := &models.MacroData{}
	var totalProtein, totalCarbs, totalFat sql.NullFloat64
	err = h.db.QueryRow(`
		SELECT AVG(protein), AVG(carbs), AVG(fat)
		FROM meals 
		WHERE user_id = ? AND date >= ? AND date <= ?
	`, userID, startDate, endDate).Scan(&totalProtein, &totalCarbs, &totalFat)
	if err == nil {
		if totalProtein.Valid {
			macros.Protein = totalProtein.Float64
		}
		if totalCarbs.Valid {
			macros.Carbs = totalCarbs.Float64
		}
		if totalFat.Valid {
			macros.Fat = totalFat.Float64
		}
	}
	analytics.Macros = macros

	// Calculate workout streak
	workoutStreak := 0
	currentDate := time.Now()
	for i := 0; i < 365; i++ { // Check up to a year back
		checkDate := currentDate.AddDate(0, 0, -i).Format("2006-01-02")
		var count int
		err := h.db.QueryRow("SELECT COUNT(*) FROM workouts WHERE DATE(date) = ?", checkDate).Scan(&count)
		if err != nil || count == 0 {
			break
		}
		workoutStreak++
	}
	analytics.WorkoutStreak = workoutStreak

	// Enhanced progress tracking
	err = h.calculateDetailedProgress(analytics, userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to calculate detailed progress: %v", err)
		// Don't fail the entire request, just skip detailed analytics
	}

	return analytics, nil
}

// createUser creates a new user and returns its ID
func (h *Handler) createUser(user models.User) (int, error) {
	query := `
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, user.Username, user.Email, user.PasswordHash, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getUserByUsername returns a user by username
func (h *Handler) getUserByUsername(username string) (models.User, error) {
	var u models.User
	
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	
	err := h.db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// createMeal creates a new meal entry and returns its ID
func (h *Handler) createMeal(meal models.Meal) (int, error) {
	query := `
		INSERT INTO meals (user_id, name, calories, protein, carbs, fat, date, meal_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, meal.UserID, meal.Name, meal.Calories, meal.Protein, meal.Carbs, meal.Fat, meal.Date, meal.MealType, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getMealsByUserAndDate returns meals for a user on a specific date
func (h *Handler) getMealsByUserAndDate(userID int, date time.Time) ([]models.Meal, error) {
	query := `
		SELECT id, user_id, name, calories, protein, carbs, fat, date, meal_type, created_at, updated_at
		FROM meals
		WHERE user_id = ? AND DATE(date) = DATE(?)
		ORDER BY created_at ASC
	`
	
	rows, err := h.db.Query(query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var m models.Meal
		err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Calories, &m.Protein, &m.Carbs, &m.Fat, &m.Date, &m.MealType, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}

	return meals, nil
}

// createBodyWeight creates a new body weight entry and returns its ID
func (h *Handler) createBodyWeight(bodyWeight models.BodyWeight) (int, error) {
	query := `
		INSERT INTO body_weights (user_id, weight, unit, date, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, bodyWeight.UserID, bodyWeight.Weight, bodyWeight.Unit, bodyWeight.Date, bodyWeight.Notes, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getBodyWeightsByUser returns body weight entries for a user
func (h *Handler) getBodyWeightsByUser(userID int) ([]models.BodyWeight, error) {
	query := `
		SELECT id, user_id, weight, unit, date, notes, created_at, updated_at
		FROM body_weights
		WHERE user_id = ?
		ORDER BY date DESC
	`
	
	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []models.BodyWeight
	for rows.Next() {
		var w models.BodyWeight
		err := rows.Scan(&w.ID, &w.UserID, &w.Weight, &w.Unit, &w.Date, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		weights = append(weights, w)
	}

	return weights, nil
}

// createBodyFat creates a new body fat entry and returns its ID
func (h *Handler) createBodyFat(bodyFat models.BodyFat) (int, error) {
	query := `
		INSERT INTO body_fats (user_id, body_fat_pct, date, measurement, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, bodyFat.UserID, bodyFat.BodyFatPct, bodyFat.Date, bodyFat.Measurement, bodyFat.Notes, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getBodyFatsByUser returns body fat entries for a user
func (h *Handler) getBodyFatsByUser(userID int) ([]models.BodyFat, error) {
	query := `
		SELECT id, user_id, body_fat_pct, date, measurement, notes, created_at, updated_at
		FROM body_fats
		WHERE user_id = ?
		ORDER BY date DESC
	`
	
	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bodyFats []models.BodyFat
	for rows.Next() {
		var bf models.BodyFat
		err := rows.Scan(&bf.ID, &bf.UserID, &bf.BodyFatPct, &bf.Date, &bf.Measurement, &bf.Notes, &bf.CreatedAt, &bf.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bodyFats = append(bodyFats, bf)
	}

	return bodyFats, nil
}

// createBodyMeasurement creates a new body measurement entry and returns its ID
func (h *Handler) createBodyMeasurement(measurement models.BodyMeasurement) (int, error) {
	query := `
		INSERT INTO body_measurements (user_id, measurement, value, unit, date, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, measurement.UserID, measurement.Measurement, measurement.Value, measurement.Unit, measurement.Date, measurement.Notes, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// getBodyMeasurementsByUser returns body measurements for a user
func (h *Handler) getBodyMeasurementsByUser(userID int) ([]models.BodyMeasurement, error) {
	query := `
		SELECT id, user_id, measurement, value, unit, date, notes, created_at, updated_at
		FROM body_measurements
		WHERE user_id = ?
		ORDER BY measurement, date DESC
	`
	
	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []models.BodyMeasurement
	for rows.Next() {
		var m models.BodyMeasurement
		err := rows.Scan(&m.ID, &m.UserID, &m.Measurement, &m.Value, &m.Unit, &m.Date, &m.Notes, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, m)
	}

	return measurements, nil
}

// User Settings Database Methods
func (h *Handler) getUserSettings(userID int) (models.UserSettings, error) {
	query := `SELECT * FROM user_settings WHERE user_id = ?`
	var settings models.UserSettings
	err := h.db.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Theme, &settings.Timezone,
		&settings.WeightUnit, &settings.DistanceUnit, &settings.DateFormat,
		&settings.Notifications, &settings.PrivacyMode, &settings.AutoLogout,
		&settings.Language, &settings.CreatedAt, &settings.UpdatedAt)
	if err != nil {
		// If no settings exist, create default settings
		defaultSettings := models.UserSettings{
			UserID:        userID,
			Theme:         "light",
			Timezone:      "UTC",
			WeightUnit:    "lbs",
			DistanceUnit:  "miles",
			DateFormat:    "MM/DD/YYYY",
			Notifications: true,
			PrivacyMode:   false,
			AutoLogout:    0,
			Language:      "en",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		_, createErr := h.createUserSettings(defaultSettings)
		if createErr != nil {
			return settings, createErr
		}
		return defaultSettings, nil
	}
	return settings, nil
}

func (h *Handler) createUserSettings(settings models.UserSettings) (int, error) {
	query := `INSERT INTO user_settings (user_id, theme, timezone, weight_unit, distance_unit, date_format, notifications, privacy_mode, auto_logout, language, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query, settings.UserID, settings.Theme, settings.Timezone, settings.WeightUnit, 
		settings.DistanceUnit, settings.DateFormat, settings.Notifications, settings.PrivacyMode, 
		settings.AutoLogout, settings.Language, settings.CreatedAt, settings.UpdatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (h *Handler) updateUserSettings(userID int, settings models.UpdateSettingsRequest) error {
	query := `UPDATE user_settings SET theme = ?, timezone = ?, weight_unit = ?, distance_unit = ?, date_format = ?, 
			notifications = ?, privacy_mode = ?, auto_logout = ?, language = ?, updated_at = ? WHERE user_id = ?`
	_, err := h.db.Exec(query, settings.Theme, settings.Timezone, settings.WeightUnit, settings.DistanceUnit,
		settings.DateFormat, settings.Notifications, settings.PrivacyMode, settings.AutoLogout, 
		settings.Language, time.Now(), userID)
	return err
}

func (h *Handler) updateUserProfile(userID int, profile models.UpdateProfileRequest) error {
	query := `UPDATE users SET username = ?, email = ?, full_name = ?, bio = ?, updated_at = ? WHERE id = ?`
	_, err := h.db.Exec(query, profile.Username, profile.Email, profile.FullName, profile.Bio, time.Now(), userID)
	return err
}

func (h *Handler) updateUserPassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := h.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

func (h *Handler) getUserByID(userID int) (models.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(full_name, '') as full_name, COALESCE(bio, '') as bio, COALESCE(avatar, '') as avatar, COALESCE(is_active, 1) as is_active, created_at, updated_at FROM users WHERE id = ?`
	var user models.User
	err := h.db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Bio, &user.Avatar, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (h *Handler) deleteUserAccount(userID int) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete user data in the correct order
	// Start with the most dependent tables first
	queries := []string{
		`DELETE FROM sets WHERE exercise_id IN (SELECT id FROM exercises WHERE workout_id IN (SELECT id FROM workouts WHERE user_id = ?))`,
		`DELETE FROM exercises WHERE workout_id IN (SELECT id FROM workouts WHERE user_id = ?)`,
		`DELETE FROM workouts WHERE user_id = ?`,
		`DELETE FROM body_measurements WHERE user_id = ?`,
		`DELETE FROM body_fats WHERE user_id = ?`,
		`DELETE FROM body_weights WHERE user_id = ?`,
		`DELETE FROM meals WHERE user_id = ?`,
		`DELETE FROM user_settings WHERE user_id = ?`,
		`DELETE FROM users WHERE id = ?`,
	}
	
	for _, query := range queries {
		_, err := tx.Exec(query, userID)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// calculateDetailedProgress computes enhanced analytics data
func (h *Handler) calculateDetailedProgress(analytics *models.AnalyticsData, userID int, startDate, endDate string) error {
	// Calculate basic set and rep statistics
	var totalSets, totalReps int
	var totalVolume, maxWeight float64

	query := `
		SELECT 
			COUNT(s.id) as total_sets,
			COALESCE(SUM(s.reps), 0) as total_reps,
			COALESCE(SUM(s.weight * s.reps), 0) as total_volume,
			COALESCE(MAX(s.weight), 0) as max_weight
		FROM sets s
		JOIN exercises e ON s.exercise_id = e.id
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ? AND s.weight > 0
	`

	err := h.db.QueryRow(query, startDate, endDate).Scan(&totalSets, &totalReps, &totalVolume, &maxWeight)
	if err != nil {
		return err
	}

	analytics.TotalSets = totalSets
	analytics.TotalReps = totalReps
	analytics.TotalVolume = totalVolume
	analytics.MaxWeight = maxWeight

	// Calculate averages
	if analytics.TotalWorkouts > 0 {
		analytics.AvgVolume = totalVolume / float64(analytics.TotalWorkouts)
		analytics.AvgSetsPerWorkout = float64(totalSets) / float64(analytics.TotalWorkouts)
	}
	if totalSets > 0 {
		analytics.AvgRepsPerSet = float64(totalReps) / float64(totalSets)
	}

	// Get personal records
	personalRecords, err := h.getPersonalRecords(userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get personal records: %v", err)
	} else {
		analytics.PersonalRecords = personalRecords
	}

	// Get strength progress for top exercises
	strengthProgress, err := h.getStrengthProgress(userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get strength progress: %v", err)
	} else {
		analytics.StrengthProgress = strengthProgress
	}

	// Get workout intensity over time
	workoutIntensity, err := h.getWorkoutIntensity(userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get workout intensity: %v", err)
	} else {
		analytics.WorkoutIntensity = workoutIntensity
	}

	// Get most popular exercises
	mostPopularExercises, err := h.getMostPopularExercises(userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get popular exercises: %v", err)
	} else {
		analytics.MostPopularExercises = mostPopularExercises
	}

	// Calculate progress trends
	progressTrends, err := h.getProgressTrends(analytics, userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get progress trends: %v", err)
	} else {
		analytics.ProgressTrends = progressTrends
	}

	return nil
}

// getPersonalRecords returns personal records for major exercises
func (h *Handler) getPersonalRecords(userID int, startDate, endDate string) ([]models.PersonalRecord, error) {
	query := `
		SELECT 
			e.name,
			MAX(s.weight) as max_weight,
			s.reps,
			(s.weight * s.reps) as volume,
			(s.weight * (1 + s.reps/30.0)) as one_rep_max,
			w.date
		FROM sets s
		JOIN exercises e ON s.exercise_id = e.id
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ? AND s.weight > 0
		GROUP BY e.name
		HAVING MAX(s.weight) = s.weight
		ORDER BY max_weight DESC
		LIMIT 10
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.PersonalRecord
	for rows.Next() {
		var record models.PersonalRecord
		err := rows.Scan(&record.ExerciseName, &record.Weight, &record.Reps, &record.Volume, &record.OneRepMax, &record.Date)
		if err != nil {
			return nil, err
		}
		// Mark as new if achieved in the last 7 days
		if time.Since(record.Date).Hours() < 168 {
			record.IsNew = true
		}
		records = append(records, record)
	}

	return records, nil
}

// getStrengthProgress returns strength progression data for top exercises
func (h *Handler) getStrengthProgress(userID int, startDate, endDate string) ([]models.StrengthProgress, error) {
	// Get top 5 exercises by frequency
	topExercisesQuery := `
		SELECT e.name, e.category, COUNT(*) as frequency
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.name, e.category
		ORDER BY frequency DESC
		LIMIT 5
	`

	rows, err := h.db.Query(topExercisesQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progressData []models.StrengthProgress
	for rows.Next() {
		var exerciseName, category string
		var frequency int
		err := rows.Scan(&exerciseName, &category, &frequency)
		if err != nil {
			continue
		}

		// Get progression data points for this exercise
		dataPointsQuery := `
			SELECT 
				w.date,
				MAX(s.weight) as max_weight,
				MAX(s.reps) as max_reps,
				SUM(s.weight * s.reps) as volume,
				MAX(s.weight * (1 + s.reps/30.0)) as one_rep_max
			FROM sets s
			JOIN exercises e ON s.exercise_id = e.id
			JOIN workouts w ON e.workout_id = w.id
			WHERE e.name = ? AND w.date >= ? AND w.date <= ? AND s.weight > 0
			GROUP BY DATE(w.date)
			ORDER BY w.date
		`

		dataRows, err := h.db.Query(dataPointsQuery, exerciseName, startDate, endDate)
		if err != nil {
			continue
		}

		var dataPoints []models.StrengthDataPoint
		for dataRows.Next() {
			var point models.StrengthDataPoint
			err := dataRows.Scan(&point.Date, &point.MaxWeight, &point.MaxReps, &point.Volume, &point.OneRepMax)
			if err != nil {
				continue
			}
			dataPoints = append(dataPoints, point)
		}
		dataRows.Close()

		// Calculate trend
		trend := "stable"
		trendPercent := 0.0
		if len(dataPoints) > 1 {
			firstWeight := dataPoints[0].MaxWeight
			lastWeight := dataPoints[len(dataPoints)-1].MaxWeight
			if firstWeight > 0 {
				trendPercent = ((lastWeight - firstWeight) / firstWeight) * 100
				if trendPercent > 5 {
					trend = "improving"
				} else if trendPercent < -5 {
					trend = "declining"
				}
			}
		}

		progressData = append(progressData, models.StrengthProgress{
			ExerciseName: exerciseName,
			Category:     category,
			DataPoints:   dataPoints,
			Trend:        trend,
			TrendPercent: trendPercent,
		})
	}

	return progressData, nil
}

// getWorkoutIntensity returns workout intensity metrics over time
func (h *Handler) getWorkoutIntensity(userID int, startDate, endDate string) ([]models.WorkoutIntensity, error) {
	query := `
		SELECT 
			w.date,
			COALESCE(SUM(s.weight * s.reps), 0) as total_volume,
			COALESCE(AVG(s.weight), 0) as avg_weight,
			COUNT(s.id) as total_sets,
			COALESCE(SUM(s.reps), 0) as total_reps
		FROM workouts w
		LEFT JOIN exercises e ON w.id = e.workout_id
		LEFT JOIN sets s ON e.id = s.exercise_id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY w.date
		ORDER BY w.date
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var intensityData []models.WorkoutIntensity
	for rows.Next() {
		var intensity models.WorkoutIntensity
		err := rows.Scan(&intensity.Date, &intensity.TotalVolume, &intensity.AvgWeight, &intensity.TotalSets, &intensity.TotalReps)
		if err != nil {
			continue
		}

		// Calculate intensity score (volume normalized by sets)
		if intensity.TotalSets > 0 {
			intensity.IntensityScore = intensity.TotalVolume / float64(intensity.TotalSets)
		}

		intensityData = append(intensityData, intensity)
	}

	return intensityData, nil
}

// getMostPopularExercises returns the most frequently performed exercises
func (h *Handler) getMostPopularExercises(userID int, startDate, endDate string) ([]models.ExerciseFrequency, error) {
	query := `
		SELECT 
			e.name,
			e.category,
			COUNT(*) as count,
			MAX(w.date) as last_performed
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.name, e.category
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// First pass: get total count
	totalExercises := 0
	var exercises []struct {
		name         string
		category     string
		count        int
		lastPerformed time.Time
	}

	for rows.Next() {
		var ex struct {
			name         string
			category     string
			count        int
			lastPerformed time.Time
		}
		err := rows.Scan(&ex.name, &ex.category, &ex.count, &ex.lastPerformed)
		if err != nil {
			continue
		}
		totalExercises += ex.count
		exercises = append(exercises, ex)
	}

	// Second pass: calculate percentages
	var frequencyData []models.ExerciseFrequency
	for _, ex := range exercises {
		percentage := 0.0
		if totalExercises > 0 {
			percentage = (float64(ex.count) / float64(totalExercises)) * 100
		}

		frequencyData = append(frequencyData, models.ExerciseFrequency{
			ExerciseName:  ex.name,
			Category:      ex.category,
			Count:         ex.count,
			Percentage:    percentage,
			LastPerformed: ex.lastPerformed,
		})
	}

	return frequencyData, nil
}

// getProgressTrends calculates overall progress trends
func (h *Handler) getProgressTrends(analytics *models.AnalyticsData, userID int, startDate, endDate string) (*models.ProgressTrends, error) {
	trends := &models.ProgressTrends{}

	// Calculate date range for comparison (previous period)
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	daysDiff := int(end.Sub(start).Hours() / 24)

	prevStartDate := start.AddDate(0, 0, -daysDiff).Format("2006-01-02")
	prevEndDate := start.AddDate(0, 0, -1).Format("2006-01-02")

	// Compare volume trends
	var currentVolume, previousVolume float64
	volumeQuery := `
		SELECT COALESCE(SUM(s.weight * s.reps), 0) as total_volume
		FROM sets s
		JOIN exercises e ON s.exercise_id = e.id
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
	`

	h.db.QueryRow(volumeQuery, startDate, endDate).Scan(&currentVolume)
	h.db.QueryRow(volumeQuery, prevStartDate, prevEndDate).Scan(&previousVolume)

	if previousVolume > 0 {
		trends.VolumeChangePercent = ((currentVolume - previousVolume) / previousVolume) * 100
		if trends.VolumeChangePercent > 5 {
			trends.Volumetrend = "up"
		} else if trends.VolumeChangePercent < -5 {
			trends.Volumetrend = "down"
		} else {
			trends.Volumetrend = "stable"
		}
	} else {
		trends.Volumetrend = "stable"
	}

	// Compare strength trends (max weight)
	var currentMaxWeight, previousMaxWeight float64
	maxWeightQuery := `
		SELECT COALESCE(MAX(s.weight), 0) as max_weight
		FROM sets s
		JOIN exercises e ON s.exercise_id = e.id
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
	`

	h.db.QueryRow(maxWeightQuery, startDate, endDate).Scan(&currentMaxWeight)
	h.db.QueryRow(maxWeightQuery, prevStartDate, prevEndDate).Scan(&previousMaxWeight)

	if previousMaxWeight > 0 {
		trends.StrengthChangePercent = ((currentMaxWeight - previousMaxWeight) / previousMaxWeight) * 100
		if trends.StrengthChangePercent > 2 {
			trends.StrengthTrend = "up"
		} else if trends.StrengthChangePercent < -2 {
			trends.StrengthTrend = "down"
		} else {
			trends.StrengthTrend = "stable"
		}
	} else {
		trends.StrengthTrend = "stable"
	}

	// Compare frequency trends
	var currentWorkouts, previousWorkouts int
	workoutQuery := `SELECT COUNT(*) FROM workouts WHERE date >= ? AND date <= ?`

	h.db.QueryRow(workoutQuery, startDate, endDate).Scan(&currentWorkouts)
	h.db.QueryRow(workoutQuery, prevStartDate, prevEndDate).Scan(&previousWorkouts)

	if previousWorkouts > 0 {
		trends.FrequencyChangePercent = ((float64(currentWorkouts) - float64(previousWorkouts)) / float64(previousWorkouts)) * 100
		if trends.FrequencyChangePercent > 10 {
			trends.FrequencyTrend = "up"
		} else if trends.FrequencyChangePercent < -10 {
			trends.FrequencyTrend = "down"
		} else {
			trends.FrequencyTrend = "stable"
		}
	} else {
		trends.FrequencyTrend = "stable"
	}

	// Calculate consistency score based on workout frequency
	expectedWorkouts := float64(daysDiff) / 7 * 3 // Assume 3 workouts per week is ideal
	if expectedWorkouts > 0 {
		consistency := math.Min((float64(currentWorkouts)/expectedWorkouts)*100, 100)
		trends.ConsistencyScore = math.Max(consistency, 0)
	}

	return trends, nil
}

// getExerciseProgressChart generates progress data for a specific exercise
func (h *Handler) getExerciseProgressChart(userID int, exerciseName string, startDate, endDate string) (*models.ExerciseProgressChart, error) {
	chart := &models.ExerciseProgressChart{
		ExerciseName: exerciseName,
		TimeRange:   fmt.Sprintf("%s to %s", startDate, endDate),
	}

	// Get basic workout stats for the exercise
	query := `
		SELECT 
			w.id as workout_id,
			w.date,
			MAX(s.weight) as max_weight,
			MAX(s.reps) as max_reps,
			SUM(s.weight * s.reps) as total_volume,
			COUNT(s.id) as sets_count
		FROM sets s
		JOIN exercises e ON s.exercise_id = e.id
		JOIN workouts w ON e.workout_id = w.id
		WHERE e.name = ? AND w.date >= ? AND w.date <= ?
		GROUP BY w.id, w.date
		ORDER BY w.date
	`

	rows, err := h.db.Query(query, exerciseName, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []models.ExerciseProgressPoint
	for rows.Next() {
		var point models.ExerciseProgressPoint
		var dateStr string
		err := rows.Scan(&point.WorkoutID, &dateStr, &point.MaxWeight, &point.MaxReps, &point.TotalVolume, &point.SetsCount)
		if err != nil {
			return nil, err
		}
		// Parse date string
		point.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			// If parsing fails, use current time
			point.Date = time.Now()
		}
		// Calculate one rep max
		point.OneRepMax = point.MaxWeight * float64(1+point.MaxReps/30.0)
		dataPoints = append(dataPoints, point)
	}
	chart.DataPoints = dataPoints

	// Calculate statistics
	stats, err := h.calculateExerciseStatistics(dataPoints)
	if err != nil {
		return nil, err
	}
	chart.Statistics = *stats

	// Identify milestones
	milestones, err := h.identifyExerciseMilestones(dataPoints)
	if err != nil {
		return nil, err
	}
	chart.Milestones = milestones

	return chart, nil
}

// calculateExerciseStatistics computes overall statistics from progress data
func (h *Handler) calculateExerciseStatistics(dataPoints []models.ExerciseProgressPoint) (*models.ExerciseStatistics, error) {
	stats := &models.ExerciseStatistics{}
	
	if len(dataPoints) == 0 {
		return stats, nil
	}

	var totalWeight, totalReps, totalVolume float64
	var firstDate, lastDate time.Time

	totalSessions := len(dataPoints)
	firstDate = dataPoints[0].Date
	lastDate = dataPoints[len(dataPoints)-1].Date

	maxWeight := 0.0
	maxReps := 0
	maxVolume := 0.0

	for _, point := range dataPoints {
		totalWeight += point.MaxWeight
		totalReps += float64(point.MaxReps)
		totalVolume += point.TotalVolume

		if point.MaxWeight > maxWeight {
			maxWeight = point.MaxWeight
		}
		if point.MaxReps > maxReps {
			maxReps = point.MaxReps
		}
		if point.TotalVolume > maxVolume {
			maxVolume = point.TotalVolume
		}
	}

	totalDays := lastDate.Sub(firstDate).Hours() / 24
	weeks := totalDays / 7
	progressRate := 0.0
	if weeks > 0 {
		progressRate = totalWeight / weeks // Simplified for illustrative purposes
	}

	stats.TotalSessions = totalSessions
	stats.TotalSets = totalSessions * 2 // Example calculation
	stats.TotalReps = int(totalReps)
	stats.TotalVolume = totalVolume
	stats.AverageWeight = totalWeight / float64(totalSessions)
	stats.AverageReps = totalReps / float64(totalSessions)
	stats.MaxWeightEver = maxWeight
	stats.MaxRepsEver = maxReps
	stats.MaxVolumeDay = maxVolume
	if weeks > 0 {
		stats.Consistency = float64(totalSessions) / weeks
	}
	stats.ProgressRate = progressRate
	stats.FirstWorkout = firstDate
	stats.LastWorkout = lastDate

	return stats, nil
}

// identifyExerciseMilestones identifies key milestones achieved for an exercise
func (h *Handler) identifyExerciseMilestones(dataPoints []models.ExerciseProgressPoint) ([]models.ExerciseMilestone, error) {
	// Example milestones - replace with real logic
	milestones := []models.ExerciseMilestone{
		{
			ID:             1,
			Type:           "weight",
			Description:    "First 100lb lift",
			Value:          100,
			Unit:           "lbs",
			AchievedAt:     dataPoints[0].Date, // Placeholder
			WorkoutID:      dataPoints[0].WorkoutID,
			IsPersonalRecord: true,
		},
	}

	return milestones, nil
}

// getWeeklySummary generates a comprehensive weekly summary for a user
func (h *Handler) getWeeklySummary(userID int, year int, week int) (*models.WeeklySummary, error) {
	// Calculate week start and end dates
	weekStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	daysToAdd := (week-1)*7 - int(weekStart.Weekday())
	weekStart = weekStart.AddDate(0, 0, daysToAdd)
	weekEnd := weekStart.AddDate(0, 0, 6)

	startDate := weekStart.Format("2006-01-02")
	endDate := weekEnd.Format("2006-01-02")

	summary := &models.WeeklySummary{
		WeekStart:  weekStart,
		WeekEnd:    weekEnd,
		WeekNumber: week,
		Year:       year,
	}

	// Get basic workout statistics
	var totalWorkouts, totalDuration int
	var avgDuration sql.NullFloat64
	query := `
		SELECT 
			COUNT(*) as total_workouts,
			COALESCE(SUM(duration), 0) as total_duration,
			COALESCE(AVG(duration), 0) as avg_duration
		FROM workouts 
		WHERE date >= ? AND date <= ?
	`
	err := h.db.QueryRow(query, startDate, endDate).Scan(&totalWorkouts, &totalDuration, &avgDuration)
	if err != nil {
		return nil, err
	}

	summary.TotalWorkouts = totalWorkouts
	summary.TotalDuration = totalDuration
	if avgDuration.Valid {
		summary.AvgDuration = avgDuration.Float64
	}

	// Get exercise and set statistics
	var totalSets, totalReps int
	var totalVolume, maxWeight float64
	var uniqueExercises int
	exerciseQuery := `
		SELECT 
			COUNT(DISTINCT e.name) as unique_exercises,
			COUNT(s.id) as total_sets,
			COALESCE(SUM(s.reps), 0) as total_reps,
			COALESCE(SUM(s.weight * s.reps), 0) as total_volume,
			COALESCE(MAX(s.weight), 0) as max_weight
		FROM workouts w
		LEFT JOIN exercises e ON w.id = e.workout_id
		LEFT JOIN sets s ON e.id = s.exercise_id
		WHERE w.date >= ? AND w.date <= ? AND s.weight > 0
	`
	err = h.db.QueryRow(exerciseQuery, startDate, endDate).Scan(&uniqueExercises, &totalSets, &totalReps, &totalVolume, &maxWeight)
	if err != nil {
		return nil, err
	}

	summary.UniqueExercises = uniqueExercises
	summary.TotalSets = totalSets
	summary.TotalReps = totalReps
	summary.TotalVolume = totalVolume
	summary.MaxWeight = maxWeight

	// Get nutrition data
	var avgCalories sql.NullFloat64
	calorieQuery := `
		SELECT AVG(calories) as avg_calories
		FROM meals 
		WHERE user_id = ? AND date >= ? AND date <= ?
	`
	err = h.db.QueryRow(calorieQuery, userID, startDate, endDate).Scan(&avgCalories)
	if err == nil && avgCalories.Valid {
		summary.AvgCalories = int(avgCalories.Float64)
	}

	// Get weight change
	weightChange, err := h.getWeightChangeForPeriod(userID, startDate, endDate)
	if err == nil {
		summary.WeightChange = weightChange
	}

	// Get top exercises
	topExercises, err := h.getTopExercisesForPeriod(userID, startDate, endDate, 3)
	if err == nil {
		summary.TopExercises = topExercises
	}

	// Count personal records achieved
	prCount, err := h.countPRsForPeriod(userID, startDate, endDate)
	if err == nil {
		summary.PRsAchieved = prCount
	}

	// Calculate consistency score (out of 7 days)
	if totalWorkouts > 0 {
		summary.ConsistencyScore = math.Min((float64(totalWorkouts)/3.0)*100, 100) // 3 workouts = 100%
	}

	// Calculate intensity score
	if totalSets > 0 {
		summary.IntensityScore = totalVolume / float64(totalSets)
	}

	return summary, nil
}

// getMonthlySummary generates a comprehensive monthly summary for a user
func (h *Handler) getMonthlySummary(userID int, year int, month int) (*models.MonthlySummary, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	summary := &models.MonthlySummary{
		Month:     month,
		Year:      year,
		MonthName: startDate.Format("January"),
	}

	// Get basic workout statistics
	var totalWorkouts, totalDuration int
	var avgDuration sql.NullFloat64
	query := `
		SELECT 
			COUNT(*) as total_workouts,
			COALESCE(SUM(duration), 0) as total_duration,
			COALESCE(AVG(duration), 0) as avg_duration
		FROM workouts 
		WHERE date >= ? AND date <= ?
	`
	err := h.db.QueryRow(query, startStr, endStr).Scan(&totalWorkouts, &totalDuration, &avgDuration)
	if err != nil {
		return nil, err
	}

	summary.TotalWorkouts = totalWorkouts
	summary.TotalDuration = totalDuration
	if avgDuration.Valid {
		summary.AvgDuration = avgDuration.Float64
	}

	// Get exercise and set statistics
	var totalSets, totalReps int
	var totalVolume, maxWeight float64
	var uniqueExercises int
	exerciseQuery := `
		SELECT 
			COUNT(DISTINCT e.name) as unique_exercises,
			COUNT(s.id) as total_sets,
			COALESCE(SUM(s.reps), 0) as total_reps,
			COALESCE(SUM(s.weight * s.reps), 0) as total_volume,
			COALESCE(MAX(s.weight), 0) as max_weight
		FROM workouts w
		LEFT JOIN exercises e ON w.id = e.workout_id
		LEFT JOIN sets s ON e.id = s.exercise_id
		WHERE w.date >= ? AND w.date <= ? AND s.weight > 0
	`
	err = h.db.QueryRow(exerciseQuery, startStr, endStr).Scan(&uniqueExercises, &totalSets, &totalReps, &totalVolume, &maxWeight)
	if err != nil {
		return nil, err
	}

	summary.UniqueExercises = uniqueExercises
	summary.TotalSets = totalSets
	summary.TotalReps = totalReps
	summary.TotalVolume = totalVolume
	summary.MaxWeight = maxWeight

	// Get nutrition data
	var avgCalories sql.NullFloat64
	calorieQuery := `
		SELECT AVG(calories) as avg_calories
		FROM meals 
		WHERE user_id = ? AND date >= ? AND date <= ?
	`
	err = h.db.QueryRow(calorieQuery, userID, startStr, endStr).Scan(&avgCalories)
	if err == nil && avgCalories.Valid {
		summary.AvgCalories = int(avgCalories.Float64)
	}

	// Get weight data
	startWeight, endWeight, weightChange, err := h.getWeightDataForPeriod(userID, startStr, endStr)
	if err == nil {
		summary.StartWeight = startWeight
		summary.EndWeight = endWeight
		summary.WeightChange = weightChange
	}

	// Get top exercises
	topExercises, err := h.getTopExercisesForPeriod(userID, startStr, endStr, 5)
	if err == nil {
		summary.TopExercises = topExercises
	}

	// Count personal records achieved
	prCount, err := h.countPRsForPeriod(userID, startStr, endStr)
	if err == nil {
		summary.PRsAchieved = prCount
	}

	// Calculate consistency score (assuming 12 workouts per month is ideal)
	if totalWorkouts > 0 {
		summary.ConsistencyScore = math.Min((float64(totalWorkouts)/12.0)*100, 100)
	}

	// Calculate intensity score
	if totalSets > 0 {
		summary.IntensityScore = totalVolume / float64(totalSets)
	}

	// Get weekly summaries for the month
	weeklySummaries, err := h.getWeeklySummariesForMonth(userID, year, month)
	if err == nil {
		summary.WeeklySummaries = weeklySummaries
	}

	// Get category breakdown
	categoryBreakdown, err := h.getCategoryBreakdownForPeriod(userID, startStr, endStr)
	if err == nil {
		summary.CategoryBreakdown = categoryBreakdown
	}

	// Generate progress highlights, goals, and recommendations
	summary.ProgressHighlights = h.generateProgressHighlights(summary)
	summary.GoalsAchieved = h.generateGoalsAchieved(summary)
	summary.Recommendations = h.generateRecommendations(summary)

	return summary, nil
}

// Helper functions for summary calculations
func (h *Handler) getWeightChangeForPeriod(userID int, startDate, endDate string) (float64, error) {
	var startWeight, endWeight sql.NullFloat64
	
	// Get first weight in period
	startQuery := `
		SELECT weight FROM body_weights 
		WHERE user_id = ? AND date >= ? 
		ORDER BY date ASC LIMIT 1
	`
	h.db.QueryRow(startQuery, userID, startDate).Scan(&startWeight)
	
	// Get last weight in period
	endQuery := `
		SELECT weight FROM body_weights 
		WHERE user_id = ? AND date <= ? 
		ORDER BY date DESC LIMIT 1
	`
	h.db.QueryRow(endQuery, userID, endDate).Scan(&endWeight)
	
	if startWeight.Valid && endWeight.Valid {
		return endWeight.Float64 - startWeight.Float64, nil
	}
	return 0, nil
}

func (h *Handler) getWeightDataForPeriod(userID int, startDate, endDate string) (float64, float64, float64, error) {
	var startWeight, endWeight sql.NullFloat64
	
	// Get first weight in period
	startQuery := `
		SELECT weight FROM body_weights 
		WHERE user_id = ? AND date >= ? 
		ORDER BY date ASC LIMIT 1
	`
	h.db.QueryRow(startQuery, userID, startDate).Scan(&startWeight)
	
	// Get last weight in period
	endQuery := `
		SELECT weight FROM body_weights 
		WHERE user_id = ? AND date <= ? 
		ORDER BY date DESC LIMIT 1
	`
	h.db.QueryRow(endQuery, userID, endDate).Scan(&endWeight)
	
	start := 0.0
	end := 0.0
	change := 0.0
	
	if startWeight.Valid {
		start = startWeight.Float64
	}
	if endWeight.Valid {
		end = endWeight.Float64
	}
	if startWeight.Valid && endWeight.Valid {
		change = end - start
	}
	
	return start, end, change, nil
}

func (h *Handler) getTopExercisesForPeriod(userID int, startDate, endDate string, limit int) ([]string, error) {
	query := `
		SELECT e.name, COUNT(*) as frequency
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.name
		ORDER BY frequency DESC
		LIMIT ?
	`
	
	rows, err := h.db.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var exercises []string
	for rows.Next() {
		var name string
		var frequency int
		err := rows.Scan(&name, &frequency)
		if err != nil {
			continue
		}
		exercises = append(exercises, name)
	}
	
	return exercises, nil
}

func (h *Handler) countPRsForPeriod(userID int, startDate, endDate string) (int, error) {
	// Count exercises where max weight in this period > previous max
	query := `
		SELECT COUNT(DISTINCT e.name) as pr_count
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		JOIN sets s ON e.id = s.exercise_id
		WHERE w.date >= ? AND w.date <= ?
		AND s.weight = (
			SELECT MAX(s2.weight)
			FROM sets s2
			JOIN exercises e2 ON s2.exercise_id = e2.id
			JOIN workouts w2 ON e2.workout_id = w2.id
			WHERE e2.name = e.name
		)
	`
	
	var count int
	err := h.db.QueryRow(query, startDate, endDate).Scan(&count)
	return count, err
}

func (h *Handler) getWeeklySummariesForMonth(userID int, year int, month int) ([]models.WeeklySummary, error) {
	// Get all weeks that overlap with this month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	
	var summaries []models.WeeklySummary
	
	// Find weeks that overlap with this month
	_, startWeek := startDate.ISOWeek()
	_, endWeek := endDate.ISOWeek()
	
	for week := startWeek; week <= endWeek; week++ {
		summary, err := h.getWeeklySummary(userID, year, week)
		if err == nil {
			summaries = append(summaries, *summary)
		}
	}
	
	return summaries, nil
}

func (h *Handler) getCategoryBreakdownForPeriod(userID int, startDate, endDate string) (map[string]int, error) {
	query := `
		SELECT e.category, COUNT(*) as count
		FROM exercises e
		JOIN workouts w ON e.workout_id = w.id
		WHERE w.date >= ? AND w.date <= ?
		GROUP BY e.category
	`
	
	rows, err := h.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	breakdown := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		if err != nil {
			continue
		}
		breakdown[category] = count
	}
	
	return breakdown, nil
}

// Generate insights and recommendations
func (h *Handler) generateProgressHighlights(summary *models.MonthlySummary) []string {
	highlights := []string{}
	
	if summary.TotalWorkouts > 12 {
		highlights = append(highlights, fmt.Sprintf("Great consistency with %d workouts this month!", summary.TotalWorkouts))
	}
	
	if summary.WeightChange < -2 {
		highlights = append(highlights, fmt.Sprintf("Excellent weight loss of %.1f lbs", math.Abs(summary.WeightChange)))
	} else if summary.WeightChange > 2 {
		highlights = append(highlights, fmt.Sprintf("Good weight gain of %.1f lbs", summary.WeightChange))
	}
	
	if summary.PRsAchieved > 5 {
		highlights = append(highlights, fmt.Sprintf("Outstanding! %d new personal records", summary.PRsAchieved))
	}
	
	if summary.TotalVolume > 50000 {
		highlights = append(highlights, fmt.Sprintf("Impressive total volume of %.0f lbs lifted", summary.TotalVolume))
	}
	
	return highlights
}

func (h *Handler) generateGoalsAchieved(summary *models.MonthlySummary) []string {
	goals := []string{}
	
	if summary.ConsistencyScore >= 80 {
		goals = append(goals, "Consistency Goal: 80%+ workout attendance")
	}
	
	if summary.PRsAchieved >= 3 {
		goals = append(goals, "Strength Goal: 3+ new personal records")
	}
	
	if summary.UniqueExercises >= 10 {
		goals = append(goals, "Variety Goal: 10+ different exercises")
	}
	
	return goals
}

func (h *Handler) generateRecommendations(summary *models.MonthlySummary) []string {
	recommendations := []string{}
	
	if summary.ConsistencyScore < 50 {
		recommendations = append(recommendations, "Try to maintain at least 3 workouts per week for better results")
	}
	
	if summary.UniqueExercises < 8 {
		recommendations = append(recommendations, "Consider adding more exercise variety to target different muscle groups")
	}
	
	if summary.PRsAchieved == 0 {
		recommendations = append(recommendations, "Focus on progressive overload to achieve new personal records")
	}
	
	// Check category balance
	if len(summary.CategoryBreakdown) > 0 {
		hasCardio := summary.CategoryBreakdown["cardio"] > 0
		hasStrength := summary.CategoryBreakdown["strength"] > 0 || summary.CategoryBreakdown["chest"] > 0 || summary.CategoryBreakdown["back"] > 0
		
		if !hasCardio {
			recommendations = append(recommendations, "Consider adding cardio exercises for cardiovascular health")
		}
		if !hasStrength {
			recommendations = append(recommendations, "Include strength training for muscle development")
		}
	}
	
	return recommendations
}
