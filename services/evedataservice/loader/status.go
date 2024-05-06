package loader

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/db"
	"github.com/andrewapj/arcturus/domain"
	"github.com/andrewapj/arcturus/esi"
	"time"
)

type statusLoader struct {
	lockName   string
	loaderName string
}

func newStatusLoader() statusLoader {
	return statusLoader{lockName: "status", loaderName: "status"}
}

func (l statusLoader) load(ctx context.Context, client *esi.Client) (time.Time, error) {

	esiStatus, err := client.FetchStatus(ctx)
	if err != nil {
		return time.Time{}, err
	}
	status, err := domain.MapStatusFromEsi(&esiStatus)
	if err != nil {
		return time.Time{}, err
	}

	tx, err := db.StartTx(ctx)
	if err != nil {
		return time.Time{}, err
	}

	err = db.InsertOrUpdateBatch(ctx, tx, "status", []db.Entity{status})
	if err != nil {
		rollbackErr := db.RollbackTx(ctx, tx)
		if rollbackErr != nil {
			return time.Time{}, fmt.Errorf("rollback transaction failed: %v", rollbackErr)
		}
		return time.Time{}, err
	}

	err = db.CommitTx(ctx, tx)
	if err != nil {
		return time.Time{}, err
	}

	expires, err := clock.FindEarliestTime([]time.Time{esiStatus.Expires()})
	if err != nil {
		return time.Time{}, err
	}

	return expires, nil
}

func (l statusLoader) name() string {
	return l.loaderName
}

func (l statusLoader) lock() string {
	return l.lockName
}
