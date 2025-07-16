package grpc_server

import (
	"BankingService/generated/protobuf"
	"context"
)

func (s *Server) AccountIsActive(ctx context.Context, in *protobuf.IsActiveRequest) (*protobuf.IsActiveResponse, error) {
	isActive, err := s.service.AccountIsActive(ctx, in.AccountId)
	if err != nil {
		s.logger.Debugf("AccountIsActive error: %v", err)
		return nil, err
	}

	return &protobuf.IsActiveResponse{
		IsActive: isActive,
	}, nil
}
