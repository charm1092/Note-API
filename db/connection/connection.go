package connection

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	connection := os.Getenv("CONNECTION")
	return pgxpool.New(ctx, connection)
}