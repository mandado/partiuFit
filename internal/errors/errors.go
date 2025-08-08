package errors

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	NoRows               = sql.ErrNoRows
	ErrUserAlreadyExists = errors.New("já existe um usuário com esse username ou email")
	InvalidCredentials   = errors.New("invalid credentials")
	Forbidden            = errors.New("você não tem permissao para realizar essa operação")
	InvalidIDParam       = errors.New("parametro de id invalido")
	InvalidIDType        = errors.New("tipo de id invalido")
)

func isPgDuplicateUserError(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr)
}

func HandleDatabaseError(err error) error {
	if isPgDuplicateUserError(err) {
		return ErrUserAlreadyExists
	}

	return err
}
