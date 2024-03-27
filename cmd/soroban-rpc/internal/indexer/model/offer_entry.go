package model

import (
	"github.com/shopspring/decimal"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfferEntry struct {
	OfferId  xdr.Int64 `gorm:"type:bigint;primaryKey;not null"`
	SellerId string    `gorm:"type:varchar(64);primaryKey;not null"`

	SellingAssetType   xdr.AssetType `gorm:"type:int"` // enum
	SellingAssetCode   []byte        `gorm:"type:bytea"`
	SellingAssetIssuer string        `gorm:"type:varchar(64);index;not null"`

	BuyingAssetType   xdr.AssetType `gorm:"type:int"` // enum
	BuyingAssetCode   []byte        `gorm:"type:bytea"`
	BuyingAssetIssuer string        `gorm:"type:varchar(64);index;not null"`

	Amount xdr.Int64       `gorm:"type:bigint;not null"`
	Price  decimal.Decimal `gorm:"type:numeric"`

	Flags                 xdr.Uint32  `gorm:"type:int;not null"`
	Ext                   interface{} `gorm:"type:jsonb;not null"`
	SponsoringId          string      `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32  `gorm:"type:int;not null"`
	util.Ts
}

func UpsertOfferEntry(db *gorm.DB, offer *OfferEntry) error {
	// Upsert operation considering composite primary keys (OfferId and SellerId).
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "offer_id"}, {Name: "seller_id"}}, // Composite primary keys for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{
			"selling_asset_type", "selling_asset_code", "selling_asset_issuer",
			"buying_asset_type", "buying_asset_code", "buying_asset_issuer",
			"amount", "price", "flags", "ext", "sponsoring_id", "last_modified_ledger_seq",
		}), // Specify fields to update on conflict, except the primary keys
	}).Create(offer).Error

	return err
}
