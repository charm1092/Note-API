package connection

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"
// postgres://postgres:1234@localhost:5432/postgres

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	connection := os.Getenv("connection")
	return pgxpool.New(ctx, connection)
}