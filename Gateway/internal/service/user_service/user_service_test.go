package user_service

import (
	"context"
	"errors"
	"gateway/generated/protobuf"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type mockUserServiceClient struct {
	mock.Mock
}

func (m *mockUserServiceClient) ValidateJWT(ctx context.Context, req *protobuf.ValidateRequest, opts ...grpc.CallOption) (*protobuf.ValidateResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*protobuf.ValidateResponse), args.Error(1)
}

func TestUserService_Validate_Success(t *testing.T) {
	mockClient := new(mockUserServiceClient)
	service := NewUserService(mockClient)

	testToken := "test.jwt.token"
	expectedUUID := "123e4567-e89b-12d3-a456-426614174000"

	mockClient.On("ValidateJWT", mock.Anything, &protobuf.ValidateRequest{Token: testToken}).
		Return(&protobuf.ValidateResponse{Uuid: expectedUUID}, nil)

	uuid, err := service.Validate(testToken)

	assert.NoError(t, err)
	assert.Equal(t, expectedUUID, uuid)
	mockClient.AssertExpectations(t)
}

func TestUserService_Validate_Error(t *testing.T) {
	mockClient := new(mockUserServiceClient)
	service := NewUserService(mockClient)

	testToken := "invalid.jwt.token"
	expectedErr := errors.New("invalid token")

	mockClient.On("ValidateJWT", mock.Anything, &protobuf.ValidateRequest{Token: testToken}).
		Return(nil, expectedErr)

	uuid, err := service.Validate(testToken)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, uuid)
	mockClient.AssertExpectations(t)
}
