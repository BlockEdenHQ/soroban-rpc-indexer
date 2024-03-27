package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	supportlog "github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer"
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

	file, err := os.Open("ledger_entries.txt")
	if err != nil {
		logger.Errorf("failed to open file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	scanner := bufio.NewScanner(file)

	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	indexerService := indexer.New(logger)

	logger.Info("start to consume")

	processed := int64(0)

	for scanner.Scan() {
		item := scanner.Text()

		decodedBytes, err := decodeFromBase64(item)
		if err != nil {
			logger.WithError(err).Error("Error decodeFromBase64: " + item)
			continue
		}

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

		if processed%10000 == 0 {
			logger.Infof("processed %d", processed)
		}
		processed += 1
	}

	logger.Infof("done! congrats")

	if err := scanner.Err(); err != nil {
		logger.Errorf("error during scan: %s", err)
	}
}

func decodeFromBase64(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
