package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MixinNetwork/surfing/session"
	"github.com/jackc/pgx/v4"
)

const ()

type Property struct {
	Key       string
	Value     string
	UpdatedAt time.Time
}

func readProperty(ctx context.Context, tx pgx.Tx, key string) (string, error) {
	var p Property
	query := "SELECT key,value,updated_at FROM properties WHERE key=$1"
	err := tx.QueryRow(ctx, query, key).Scan(&p.Key, &p.Value, &p.UpdatedAt)
	return p.Value, err
}

func ReadProperty(ctx context.Context, key string) (string, error) {
	var v string
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var err error
		v, err = readProperty(ctx, tx, key)
		return err
	})
	if err != nil {
		return "", session.TransactionError(ctx, err)
	}
	return v, nil
}

func ReadProperties(ctx context.Context, keys []string) (map[string]string, error) {
	rows, err := session.Database(ctx).Query(ctx, fmt.Sprintf("SELECT key,value FROM properties WHERE key IN ('%s')", strings.Join(keys, "','")))
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}

	set := make(map[string]string)
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		set[key] = value
	}
	return set, nil
}

func writeProperty(ctx context.Context, tx pgx.Tx, key, value string) error {
	query := "INSERT INTO properties (key,value,updated_at) VALUES($1,$2,$3) ON CONFLICT (key) DO UPDATE SET (value,updated_at)=(EXCLUDED.value, EXCLUDED.updated_at)"
	_, err := tx.Exec(ctx, query, key, value, time.Now())
	return err
}

func WriteProperty(ctx context.Context, key, value string) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return writeProperty(ctx, tx, key, value)
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}
