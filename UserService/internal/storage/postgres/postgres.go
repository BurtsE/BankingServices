package postgres

import (
	"UserService/internal/config"
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/metric"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, cfg *config.Config) (*PostgresRepository, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)

	pgxConfig, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithDisableConnectionDetailsInAttributes())

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	repo := &PostgresRepository{
		pool: pool,
	}

	return repo, nil
}

func (p *PostgresRepository) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}
