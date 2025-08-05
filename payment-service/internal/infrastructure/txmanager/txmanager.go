package txmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
)

type TxManager struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewTxManager(db *sqlx.DB, logger logger.Logger) *TxManager {
	return &TxManager{db: db, logger: logger}
}

func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	const op = "txmanager.WithinTx"

	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			m.logger.WithOp(op).Warn(fmt.Sprintf("%s: failed to rollback transaction: %w", op, err))
		}
	}()

	txCtx := InjectTx(ctx, tx)

	if err := fn(txCtx); err != nil {
		return fmt.Errorf("%s: failed to execute transaction function: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
