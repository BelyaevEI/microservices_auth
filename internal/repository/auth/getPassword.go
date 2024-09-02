package auth

import (
	"context"

	"github.com/BelyaevEI/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

func (r *repo) GetPassword(ctx context.Context, id int64) (string, error) {
	builderSelectOne := sq.Select(passHashColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})
	query, args, err := builderSelectOne.ToSql()
	if err != nil {

		return "", err
	}
	q := db.Query{
		Name:     "user_repository.GetPassword",
		QueryRaw: query,
	}

	var password string
	err = r.db.DB().ScanOneContext(ctx, &password, q, args...)
	if err != nil {

		return "", err
	}

	return password, nil
}
