package user

import (
	"context"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	token := req.GetRefreshToken()
	accessToken, err := i.userService.GetAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
