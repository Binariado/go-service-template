package database

import (
	"context"
	"database/sql"
)

type transactionContextKey struct{}

// WithTransaction returns a context that carries tx so repositories can join
// an ongoing transaction instead of opening their own.
func WithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

// TransactionFromContext returns the transaction stored by WithTransaction.
func TransactionFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(transactionContextKey{}).(*sql.Tx)
	return tx, ok && tx != nil
}
