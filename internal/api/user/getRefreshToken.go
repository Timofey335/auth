package user

import (
	"context"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	token := req.GetRefreshToken()
	refreshToken, err := i.userService.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
