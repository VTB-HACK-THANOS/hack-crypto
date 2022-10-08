package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/VTB-HACK-THANOS/hack-crypto/storage/migrations"
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

const (
	ctxTransaction = "tx"
)

// Commit commits transaction, takes it from context.
func (s *Storage) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(ctxTransaction).(*bun.Tx)
	if !ok {
		return errors.New("failed to get transaction from context")
	}

	return tx.Commit()
}

// Rollback rollbacks transaction, takes it from context.
func (s *Storage) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(ctxTransaction).(*bun.Tx)
	if !ok {
		return errors.New("failed to get transaction from context")
	}

	err := tx.Rollback()
	if err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return err
	}

	return nil
}

// contextTransaction takes in context and check if ctx has a pointer to a transaction.
// If it is - returns the transaction, otherwise creates new one.
func (s *Storage) contextTransaction(ctx context.Context) (*bun.Tx, error) {
	var (
		tx  *bun.Tx
		ok  bool
		err error
	)
	tx, ok = ctx.Value(ctxTransaction).(*bun.Tx)
	if !ok { // start a transaction if context doesn't have a pointer to tx.
		tx, err = s.beginTx(ctx)
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// BeginTX returns pointer to a new tx.
func (s *Storage) beginTx(ctx context.Context) (*bun.Tx, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

// BeginTx takes in context, returns given context with a pointer to a transaction.
func (s *Storage) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, ctxTransaction, &tx), nil
}
