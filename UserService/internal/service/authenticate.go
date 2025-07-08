package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, time.Time, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("database error: %v", err)
	}

	if user == nil {
		return "", time.Time{}, errors.New("user with such an email does not exist")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("wrong password")
	}

	// Generate JWT
	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, exp, nil
}
