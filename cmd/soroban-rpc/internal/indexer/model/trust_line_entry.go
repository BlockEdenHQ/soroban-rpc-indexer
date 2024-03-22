package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
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
