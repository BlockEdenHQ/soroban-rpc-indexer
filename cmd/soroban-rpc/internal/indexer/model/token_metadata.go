package model

import "github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"

type TokenMetadata struct {
	ContractID   string `gorm:"primaryKey"`
	AdminAddress string `gorm:"column:admin_address"`
	Decimal      uint32 `gorm:"column:decimal"`
	Name         string `gorm:"column:name"`
	Symbol       string `gorm:"column:symbol"`
	util.Ts
}
