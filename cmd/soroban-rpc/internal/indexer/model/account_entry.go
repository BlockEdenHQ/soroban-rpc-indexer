package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
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
