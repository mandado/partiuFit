package store

import (
	"database/sql"
	tokens "partiuFit/internal/tokens"
	"time"
)

type TokensStore interface {
	InsertToken(tokens *tokens.Token) error
	CreateToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userId int) error
}

type PostgresTokensStore struct {
	db *sql.DB
}

func NewPostgresTokensStore(db *sql.DB) *PostgresTokensStore {
	return &PostgresTokensStore{
		db: db,
	}
}

func (s *PostgresTokensStore) CreateToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = s.InsertToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *PostgresTokensStore) InsertToken(token *tokens.Token) error {
	query := `
		insert into tokens (hash, user_id, scope, expires_at)
		values ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Scope, token.ExpiresAt)

	return err
}

func (s *PostgresTokensStore) DeleteAllTokensForUser(userId int) error {
	query := `
		delete from tokens
		where user_id = $1
	`

	_, err := s.db.Exec(query, userId)

	return err
}
