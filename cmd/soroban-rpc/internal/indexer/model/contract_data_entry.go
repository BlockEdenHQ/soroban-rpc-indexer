package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
)

type ContractDataEntry struct {
	KeyHash             string     `gorm:"type:text;primaryKey:not null"`
	ContractId          string     `gorm:"type:varchar(64)"`
	KeyXdr              string     `gorm:"type:text"`
	ExpirationLedgerSeq xdr.Uint32 `gorm:"type:int"`

	Key        interface{} `gorm:"type:jsonb;"` // native json of KeyXdr for better readability
	Durability string      `gorm:"varchar(64)"` // enum persistent or temporary
	Flags      xdr.Uint32  `gorm:"type:int"`    // contract data body flags
	ValXdr     string      `gorm:"type:text"`
	Val        interface{} `gorm:"type:jsonb"` // native json of contract data body val

	util.Ts
}
