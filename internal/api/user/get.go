package user

import (
	"context"
	"log"

	"github.com/Timofey335/auth/internal/converter"
	"github.com/fatih/color"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (s *Implementation) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := s.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("Get user by id: %d", userObj.ID))

	userObjConvert := converter.ToUserFromService(userObj)

	return &desc.GetUserResponse{
		Id:        userObjConvert.Id,
		Name:      userObj.Name,
		Email:     userObj.Email,
		Role:      userObjConvert.Role,
		CreatedAt: userObjConvert.CreatedAt,
		UpdatedAt: userObjConvert.UpdatedAt,
	}, nil
}
