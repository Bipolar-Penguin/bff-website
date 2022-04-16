package service

import (
	"errors"
	"time"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/Bipolar-Penguin/bff-website/pkg/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "foobar"
	tokenTTL   = 999 * time.Hour
	signingKey = "foobar"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

type Service struct {
	rep *repository.Repositories
}

func NewService(rep *repository.Repositories) *Service {
	return &Service{rep}
}

// User features

func (s *Service) SaveUser(user domain.User) (domain.User, error) {
	return s.rep.User.Save(user)
}

func (s *Service) Authenticate(authHeader string) (string, error) {
	token, err := jwt.ParseWithClaims(authHeader, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func (s *Service) GenerateToken(userID string) (string, error) {

	user, err := s.rep.User.Find(userID)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}
