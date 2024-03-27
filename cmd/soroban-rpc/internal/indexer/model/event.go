package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Event struct {
	ID                       string      `gorm:"column:id;primaryKey"`
	TxIndex                  int32       `gorm:"column:tx_index"`
	EventType                string      `gorm:"column:type"`
	Ledger                   int32       `gorm:"column:ledger"`
	LedgerClosedAt           string      `gorm:"column:ledger_closed_at"`
	ContractID               string      `gorm:"column:contract_id"`
	PagingToken              string      `gorm:"column:paging_token"`
	Topic                    interface{} `gorm:"column:topic;type:jsonb"`
	Value                    interface{} `gorm:"column:value;type:jsonb"`
	InSuccessfulContractCall bool        `gorm:"column:in_successful_contract_call"`
	LastModifiedLedgerSeq    xdr.Uint32  `gorm:"type:int;not null"`
	util.Ts
}

func UpsertEvent(db *gorm.DB, event *Event) error {
	// Upsert operation considering 'ID' as the primary key.
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // Primary key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{
			"tx_index", "type", "ledger", "ledger_closed_at", "contract_id",
			"paging_token", "topic", "value", "in_successful_contract_call", "last_modified_ledger_seq",
		}), // Specify fields to update on conflict, except the primary key
	}).Create(event).Error

	return err
}
