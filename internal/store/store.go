package store

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	WorkoutStore WorkoutStore
	UserStore    UserStore
	TokensStore  TokensStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		WorkoutStore: NewPostgresWorkoutStore(db),
		UserStore:    NewPostgresUserStore(db),
		TokensStore:  NewPostgresTokensStore(db),
	}
}
