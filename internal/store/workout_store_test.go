package store

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	testingUtils "partiuFit/internal/database/testing_utils"
	"partiuFit/internal/utils"
	"testing"
)

func TestWorkoutStore(t *testing.T) {
	db := utils.Must(testingUtils.SetupTestDB())
	utils.MustIfError(testingUtils.SeedDB(db))
	user, _ := testingUtils.CreateToken(db, "johndoe")
	defer testingUtils.TeardownTestDB(db)

	workutStore := NewPostgresWorkoutStore(db)

	t.Run("Create with valid data", func(t *testing.T) {
		workout := &Workout{
			Title:           "Test Workout",
			Description:     "This is a test workout",
			DurationMinutes: 30,
			CaloriesBurned:  200,
			UserID:          user.ID,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench Press",
					Sets:         3,
					Reps:         utils.ValueToPointer(10),
					Weight:       50,
					Notes:        "Use a bench with 45-degree incline",
					OrderIndex:   1,
					UserID:       user.ID,
				},
			},
		}

		createdWorkout, err := workutStore.CreateWorkout(workout)

		assert.NoError(t, err)
		assert.NotNil(t, createdWorkout)
		assert.Equal(t, "Test Workout", createdWorkout.Title)
		assert.Equal(t, "This is a test workout", createdWorkout.Description)
		assert.Equal(t, 30, createdWorkout.DurationMinutes)
		assert.Equal(t, 200, createdWorkout.CaloriesBurned)
		assert.Equal(t, 1, len(createdWorkout.Entries))
		assert.Equal(t, "Bench Press", createdWorkout.Entries[0].ExerciseName)
		assert.NotEmpty(t, createdWorkout.Entries)
	})

	t.Run("Create with invalid", func(t *testing.T) {
		workout := &Workout{
			Title:           "Test Workout",
			Description:     "This is a test workout",
			DurationMinutes: 30,
			CaloriesBurned:  200,
			UserID:          user.ID,
			Entries: []WorkoutEntry{
				{
					ExerciseName:    "Bench Press",
					Sets:            3,
					Reps:            utils.ValueToPointer(3),
					DurationSeconds: utils.ValueToPointer(60),
					Weight:          50,
					Notes:           "Use a bench with 45-degree incline",
					OrderIndex:      1,
					UserID:          user.ID,
				},
			},
		}

		createdWorkout, err := workutStore.CreateWorkout(workout)

		assert.Error(t, err)
		assert.Equal(t, `ERROR: new row for relation "workout_entries" violates check constraint "valid_workout_entry" (SQLSTATE 23514)`, err.Error())
		assert.Nil(t, createdWorkout)
	})

	t.Run("Update with valid data", func(t *testing.T) {
		workout := &Workout{
			Title:           "Test Workout",
			Description:     "This is a test workout",
			DurationMinutes: 30,
			CaloriesBurned:  200,
			UserID:          user.ID,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench Press",
					Sets:         3,
					Reps:         utils.ValueToPointer(10),
					Weight:       50,
					Notes:        "Use a bench with 45-degree incline",
					OrderIndex:   1,
					UserID:       user.ID,
				},
			},
		}

		createdWorkout := utils.Must(workutStore.CreateWorkout(workout))

		createdWorkout.Title = "Updated Test Workout"

		updatedWorkout, err := workutStore.UpdateWorkout(createdWorkout.ID, createdWorkout)

		assert.NoError(t, err)
		assert.NotNil(t, updatedWorkout)
		assert.Equal(t, "Updated Test Workout", updatedWorkout.Title)
	})

	t.Run("Get by id", func(t *testing.T) {
		workout := &Workout{
			Title:           "Test Workout 3",
			Description:     "This is a test workout",
			DurationMinutes: 30,
			CaloriesBurned:  200,
			UserID:          user.ID,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench Press",
					Sets:         3,
					Reps:         utils.ValueToPointer(10),
					Weight:       50,
					Notes:        "Use a bench with 45-degree incline",
					OrderIndex:   1,
					UserID:       user.ID,
				},
			},
		}

		createdWorkout := utils.Must(workutStore.CreateWorkout(workout))

		retrievedWorkout, err := workutStore.GetWorkoutById(createdWorkout.ID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedWorkout)
		assert.Equal(t, "Test Workout 3", retrievedWorkout.Title)
	})

	t.Run("Delete by id", func(t *testing.T) {
		workout := &Workout{
			Title:           "Test Workout",
			Description:     "This is a test workout",
			DurationMinutes: 30,
			CaloriesBurned:  200,
			UserID:          user.ID,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench Press",
					Sets:         3,
					Reps:         utils.ValueToPointer(10),
					Weight:       50,
					Notes:        "Use a bench with 45-degree incline",
					OrderIndex:   1,
					UserID:       user.ID,
				},
			},
		}

		createdWorkout := utils.Must(workutStore.CreateWorkout(workout))

		err := workutStore.DeleteWorkout(createdWorkout.ID)

		assert.NoError(t, err)
		_, err = workutStore.GetWorkoutById(createdWorkout.ID)
		assert.Error(t, err)
	})

	t.Run("Get all", func(t *testing.T) {
		workouts := utils.Must(workutStore.GetAllWorkouts(user.ID))

		assert.Len(t, workouts, 3)

		assert.Equal(t, "Test Workout 3", workouts[0].Title)
		assert.Equal(t, "Updated Test Workout", workouts[1].Title)
		assert.Equal(t, "Test Workout", workouts[2].Title)
	})
}
