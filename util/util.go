package util

import (
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomBytes ...
// Generate random bytes of len x
func RandomBytes(strlen int) []byte {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return result
}

// ParseJwtToken ...
// Parse a jwt tokenstring and return a *jwt.Token
func ParseJwtToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	parseFunc := func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(jwtSecret), nil
	}

	token, err := jwt.Parse(tokenString, parseFunc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse jwt token")
	}

	if !token.Valid {
		return nil, errors.New("non-valid jwt token")
	}
	return token, nil
}
