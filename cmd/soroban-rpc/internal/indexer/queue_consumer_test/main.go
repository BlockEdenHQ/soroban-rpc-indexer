package main

import (
	"encoding/json"
	"fmt"
	"github.com/stellar/go/support/log"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"time"
)

func main() {
	// Example TokenOperation struct with an Int128 amount
	amount := new(util.Int128)
	amount.SetString("123456789012345678901234567890", 10)

	tokenOp := model.TokenOperation{
		ID:               "1",
		Type:             "transfer",
		TxIndex:          0,
		Ledger:           123456,
		LedgerClosedAt:   "2021-01-01T15:04:05Z",
		ContractID:       "contract_123",
		From:             "user_from",
		To:               ptrToString("user_to"),
		Amount:           amount,
		Authorized:       ptrToBool(true),
		ExpirationLedger: ptrToInt32(int32(111)),
		Ts: util.Ts{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// before
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(tokenOp)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %s", err)
	}
	fmt.Printf("Marshalled JSON:\n%s\n", string(jsonData))

	// Unmarshal the JSON back to a struct
	var unmarshalledTokenOp model.TokenOperation
	err = json.Unmarshal(jsonData, &unmarshalledTokenOp)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// after
	jsonData, err = json.Marshal(unmarshalledTokenOp)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %s", err)
	}
	fmt.Printf("Marshalled JSON:\n%s\n", string(jsonData))

}

func ptrToString(s string) *string {
	return &s
}

func ptrToBool(b bool) *bool {
	return &b
}

func ptrToInt32(b int32) *int32 {
	return &b
}
