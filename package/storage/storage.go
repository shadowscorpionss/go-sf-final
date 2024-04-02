// Пакет для работы с БД приложения GoNews.
package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// database structure
type DB struct {
	pool *pgxpool.Pool
}

// database constructor
func New(connstr string) (*DB, error) {
	if connstr == "" {
		return nil, errors.New("database connection string is empty")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

func NullInt32(i int32) sql.NullInt32 {
	if i < 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

func NullInt64(i int64) sql.NullInt64 {
	if i < 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
