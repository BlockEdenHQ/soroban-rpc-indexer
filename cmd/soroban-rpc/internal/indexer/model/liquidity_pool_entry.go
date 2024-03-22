package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
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
