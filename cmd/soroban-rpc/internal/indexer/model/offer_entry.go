package model

import (
	"github.com/shopspring/decimal"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
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
