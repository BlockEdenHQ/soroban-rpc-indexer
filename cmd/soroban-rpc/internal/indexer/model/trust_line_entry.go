package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrustLineEntry struct {
	AccountId             string        `gorm:"type:varchar(64);primaryKey;not null"`
	Balance               xdr.Int64     `gorm:"type:bigint;not null"`
	Limit                 xdr.Int64     `gorm:"type:bigint;not null"`
	AssetType             xdr.AssetType `gorm:"type:int;primaryKey"`
	AssetCode             []byte        `gorm:"type:bytea;primaryKey"`
	AssetIssuer           string        `gorm:"type:varchar(64);primaryKey"`
	LiquidityPoolId       []byte        `gorm:"type:bytea;index"`
	Flags                 xdr.Uint32    `gorm:"type:int;not null"`
	Ext                   interface{}   `gorm:"type:jsonb;not null"`
	SponsoringId          string        `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32    `gorm:"type:int;not null"`
	util.Ts
}

func UpsertTrustLineEntry(db *gorm.DB, entry *TrustLineEntry) error {
	// Upsert operation considering composite primary keys.
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "account_id"},
			{Name: "asset_type"},
			{Name: "asset_code"},
			{Name: "asset_issuer"},
		}, // Composite primary keys for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{
			"balance", "limit", "liquidity_pool_id", "flags", "ext",
			"sponsoring_id", "last_modified_ledger_seq",
		}), // Specify fields to update on conflict, excluding primary keys
	}).Create(entry).Error

	return err
}
