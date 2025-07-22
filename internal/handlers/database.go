package handlers

import (
	"database/sql"
	"fmt"
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
		SELECT id, name, category, description, created_at, updated_at
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
		err := rows.Scan(&e.ID, &e.Name, &e.Category, &e.Description, &e.CreatedAt, &e.UpdatedAt)
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
		SELECT id, name, category, description, created_at, updated_at
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
		err := rows.Scan(&e.ID, &e.Name, &e.Category, &e.Description, &e.CreatedAt, &e.UpdatedAt)
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
		INSERT INTO predefined_exercises (name, category, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := h.db.Exec(query, exercise.Name, exercise.Category, exercise.Description, time.Now(), time.Now())
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
