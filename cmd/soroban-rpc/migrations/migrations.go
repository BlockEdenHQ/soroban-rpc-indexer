package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Migration struct {
	ID       string
	Migrate  func(tx *gorm.DB) error
	Rollback func(tx *gorm.DB) error
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		panic("POSTGRES_DSN is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	ms := []Migration{
		{
			ID: "create inital tables",
			Migrate: func(tx *gorm.DB) error {
				err = tx.AutoMigrate(
					&model.Event{},
					&model.Transaction{},
					&model.AccountEntry{},
					&model.TrustLineEntry{},
					&model.OfferEntry{},
					&model.DataEntry{},
					&model.ClaimableBalanceEntry{},
					&model.LiquidityPoolEntry{},
					&model.TokenOperation{},
					&model.TokenMetadata{},
					&model.ContractDataEntry{},
					&model.TokenBalance{},
				)

				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	}

	for _, m := range ms {
		fmt.Println("Migrate DB: " + m.ID)
		err := m.Migrate(db)
		if err != nil {
			panic("failed to migrate: " + m.ID)
		}
	}
}
