package auth

import (
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	db        *sql.DB
	jwtSecret string
)

func SetDB(database *sql.DB) {
	db = database
}

func SetJWTSecret(secret string) {
	jwtSecret = secret
}

func Register(login, password string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", login).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", login, hashedPassword)
	return err
}

func Login(login, password string) (string, error) {
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE login = ?", login).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,
	})

	return token.SignedString([]byte(jwtSecret))
}
