package service

import (
	"UserService/internal/config"
	model "UserService/internal/domain"
	"UserService/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	Register(ctx context.Context, email, username, password, fullName string) (*model.User, error)
	Authenticate(ctx context.Context, email, password string) (jwtToken string, expiresAt time.Time, err error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

var _ UserService = (*Service)(nil)

type Service struct {
	repo      storage.UserStorage
	jwtSecret []byte
}

func NewUserService(repo storage.UserStorage, cfg *config.Config) *Service {
	return &Service{repo: repo, jwtSecret: []byte(config.GetJWTSecretKey())}
}

func (s *Service) Register(ctx context.Context, email, username, password, fullName string) (*model.User, error) {
	// Проверка уникальности email и username
	if user, err := s.repo.FindByEmail(ctx, email); user != nil {
		return nil, errors.New("email уже зарегистрирован")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	if user, err := s.repo.FindByUsername(ctx, username); user != nil {
		return nil, errors.New("username уже занят")
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

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, time.Time, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", time.Time{}, errors.New("пользователя с таким email не существует")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("неверный пароль")
	}

	// Генерируем JWT
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

func (s *Service) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.FindByID(ctx, userID)
}
