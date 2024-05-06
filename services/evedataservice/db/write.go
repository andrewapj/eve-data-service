package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

const (
	insertOrUpdateBatch = `
	INSERT INTO %s (id, data)
	VALUES ($1, $2) ON CONFLICT (id) DO 
	UPDATE SET data = $2;`
)

func StartTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	return tx, nil
}

func InsertOrUpdateBatch(ctx context.Context, tx pgx.Tx, kind string, entities []Entity) error {

	batch := &pgx.Batch{}
	for _, entity := range entities {
		id := entity.Id()
		data, err := entity.Data()
		if err != nil {
			return fmt.Errorf("unable to insert batch: %w", err)
		}

		batch.Queue(fmt.Sprintf(insertOrUpdateBatch, kind), id, data)
	}
	res := tx.SendBatch(ctx, batch)

	batchErr := res.Close()
	if batchErr != nil {
		return fmt.Errorf("unable to close batch: %w", batchErr)
	}
	return nil
}

func CommitTx(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}
	return nil
}

func RollbackTx(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("unable to rollback transaction: %w", err)
	}
	return nil
}
