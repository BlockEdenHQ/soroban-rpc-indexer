package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/log"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	golg "gorm.io/gorm/logger"
	"os"
	"time"
)

func main() {
	tx := model.Transaction{
		ID:               "tx_123ABC",
		Status:           "Completed",
		Ledger:           new(uint32),
		CreatedAt:        new(int64),
		ApplicationOrder: new(int32),
		FeeBump:          new(bool),
		FeeBumpInfo:      &util.FeeBumpInfo{Fee: 1222},
		Fee:              new(int32),
		FeeCharged:       new(int32),
		Sequence:         new(int64),
		SourceAccount:    new(string),
		MuxedAccountId:   new(int64),
		Memo:             &util.TypeItem{Type: "Text", Value: "Sample Memo"},
		Preconditions:    nil,
		Signatures:       &[]util.Signature{{Hint: "Key1", Signature: "Signature1"}},
		Ts:               util.Ts{CreatedAt: time.Date(2024, time.April, 2, 23, 13, 42, 665153413, time.Local), UpdatedAt: time.Now()},
	}

	// before
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(tx)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %s", err)
	}
	fmt.Printf("Marshalled JSON:\n%s\n", string(jsonData))

	// Unmarshal the JSON back to a struct
	var unmarshalledTokenOp model.Transaction
	err = json.Unmarshal(jsonData, &unmarshalledTokenOp)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// after
	jsonData, err = json.Marshal(unmarshalledTokenOp)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %s", err)
	}
	fmt.Printf("Marshalled JSON:\n%s\n", string(jsonData))

	err = godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file") // Use logrus directly for startup errors
	}
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		logrus.Fatal("POSTGRES_DSN is empty") // Use logrus directly for startup errors
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: golg.Default.LogMode(golg.Error),
	})

	err = model.UpsertTransaction(db, &unmarshalledTokenOp)
	if err != nil {
		fmt.Printf("error upsert :\n%s\n", string(jsonData))
	}
}

func ptrToString(s string) *string {
	return &s
}

func ptrToBool(b bool) *bool {
	return &b
}

func ptrToInt32(b int32) *int32 {
	return &b
}
