package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
)

type DataEntry struct {
	AccountId             string       `gorm:"type:varchar(64);primaryKey:not null"`
	DataName              xdr.String64 `gorm:"type:varchar(64);primaryKey;not null"`
	DataValue             interface{}  `gorm:"type:jsonb;not null"`
	Ext                   interface{}  `gorm:"type:jsonb;not null"`
	SponsoringId          string       `gorm:"type:varchar(64);index"`
	LastModifiedLedgerSeq xdr.Uint32   `gorm:"type:int;not null"`
	util.Ts
}
