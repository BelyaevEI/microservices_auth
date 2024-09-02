package auth

import (
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	"github.com/BelyaevEI/platform_common/pkg/db"
)

const (
	tableName       = "user"
	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passHashColumn  = "pass_hash"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates a new user repository.
func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}
