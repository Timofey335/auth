package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Timofey335/auth/internal/cache"
	cacheMocks "github.com/Timofey335/auth/internal/cache/mocks"
	"github.com/Timofey335/auth/internal/client/db"
	dbMocks "github.com/Timofey335/auth/internal/client/db/mocks"
	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/repository"
	repoMocks "github.com/Timofey335/auth/internal/repository/mocks"
	"github.com/Timofey335/auth/internal/service/user"
)

func TestCreateUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type userCacheMockFunc func(mc *minimock.Controller) cache.UserCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UserModel
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = int64(gofakeit.Number(1, 2))

		serviceErr = fmt.Errorf("service error")

		req = &model.UserModel{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		userCacheMock      userCacheMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				// mock.CreateUserMock.Expect(ctx, req).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(0, serviceErr)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				// mock.CreateUserMock.Expect(ctx, req).Return(serviceErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepositoryMock(mc)
			userCacheMock := tt.userCacheMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := user.NewMockService(userRepoMock, userCacheMock, txManagerMock)

			resHandler, err := service.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
