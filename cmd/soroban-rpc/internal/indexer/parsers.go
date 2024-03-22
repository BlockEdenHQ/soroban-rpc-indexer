package indexer

import (
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model/util"
)

func getInt128FromString(str string) (int128 util.Int128) {
	int128.SetString(str, 10)
	return int128
}

func setSignatures(t *model.Transaction, sigs []xdr.DecoratedSignature) {
	var signatures []util.Signature
	for _, sig := range sigs {
		signatures = append(signatures, util.Signature{
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
			Hint:      base64.StdEncoding.EncodeToString(sig.Hint[:]),
		})
	}
	t.Signatures = &signatures
}

func parseSourceAccount(account xdr.MuxedAccount) (*string, *int64) {
	if res, ok := account.GetEd25519(); ok {
		raw := make([]byte, 32)
		copy(raw, res[:])
		var str, _ = strkey.Encode(strkey.VersionByteAccountID, raw)
		return &str, nil
	} else if res, ok := account.GetMed25519(); ok {
		var id = int64(res.Id)
		raw := make([]byte, 32)
		copy(raw, res.Ed25519[:])
		var str, _ = strkey.Encode(strkey.VersionByteMuxedAccount, raw)
		return &str, &id
	}
	return nil, nil
}

func setMemo(t *model.Transaction, memo xdr.Memo) {
	var item = util.TypeItem{}
	if hash, ok := memo.GetHash(); ok {
		item.Type = "hash"
		item.Value = base64.StdEncoding.EncodeToString(hash[:])
	} else if text, ok := memo.GetText(); ok {
		item.Type = "text"
		item.Value = text
	} else if retHash, ok := memo.GetRetHash(); ok {
		item.Type = "ret_hash"
		item.Value = base64.StdEncoding.EncodeToString(retHash[:])
	} else if id, ok := memo.GetId(); ok {
		item.Type = "id"
		item.Value = strconv.FormatUint(uint64(id), 10)
	} else {
		return
	}
	t.Memo = &item
}

func parseSignerKey(signer xdr.SignerKey) util.SignerKey {
	var signerKey = util.SignerKey{}
	if res, ok := signer.GetEd25519(); ok {
		raw := make([]byte, 32)
		copy(raw, res[:])
		var str, _ = strkey.Encode(strkey.VersionByteAccountID, raw)
		signerKey.Type = "ed25519"
		signerKey.Ed25519 = str
	} else if res, ok := signer.GetPreAuthTx(); ok {
		raw := make([]byte, 32)
		copy(raw, res[:])
		var str, _ = strkey.Encode(strkey.VersionByteHashTx, raw)
		signerKey.Type = "pre_auth_tx"
		signerKey.PreAuthTx = str
	} else if res, ok := signer.GetHashX(); ok {
		raw := make([]byte, 32)
		copy(raw, res[:])
		var str, _ = strkey.Encode(strkey.VersionByteHashX, raw)
		signerKey.Type = "hash_x"
		signerKey.HashX = str
	} else if res, ok := signer.GetEd25519SignedPayload(); ok {
		raw := make([]byte, 32)
		copy(raw, res.Ed25519[:])
		var ed25519, _ = strkey.Encode(strkey.VersionByteHashX, raw)
		var payload, _ = strkey.Encode(strkey.VersionByteSignedPayload, res.Payload)
		signerKey.Type = "ed25519_signed_payload"
		signerKey.Ed25519 = ed25519
		signerKey.Ed25519SignedPayload = payload
	}
	return signerKey
}

