package wom

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"
	"strconv"
	"time"

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
	err := checkIntMigration("storage", "SELECT MAX(id) as MaxVersion FROM storage.migrations", 10)
	if err != nil {
		return err
	}
	err = checkIntMigration("realtime", "select max(version) as MaxVersion from realtime.schema_migrations", 20220712093339)
	if err != nil {
		return err
	}
	err = checkStringMigration("auth", "select max(version) as MaxVersion from auth.schema_migrations", "20221114143410")
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS supabase_migrations")
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
	version, dirty, err := m.Version()
	if err != nil {
		return err
	}
	log.Printf("DB Version: %d (%t).", version, dirty)
	if dirty {
		return errors.New("database is dirty")
	}
	return nil
}

func checkIntMigration(name, sql string, max int) error {
	var maxVersion int
	for maxVersion < max {
		row := pool.QueryRow(context.Background(), sql)
		err := row.Scan(&maxVersion)
		if err != nil {
			return fmt.Errorf("error checking %s migrations", name)
		}
		if maxVersion < max {
			log.Printf("Waiting for %s migration", name)
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func checkStringMigration(name, sql string, max string) error {
	maxInt, _ := strconv.Atoi(max)
	var maxVersion int
	for maxVersion < maxInt {
		var result string
		row := pool.QueryRow(context.Background(), sql)
		err := row.Scan(&result)
		if err != nil {
			return fmt.Errorf("error checking %s migrations", name)
		}
		maxVersion, err = strconv.Atoi(result)
		if err != nil {
			return fmt.Errorf("error checking %s migrations", name)
		}
		if maxVersion < maxInt {
			log.Printf("Waiting for %s migration", name)
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func CloseDatabase() {
	pool.Close()
}
