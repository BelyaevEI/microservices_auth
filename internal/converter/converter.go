package converter

import (
	"github.com/BelyaevEI/microservices_auth/internal/model"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserCreateFromDesc converts a UserCreate object from the desc package to a UserCreate object from the model package.
func ToUserCreateFromDesc(userCreate *desc.CreateRequest) *model.UserCreate {
	return &model.UserCreate{
		Name:     userCreate.GetInfo().GetName(),
		Email:    userCreate.GetInfo().GetEmail(),
		Role:     model.Role(userCreate.GetInfo().GetRole()),
		Password: userCreate.GetPassword(),
	}
}

// ToUserFromService converts a User object from the model package to a User object from the desc package.
func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp

	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id: user.ID,
		Info: &desc.UserInfo{
			Name:  user.Info.Name,
			Email: user.Info.Email,
			Role:  desc.Role(user.Info.Role),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserUpdateFromDesc converts a UserUpdate object from the desc package to a UserUpdate object from the model package.
func ToUserUpdateFromDesc(userUpdate *desc.UpdateRequest) *model.UserUpdate {
	var (
		email *string
		name  *string
		role  model.Role
	)

	if userUpdate.Info.GetEmail() != nil {
		email = &userUpdate.Info.GetEmail().Value
	}
	if userUpdate.Info.GetName() != nil {
		name = &userUpdate.Info.GetName().Value
	}
	if userUpdate.Info.GetRole() != desc.Role_UNKNOWN {
		role = (model.Role)(userUpdate.Info.GetRole())
	}
	return &model.UserUpdate{
		Email: email,
		Name:  name,
		Role:  role,
	}
}
