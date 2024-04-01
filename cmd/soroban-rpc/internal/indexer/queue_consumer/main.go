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
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/clients"
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
	indexerService := indexer.New(logger)
	processed := int64(0)
	logger.Info("start to consume")

	ctx := context.Background() // Assuming context is defined

	for {
		result, err := rdb.BLPop(ctx, 0*time.Second, "change_queue").Result()
		if err != nil {
			logger.WithError(err).Error("Error dequeuing")
			time.Sleep(1 * time.Second) // Example: simple backoff
			continue
		}

		if len(result) == 2 {
			item := result[1]

			decodedBytes, err := decodeFromBase64(item) // Assume util.DecodeFromBase64 exists
			if err != nil {
				logger.WithError(err).Error("Error decodeFromBase64")
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

			processed++
			if processed%10000 == 0 {
				logger.Infof("Processed %d items", processed)
			}
		}
	}
}

func processFile() {
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
