// Пакет для работы с БД приложения GoNews.
package storage

import (
	"context"
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
