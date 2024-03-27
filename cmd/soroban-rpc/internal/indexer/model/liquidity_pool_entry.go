package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LiquidityPoolEntry struct {
	LiquidityPoolId          []byte                `gorm:"type:bytea;primaryKey"`
	Type                     xdr.LiquidityPoolType `gorm:"type:int"`
	AssetAType               xdr.AssetType         `gorm:"type:int"`
	AssetACode               []byte                `gorm:"type:bytea"`
	AssetAIssuer             string                `gorm:"type:varchar(64);index;not null"`
	AssetBType               xdr.AssetType         `gorm:"type:int"`
	AssetBCode               []byte                `gorm:"type:bytea"`
	AssetBIssuer             string                `gorm:"type:varchar(64);index;not null"`
	Fee                      xdr.Int32             `gorm:"type:int"`
	ReserveA                 xdr.Int64             `gorm:"type:bigint"`
	ReserveB                 xdr.Int64             `gorm:"type:bigint"`
	TotalPoolShares          xdr.Int64             `gorm:"type:bigint"`
	PoolSharesTrustLineCount xdr.Int64             `gorm:"type:bigint"`
	LastModifiedLedgerSeq    xdr.Uint32            `gorm:"type:int;not null"`
	SponsoringId             string                `gorm:"type:varchar(64);index"`
	util.Ts
}

func UpsertLiquidityPoolEntry(db *gorm.DB, entry *LiquidityPoolEntry) error {
	// Assuming `LiquidityPoolId` is the unique identifier for the upsert operation
	// Adjust the clause below based on your actual unique constraint or conflict target
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "liquidity_pool_id"}}, // Unique or primary key column names
		DoUpdates: clause.AssignmentColumns([]string{
			"type", "asset_a_type", "asset_a_code", "asset_a_issuer",
			"asset_b_type", "asset_b_code", "asset_b_issuer", "fee",
			"reserve_a", "reserve_b", "total_pool_shares",
			"pool_shares_trust_line_count", "last_modified_ledger_seq", "sponsoring_id",
		}), // Specify columns to be updated on conflict
	}).Create(entry).Error

	if err != nil {
		return err
	}

	return nil
}
