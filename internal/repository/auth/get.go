package auth

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/repository/user/converter"
	modelRepo "github.com/BelyaevEI/microservices_auth/internal/repository/user/model"

	"github.com/BelyaevEI/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {

	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})
	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}
	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}
	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
