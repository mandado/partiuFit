package store

import (
	"errors"
	testingUtils "partiuFit/internal/database/testing_utils"
	internalErrors "partiuFit/internal/errors"
	"partiuFit/internal/utils"
	"partiuFit/internal/valueObjects"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	db := utils.Must(testingUtils.SetupTestDB())
	utils.MustIfError(testingUtils.SeedDB(db))
	defer func() { _ = testingUtils.TeardownTestDB(db) }()

	userStore := NewPostgresUserStore(db)

	t.Run("CreateUser with valid data", func(t *testing.T) {
		user := &User{
			Name:     "Test User",
			Username: "testuser",
			Password: &valueObjects.Password{
				PlainText: "password123",
			},
			Email: "test@example.com",
		}

		err := user.HashPassword()
		assert.NoError(t, err)

		err = userStore.CreateUser(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("CreateUser with duplicate username", func(t *testing.T) {
		user := &User{
			Name:     "Duplicate User",
			Username: "johndoe",
			Password: &valueObjects.Password{
				PlainText: "password123",
			},
			Email: "duplicate@example.com",
		}

		err := user.HashPassword()
		assert.NoError(t, err)

		err = userStore.CreateUser(user)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrUserAlreadyExists))
	})

	t.Run("CreateUser with duplicate email", func(t *testing.T) {
		user := &User{
			Name:     "Another User",
			Username: "anotheruser",
			Password: &valueObjects.Password{
				PlainText: "password123",
			},
			Email: "john@example.com",
		}

		err := user.HashPassword()
		assert.NoError(t, err)

		err = userStore.CreateUser(user)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrUserAlreadyExists))
	})

	t.Run("GetUserByUsername with existing user", func(t *testing.T) {
		user, err := userStore.GetUserByUsername("johndoe")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "johndoe", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
		assert.NotEmpty(t, user.Password.Hash)
	})

	t.Run("GetUserByUsername with non-existing user", func(t *testing.T) {
		user, err := userStore.GetUserByUsername("nonexistent")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrNoRows))
		assert.Nil(t, user)
	})

	t.Run("UpdateUser with valid data", func(t *testing.T) {
		existingUser, _ := userStore.GetUserByUsername("johndoe")

		updatedUser := &User{
			Name:     "John Updated",
			Username: "johnupdated",
			Password: &valueObjects.Password{
				PlainText: "newpassword123",
			},
			Email: "johnupdated@example.com",
		}

		err := updatedUser.HashPassword()
		assert.NoError(t, err)

		err = userStore.UpdateUser(existingUser.ID, updatedUser)
		assert.NoError(t, err)

		retrievedUser, err := userStore.GetUserByUsername("johnupdated")
		assert.NoError(t, err)
		assert.Equal(t, "John Updated", retrievedUser.Name)
		assert.Equal(t, "johnupdated", retrievedUser.Username)
		assert.Equal(t, "johnupdated@example.com", retrievedUser.Email)
	})

	t.Run("UpdateUser with non-existing ID", func(t *testing.T) {
		user := &User{
			Name:     "Non Existing",
			Username: "nonexisting",
			Password: &valueObjects.Password{
				PlainText: "password123",
			},
			Email: "nonexisting@example.com",
		}

		err := user.HashPassword()
		assert.NoError(t, err)

		err = userStore.UpdateUser(99999, user)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrNoRows))
	})

	t.Run("GetUserFromToken with valid token", func(t *testing.T) {
		existingUser, token := testingUtils.CreateToken(db, "janedoe")

		retrievedUser, err := userStore.GetUserFromToken("authentication", token.Plaintext)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, existingUser.ID, retrievedUser.ID)
		assert.Equal(t, existingUser.Username, retrievedUser.Username)
	})

	t.Run("GetUserFromToken with invalid token", func(t *testing.T) {
		user, err := userStore.GetUserFromToken("authentication", "invalidtoken")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrNoRows))
		assert.Nil(t, user)
	})

	t.Run("GetUserFromToken with wrong scope", func(t *testing.T) {
		_, token := testingUtils.CreateToken(db, "janedoe")

		user, err := userStore.GetUserFromToken("wrongscope", token.Plaintext)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, internalErrors.ErrNoRows))
		assert.Nil(t, user)
	})

}
