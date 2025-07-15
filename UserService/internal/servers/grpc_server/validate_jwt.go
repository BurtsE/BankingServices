package grpc_server

import (
	"UserService/generated/protobuf"
	"context"
)

func (s *Server) ValidateJWT(ctx context.Context, in *protobuf.ValidateRequest) (*protobuf.ValidateResponse, error) {
	uuid, err := s.service.ValidateJWT(ctx, in.Token)
	if err != nil {
		s.logger.Debugf("ValidateJWT error: %v", err)
		return nil, err
	}

	return &protobuf.ValidateResponse{
		Uuid: uuid,
	}, nil
}
