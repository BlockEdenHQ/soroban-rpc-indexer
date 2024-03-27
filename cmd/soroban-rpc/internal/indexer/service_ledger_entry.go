package indexer

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/parser"
)

func (s *Service) UpsertLedgerEntry(entry xdr.LedgerEntry) error {
	// We can do a little extra validation to ensure the entry and key match,
	// because the key can be derived from the entry.
	key, err := entry.LedgerKey()
	if err != nil {
		return errors.Wrap(err, "could not get ledger key from entry")
	}

	bin, err := key.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "could not marshal ledger key binary")
	}
	keyHash := sha256.Sum256(bin)
	hexKey := hex.EncodeToString(keyHash[:])
	// upsert according to
	// xdr/xdr_generated.go
	// LedgerKey is an XDR Union defines as:
	//
	//	union LedgerKey switch (LedgerEntryType type)
	//	 {
	//	 case ACCOUNT:
	//	     struct
	//	     {
	//	         AccountID accountID;
	//	     } account;
	//
	//	 case TRUSTLINE:
	//	     struct
	//	     {
	//	         AccountID accountID;
	//	         TrustLineAsset asset;
	//	     } trustLine;
	//
	//	 case OFFER:
	//	     struct
	//	     {
	//	         AccountID sellerID;
	//	         int64 offerID;
	//	     } offer;
	//
	//	 case DATA:
	//	     struct
	//	     {
	//	         AccountID accountID;
	//	         string64 dataName;
	//	     } data;
	//
	//	 case CLAIMABLE_BALANCE:
	//	     struct
	//	     {
	//	         ClaimableBalanceID balanceID;
	//	     } claimableBalance;
	//
	//	 case LIQUIDITY_POOL:
	//	     struct
	//	     {
	//	         PoolID liquidityPoolID;
	//	     } liquidityPool;
	//	 case CONTRACT_DATA:
	//	     struct
	//	     {
	//	         SCAddress contract;
	//	         SCVal key;
	//	         ContractDataDurability durability;
	//	         ContractEntryBodyType bodyType;
	//	     } contractData;
	//	 case CONTRACT_CODE:
	//	     struct
	//	     {
	//	         Hash hash;
	//	         ContractEntryBodyType bodyType;
	//	     } contractCode;
	//	 case CONFIG_SETTING:
	//	     struct
	//	     {
	//	         ConfigSettingID configSettingID;
	//	     } configSetting;
	//	 };

	if key.Ttl != nil {
		searchKey := hex.EncodeToString(entry.Data.Ttl.KeyHash[:])

		// find and update existing contract data
		// to do: decide whether we should store data entry in another table
		if err := s.indexerDB.Model(&model.ContractDataEntry{}).Where(&model.ContractDataEntry{KeyHash: searchKey}).UpdateColumn("expiration_ledger_seq", entry.Data.Ttl.LiveUntilLedgerSeq).Error; err != nil {
			return errors.Wrap(err, "failed to update ContractData Expiry: "+searchKey)
		}
	}

	if key.ContractData != nil {
		oldEm := model.ContractDataEntry{}
		errFindOldEm := s.indexerDB.Where(&model.ContractDataEntry{KeyHash: hexKey}).First(&oldEm).Error
		em := parser.GetContractDataModel(entry)
		if em != nil {
			if errFindOldEm != nil {
				em.ExpirationLedgerSeq = oldEm.ExpirationLedgerSeq
			}
			em.KeyHash = hexKey
			key, _ := s.scValToJSON(entry.Data.ContractData.Key)
			em.Key = key
			val, _ := s.scValToJSON(entry.Data.ContractData.Val)
			em.Val = val

			s.UpsertTokenBalance(em.ContractId, key, val)
			s.UpsertTokenMetadata(em.ContractId, key, val)

			if em.CreatedAt == (time.Time{}) {
				em.CreatedAt = time.Now()
			}

			if err := model.UpsertContractDataEntry(s.indexerDB, em); err != nil {
				return errors.Wrap(err, "failed to upsert or update ContractDataEntry")
			}
		}
	}

	if key.Account != nil {
		em := parser.GetAccountEntryModel(entry)
		if err := model.UpsertAccountEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert AccountEntry")
		}
	}

	if key.TrustLine != nil {
		em := parser.GetTrustLineEntryModel(entry)
		if err := model.UpsertTrustLineEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert TrustLineEntry")
		}
	}

	if key.Offer != nil {
		em := parser.GetOfferEntryModel(entry)
		if err := model.UpsertOfferEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert OfferEntry")
		}
	}

	if key.Data != nil {
		em := parser.GetDataEntryModel(entry)
		if err := model.UpsertDataEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert DataEntry")
		}
	}

	if key.ClaimableBalance != nil {
		em := parser.GetClaimableBalanceEntryModel(entry)
		if err := model.UpsertClaimableBalanceEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert ClaimableBalanceEntry")
		}
	}

	if key.LiquidityPool != nil {
		em := parser.GetLiquidityPoolEntryModel(entry)
		if err := model.UpsertLiquidityPoolEntry(s.indexerDB, em); err != nil {
			return errors.Wrap(err, "failed to upsert LiquidityPoolEntry")
		}
	}

	return nil
}
