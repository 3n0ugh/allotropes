package token

import (
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email string
	Role  int8
	jwt.StandardClaims
}

func NewToken(cfg config.Config, claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.Application.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *Claims) ParseToken(token string) error {
	tkn, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	if err != nil {
		return errors.Wrap(err, "token parse")
	}

	if !tkn.Valid {
		return errors.Wrap(err, "token validation")
	}

	return nil
}
