package store

import (
	"crypto/sha256"
	"database/sql"
	internalErrors "partiuFit/internal/errors"
	"partiuFit/internal/requests"
	"partiuFit/internal/valueObjects"
	"time"

	_ "golang.org/x/crypto/bcrypt"
)

var (
	AnonymousUser = &User{}
)

type User struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Username  string                 `json:"username"`
	Password  *valueObjects.Password `json:"-"`
	Email     string                 `json:"email"`
	CreatedAt *time.Time             `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) FromUserRequest(userRequest *requests.UserRequest) *User {
	u.Name = userRequest.Name
	u.Username = userRequest.Username
	u.Email = userRequest.Email
	u.Password = valueObjects.NewPassword(userRequest.Password)

	return u
}

func (u *User) HashPassword() error {
	return u.Password.HashPassword()
}

type UserStore interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(id int, user *User) error
	GetUserFromToken(scope, token string) (*User, error)
}

type UserPostgresStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *UserPostgresStore {
	return &UserPostgresStore{
		db: db,
	}
}

func (s *UserPostgresStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		Password: &valueObjects.Password{},
	}
	query := `
		select id, name, username, password, email
		from users
		where username = $1
	`
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Username, &user.Password.Hash, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserPostgresStore) CreateUser(user *User) error {
	query := `
		insert into users (name, username, password, email)
		values ($1, $2, $3, $4)
		returning id, created_at, updated_at
	`

	err := s.db.QueryRow(query, user.Name, user.Username, user.Password.GetHash(), user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return internalErrors.HandleDatabaseError(err)
	}

	return nil
}

func (s *UserPostgresStore) UpdateUser(id int, user *User) error {
	query := `
		update users
		set name = $2, username = $3, password = $4, email = $5
		where id = $1
	`

	result, err := s.db.Exec(query, id, user.Name, user.Username, user.Password.GetHash(), user.Email)

	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return internalErrors.ErrNoRows
	}

	return nil
}

func (s *UserPostgresStore) GetUserFromToken(scope, token string) (*User, error) {
	hash := sha256.Sum256([]byte(token))

	query := `
		select id, name, username, password, email
		from users where exists (
		    select 1 from tokens 
		             where 
		                 tokens.user_id = users.id and 
		                 tokens.scope = $1 and 
		                 tokens.hash = $2 and 
		                 tokens.expires_at > $3
		)
	`
	user := &User{
		Password: &valueObjects.Password{},
	}

	err := s.db.QueryRow(query, scope, hash[:], time.Now()).Scan(&user.ID, &user.Name, &user.Username, &user.Password.Hash, &user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}
