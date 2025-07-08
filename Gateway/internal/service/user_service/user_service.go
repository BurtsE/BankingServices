package user_service

import (
	"context"
	"gateway/generated/protobuf"
	"gateway/internal/service"
	"time"
)

var _ service.IUserService = (*UserService)(nil)

type UserService struct {
	client protobuf.UserServiceClient
}

func NewUserService(client protobuf.UserServiceClient) *UserService {
	return &UserService{client: client}
}

func (u *UserService) Validate(jwtToken string) (uuid string, err error) {

	validateRequest := &protobuf.ValidateRequest{Token: jwtToken}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := u.client.ValidateJWT(ctx, validateRequest)
	if err != nil {
		return "", err
	}

	return resp.GetUuid(), nil
}
