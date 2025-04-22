package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx"
)

type ctxKey string

const TxKey ctxKey = "tx"

func WithTransaction(db *pgx.Conn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.Begin()
			if err != nil {
				http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), TxKey, tx)

			defer func() {
				if rec := recover(); rec != nil {
					tx.Rollback()
					panic(rec)
				} else if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()

			// Pass the request with the transaction context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTxFromContext(ctx context.Context) (*pgx.Tx, error) {
	tx, ok := ctx.Value(TxKey).(*pgx.Tx)
	if !ok {
		return nil, errors.New("no transaction in context")
	}
	return tx, nil
}
