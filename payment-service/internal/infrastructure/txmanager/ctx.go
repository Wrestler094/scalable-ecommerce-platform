package txmanager

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ctxKeyTx struct{}

func InjectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

func ExtractTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(ctxKeyTx{}).(*sqlx.Tx)
	return tx, ok
}
