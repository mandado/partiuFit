package testing_utils

import (
	"database/sql"
	"partiuFit/internal/tokens"
	"partiuFit/internal/utils"
	"partiuFit/internal/valueObjects"
	"time"
)

type User struct {
	ID       int
	Name     string
	Username string
	Password *valueObjects.Password
	Email    string
}

func SeedDB(db *sql.DB) error {
	err := createUsers(db)

	if err != nil {
		return err
	}

	return nil
}

func createUsers(db *sql.DB) error {
	users := []struct {
		Name     string
		Username string
		Email    string
		Password string
	}{
		{Name: "John Doe", Username: "johndoe", Email: "john@example.com", Password: "password"},
		{Name: "Jane Doe", Username: "janedoe", Email: "jonM@example.com", Password: "password"},
	}

	for _, user := range users {
		password := &valueObjects.Password{
			PlainText: user.Password,
		}
		password.HashPassword()
		_, err := db.Exec("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)", user.Name, user.Username, user.Email, password.GetHash())

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateToken(db *sql.DB, username string) (*User, *tokens.Token) {
	user := &User{
		Password: &valueObjects.Password{},
	}
	query := `
		select id, name, username, password, email
		from users
		where username = $1
	`
	utils.MustIfError(db.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Username, &user.Password.Hash, &user.Email))
	token := utils.Must(tokens.GenerateToken(user.ID, time.Hour, tokens.ScopeAuthentication))

	// Insert token into database
	_, err := db.Exec(`
		INSERT INTO tokens (hash, user_id, expires_at, scope) 
		VALUES ($1, $2, $3, $4)
	`, token.Hash, token.UserID, token.ExpiresAt, token.Scope)
	utils.MustIfError(err)

	return user, token
}
