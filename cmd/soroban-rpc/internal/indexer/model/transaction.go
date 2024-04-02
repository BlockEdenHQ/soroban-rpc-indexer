package model

import (
	"encoding/json"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	ID               string              `gorm:"column:id;primaryKey"`
	Status           string              `gorm:"column:status"`
	Ledger           *uint32             `gorm:"column:ledger"`
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

func UpsertTransaction(db *gorm.DB, tx *Transaction) error {
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "ledger", "created_at", "application_order", "fee_bump",
			"fee_bump_info", "fee", "fee_charged", "sequence", "source_account", "muxed_account_id", "memo", "preconditions",
			"signatures"}), // List columns to update
	}).Create(&tx).Error

	return err
}

func NewTransaction(inp []byte) (Transaction, error) {
	var unmarshalledTokenOp Transaction
	err := json.Unmarshal(inp, &unmarshalledTokenOp)
	return unmarshalledTokenOp, err
}
