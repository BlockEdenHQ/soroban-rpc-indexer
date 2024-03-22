package indexer

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/db"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
)

type ScValMapSimplePair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ScValContractInstance struct {
	Storage    []ScValMapSimplePair `json:"storage"`
	Executable interface{}          `json:"executable"`
}

func (s *Service) UpsertTokenBalance(contractId string, key string, val string) {
	var data []string
	if err := json.Unmarshal([]byte(key), &data); err != nil {
		errors.Wrap(err, "failed to unmarshal the key")
		return
	}
	if len(data) == 2 && data[0] == "Balance" {
		// there are two val format, need to integrate both:
		// 1. 6600000000
		// 2. [{"key":"amount","value":6600000000},{"key":"authorized","value":true},{"key":"clawback","value":false}]
		var balance string
		if err1 := json.Unmarshal([]byte(val), &balance); err1 != nil {
			var pairs []ScValMapSimplePair
			if err2 := json.Unmarshal([]byte(val), &pairs); err2 != nil {
				errors.Wrap(err2, "failed to unmarshal the val")
				return
			}
			for _, item := range pairs {
				if item.Key == "amount" {
					value, ok := item.Value.(string)
					if ok {
						balance = value
					}
					break
				}
			}
		}

		tokenBalance := model.TokenBalance{
			ContractID: contractId,
			Address:    data[1],
			Balance:    balance,
		}
		if err := s.indexerDB.Create(&tokenBalance).Error; err != nil {
			errors.Wrap(err, "failed to update the token balance")
		}
	}
}

func (s *Service) UpsertTokenMetadata(contractId string, key string, val string) {
	if key == "\"ScvLedgerKeyContractInstance\"" {
		var instance ScValContractInstance
		if err := json.Unmarshal([]byte(val), &instance); err != nil {
			errors.Wrap(err, "failed to unmarshal the contract instance value")
			return
		}
		var tokenMeta model.TokenMetadata
		for _, item := range instance.Storage {
			if item.Key == "METADATA" {
				meta, _ := json.Marshal(item.Value)
				var data []ScValMapSimplePair
				json.Unmarshal([]byte(meta), &data)
				simpleMap := make(map[string]interface{})
				for _, pair := range data {
					simpleMap[pair.Key] = pair.Value
				}
				simpleStr, _ := json.Marshal(simpleMap)
				json.Unmarshal([]byte(simpleStr), &tokenMeta)
				tokenMeta.ContractID = contractId
			}
			if item.Key == "Admin" {
				val, ok := item.Value.(string)
				if ok {
					tokenMeta.AdminAddress = val
				}
			}
		}
		if tokenMeta.Name != "" {
			if tokenMeta.CreatedAt == (time.Time{}) {
				tokenMeta.CreatedAt = time.Now()
			}
			s.indexerDB.Save(&tokenMeta)
		}
	}
}

func initTokenOpFromEvent(opType string, from string, event model.Event) model.TokenOperation {
	return model.TokenOperation{
		ID:             event.ID,
		Type:           opType,
		TxIndex:        event.TxIndex,
		Ledger:         event.Ledger,
		LedgerClosedAt: event.LedgerClosedAt,
		ContractID:     event.ContractID,
		From:           from,
	}
}

func (s *Service) UpdateTokenMetadata(tx db.LedgerEntryReadTx, contractId string, admin string) {
	var meta model.TokenMetadata
	if err := s.indexerDB.First(meta, contractId).Error; err != nil {
		errors.Wrap(err, "failed to find the record")
	}
	if meta.AdminAddress != admin {
		meta.AdminAddress = admin
		if meta.CreatedAt == (time.Time{}) {
			meta.CreatedAt = time.Now()
		}
		if err := s.indexerDB.Save(&meta).Error; err != nil {
			errors.Wrap(err, "failed to update the record")
		}
	}
}

func (s *Service) CreateTokenOperation(tx db.LedgerEntryReadTx, topicRaw []string, value string, event model.Event) {
	topic := make([]string, 0, 4)
	for _, t := range topicRaw {
		t = strings.TrimPrefix(t, "\"")
		t = strings.TrimSuffix(t, "\"")
		topic = append(topic, t)
	}
	switch topic[0] {
	case "set_admin":
		if len(topic) < 2 {
			return
		}
		tokenOp := initTokenOpFromEvent(topic[0], topic[1], event)
		var admin = value[1 : len(value)-2]
		tokenOp.To = &admin
		s.indexerDB.Create(&tokenOp)

		s.UpdateTokenMetadata(tx, event.ContractID, admin)
	case "set_authorized":
		if len(topic) < 3 {
			return
		}
		authorized, _ := strconv.ParseBool(value)
		tokenOp := initTokenOpFromEvent(topic[0], topic[1], event)
		tokenOp.From = topic[1]
		tokenOp.To = &topic[2]
		tokenOp.Authorized = &authorized
		s.indexerDB.Create(&tokenOp)
	case "approve":
		if len(topic) < 3 {
			return
		}
		var data []json.RawMessage
		json.Unmarshal([]byte(value), &data)
		var expiration int32
		num, _ := strconv.Atoi(string(data[1]))
		expiration = int32(num)
		amount := getInt128FromString(string(data[0]))
		tokenOp := initTokenOpFromEvent(topic[0], topic[1], event)
		tokenOp.To = &topic[2]
		tokenOp.Amount = &amount
		tokenOp.ExpirationLedger = &expiration
		s.indexerDB.Create(&tokenOp)
	case "mint", "transfer":
		if len(topic) < 3 {
			return
		}
		amount := getInt128FromString(value)
		tokenOp := initTokenOpFromEvent(topic[0], topic[1], event)
		tokenOp.To = &topic[2]
		tokenOp.Amount = &amount
		s.indexerDB.Create(&tokenOp)
	case "clawback":
		if len(topic) < 3 {
			return
		}
		amount := getInt128FromString(value)
		tokenOp := initTokenOpFromEvent(topic[0], topic[2], event)
		tokenOp.To = &topic[1]
		tokenOp.Amount = &amount
		s.indexerDB.Create(&tokenOp)
	case "burn":
		if len(topic) < 2 {
			return
		}
		amount := getInt128FromString(value)
		tokenOp := initTokenOpFromEvent(topic[0], topic[1], event)
		tokenOp.Amount = &amount
		s.indexerDB.Create(&tokenOp)
	}
}
