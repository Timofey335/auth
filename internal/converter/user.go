package converter

import (
	"github.com/Timofey335/auth/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func ToAuthToService(user *model.User) *desc.GetUserResponse {
	var updated_at *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updated_at = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updated_at,
	}
}
