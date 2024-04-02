package ingest

import (
	"context"
	"encoding/base64"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/methods"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/transactions"
)

func (s *Service) enqueueChangePost(changePost xdr.LedgerEntry) {
	bytes, err := changePost.MarshalBinary()
	if err != nil {
		s.logger.WithError(err).Error("error cannot marshal LedgerEntry")
	}
	encodedEntry := base64.StdEncoding.EncodeToString(bytes)
	err = s.rdb.RPush(context.Background(), indexer.QueueKey, indexer.LedgerEntry, encodedEntry).Err()
	if err != nil {
		s.logger.WithError(err).Error("error push change_queue")
	}
}

func (s *Service) enqueueTransaction(hash string, info methods.GetTransactionResponse, tx transactions.Transaction) {
	marshaledTx := s.indexerService.MarshalTransaction(hash, info, tx)
	err := s.rdb.RPush(context.Background(), indexer.QueueKey, indexer.Tx, marshaledTx).Err()
	if err != nil {
		s.logger.WithError(err).Error("error push tx_queue")
	}
}
