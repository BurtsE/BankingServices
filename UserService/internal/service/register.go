package service

import (
	model "UserService/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) Register(ctx context.Context, email, username, password, fullName string) (*model.User, error) {
	// Check for unique email and username
	if user, err := s.repo.FindByEmail(ctx, email); user != nil {
		return nil, errors.New("email already registered")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user, err := s.repo.FindByUsername(ctx, username); user != nil {
		return nil, errors.New("username taken")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UUID:         uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
