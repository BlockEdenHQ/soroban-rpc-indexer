package model

import "github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"

type Transaction struct {
	ID               string              `gorm:"column:id;primaryKey"`
	Status           string              `gorm:"column:status"`
	Ledger           *uint32             `gorm:"column:ledger"`
	CreatedAt        *int64              `gorm:"column:created_at"`
	ApplicationOrder *int32              `gorm:"column:application_order"`
	FeeBump          *bool               `gorm:"column:fee_bump"`
	FeeBumpInfo      *util.FeeBumpInfo   `gorm:"column:fee_bump_info;type:jsonb"`
	Fee              *int32              `gorm:"column:fee"`
	FeeCharged       *int32              `gorm:"column:fee_charged"`
	Sequence         *int64              `gorm:"column:sequence"`
	SourceAccount    *string             `gorm:"column:source_account"`
	MuxedAccountId   *int64              `gorm:"column:muxed_account_id"` // only set for muxed account
	Memo             *util.TypeItem      `gorm:"column:memo;type:jsonb"`
	Preconditions    *util.Preconditions `gorm:"column:preconditions;type:jsonb"`
	Signatures       *[]util.Signature   `gorm:"column:signatures;type:jsonb"`
	util.Ts
}
