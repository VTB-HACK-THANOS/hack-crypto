package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/anthonyaspen/emlvid-back/storage/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

type Storage struct {
	db *bun.DB
} // New constructs a new storage with bun driver to postgres database.
func New(address, database, user, password string, poolSize int, envName, appName string) (*Storage, error) {
	db := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(address),
		pgdriver.WithInsecure(true),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database), pgdriver.WithApplicationName(appName),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	))
	store := bun.NewDB(db, pgdialect.New(), bun.WithDiscardUnknownColumns())
	store.SetMaxOpenConns(poolSize)
	store.AddQueryHook(bundebug.NewQueryHook())

	// if envName == "dev" {
	// 	store.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	// }

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return &Storage{db: store}, nil
}

// Close closes all connections to a storage.
func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

// Migrate runs migrations to a database.
func (s *Storage) Migrate() (int64, error) {
	if s.db == nil {
		return 0, errors.New("database is not initialized")
	}

	migrator, err := runInitCommand(s.db)
	if err != nil {
		return 0, err
	}

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		return 0, err
	}

	return group.Migrations.LastGroupID(), nil

}

func runInitCommand(db *bun.DB) (*migrate.Migrator, error) {
	migrator := migrate.NewMigrator(db, migrations.Migration)
	if err := migrator.Init(context.Background()); err != nil {
		return nil, err
	}
	return migrator, nil
}
