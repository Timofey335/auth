package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Timofey335/platform_common/pkg/db"
	dbMocks "github.com/Timofey335/platform_common/pkg/db/mocks"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/cache"
	cacheMocks "github.com/Timofey335/auth/internal/cache/mocks"
	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/repository"
	repoMocks "github.com/Timofey335/auth/internal/repository/mocks"
	"github.com/Timofey335/auth/internal/service/user"
)

func TestUpdateUser(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type userCacheMockFunc func(mc *minimock.Controller) cache.UserCache
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UserUpdateModel
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = int64(gofakeit.Number(1, 2))

		serviceErr = fmt.Errorf("service error")

		req = &model.UserUpdateModel{
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
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		userRepositoryMock userRepositoryMockFunc
		userCacheMock      userCacheMockFunc
		txManagaerMock     txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
			txManagaerMock: func(mc *minimock.Controller) db.TxManager {
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
			want: nil,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(nil, serviceErr)
				return mock
			},
			userCacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				return mock
			},
			txManagaerMock: func(mc *minimock.Controller) db.TxManager {
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
			txManagerMock := tt.txManagaerMock(mc)

			service := user.NewService(userRepoMock, userCacheMock, txManagerMock)

			resHandler, err := service.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}

}
