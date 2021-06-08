package durable

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	*pgxpool.Pool
}

func NewDatabase(ctx context.Context, db *pgxpool.Pool) (*Database, error) {
	return &Database{db}, db.Ping(ctx)
}

func (d *Database) RunInTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := d.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		return tx.Rollback(ctx)
	}
	return tx.Commit(ctx)
}

type Row interface {
	Scan(dest ...interface{}) error
}
