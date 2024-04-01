package ingest

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"

	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/db"
)

func (s *Service) ingestLedgerEntryChanges(ctx context.Context, reader ingest.ChangeReader, tx db.WriteTx, progressLogPeriod int, fillingFromCheckpoint bool) error {
	entryCount := 0
	startTime := time.Now()
	writer := tx.LedgerEntryWriter()

	changeStatsProcessor := ingest.StatsChangeProcessor{}
	for ctx.Err() == nil {
		change, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if ENQUEUE_LEDGER_ENTRIES_ENABLED && fillingFromCheckpoint && entryCount <= 40320000 {
			// write to file
			s.changeQueue <- *change.Post
		} else {
			if change.Post != nil && ((fillingFromCheckpoint && entryCount > 40320000) || !fillingFromCheckpoint) {
				s.changeQueue <- *change.Post
			}
			// write to sqlite
			err = ingestLedgerEntryChange(writer, change)
			if err != nil {
				return err
			}
		}

		if err = changeStatsProcessor.ProcessChange(ctx, change); err != nil {
			return err
		}
		entryCount++
		if progressLogPeriod > 0 && entryCount%progressLogPeriod == 0 {
			s.logger.Infof("processed %d ledger entry changes", entryCount)
		}
	}

	results := changeStatsProcessor.GetResults()
	for stat, value := range results.Map() {
		stat = strings.Replace(stat, "stats_", "change_", 1)
		s.metrics.ledgerStatsMetric.
			With(prometheus.Labels{"type": stat}).Add(float64(value.(int64)))
	}
	s.metrics.ingestionDurationMetric.
		With(prometheus.Labels{"type": "ledger_entries"}).Observe(time.Since(startTime).Seconds())
	return ctx.Err()
}

func (s *Service) ingestTempLedgerEntryEvictions(
	ctx context.Context,
	evictedTempLedgerKeys []xdr.LedgerKey,
	tx db.WriteTx,
) error {
	startTime := time.Now()
	writer := tx.LedgerEntryWriter()
	counts := map[string]int{}

	for _, key := range evictedTempLedgerKeys {
		if err := writer.DeleteLedgerEntry(key); err != nil {
			return err
		}
		counts["evicted_"+key.Type.String()]++
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	for evictionType, count := range counts {
		s.metrics.ledgerStatsMetric.
			With(prometheus.Labels{"type": evictionType}).Add(float64(count))
	}
	s.metrics.ingestionDurationMetric.
		With(prometheus.Labels{"type": "evicted_temp_ledger_entries"}).Observe(time.Since(startTime).Seconds())
	return ctx.Err()
}

func ingestLedgerEntryChange(writer db.LedgerEntryWriter, change ingest.Change) error {
	if change.Post == nil {
		ledgerKey, err := xdr.GetLedgerKeyFromData(change.Pre.Data)
		if err != nil {
			return err
		}
		return writer.DeleteLedgerEntry(ledgerKey)
	} else {
		return writer.UpsertLedgerEntry(*change.Post)
	}
}
