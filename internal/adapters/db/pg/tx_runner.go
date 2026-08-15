package pg

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxRunner struct {
	pool *pgxpool.Pool
}

func NewTxRunner(pool *pgxpool.Pool) *TxRunner {
	return &TxRunner{pool: pool}
}

// RunTx executes fn as an atomic transaction on the database.
func (r *TxRunner) RunTx(ctx context.Context, fn func(ports.Querier) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
