package model

import (
	"encoding/json"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenMetadata struct {
	ContractID   string `gorm:"primaryKey"`
	AdminAddress string `gorm:"column:admin_address"`
	Decimal      uint32 `gorm:"column:decimal"`
	Name         string `gorm:"column:name"`
	Symbol       string `gorm:"column:symbol"`
	util.Ts
}

func UpsertTokenMetadata(db *gorm.DB, metadata *TokenMetadata) error {
	// Upsert operation using 'ContractID' as the primary key for conflict resolution.
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "contract_id"}},                                           // Primary key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{"admin_address", "decimal", "name", "symbol"}), // Fields to update on conflict
	}).Create(metadata).Error

	return err
}

func NewTokenMetadata(inp []byte) (TokenMetadata, error) {
	var tm TokenMetadata
	err := json.Unmarshal(inp, &tm)
	return tm, err
}
