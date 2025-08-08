package store

import (
	"database/sql"
	internalErrors "partiuFit/internal/errors"
	"time"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	UserID          int            `json:"user_id"`
}

type WorkoutEntry struct {
	ID              int     `json:"id"`
	ExerciseName    string  `json:"exercise_name"`
	Reps            *int    `json:"reps"`
	Sets            int     `json:"sets"`
	Weight          float64 `json:"weight"`
	DurationSeconds *int    `json:"duration_seconds"`
	Notes           string  `json:"notes"`
	OrderIndex      int     `json:"order_index"`
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	UserID          int
}

type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	UpdateWorkout(id int, workout *Workout) (*Workout, error)
	GetWorkoutById(id int) (*Workout, error)
	DeleteWorkout(id int) error
	GetAllWorkouts(userID int) ([]Workout, error)
	OwnsWorkout(id int, userID int) (bool, error)
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

func (s *PostgresWorkoutStore) GetAllWorkouts(userID int) ([]Workout, error) {
	query := `
		select id, title, description, duration_minutes, calories_burned, created_at, updated_at, user_id
		from workouts
		where user_id = $1
		order by created_at desc
	`

	rows, err := s.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var workouts = make([]Workout, 0)

	for rows.Next() {
		workout := &Workout{}

		err := rows.Scan(
			&workout.ID,
			&workout.Title,
			&workout.Description,
			&workout.DurationMinutes,
			&workout.CaloriesBurned,
			&workout.CreatedAt,
			&workout.UpdatedAt,
			&workout.UserID,
		)

		if err != nil {
			return nil, err
		}

		workouts = append(workouts, *workout)
	}

	return workouts, nil
}

func (s *PostgresWorkoutStore) GetWorkoutById(id int) (*Workout, error) {
	workout := &Workout{}
	query := `
		select id, title, description, duration_minutes, calories_burned, created_at, updated_at, user_id
		from workouts
		where id = $1
	`

	err := s.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Description,
		&workout.DurationMinutes,
		&workout.CaloriesBurned,
		&workout.CreatedAt,
		&workout.UpdatedAt,
		&workout.UserID,
	)

	if err != nil {
		return nil, err
	}

	entriesQuery := `
		select id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index, created_at, updated_at, user_id
		from workout_entries
		where workout_id = $1
		order by order_index
	`

	rows, err := s.db.Query(entriesQuery, id)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		entry := WorkoutEntry{}

		err := rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.UserID,
		)

		if err != nil {
			return nil, err
		}

		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (s *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `
			insert into workouts (title, description, duration_minutes, calories_burned, user_id)
			values ($1, $2, $3, $4, $5)
			returning id, created_at, updated_at
	`

	err = tx.QueryRow(
		query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
		workout.UserID).Scan(&workout.ID, &workout.CreatedAt, &workout.UpdatedAt)

	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query := `
			insert into workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index, user_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			returning id, created_at, updated_at
		`

		err := tx.QueryRow(
			query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex, entry.UserID).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)

		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *PostgresWorkoutStore) UpdateWorkout(id int, workout *Workout) (*Workout, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		update workouts
		set title = $2, description = $3, duration_minutes = $4, calories_burned = $5
		where id = $1
	`

	workout.ID = int(id)

	result, err := tx.Exec(query,
		id,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned)

	if err != nil {
		return nil, err
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, internalErrors.ErrNoRows
	}

	_, err = tx.Exec("delete from workout_entries where workout_id = $1", id)

	if err != nil {
		return nil, err
	}

	for i, entry := range workout.Entries {
		query := `
			insert into workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index, user_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			returning id
		`

		err := tx.QueryRow(
			query,
			id,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
			entry.UserID).Scan(&workout.Entries[i].ID)

		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *PostgresWorkoutStore) DeleteWorkout(id int) error {
	result, err := s.db.Exec("delete from workouts where id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return internalErrors.ErrNoRows
	}

	return err
}

func (s *PostgresWorkoutStore) OwnsWorkout(id int, userID int) (bool, error) {
	var owns bool

	query := `
		select exists(select 1 from workouts where id = $1 and user_id = $2)
	`

	err := s.db.QueryRow(query, id, userID).Scan(&owns)

	return owns, err
}