func setCond(t *model.Transaction, cond xdr.Preconditions) {
	var precond = util.Preconditions{}
	if timeBonds, ok := cond.GetTimeBounds(); ok {
		var bonds = util.Bonds{
			Min: int64(timeBonds.MinTime),
			Max: int64(timeBonds.MaxTime),
		}
		precond.TimeBonds = &bonds
	} else if v2, ok := cond.GetV2(); ok {
		if v2.LedgerBounds != nil {
			precond.LedgerBounds = &util.Bonds{
				Min: int64(v2.LedgerBounds.MinLedger),
				Max: int64(v2.LedgerBounds.MaxLedger),
			}
		}
		if v2.MinSeqNum != nil {
			precond.MinSeqNum = (*int64)(v2.MinSeqNum)
		}
		var minSeqAge = int64(v2.MinSeqAge)
		precond.MinSeqAge = &minSeqAge
		var minSeqLedgerGap = int32(v2.MinSeqLedgerGap)
		precond.MinSeqLedgerGap = &minSeqLedgerGap

		var extraSigners = []util.SignerKey{}
		for _, signer := range v2.ExtraSigners {
			extraSigners = append(extraSigners, parseSignerKey(signer))
		}
		precond.ExtraSigners = &extraSigners
	} else {
		return
	}
	t.Preconditions = &precond
}

func processTxV0Envelope(transaction *model.Transaction, txv0 xdr.TransactionV0Envelope) {
	var fee = int32(txv0.Tx.Fee)
	transaction.Fee = &fee
	var sequence = int64(txv0.Tx.SeqNum)
	transaction.Sequence = &sequence

	setSignatures(transaction, txv0.Signatures)
	setMemo(transaction, txv0.Tx.Memo)

	raw := make([]byte, 32)
	copy(raw, txv0.Tx.SourceAccountEd25519[:])
	var str, _ = strkey.Encode(strkey.VersionByteAccountID, raw)
	transaction.SourceAccount = &str

	if txv0.Tx.TimeBounds != nil {
		var precond = util.Preconditions{}
		var bonds = util.Bonds{
			Min: int64(txv0.Tx.TimeBounds.MinTime),
			Max: int64(txv0.Tx.TimeBounds.MaxTime),
		}
		precond.TimeBonds = &bonds
		transaction.Preconditions = &precond
	}
}

func processTxV1Envelope(transaction *model.Transaction, txv1 xdr.TransactionV1Envelope) {
	var fee = int32(txv1.Tx.Fee)
	transaction.Fee = &fee
	var sequence = int64(txv1.Tx.SeqNum)
	transaction.Sequence = &sequence

	setSignatures(transaction, txv1.Signatures)
	setMemo(transaction, txv1.Tx.Memo)

	str, num := parseSourceAccount(txv1.Tx.SourceAccount)
	transaction.SourceAccount = str
	transaction.MuxedAccountId = num

	setCond(transaction, txv1.Tx.Cond)
}

func processTxFeeBumpEnvelope(transaction *model.Transaction, txFeeBump xdr.FeeBumpTransactionEnvelope) {
	processTxV1Envelope(transaction, *txFeeBump.Tx.InnerTx.V1)

	var feeBumpInfo = util.FeeBumpInfo{
		Fee: int32(txFeeBump.Tx.Fee),
	}

	str, num := parseSourceAccount(txFeeBump.Tx.FeeSource)
	feeBumpInfo.SourceAccount = str
	feeBumpInfo.MuxedAccountId = num

	transaction.FeeBumpInfo = &feeBumpInfo
}

func scMapToGo(scMap xdr.ScMap) interface{} {
	var data []interface{}
	for _, pair := range scMap {
		item := make(map[string]interface{})
		key := scValToGo(pair.Key)
		// flatten single value slice
		if val, ok := key.([]interface{}); ok {
			if len(val) == 1 {
				key = val[0]
			}
		}
		value := scValToGo(pair.Val)
		item["key"] = key
		item["value"] = value
		data = append(data, item)
	}
	return data
}

