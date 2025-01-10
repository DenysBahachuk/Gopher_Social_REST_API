package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "testSecret"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-iss",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
}

type TestAuthenticator struct{}

func (a *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString, nil
}
func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	},
	)
}
