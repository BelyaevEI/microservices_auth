package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/BelyaevEI/microservices_auth/internal/cache"
	cacheMocks "github.com/BelyaevEI/microservices_auth/internal/cache/mocks"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	repoMocks "github.com/BelyaevEI/microservices_auth/internal/repository/mocks"
	userService "github.com/BelyaevEI/microservices_auth/internal/service/user"
	"github.com/BelyaevEI/platform_common/pkg/db"
	"github.com/BelyaevEI/platform_common/pkg/db/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type cacheMockFunc func(mc *minimock.Controller) cache.UserCache
	type txManagerMockFunc func(f func(context.Context) error, mc *minimock.Controller) db.TxManager

	type args struct {
		ctx         context.Context
		userRepoReq *model.UserCreate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Animal()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = gofakeit.IntRange(1, 2)

		repoErr = fmt.Errorf("repo error")

		userRepoReq = &model.UserCreate{
			Name:     name,
			Email:    email,
			Role:     model.Role(role),
			Password: password,
		}

		cacheUser = &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:  userRepoReq.Name,
				Email: userRepoReq.Email,
				Role:  userRepoReq.Role,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		cacheMock          cacheMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case 1",
			args: args{
				ctx:         ctx,
				userRepoReq: userRepoReq,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, userRepoReq).Return(id, nil)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.CreateUserMock.Expect(ctx, cacheUser).Return(nil)
				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, f).Return(nil)
				return mock
			},
		},

		{
			name: "service error case",
			args: args{
				ctx:         ctx,
				userRepoReq: userRepoReq,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, userRepoReq).Return(0, repoErr)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.CreateUserMock.Expect(ctx, cacheUser).Return(nil)
				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, f).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			cacheMock := tt.cacheMock(mc)
			txManagerMock := tt.txManagerMock(func(ctx context.Context) error {
				return nil
			}, mc)
			service := userService.NewService(userRepoMock, txManagerMock, cacheMock)
			newID, err := service.CreateUser(tt.args.ctx, tt.args.userRepoReq)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}

}
