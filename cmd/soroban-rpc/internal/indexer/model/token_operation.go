package model

import "github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"

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
