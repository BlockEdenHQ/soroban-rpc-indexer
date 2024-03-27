package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClaimableBalanceEntry struct {
	BalanceId string `gorm:"type:varchar(64);primaryKey;not null"`
	// TODO: should parse more?
	Claimants             interface{}   `gorm:"type:jsonb"`
	AssetType             xdr.AssetType `gorm:"type:int"`
	AssetCode             []byte        `gorm:"type:bytea"`
	AssetIssuer           string        `gorm:"type:varchar(64);index;not null"`
	Amount                xdr.Int64     `gorm:"type:bigint;not null"`
	Ext                   interface{}   `gorm:"type:jsonb"`
	SponsoringId          string        `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32    `gorm:"type:int;not null"`
	util.Ts
}

func UpsertClaimableBalanceEntry(db *gorm.DB, entry *ClaimableBalanceEntry) error {
	// Use the Clauses method with ON CONFLICT directive for PostgreSQL,
	// targeting the primary key (BalanceId) for conflict detection.
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "balance_id"}}, // Conflict target
		DoUpdates: clause.AssignmentColumns([]string{
			"claimants", "asset_type", "asset_code", "asset_issuer", "amount",
			"ext", "sponsoring_id", "last_modified_ledger_seq",
		}), // Fields to update in case of conflict, excluding primary key
	}).Create(entry).Error

	return err
}
