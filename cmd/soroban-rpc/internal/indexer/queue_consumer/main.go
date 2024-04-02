package main

import (
	"context"
	"encoding/base64"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/clients"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file") // Use logrus directly for startup errors
	}
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		logrus.Fatal("POSTGRES_DSN is empty") // Use logrus directly for startup errors
	}

	logger := supportlog.New()
	logger.SetLevel(logrus.DebugLevel)

	logger.Info("init services")

	rdb := clients.NewRedis(logger)
	indexerService := indexer.New(logger, rdb)
	processed := int64(0)
	logger.Info("start to consume")

	ctx := context.Background() // Assuming context is defined

	for {
		result, err := rdb.BLPop(ctx, 0*time.Second, indexer.QueueKey).Result()
		if err != nil {
			logger.WithError(err).Error("Error dequeuing")
			time.Sleep(1 * time.Second) // Example: simple backoff
			continue
		}

		if len(result) != 2 {
			logger.WithError(err).Error("result length incorrect")
			continue
		}

		rawValue := result[1]
		itemKey := rawValue[0:1]
		item := rawValue[2:]

		decodedBytes, err := decodeFromBase64(item) // Assume util.DecodeFromBase64 exists
		if err != nil {
			logger.WithError(err).Error("Error decodeFromBase64: " + item)
			continue
		}

		switch itemKey {
		case indexer.LedgerEntry:
			entry := xdr.LedgerEntry{}
			err = entry.UnmarshalBinary(decodedBytes)
			if err != nil {
				logger.WithError(err).Error("Error UnmarshalBinary")
				continue
			}

			err = indexerService.UpsertLedgerEntry(entry)
			if err != nil {
				logger.WithError(err).Error("Error UpsertLedgerEntry")
			}
			break
		case indexer.Tx:
			tx, err := model.NewTransaction(decodedBytes)
			if err != nil {
				logger.WithError(err).Error("Error NewTransaction")
			}
			err = indexerService.UpsertTransaction(&tx)
			if err != nil {
				logger.WithError(err).Error("Error Tx")
			}
			break
		case indexer.TokenMetadata:
			tm, err := model.NewTokenMetadata(decodedBytes)
			if err != nil {
				logger.WithError(err).Error("Error NewTokenMetadata")
			}
			err = indexerService.UpsertTokenMetadataFromStruct(&tm)
			if err != nil {
				logger.WithError(err).Error("Error NewTransaction")
			}
			break
		case indexer.Event:
			ev, err := model.NewEvent(decodedBytes)
			if err != nil {
				logger.WithError(err).Error("Error NewEvent")
			}
			err = indexerService.UpsertEvent(&ev)
			if err != nil {
				logger.WithError(err).Error("Error UpsertEvent")
			}
			break
		case indexer.TokenOperation:
			to, err := model.NewTokenOperation(decodedBytes)
			if err != nil {
				logger.WithError(err).Error("Error NewTokenOperation")
			}
			err = indexerService.UpsertTokenOperation(&to)
			if err != nil {
				logger.WithError(err).Error("Error UpsertTokenOperation")
			}
		}

		processed++
		if processed%10000 == 0 {
			logger.Infof("Processed %d items", processed)
		}
	}
}

func decodeFromBase64(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
