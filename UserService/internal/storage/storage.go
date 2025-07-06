package storage

import (
	model "UserService/internal/domain"
	"context"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)       // Пустой указатель, если пользователь не найден
	FindByUsername(ctx context.Context, username string) (*model.User, error) // Пустой указатель, если пользователь не найден
	FindByID(ctx context.Context, userID string) (*model.User, error)
}
