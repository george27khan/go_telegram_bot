package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPostgresPool отдельный конструктор для пула, т.к. он должен быть общим для всех репозиториев
func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	//user := os.Getenv("POSTGRES_USER")
	//pass := os.Getenv("POSTGRES_PASSWORD")
	//db := os.Getenv("POSTGRES_DB")
	//host := os.Getenv("POSTGRES_HOST")
	//port := os.Getenv("POSTGRES_PORT")
	user := "postgres"
	pass := "postgres"
	db := "postgres"
	host := "localhost"
	port := "50270"
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, db,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

// NewRepository конструктор для репозитория
func NewRepository(pool *pgxpool.Pool, logger *slog.Logger) *DB {
	return &DB{Pool: pool, logger: logger}
}

// Close закрытие пула соединений
func (r *DB) Close() {
	r.Pool.Close()
}
