package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/BelyaevEI/microservices_auth/internal/api/user"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	serviceMocks "github.com/BelyaevEI/microservices_auth/internal/service/mocks"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx      = context.Background()
		mc       = minimock.NewController(t)
		id       = gofakeit.Int64()
		name     = gofakeit.Animal()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = gofakeit.IntRange(1, 2)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			Password:        password,
			PasswordConfirm: password,
		}

		reqBadConfirmPass = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			Password:        password,
			PasswordConfirm: password + "a",
		}

		userCreate = model.UserCreate{
			Name:     name,
			Email:    email,
			Role:     model.Role(role),
			Password: password,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
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
				mock.CreateUserMock.Expect(ctx, &userCreate).Return(id, nil)
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
				mock.CreateUserMock.Expect(ctx, &userCreate).Return(0, serviceErr)
				return mock
			},
		},

		{
			name: "password and password confirm do not match",
			args: args{
				ctx: ctx,
				req: reqBadConfirmPass,
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "password and password confirm do not match"),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			createServiceMock := test.userServiceMock(mc)
			api := user.NewImplementation(createServiceMock)

			newID, err := api.CreateUser(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, newID)
		})
	}

}
