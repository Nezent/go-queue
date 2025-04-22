package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const TxKey ctxKey = "tx"

func WithTransaction(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.Begin(r.Context())
			if err != nil {
				http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), TxKey, tx)

			defer func() {
				if rec := recover(); rec != nil {
					tx.Rollback(ctx)
					panic(rec)
				} else if err != nil {
					tx.Rollback(ctx)
				} else {
					tx.Commit(ctx)
				}
			}()

			// Pass the request with the transaction context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTxFromContext(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}
