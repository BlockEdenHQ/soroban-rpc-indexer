package model

import (
	"encoding/json"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenOperation struct {
	ID               string       `gorm:"column:id;primaryKey"`
	Type             string       `gorm:"column:type"`
	TxIndex          int32        `gorm:"column:tx_index"`
	Ledger           int32        `gorm:"column:ledger"`
	LedgerClosedAt   string       `gorm:"column:ledger_closed_at"`
	ContractID       string       `gorm:"column:contract_id"`
	From             string       `gorm:"column:from"`
	To               *string      `gorm:"column:to"`
	Amount           *util.Int128 `gorm:"column:amount"`
	Authorized       *bool        `gorm:"column:authorized"`
	ExpirationLedger *int32       `gorm:"column:expiration_ledger"`
	util.Ts
}

func UpsertTokenOperation(db *gorm.DB, tokenOp *TokenOperation) error {
	// Assuming `ID` is the field that should uniquely identify the record,
	// and you want to update all fields on conflict.
	err := db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // Primary Key
			DoUpdates: clause.AssignmentColumns([]string{"type", "tx_index", "ledger", "ledger_closed_at", "contract_id", "from", "to", "amount", "authorized", "expiration_ledger", "created_at", "updated_at"}),
		}).Create(tokenOp).Error
	return err
}

func NewTokenOperation(inp []byte) (TokenOperation, error) {
	var to TokenOperation
	err := json.Unmarshal(inp, &to)
	return to, err
}
