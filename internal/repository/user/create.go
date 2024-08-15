package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/repository/user/converter"
	"github.com/BelyaevEI/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

func (r *repo) CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	userCreateRepo := converter.ToUserCreateFromService(userCreate)

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passHashColumn, roleColumn).
		Values(userCreateRepo.Name, userCreateRepo.Email, userCreateRepo.PasswordHash, userCreateRepo.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
