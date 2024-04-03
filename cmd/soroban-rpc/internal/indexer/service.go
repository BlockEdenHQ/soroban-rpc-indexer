package indexer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/events"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	golg "gorm.io/gorm/logger"

	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/methods"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/transactions"
)

type Service struct {
	logger    *log.Entry
	indexerDB *gorm.DB
	rdb       *redis.Client
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

func New(logger *log.Entry, rdb *redis.Client) *Service {
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
		rdb:       rdb,
	}

	return s
}

func (s *Service) EnqueueEvent(ev events.EventInfoRaw) {
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
	s.enqueueTokenOperation(topic, value, event)
	s.enqueueEvent(event)
}

func (s *Service) UpsertEvent(event *model.Event) error {
	return model.UpsertEvent(s.indexerDB, event)
}

func (s *Service) MarshalTransaction(hash string, info methods.GetTransactionResponse, tx transactions.Transaction) string {
	transaction := model.Transaction{
		ID:     hash,
		Status: info.Status,
		Ledger: &info.Ledger,
		Ts: util.Ts{
			CreatedAt: time.Unix(info.LedgerCloseTime, 0).UTC(),
			UpdatedAt: time.Now(),
		},
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

	jsonDataPretty, err := json.Marshal(transaction)
	if err != nil {
		s.logger.WithError(err).Error("error cannot marshal tx")
	}
	return base64.StdEncoding.EncodeToString(jsonDataPretty)
}

func (s *Service) UpsertTransaction(transaction *model.Transaction) error {
	return model.UpsertTransaction(s.indexerDB, transaction)
}

func (s *Service) enqueueTokenMetadata(tm model.TokenMetadata) {
	jsonData, err := json.Marshal(tm)
	if err != nil {
		s.logger.WithError(err).Error("error cannot marshal TokenMetadata")
	}
	marshaled := base64.StdEncoding.EncodeToString(jsonData)
	err = s.rdb.RPush(context.Background(), QueueKey, TokenMetadata+":"+marshaled).Err()
	if err != nil {
		s.logger.WithError(err).Error("error push event_token_metadata")
	}
}

func (s *Service) enqueueEvent(event model.Event) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		s.logger.WithError(err).Error("error cannot marshal event")
	}
	marshaledEvent := base64.StdEncoding.EncodeToString(jsonData)
	err = s.rdb.RPush(context.Background(), QueueKey, Event+":"+marshaledEvent).Err()
	if err != nil {
		s.logger.WithError(err).Error("error push event_queue")
	}
}

func (s *Service) enqueueTokenOp(op model.TokenOperation) {
	jsonData, err := json.Marshal(op)
	if err != nil {
		s.logger.WithError(err).Error("error cannot marshal token op")
	}
	marshaledTokenOp := base64.StdEncoding.EncodeToString(jsonData)
	err = s.rdb.RPush(context.Background(), QueueKey, TokenOperation+":"+marshaledTokenOp).Err()
	if err != nil {
		s.logger.WithError(err).Error("error push token_op")
	}
}

func (s *Service) UpsertTokenOperation(to *model.TokenOperation) error {
	return model.UpsertTokenOperation(s.indexerDB, to)
}

// key: "change_queue" value: "${number}:${base64encoded}"
const (
	QueueKey       = "change_queue"
	LedgerEntry    = "1"
	Tx             = "2"
	TokenMetadata  = "3"
	Event          = "4"
	TokenOperation = "5"
)
