package helper

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RollbackOrCommit(ctx context.Context, tx *pgxpool.Tx) {
	err := recover()
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			panic(errRollback)
		}
	} else {
		errCommit := tx.Commit(ctx)
		if errCommit != nil {
			panic(errCommit)
		}
	}
}
