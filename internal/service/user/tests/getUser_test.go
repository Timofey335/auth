package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Timofey335/platform_common/pkg/db"
	dbMocks "github.com/Timofey335/platform_common/pkg/db/mocks"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Timofey335/auth/internal/cache"
	cacheMocks "github.com/Timofey335/auth/internal/cache/mocks"
	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/repository"
	repoMocks "github.com/Timofey335/auth/internal/repository/mocks"
	"github.com/Timofey335/auth/internal/service/user"
)

func TestGetUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type userCacheMockFunc func(mc *minimock.Controller) cache.UserCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = int64(gofakeit.Number(1, 2))

		serviceErr = fmt.Errorf("service error")

		res = &model.UserModel{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.UserModel
		err                error
		userRepositoryMock userRepositoryMockFunc
		userCacheMock      userCacheMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
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

			service := user.NewService(userRepoMock, userCacheMock, txManagerMock)

			resHandler, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
