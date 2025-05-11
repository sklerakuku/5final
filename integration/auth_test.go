package integration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sklerakuku/5final/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	// Setup
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE users (login TEXT PRIMARY KEY, password_hash TEXT)")
	assert.NoError(t, err)

	auth.SetDB(db)
	auth.SetJWTSecret("testsecret")

	// Test registration
	t.Run("Successful registration", func(t *testing.T) {
		err := auth.Register("user1", "password123")
		assert.NoError(t, err)
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		err := auth.Register("user1", "password123")
		assert.Error(t, err)
	})

	// Test login
	t.Run("Successful login", func(t *testing.T) {
		token, err := auth.Login("user1", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		_, err := auth.Login("user1", "wrongpassword")
		assert.Error(t, err)
	})
}
