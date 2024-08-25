package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	userApi "github.com/BelyaevEI/microservices_auth/internal/api/user"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	serviceMocks "github.com/BelyaevEI/microservices_auth/internal/service/mocks"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Animal()
		email     = gofakeit.Email()
		role      = gofakeit.IntRange(1, 2)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetRequest{
			Id: id,
		}

		res = &desc.GetResponse{
			User: &desc.User{
				Id: id,
				Info: &desc.UserInfo{
					Name:  name,
					Email: email,
					Role:  desc.Role(role),
				},
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
		user = model.User{
			ID: id,
			Info: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  model.Role(role),
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
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
				mock.GetUserByIDMock.Expect(ctx, id).Return(&user, nil)
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
				mock.GetUserByIDMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			getServiceMock := test.userServiceMock(mc)
			api := userApi.NewImplementation(getServiceMock)

			newID, err := api.GetUserByID(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, newID)
		})
	}

}
