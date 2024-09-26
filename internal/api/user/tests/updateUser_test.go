package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Timofey335/auth/internal/api/user"
	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/service"
	serviceMocks "github.com/Timofey335/auth/internal/service/mocks"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func TestUpadateUser(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = int64(gofakeit.Number(1, 2))
		roleReq  = desc.Role(role)

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateUserRequest{
			Id:              id,
			Name:            wrapperspb.String(name),
			Password:        wrapperspb.String(password),
			PasswordConfirm: wrapperspb.String(password),
			Role:            &roleReq,
		}

		userData = &model.UserUpdateModel{
			ID:              id,
			Name:            &name,
			Password:        &password,
			PasswordConfirm: &password,
			Role:            &role,
		}

		res = &emptypb.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, userData).Return(res, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, userData).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			resHandler, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})

	}
}
