package postgresql

import (
	"api-gateway/internal/pkg/postgres"
	"api-gateway/internal/usecase/refresh_token"
)

func NewRefreshTokenRepo(db *postgres.PostgresDB) refresh_token.RefreshTokenRepo {
	return nil
}
