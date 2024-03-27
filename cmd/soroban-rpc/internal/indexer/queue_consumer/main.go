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
	"os"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file") // Critical startup error
	}
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		panic("POSTGRES_DSN is empty") // Critical startup error
	}

	logger := supportlog.New()
	logger.SetLevel(logrus.DebugLevel)

	logger.Info("init services")

	queue := clients.NewFileQueue("ledger_entries.txt", "ledger_entries.position.txt", logger)

	indexerService := indexer.New(logger)

	logger.Info("start to consume")

	for {
		item := queue.Dequeue()
		if item == "" {
			logger.Info("Reached the end of the queue or encountered an error")
			break
		}

		decodedBytes, err := decodeFromBase64(item)
		if err != nil {
			logger.WithError(err).Error("Error decodeFromBase64")
		}

		entry := xdr.LedgerEntry{}
		err = entry.UnmarshalBinary(decodedBytes)
		if err != nil {
			logger.WithError(err).Error("Error UnmarshalBinary")
		}

		err = indexerService.UpsertLedgerEntry(entry)
		if err != nil {
			logger.WithError(err).Error("Error UpsertLedgerEntry")
		}
	}

	queue.Close()
}

func decodeFromBase64(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
