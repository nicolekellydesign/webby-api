package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var signingKey string

type Auth struct {
	jwt.StandardClaims
	UserID uint `json:"id" db:"user_id"`
}

func GenerateToken(user User) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Auth{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiresAt,
		},
		UserID: user.ID,
	})

	return token.SignedString(signingKey)
}

func ValidateToken(tokenString string) (uint, string, error) {
	var claims Auth
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	id := claims.UserID
	username := claims.Subject
	return id, username, nil
}

// User represents a user that can manage Webby.
type User struct {
	ID       uint   `json:"id,omitempty" db:"id"`
	Username string `json:"username" db:"user_name"`
	Password string `json:"password,omitempty" db:"pwdhash"`
}
