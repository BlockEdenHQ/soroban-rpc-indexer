package model

import "github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"

type TokenBalance struct {
	ContractID string `gorm:"primaryKey:not null"`
	Address    string `gorm:"primaryKey:not null"`
	Balance    string `gorm:""`
	util.Ts
}
