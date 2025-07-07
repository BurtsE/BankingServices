package user_service

import "gateway/internal/service"

var _ service.IUserService = (*UserService)(nil)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Validate(jwtToken string) (uuid string, err error) {
	return "", nil
}