func scValToGo(val xdr.ScVal) interface{} {
	switch val.Type {
	case xdr.ScValTypeScvBool:
		return *val.B
	case xdr.ScValTypeScvVoid:
		return ""
	case xdr.ScValTypeScvError:
		data := make(map[string]interface{})
		data["type"] = val.Error.Type.String()
		data["code"] = val.Error.Code.String()
		return data
	case xdr.ScValTypeScvU32:
		return uint32(*val.U32)
	case xdr.ScValTypeScvI32:
		return int32(*val.I32)
	case xdr.ScValTypeScvU64:
		return uint64(*val.U64)
	case xdr.ScValTypeScvI64:
		return int64(*val.I64)
	case xdr.ScValTypeScvTimepoint:
		return uint64(*val.Timepoint)
	case xdr.ScValTypeScvDuration:
		return uint64(*val.Duration)
	case xdr.ScValTypeScvI128:
		hi := big.NewInt(int64(val.I128.Hi))
		lo := big.NewInt(0).SetUint64(uint64(val.I128.Lo))
		hi.Lsh(hi, 64)
		hi.Add(hi, lo)
		return hi.String()
	case xdr.ScValTypeScvU128:
		hi := big.NewInt(0).SetUint64(uint64(val.U128.Hi))
		lo := big.NewInt(0).SetUint64(uint64(val.U128.Lo))
		hi.Lsh(hi, 64)
		hi.Add(hi, lo)
		return hi.String()
	case xdr.ScValTypeScvI256:
		hihi := big.NewInt(int64(val.I256.HiHi))
		hilo := big.NewInt(0).SetUint64(uint64(val.I256.HiLo))
		lohi := big.NewInt(0).SetUint64(uint64(val.I256.LoHi))
		lolo := big.NewInt(0).SetUint64(uint64(val.I256.LoLo))
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, hilo)
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, lohi)
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, lolo)
		return hihi.String()
	case xdr.ScValTypeScvU256:
		hihi := big.NewInt(0).SetUint64(uint64(val.U256.HiHi))
		hilo := big.NewInt(0).SetUint64(uint64(val.U256.HiLo))
		lohi := big.NewInt(0).SetUint64(uint64(val.U256.LoHi))
		lolo := big.NewInt(0).SetUint64(uint64(val.U256.LoLo))
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, hilo)
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, lohi)
		hihi.Lsh(hihi, 64)
		hihi.Add(hihi, lolo)
		return hihi.String()
	case xdr.ScValTypeScvBytes:
		return *val.Bytes
	case xdr.ScValTypeScvString:
		return string(*val.Str)
	case xdr.ScValTypeScvSymbol:
		return string(*val.Sym)
	case xdr.ScValTypeScvVec:
		scVec, _ := val.GetVec()
		var list [](interface{})
		// list := make([]interface{}, len(*scVec))
		for _, value := range *scVec {
			list = append(list, scValToGo(value))
		}
		return list
	case xdr.ScValTypeScvMap:
		scMap, _ := val.GetMap()
		return scMapToGo(*scMap)
	case xdr.ScValTypeScvAddress:
		var address = ""
		if val.Address.AccountId != nil {
			address = val.Address.AccountId.Address()
		} else if val.Address.ContractId != nil {
			address = hex.EncodeToString(val.Address.ContractId[:])
		}
		return address
	case xdr.ScValTypeScvLedgerKeyContractInstance:
		return "ScvLedgerKeyContractInstance"
	case xdr.ScValTypeScvLedgerKeyNonce:
		return int64(val.NonceKey.Nonce)
	case xdr.ScValTypeScvContractInstance:
		if val.Instance == nil {
			return nil
		}
		instance := *val.Instance
		data := make(map[string]interface{})
		executable := make(map[string]interface{})
		data["executable"] = executable
		executable["type"] = instance.Executable.Type.String()
		if instance.Executable.WasmHash != nil {
			executable["wasmHash"] = hex.EncodeToString(instance.Executable.WasmHash[:])
		}
		if instance.Storage != nil {
			storage := scMapToGo(*instance.Storage)
			data["storage"] = storage
		}

		return data
	default:
		return nil
	}
}
