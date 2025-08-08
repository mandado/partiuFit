package store

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	testingUtils "partiuFit/internal/database/testing_utils"
	"partiuFit/internal/tokens"
	"partiuFit/internal/utils"
	"testing"
	"time"
)

func TestTokensStore(t *testing.T) {
	db := utils.Must(testingUtils.SetupTestDB())
	utils.MustIfError(testingUtils.SeedDB(db))
	defer testingUtils.TeardownTestDB(db)

	tokensStore := NewPostgresTokensStore(db)

	// Get actual user IDs from seeded data
	var johnID, janeID int
	err := db.QueryRow("SELECT id FROM users WHERE username = 'johndoe'").Scan(&johnID)
	assert.NoError(t, err)
	err = db.QueryRow("SELECT id FROM users WHERE username = 'janedoe'").Scan(&janeID)
	assert.NoError(t, err)

	t.Run("CreateToken with valid data", func(t *testing.T) {
		token, err := tokensStore.CreateToken(johnID, 24*time.Hour, "authentication")

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, johnID, token.UserID)
		assert.Equal(t, "authentication", token.Scope)
		assert.NotEmpty(t, token.Hash)
		assert.NotEmpty(t, token.Plaintext)
		assert.True(t, token.ExpiresAt.After(time.Now()))
	})

	t.Run("InsertToken with valid token", func(t *testing.T) {
		token, err := tokens.GenerateToken(janeID, 12*time.Hour, "password-reset")
		assert.NoError(t, err)

		err = tokensStore.InsertToken(token)
		assert.NoError(t, err)
	})

	t.Run("DeleteAllTokensForUser", func(t *testing.T) {
		// Create some tokens for user
		_, err := tokensStore.CreateToken(johnID, 24*time.Hour, "authentication")
		assert.NoError(t, err)

		_, err = tokensStore.CreateToken(johnID, 24*time.Hour, "password-reset")
		assert.NoError(t, err)

		// Delete all tokens for user
		err = tokensStore.DeleteAllTokensForUser(johnID)
		assert.NoError(t, err)

		// Verify tokens were deleted by trying to use them
		// This would typically be tested through the user store's GetUserFromToken method
	})

	t.Run("CreateToken with invalid user ID", func(t *testing.T) {
		token, err := tokensStore.CreateToken(99999, 24*time.Hour, "authentication")

		// The error comes from foreign key constraint when inserting
		// But the token is still generatedX and returned even if insertion fails
		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("DeleteAllTokensForUser with non-existing user", func(t *testing.T) {
		err := tokensStore.DeleteAllTokensForUser(99999)

		// Should not error even if no tokens exist for the user
		assert.NoError(t, err)
	})

	t.Run("CreateToken with zero TTL", func(t *testing.T) {
		token, err := tokensStore.CreateToken(janeID, 0, "authentication")

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.ExpiresAt.Before(time.Now()) || token.ExpiresAt.Equal(time.Now()))
	})

	t.Run("CreateToken with negative TTL", func(t *testing.T) {
		token, err := tokensStore.CreateToken(janeID, -1*time.Hour, "authentication")

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.ExpiresAt.Before(time.Now()))
	})
}
