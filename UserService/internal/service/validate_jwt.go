package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) ValidateJWT(ctx context.Context, jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		UUID, err := claims.GetSubject()
		if err != nil {
			return "", errors.New("user_id missing in token")
		}

		_, err = s.repo.FindByID(ctx, UUID)
		if err != nil {
			return "", err
		}
		
		return UUID, nil
	}
	return "", errors.New("invalid token")
}
