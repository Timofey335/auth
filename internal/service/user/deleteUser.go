package user

import (
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error) {
	_, err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
