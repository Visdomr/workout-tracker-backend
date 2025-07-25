package main

import (
	"encoding/json"
	"fmt"
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

	// Sample exercises with form videos and instructions
	exercises := []models.PredefinedExercise{
		{
			Name:         "Barbell Squat",
			Category:     "strength",
			Description:  "A compound exercise that targets the quadriceps, hamstrings, and glutes.",
			VideoURL:     "https://www.youtube.com/watch?v=nEQQle9-0NA",
			Instructions: "1. Stand with feet shoulder-width apart\n2. Place barbell on your upper back\n3. Lower your body by pushing hips back and bending knees\n4. Keep chest up and knees tracking over toes\n5. Lower until thighs are parallel to ground\n6. Drive through heels to return to starting position",
			Tips:         "- Keep your core tight throughout the movement\n- Don't let knees collapse inward\n- Maintain neutral spine\n- Start with bodyweight before adding load",
			MuscleGroups: "Quadriceps, Hamstrings, Glutes, Core",
			Equipment:    "Barbell, Squat Rack",
			Difficulty:   "intermediate",
			ImageURL:     "https://example.com/images/barbell-squat.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Deadlift",
			Category:     "strength", 
			Description:  "A fundamental compound movement that works the entire posterior chain.",
			VideoURL:     "https://www.youtube.com/watch?v=op9kVnSso6Q",
			Instructions: "1. Stand with feet hip-width apart, bar over mid-foot\n2. Hinge at hips and bend down to grip bar\n3. Keep chest up, shoulders back\n4. Drive through heels and stand up straight\n5. Keep bar close to body throughout\n6. Lower bar by pushing hips back first",
			Tips:         "- Keep the bar close to your body\n- Don't round your back\n- Engage your lats to keep bar close\n- Start light and focus on form",
			MuscleGroups: "Hamstrings, Glutes, Erector Spinae, Traps, Lats",
			Equipment:    "Barbell, Weight plates",
			Difficulty:   "intermediate",
			ImageURL:     "https://example.com/images/deadlift.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Bench Press",
			Category:     "strength",
			Description:  "Upper body compound exercise targeting chest, shoulders, and triceps.",
			VideoURL:     "https://www.youtube.com/watch?v=gRVjAtPip0Y",
			Instructions: "1. Lie on bench with eyes under the bar\n2. Grip bar slightly wider than shoulder-width\n3. Plant feet firmly on ground\n4. Unrack bar and lower to chest\n5. Press bar up in straight line\n6. Keep shoulders retracted throughout",
			Tips:         "- Squeeze shoulder blades together\n- Keep wrists straight\n- Control the weight down slowly\n- Don't bounce bar off chest",
			MuscleGroups: "Chest, Anterior Deltoids, Triceps",
			Equipment:    "Barbell, Bench, Weight plates",
			Difficulty:   "intermediate",
			ImageURL:     "https://example.com/images/bench-press.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Pull-ups",
			Category:     "strength",
			Description:  "Bodyweight exercise for upper body pulling strength.",
			VideoURL:     "https://www.youtube.com/watch?v=eGo4IYlbE5g",
			Instructions: "1. Hang from bar with hands shoulder-width apart\n2. Pull yourself up until chin clears bar\n3. Lower yourself with control\n4. Keep core tight throughout\n5. Avoid swinging or kipping",
			Tips:         "- Start with assisted variations if needed\n- Focus on pulling with your back muscles\n- Keep shoulders down and back\n- Full range of motion is key",
			MuscleGroups: "Lats, Rhomboids, Biceps, Rear Delts",
			Equipment:    "Pull-up bar",
			Difficulty:   "intermediate",
			ImageURL:     "https://example.com/images/pull-ups.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Push-ups",
			Category:     "strength",
			Description:  "Bodyweight exercise for chest, shoulders, and triceps.",
			VideoURL:     "https://www.youtube.com/watch?v=IODxDxX7oi4",
			Instructions: "1. Start in plank position with hands under shoulders\n2. Lower body until chest nearly touches ground\n3. Push back up to starting position\n4. Keep body in straight line throughout\n5. Don't let hips sag or pike up",
			Tips:         "- Modify on knees if needed\n- Keep core tight\n- Control the descent\n- Full range of motion",
			MuscleGroups: "Chest, Anterior Deltoids, Triceps, Core",
			Equipment:    "None (bodyweight)",
			Difficulty:   "beginner",
			ImageURL:     "https://example.com/images/push-ups.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Overhead Press",
			Category:     "strength",
			Description:  "Vertical pressing movement for shoulders and triceps.",
			VideoURL:     "https://www.youtube.com/watch?v=2yjwXTZQDDI",
			Instructions: "1. Stand with feet hip-width apart\n2. Hold bar at shoulder height\n3. Press bar straight up overhead\n4. Keep core tight and avoid arching back\n5. Lower bar with control to shoulders",
			Tips:         "- Don't press behind neck\n- Keep elbows slightly forward\n- Engage your core\n- Start with lighter weight",
			MuscleGroups: "Anterior Deltoids, Triceps, Core",
			Equipment:    "Barbell or Dumbbells",
			Difficulty:   "intermediate",
			ImageURL:     "https://example.com/images/overhead-press.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Insert exercises into database
	for _, exercise := range exercises {
		query := `
			INSERT INTO predefined_exercises (name, category, description, video_url, instructions, tips, muscle_groups, equipment, difficulty, image_url, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		
		_, err := db.Exec(query, exercise.Name, exercise.Category, exercise.Description, exercise.VideoURL, exercise.Instructions, exercise.Tips, exercise.MuscleGroups, exercise.Equipment, exercise.Difficulty, exercise.ImageURL, exercise.CreatedAt, exercise.UpdatedAt)
		if err != nil {
			log.Printf("Failed to insert exercise %s: %v", exercise.Name, err)
		} else {
			fmt.Printf("Successfully inserted exercise: %s\n", exercise.Name)
		}
	}

	// Query and display the results
	fmt.Println("\n=== Predefined Exercises with Video URLs ===")
	rows, err := db.Query(`
		SELECT id, name, category, description, video_url, instructions, tips, muscle_groups, equipment, difficulty, image_url
		FROM predefined_exercises 
		ORDER BY category, name
	`)
	if err != nil {
		log.Fatal("Failed to query exercises:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exercise models.PredefinedExercise
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Category, &exercise.Description, &exercise.VideoURL, &exercise.Instructions, &exercise.Tips, &exercise.MuscleGroups, &exercise.Equipment, &exercise.Difficulty, &exercise.ImageURL)
		if err != nil {
			log.Printf("Failed to scan exercise: %v", err)
			continue
		}

		// Pretty print the exercise as JSON
		exerciseJSON, _ := json.MarshalIndent(exercise, "", "  ")
		fmt.Printf("%s\n\n", exerciseJSON)
	}

	fmt.Println("Exercise seeding completed!")
}
