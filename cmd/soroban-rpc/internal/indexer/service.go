package indexer

import (
	"encoding/json"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/events"
	"os"
	"strings"

	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	golg "gorm.io/gorm/logger"

	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/db"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/methods"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/transactions"
)

type Service struct {
	logger    *log.Entry
	indexerDB *gorm.DB
}

func (s *Service) scValXdrToJSON(str string) (string, error) {
	var scVal = xdr.ScVal{}
	xdr.SafeUnmarshalBase64(str, &scVal)
	return s.scValToJSON(scVal)
}

func (s *Service) scValToJSON(scVal xdr.ScVal) (string, error) {
	data := scValToGo(scVal)
	jsonData, err := json.Marshal(data)
	return strings.Replace(string(jsonData), "\\u0000", "", -1), err
}

func New(logger *log.Entry) *Service {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		panic("POSTGRES_DSN is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: golg.Default.LogMode(golg.Error),
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	s := &Service{
		indexerDB: db,
		logger:    logger,
	}

	return s
}

func (s *Service) CreateEvent(tx db.LedgerEntryReadTx, ev events.EventInfoRaw) {
	info, err := methods.NewEventInfoForEvent(ev.Event, ev.Cursor, ev.LedgerClosedAt, "") // don't need hash info
	if err != nil {
		return
	}

	topic := make([]string, 0, 4)
	for _, segment := range info.Topic {
		t, err := s.scValXdrToJSON(segment)
		if err != nil {
			s.logger.WithError(err).Error("error failed to parse segment " + segment)
			topic = append(topic, "\""+segment+"\"")
		} else {
			topic = append(topic, string(t))
		}
	}
	topicData := "[" + strings.Join(topic, ",") + "]"
	value, _ := s.scValXdrToJSON(info.Value)

	event := model.Event{
		ID:                       info.ID,
		TxIndex:                  int32(ev.Cursor.Tx),
		EventType:                info.EventType,
		Ledger:                   info.Ledger,
		LedgerClosedAt:           info.LedgerClosedAt,
		ContractID:               info.ContractID,
		PagingToken:              info.PagingToken,
		Topic:                    topicData,
		Value:                    value,
		InSuccessfulContractCall: info.InSuccessfulContractCall,
	}
	s.CreateTokenOperation(tx, topic, string(value), event)
	s.indexerDB.Create(&event)
}

func (s *Service) CreateTransaction(hash string, info methods.GetTransactionResponse, tx transactions.Transaction) {
	transaction := model.Transaction{
		ID:               hash,
		Status:           info.Status,
		Ledger:           &info.Ledger,
		CreatedAt:        &info.LedgerCloseTime,
		ApplicationOrder: &info.ApplicationOrder,
		FeeBump:          &info.FeeBump,
	}

	envelope := xdr.TransactionEnvelope{}
	envelope.UnmarshalBinary(tx.Envelope)

	if txv0, ok := envelope.GetV0(); ok {
		processTxV0Envelope(&transaction, txv0)
	} else if txv1, ok := envelope.GetV1(); ok {
		processTxV1Envelope(&transaction, txv1)
	} else if txFeeBump, ok := envelope.GetFeeBump(); ok {
		processTxFeeBumpEnvelope(&transaction, txFeeBump)
	}

	result := xdr.TransactionResult{}
	result.UnmarshalBinary(tx.Result)

	var feeCharged = int32(result.FeeCharged)
	transaction.FeeCharged = &feeCharged
	s.indexerDB.Create(&transaction)
}
