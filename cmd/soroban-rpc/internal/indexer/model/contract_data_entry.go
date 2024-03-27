package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContractDataEntry struct {
	KeyHash             string     `gorm:"type:text;primaryKey;not null"`
	ContractId          string     `gorm:"type:varchar(64)"`
	KeyXdr              string     `gorm:"type:text"`
	ExpirationLedgerSeq xdr.Uint32 `gorm:"type:int"`

	Key        interface{} `gorm:"type:jsonb;"` // native json of KeyXdr for better readability
	Durability string      `gorm:"varchar(64)"` // enum persistent or temporary
	Flags      xdr.Uint32  `gorm:"type:int"`    // contract data body flags
	ValXdr     string      `gorm:"type:text"`
	Val        interface{} `gorm:"type:jsonb"` // native json of contract data body val

	util.Ts
}

func UpsertContractDataEntry(db *gorm.DB, entry *ContractDataEntry) error {
	// Use the Clauses method with ON CONFLICT directive for PostgreSQL,
	// targeting the primary key (KeyHash) for conflict detection.
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key_hash"}}, // Conflict target
		DoUpdates: clause.AssignmentColumns([]string{
			"contract_id", "key_xdr", "expiration_ledger_seq",
			"key", "durability", "flags", "val_xdr", "val",
		}), // Fields to update in case of conflict, excluding primary key
	}).Create(entry).Error

	return err
}
