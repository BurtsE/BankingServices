package postgres

import (
	"BankingService/internal/config"
	"BankingService/internal/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/openpgp/packet"
)

var _ storage.BankingStorage = (*PostgresRepository)(nil)

type PostgresRepository struct {
	config *packet.Config
	pool   *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, cfg *config.Config) (*PostgresRepository, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	repo := &PostgresRepository{
		pool: pool,
	}

	return repo, nil
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}

func (p *PostgresRepository) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTransaction: %w", err)
	}
	return tx, nil
}
