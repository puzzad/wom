package wom

import (
	"context"
	"embed"
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbUrl = flag.String("database-url", "", "DSN to use to connect to the database")

	pool *pgxpool.Pool

	//go:embed migrations
	migrations embed.FS
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

func EnsureMigrations() error {
	_, err := pool.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS supabase_migrations")
	if err != nil {
		return err
	}
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d,
		*dbUrl+"?x-migrations-table=\"supabase_migrations\".\"schema_migrations\"&x-migrations-table-quoted=true")
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func CloseDatabase() {
	pool.Close()
}
