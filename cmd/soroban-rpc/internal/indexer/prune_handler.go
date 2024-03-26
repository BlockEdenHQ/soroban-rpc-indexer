package indexer

import (
	"context"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/stellar/go/support/log"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"gorm.io/gorm"
	"time"
)

type PruneDatabaseResult struct {
	Status string `json:"status"`
}

func deleteOldRecords[T any](db *gorm.DB, model *T, logger *log.Entry) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	result := db.Where("updated_at < ?", thirtyDaysAgo).Delete(model)
	if result.Error != nil {
		logger.Errorf("Error deleting records for %T: %v", model, result.Error)
		return
	}

	logger.Debugf("Deleted %v records for %T.", result.RowsAffected, model)
}

func NewPruneDatabaseHandler(service *Service, logger *log.Entry) jrpc2.Handler {
	return handler.New(func(ctx context.Context) (PruneDatabaseResult, error) {
		deleteOldRecords(service.indexerDB, &model.AccountEntry{}, logger)
		deleteOldRecords(service.indexerDB, &model.ClaimableBalanceEntry{}, logger)
		deleteOldRecords(service.indexerDB, &model.ContractDataEntry{}, logger)
		deleteOldRecords(service.indexerDB, &model.DataEntry{}, logger)
		deleteOldRecords(service.indexerDB, &model.Event{}, logger)
		deleteOldRecords(service.indexerDB, &model.LiquidityPoolEntry{}, logger)
		deleteOldRecords(service.indexerDB, &model.TrustLineEntry{}, logger)
		return PruneDatabaseResult{Status: "OK"}, nil
	})
}
