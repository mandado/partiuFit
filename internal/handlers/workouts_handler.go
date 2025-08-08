package handlers

import (
	"net/http"
	internalErrors "partiuFit/internal/errors"
	"partiuFit/internal/middlewares"
	"partiuFit/internal/store"
	"partiuFit/internal/utils"

	"go.uber.org/zap"
)

type WorkoutsHandlers struct {
	Store  *store.Store
	Logger *zap.SugaredLogger
}

type UpdateWorkoutRequest struct {
	Title           *string              `json:"title"`
	Description     *string              `json:"description"`
	DurationMinutes *int                 `json:"duration_minutes"`
	CaloriesBurned  *int                 `json:"calories_burned"`
	Entries         []store.WorkoutEntry `json:"entries"`
}

func NewWorkoutsHandlers(store *store.Store, logger *zap.SugaredLogger) *WorkoutsHandlers {
	return &WorkoutsHandlers{
		Store:  store,
		Logger: logger,
	}
}

func (wh *WorkoutsHandlers) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := middlewares.GetUser(r)
	workouts, err := wh.Store.WorkoutStore.GetAllWorkouts(user.ID)

	if err != nil {
		wh.Logger.Error("failed to get workouts", zap.Error(err))
		utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to get workouts"})
		return
	}

	utils.MustWriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}

func (wh *WorkoutsHandlers) GetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	user := middlewares.GetUser(r)

	if err != nil {
		wh.Logger.Error("failed to parse workout id", zap.Error(err))
		utils.MustWriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout id"})
		return
	}

	utils.MustIfError(checkOwnerOfWorkout(wh, w, user, workoutID))
	workout := utils.Must(wh.Store.WorkoutStore.GetWorkoutById(workoutID))

	utils.MustWriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutsHandlers) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	workout := &store.Workout{}
	utils.MustReadJSON(w, r, workout)

	wh.Logger.Info("creating workout", zap.String("title", workout.Title))
	createdWorkout, err := wh.Store.WorkoutStore.CreateWorkout(workout)

	if err != nil {
		wh.Logger.Errorf("failed to create workout: %v", err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	utils.MustWriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutsHandlers) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUser(r)
	workoutID := utils.Must(utils.ReadIDParam(r))
	existingWorkout := utils.Must(wh.Store.WorkoutStore.GetWorkoutById(workoutID))

	if existingWorkout == nil {
		wh.Logger.Error("workout not found")
		utils.MustWriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	utils.MustIfError(checkOwnerOfWorkout(wh, w, user, workoutID))

	workout := &UpdateWorkoutRequest{}
	utils.MustReadJSON(w, r, workout)

	if workout.Title != nil {
		existingWorkout.Title = *workout.Title
	}

	if workout.Description != nil {
		existingWorkout.Description = *workout.Description
	}

	if workout.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *workout.DurationMinutes
	}

	if workout.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *workout.CaloriesBurned
	}

	if workout.Entries != nil {
		existingWorkout.Entries = workout.Entries
	}

	updatedWorkout, err := wh.Store.WorkoutStore.UpdateWorkout(workoutID, existingWorkout)

	if err != nil {
		wh.Logger.Errorf("failed to update workout: %v", err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}

	utils.MustWriteJSON(w, http.StatusOK, utils.Envelope{"workout": updatedWorkout})
}

func (wh *WorkoutsHandlers) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID := utils.Must(utils.ReadIDParam(r))
	user := middlewares.GetUser(r)

	utils.MustIfError(checkOwnerOfWorkout(wh, w, user, workoutID))
	utils.MustIfError(wh.Store.WorkoutStore.DeleteWorkout(workoutID))

	w.WriteHeader(http.StatusNoContent)
}

func checkOwnerOfWorkout(wh *WorkoutsHandlers, w http.ResponseWriter, user *store.User, workoutID int) error {
	isWorkoutOwner := utils.Must(wh.Store.WorkoutStore.OwnsWorkout(workoutID, user.ID))

	if !isWorkoutOwner {
		wh.Logger.Error("user does not own this workout")
		return internalErrors.ErrForbidden
	}

	return nil
}
