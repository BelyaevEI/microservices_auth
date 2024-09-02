package user

import (
	"context"
	"time"

	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/pkg/user_v1"
	"github.com/BelyaevEI/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

func (r *repo) UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error {

	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar)

	if userUpdate.Name != nil {
		builderUpdate = builderUpdate.
			Set(nameColumn, userUpdate.Name)
	}
	if userUpdate.Email != nil {
		builderUpdate = builderUpdate.
			Set(emailColumn, userUpdate.Email)
	}
	if userUpdate.Role != model.Role(user_v1.Role_UNKNOWN) {
		builderUpdate = builderUpdate.
			Set(roleColumn, userUpdate.Role)
	}

	builderUpdate = builderUpdate.
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderUpdate.ToSql()

	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
