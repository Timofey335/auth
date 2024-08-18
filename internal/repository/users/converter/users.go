package converter

import (
	"github.com/Timofey335/auth/internal/repository/users/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func ToUsersFromRepo(users *model.Users) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if users.UpdatedAt.Valid {
		updatedAt = timestamppb.New(users.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        users.ID,
		Name:      users.Name,
		Email:     users.Email,
		Role:      desc.Role(users.Role),
		CreatedAt: timestamppb.New(users.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
