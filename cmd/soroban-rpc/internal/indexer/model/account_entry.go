package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountEntry struct {
	AccountId             string             `gorm:"type:varchar(64);primaryKey;not null"`
	Balance               xdr.Int64          `gorm:"type:bigint;not null"`
	SeqNum                xdr.SequenceNumber `gorm:"type:bigint;not null"`
	NumSubEntries         xdr.Uint32         `gorm:"type:int;not null"`
	Flags                 xdr.Uint32         `gorm:"type:int;not null"`
	HomeDomain            xdr.String32       `gorm:"type:varchar(32);index"`
	Signers               *[]byte            `gorm:"type:jsonb"`
	Ext                   interface{}        `gorm:"type:jsonb;not null"`
	InflationDest         string             `gorm:"type:varchar(64);not null"`
	Thresholds            []byte             `gorm:"type:bytea"`
	SponsoringId          string             `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32         `gorm:"type:int;not null"`
	util.Ts
}

func UpsertAccountEntry(db *gorm.DB, entry *AccountEntry) error {
	// Use the Clauses method to specify the ON CONFLICT behavior.
	// Since AccountId is the primary key, the conflict will be based on it.
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "account_id"}}, // Primary key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{
			"balance", "seq_num", "num_sub_entries", "flags",
			"home_domain", "signers", "ext", "inflation_dest",
			"thresholds", "sponsoring_id", "last_modified_ledger_seq",
		}), // Columns to update in case of conflict, excluding the primary key
	}).Create(entry).Error

	return err
}
