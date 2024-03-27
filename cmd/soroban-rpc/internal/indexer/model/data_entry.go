package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataEntry struct {
	AccountId             string       `gorm:"type:varchar(64);primaryKey;not null"`
	DataName              xdr.String64 `gorm:"type:varchar(64);primaryKey;not null"`
	DataValue             interface{}  `gorm:"type:jsonb;not null"`
	Ext                   interface{}  `gorm:"type:jsonb;not null"`
	SponsoringId          string       `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32   `gorm:"type:int;not null"`
	util.Ts
}

func UpsertDataEntry(db *gorm.DB, entry *DataEntry) error {
	// Upsert operation considering composite primary keys (AccountId, DataName)
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "data_name"}},                                           // Composite primary keys
		DoUpdates: clause.AssignmentColumns([]string{"data_value", "ext", "sponsoring_id", "last_modified_ledger_seq"}), // Fields to update on conflict
	}).Create(entry).Error

	return err
}
