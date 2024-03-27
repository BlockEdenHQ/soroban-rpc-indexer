package model

import (
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenBalance struct {
	ContractID string `gorm:"primaryKey;not null"`
	Address    string `gorm:"primaryKey;not null"`
	Balance    string `gorm:""`
	util.Ts
}

func UpsertTokenBalance(db *gorm.DB, tokenBalance *TokenBalance) error {
	// Upsert operation considering composite primary keys (ContractID and Address).
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "contract_id"}, {Name: "address"}}, // Composite primary keys for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{"balance"}),             // Fields to update on conflict, excluding the primary keys
	}).Create(tokenBalance).Error

	return err
}
