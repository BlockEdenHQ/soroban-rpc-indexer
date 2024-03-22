package model

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
)

type Event struct {
	ID                       string      `gorm:"column:id;primaryKey"`
	TxIndex                  int32       `gorm:"column:tx_index"`
	EventType                string      `gorm:"column:type"`
	Ledger                   int32       `gorm:"column:ledger"`
	LedgerClosedAt           string      `gorm:"column:ledger_closed_at"`
	ContractID               string      `gorm:"column:contract_id"`
	PagingToken              string      `gorm:"column:paging_token"`
	Topic                    interface{} `gorm:"column:topic;type:jsonb"`
	Value                    interface{} `gorm:"column:value;type:jsonb"`
	InSuccessfulContractCall bool        `gorm:"column:in_successful_contract_call"`
	LastModifiedLedgerSeq    xdr.Uint32  `gorm:"type:int;not null"`
	util.Ts
}
