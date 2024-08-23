package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

func (s *serv) UpdateUser(ctx context.Context, user *model.User) (*emptypb.Empty, error) {
	_, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
