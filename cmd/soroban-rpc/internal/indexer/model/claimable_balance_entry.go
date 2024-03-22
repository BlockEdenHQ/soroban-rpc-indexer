package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
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
