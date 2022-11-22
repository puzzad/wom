package wom

import (
	"context"
	"flag"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbUrl = flag.String("database-url", "", "DSN to use to connect to the database")

	pool *pgxpool.Pool
)

func ConnectToDatabase() error {
	var err error
	pool, err = pgxpool.New(context.Background(), *dbUrl)
	if err != nil {
		return err
	}
	return pool.Ping(context.Background())
}

func addEmailToMailingList(ctx context.Context, email string) error {
	_, err := pool.Exec(ctx, "INSERT INTO internal.mailinglist (email) VALUES ($1) ON CONFLICT DO NOTHING", email)
	return err
}

func removeEmailFromMailingList(ctx context.Context, email string) error {
	_, err := pool.Exec(ctx, "DELETE FROM internal.mailinglist WHERE email = $1", email)
	return err
}

func CloseDatabase() {
	pool.Close()
}
