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


	// Define common exercises organized by category
	exercises := []models.PredefinedExercise{
		// Chest exercises
		{Name: "Bench Press", Category: "chest", Description: "Classic chest exercise using a barbell on a bench"},
		{Name: "Dumbbell Flyes", Category: "chest", Description: "Chest isolation exercise using dumbbells"},
		{Name: "Push-ups", Category: "chest", Description: "Bodyweight chest exercise"},
		{Name: "Incline Bench Press", Category: "chest", Description: "Upper chest focused bench press on inclined bench"},
		{Name: "Dips", Category: "chest", Description: "Bodyweight exercise targeting chest and triceps"},
		{Name: "Chest Press Machine", Category: "chest", Description: "Machine-based chest exercise"},

		// Back exercises
		{Name: "Pull-ups", Category: "back", Description: "Bodyweight exercise targeting the entire back"},
		{Name: "Deadlifts", Category: "back", Description: "Compound exercise for back, glutes, and hamstrings"},
		{Name: "Bent-over Rows", Category: "back", Description: "Barbell rowing exercise for mid-back"},
		{Name: "Lat Pulldowns", Category: "back", Description: "Cable exercise targeting latissimus dorsi"},
		{Name: "T-Bar Rows", Category: "back", Description: "Rowing exercise using T-bar apparatus"},
		{Name: "Single-Arm Dumbbell Row", Category: "back", Description: "Unilateral back exercise using dumbbells"},

		// Legs exercises
		{Name: "Squats", Category: "legs", Description: "Fundamental compound leg exercise"},
		{Name: "Lunges", Category: "legs", Description: "Single-leg exercise for quads, glutes, and hamstrings"},
		{Name: "Romanian Deadlifts", Category: "legs", Description: "Hamstring and glute focused deadlift variation"},
		{Name: "Leg Press", Category: "legs", Description: "Machine-based leg exercise"},
		{Name: "Calf Raises", Category: "legs", Description: "Exercise targeting calf muscles"},
		{Name: "Leg Curls", Category: "legs", Description: "Hamstring isolation exercise"},
		{Name: "Leg Extensions", Category: "legs", Description: "Quadriceps isolation exercise"},

		// Shoulders exercises
		{Name: "Overhead Press", Category: "shoulders", Description: "Compound shoulder exercise using barbell or dumbbells"},
		{Name: "Lateral Raises", Category: "shoulders", Description: "Side deltoid isolation exercise"},
		{Name: "Front Raises", Category: "shoulders", Description: "Front deltoid isolation exercise"},
		{Name: "Rear Delt Flyes", Category: "shoulders", Description: "Posterior deltoid targeting exercise"},
		{Name: "Upright Rows", Category: "shoulders", Description: "Compound exercise for shoulders and traps"},
		{Name: "Arnold Press", Category: "shoulders", Description: "Dumbbell press variation for complete shoulder development"},

		// Arms exercises
		{Name: "Bicep Curls", Category: "arms", Description: "Basic bicep isolation exercise"},
		{Name: "Tricep Dips", Category: "arms", Description: "Bodyweight tricep exercise"},
		{Name: "Hammer Curls", Category: "arms", Description: "Neutral grip bicep exercise"},
		{Name: "Tricep Pushdowns", Category: "arms", Description: "Cable tricep isolation exercise"},
		{Name: "Close-Grip Bench Press", Category: "arms", Description: "Compound exercise emphasizing triceps"},
		{Name: "Preacher Curls", Category: "arms", Description: "Bicep exercise using preacher bench"},

		// Core exercises
		{Name: "Plank", Category: "core", Description: "Isometric core strengthening exercise"},
		{Name: "Crunches", Category: "core", Description: "Basic abdominal exercise"},
		{Name: "Russian Twists", Category: "core", Description: "Rotational core exercise"},
		{Name: "Mountain Climbers", Category: "core", Description: "Dynamic core and cardio exercise"},
		{Name: "Dead Bug", Category: "core", Description: "Core stability exercise"},
		{Name: "Bicycle Crunches", Category: "core", Description: "Dynamic abdominal exercise"},

		// Cardio exercises
		{Name: "Treadmill Running", Category: "cardio", Description: "Running on a treadmill"},
		{Name: "Stationary Bike", Category: "cardio", Description: "Cycling on a stationary bike"},
		{Name: "Elliptical", Category: "cardio", Description: "Low-impact cardio exercise"},
		{Name: "Rowing Machine", Category: "cardio", Description: "Full-body cardio exercise"},
		{Name: "Jump Rope", Category: "cardio", Description: "High-intensity cardio exercise"},
		{Name: "Burpees", Category: "cardio", Description: "High-intensity full-body exercise"},

		// Flexibility exercises
		{Name: "Yoga Flow", Category: "flexibility", Description: "Dynamic yoga sequence"},
		{Name: "Static Stretching", Category: "flexibility", Description: "Holding stretches for extended periods"},
		{Name: "Foam Rolling", Category: "flexibility", Description: "Self-myofascial release technique"},
		{Name: "Dynamic Warm-up", Category: "flexibility", Description: "Movement-based warm-up routine"},
	}

	// Insert exercises into database
	query := `
		INSERT INTO predefined_exercises (name, category, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	for _, exercise := range exercises {
		result, err := db.Exec(query, exercise.Name, exercise.Category, exercise.Description, time.Now(), time.Now())
		if err != nil {
			log.Printf("Failed to create exercise %s: %v", exercise.Name, err)
			continue
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Failed to get ID for exercise %s: %v", exercise.Name, err)
			continue
		}
		
		log.Printf("Created exercise: %s (%s) with ID %d", exercise.Name, exercise.Category, id)
	}

	log.Println("Exercise seeding completed!")
}
