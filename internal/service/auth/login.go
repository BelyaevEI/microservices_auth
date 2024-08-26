package auth

import (
	"context"

	errorAuth "github.com/BelyaevEI/microservices_auth/internal/errors"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/utils"
)

func (s serv) Login(ctx context.Context, login *model.UserLogin) (string, error) {
	var refreshToken string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		passHash, errTx := s.authRepository.GetPassword(ctx, login.ID)
		if errTx != nil {
			return errTx
		}

		if !utils.VerifyPassword(string(passHash), login.Password) {
			return errorAuth.ErrWrongPassword
		}
		user, errTx := s.authRepository.Get(ctx, login.ID)
		if errTx != nil {
			return errTx
		}

		refreshToken, errTx = utils.GenerateToken(*user,
			[]byte(s.refreshTokenSecretKey),
			s.refreshTokenExpiration,
		)
		if errTx != nil {
			return errorAuth.ErrGenerateToken
		}

		return nil
	})
	return refreshToken, err
}
