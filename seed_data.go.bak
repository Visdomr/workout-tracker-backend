package main

import (
	"log"
	"time"

	"workout-tracker/internal/database"
	"workout-tracker/internal/models"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Sample workouts
	workouts := []models.Workout{
		{
			Name:     "Morning Push Day",
			Date:     time.Now().AddDate(0, 0, -1),
			Duration: 45,
			Notes:    "Great chest and triceps workout",
		},
		{
			Name:     "Leg Day",
			Date:     time.Now().AddDate(0, 0, -3),
			Duration: 60,
			Notes:    "Focused on squats and deadlifts",
		},
		{
			Name:     "Full Body Circuit",
			Date:     time.Now().AddDate(0, 0, -5),
			Duration: 35,
			Notes:    "High intensity interval training",
		},
		{
			Name:     "Upper Body Strength",
			Date:     time.Now().AddDate(0, 0, -7),
			Duration: 50,
			Notes:    "Back and biceps focused session",
		},
		{
			Name:     "Cardio Session",
			Date:     time.Now().AddDate(0, 0, -2),
			Duration: 30,
			Notes:    "Running and cycling",
		},
	}

	// Insert workouts
	for _, workout := range workouts {
		result, err := db.Exec(`
			INSERT INTO workouts (name, date, duration, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, workout.Name, workout.Date, workout.Duration, workout.Notes, time.Now(), time.Now())
		
		if err != nil {
			log.Printf("Error inserting workout %s: %v", workout.Name, err)
			continue
		}

		workoutID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting workout ID: %v", err)
			continue
		}

		log.Printf("Created workout: %s (ID: %d)", workout.Name, workoutID)

		// Add some sample exercises for the first workout
		if workout.Name == "Morning Push Day" {
			exercises := []struct {
				Name     string
				Category string
				Sets     []struct {
					SetNum int
					Reps   int
					Weight float64
				}
			}{
				{
					Name:     "Bench Press",
					Category: "strength",
					Sets: []struct {
						SetNum int
						Reps   int
						Weight float64
					}{
						{1, 10, 135},
						{2, 8, 155},
						{3, 6, 175},
					},
				},
				{
					Name:     "Push-ups",
					Category: "strength",
					Sets: []struct {
						SetNum int
						Reps   int
						Weight float64
					}{
						{1, 15, 0},
						{2, 12, 0},
						{3, 10, 0},
					},
				},
				{
					Name:     "Tricep Dips",
					Category: "strength",
					Sets: []struct {
						SetNum int
						Reps   int
						Weight float64
					}{
						{1, 12, 0},
						{2, 10, 0},
						{3, 8, 0},
					},
				},
			}

			for _, exercise := range exercises {
				result, err := db.Exec(`
					INSERT INTO exercises (workout_id, name, category, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?)
				`, workoutID, exercise.Name, exercise.Category, time.Now(), time.Now())
				
				if err != nil {
					log.Printf("Error inserting exercise %s: %v", exercise.Name, err)
					continue
				}

				exerciseID, err := result.LastInsertId()
				if err != nil {
					log.Printf("Error getting exercise ID: %v", err)
					continue
				}

				log.Printf("  Created exercise: %s (ID: %d)", exercise.Name, exerciseID)

				// Add sets
				for _, set := range exercise.Sets {
					_, err := db.Exec(`
						INSERT INTO sets (exercise_id, set_number, reps, weight, distance, duration, rest_time, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
					`, exerciseID, set.SetNum, set.Reps, set.Weight, 0, 0, 60, time.Now(), time.Now())
					
					if err != nil {
						log.Printf("    Error inserting set %d: %v", set.SetNum, err)
						continue
					}

					log.Printf("    Created set %d: %d reps @ %.1f lbs", set.SetNum, set.Reps, set.Weight)
				}
			}
		}
	}

	log.Println("Sample data seeding completed!")
}
