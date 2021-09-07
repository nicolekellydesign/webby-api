package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nicolekellydesign/webby-api/entities"
)

// signingKey is used to generate tokens. It is set at build time with ldflags.
var signingKey string

var errInvalidToken = errors.New("invalid token")

// authClaims holds authenitcation data for a user.
type authClaims struct {
	jwt.StandardClaims
	UserID uint `json:"id"`
}

// generateToken creates a new authentication token that expires after 24 hours.
func generateToken(user entities.User) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiresAt,
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

// validateToken parses a token string and checks if it is valid.
func validateToken(tokenString string) (uint, string, error) {
	if tokenString == "" {
		return 0, "", errInvalidToken
	}

	var claims authClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errInvalidToken
	}

	id := claims.UserID
	username := claims.Subject
	return id, username, nil
}
