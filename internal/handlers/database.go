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
