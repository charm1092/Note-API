package connection

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(
		ctx, "postgres://postgres:1234@localhost:5432/postgres")
}