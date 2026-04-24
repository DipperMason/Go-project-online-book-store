package jwt

import (
	"errors"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret         string
	AccessTokenTTL int
}

type TokenPair struct {
	Token string
}

type Service struct {
	config Config
}

type UserClaims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`

	golangjwt.RegisteredClaims
}

func NewService(config Config) Service {
	return Service{config: config}
}

func (j *Service) GenerateToken(claims UserClaims) (TokenPair, error) {
	claims.RegisteredClaims.ExpiresAt = golangjwt.NewNumericDate(
		time.Now().Add(time.Second * time.Duration(j.config.AccessTokenTTL)),
	)

	token := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		Token: tokenString,
	}, nil
}

func (j *Service) ValidateAndParseToken(pair TokenPair) (UserClaims, error) {
	claims := UserClaims{}
	_, err := golangjwt.ParseWithClaims(pair.Token, &claims, func(token *golangjwt.Token) (any, error) {
		if token.Method != golangjwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return UserClaims{}, err
	}

	return claims, err
}
