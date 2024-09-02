package auth

import (
	"context"

	errorAuth "github.com/BelyaevEI/microservices_auth/internal/errors"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/utils"
)

func (s serv) GetRefreshToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.refreshTokenSecretKey))
	if err != nil {
		return "", errorAuth.ErrInvalidRefreshToken
	}

	refreshToken, err = utils.GenerateToken(model.User{
		ID: claims.ID,
		Info: model.UserInfo{
			Role: claims.Role,
		},
	},
		[]byte(s.refreshTokenSecretKey),
		s.refreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
