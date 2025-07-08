package service

import (
	"UserService/internal/config"
	model "UserService/internal/domain"
	"UserService/internal/storage"
	"context"
	"time"
)

type UserService interface {
	Register(ctx context.Context, email, username, password, fullName string) (*model.User, error)
	Authenticate(ctx context.Context, email, password string) (jwtToken string, expiresAt time.Time, err error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	ValidateJWT(ctx context.Context, jwtToken string) (uuid string, err error)
}

var _ UserService = (*Service)(nil)

type Service struct {
	repo      storage.UserStorage
	jwtSecret []byte
}

func NewUserService(repo storage.UserStorage, cfg *config.Config) *Service {
	return &Service{repo: repo, jwtSecret: []byte(config.GetJWTSecretKey())}
}

func (s *Service) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.FindByID(ctx, userID)
}
