package parser

import (
	"encoding/json"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/stellar/soroban-rpc/cmd/soroban-rpc/internal/indexer/model"
)

func GetTrustLineEntryModel(entry xdr.LedgerEntry) *model.TrustLineEntry {
	trustLine := entry.Data.TrustLine
	var poolId []byte
	if trustLine.Asset.LiquidityPoolId != nil {
		poolId = trustLine.Asset.LiquidityPoolId[:]
	}

	em := &model.TrustLineEntry{
		AccountId:       trustLine.AccountId.Address(),
		AssetType:       trustLine.Asset.Type,
		LiquidityPoolId: poolId,
		Balance:         trustLine.Balance,
		Limit:           trustLine.Limit,
		Flags:           trustLine.Flags,
		Ext:             trustLine.Ext,
	}

	switch trustLine.Asset.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.AssetCode = trustLine.Asset.AlphaNum4.AssetCode[:]
		em.AssetIssuer = trustLine.Asset.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.AssetCode = trustLine.Asset.AlphaNum12.AssetCode[:]
		em.AssetIssuer = trustLine.Asset.AlphaNum12.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypePoolShare:
		em.AssetCode = trustLine.Asset.LiquidityPoolId[:]
		break
	}

	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}

	return em
}

func GetOfferEntryModel(entry xdr.LedgerEntry) *model.OfferEntry {
	offer := entry.Data.Offer
	em := &model.OfferEntry{
		OfferId:          offer.OfferId,
		SellerId:         offer.SellerId.Address(),
		SellingAssetType: offer.Selling.Type,
		BuyingAssetType:  offer.Buying.Type,
		Amount:           offer.Amount,
		Price:            toDecimal(offer.Price),
		Flags:            offer.Flags,
		Ext:              offer.Ext,
	}

	switch offer.Selling.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.SellingAssetCode = offer.Selling.AlphaNum4.AssetCode[:]
		em.SellingAssetIssuer = offer.Selling.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.SellingAssetCode = offer.Selling.AlphaNum12.AssetCode[:]
		em.SellingAssetIssuer = offer.Selling.AlphaNum12.Issuer.Address()
		break
	}

	switch offer.Buying.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.BuyingAssetCode = offer.Buying.AlphaNum4.AssetCode[:]
		em.BuyingAssetIssuer = offer.Buying.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.BuyingAssetCode = offer.Buying.AlphaNum12.AssetCode[:]
		em.BuyingAssetIssuer = offer.Buying.AlphaNum12.Issuer.Address()
		break
	}

	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}

	return em
}

func GetContractDataModel(entry xdr.LedgerEntry) *model.ContractDataEntry {
	contractData := entry.Data.ContractData

	contractId := ""
	if contractData.Contract.ContractId != nil {
		contractId = strkey.MustEncode(strkey.VersionByteContract, contractData.Contract.ContractId[:])
	} else {
		return nil
	}
	keyXdr, _ := xdr.MarshalBase64(contractData.Key)
	valXdr, _ := xdr.MarshalBase64(contractData.Val)

	durability := "persistent"
	if contractData.Durability == xdr.ContractDataDurabilityTemporary {
		durability = "temporary"
	}

	em := &model.ContractDataEntry{
		ContractId: contractId,
		KeyXdr:     keyXdr,
		Durability: durability,
		ValXdr:     valXdr,
	}

	return em
}

func GetAccountEntryModel(entry xdr.LedgerEntry) *model.AccountEntry {
	inflationDest, _ := entry.Data.Account.InflationDest.GetAddress()
	em := &model.AccountEntry{
		AccountId:     entry.Data.Account.AccountId.Address(),
		Balance:       entry.Data.Account.Balance,
		SeqNum:        entry.Data.Account.SeqNum,
		NumSubEntries: entry.Data.Account.NumSubEntries,
		InflationDest: inflationDest,
		Flags:         entry.Data.Account.Flags,
		HomeDomain:    entry.Data.Account.HomeDomain,
		Thresholds:    entry.Data.Account.Thresholds[:],
		Signers:       getSignerIds(entry.Data.Account.Signers),
		Ext:           entry.Data.Account.Ext.V,
	}
	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}
	return em
}

func GetDataEntryModel(entry xdr.LedgerEntry) *model.DataEntry {
	data := entry.Data.Data
	em := &model.DataEntry{
		AccountId: data.AccountId.Address(),
		DataName:  data.DataName,
		Ext:       data.Ext,
		DataValue: data.DataValue,
	}
	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}
	return em
}

func GetClaimableBalanceEntryModel(entry xdr.LedgerEntry) *model.ClaimableBalanceEntry {
	claimableBalance := entry.Data.ClaimableBalance
	jsonData, _ := json.Marshal(claimableBalance.Claimants)
	em := &model.ClaimableBalanceEntry{
		BalanceId: claimableBalance.BalanceId.V0.HexString(),
		Claimants: jsonData,
		AssetType: claimableBalance.Asset.Type,
	}
	switch claimableBalance.Asset.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.AssetCode = claimableBalance.Asset.AlphaNum4.AssetCode[:]
		em.AssetIssuer = claimableBalance.Asset.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.AssetCode = claimableBalance.Asset.AlphaNum12.AssetCode[:]
		em.AssetIssuer = claimableBalance.Asset.AlphaNum12.Issuer.Address()
		break
	}
	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}
	return em
}

func GetLiquidityPoolEntryModel(entry xdr.LedgerEntry) *model.LiquidityPoolEntry {
	lp := entry.Data.LiquidityPool
	em := &model.LiquidityPoolEntry{
		LiquidityPoolId:          lp.LiquidityPoolId[:],
		Type:                     lp.Body.Type,
		Fee:                      lp.Body.ConstantProduct.Params.Fee,
		ReserveA:                 lp.Body.ConstantProduct.ReserveA,
		ReserveB:                 lp.Body.ConstantProduct.ReserveB,
		TotalPoolShares:          lp.Body.ConstantProduct.TotalPoolShares,
		PoolSharesTrustLineCount: lp.Body.ConstantProduct.PoolSharesTrustLineCount,
	}

	assetA := lp.Body.ConstantProduct.Params.AssetA
	em.AssetAType = assetA.Type
	switch assetA.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.AssetACode = assetA.AlphaNum4.AssetCode[:]
		em.AssetAIssuer = assetA.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.AssetACode = assetA.AlphaNum12.AssetCode[:]
		em.AssetAIssuer = assetA.AlphaNum12.Issuer.Address()
		break
	}

	assetB := lp.Body.ConstantProduct.Params.AssetB
	em.AssetBType = assetB.Type
	switch assetB.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		em.AssetBCode = assetB.AlphaNum4.AssetCode[:]
		em.AssetBIssuer = assetB.AlphaNum4.Issuer.Address()
		break
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		em.AssetBCode = assetB.AlphaNum12.AssetCode[:]
		em.AssetBIssuer = assetB.AlphaNum12.Issuer.Address()
		break
	}

	em.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq
	if entry.SponsoringID() != nil {
		em.SponsoringId = entry.SponsoringID().Address()
	}

	return em
}

func toDecimal(p xdr.Price) decimal.Decimal {
	if p.D == 0 {
		// handle division by zero
		return decimal.NewFromInt32(0)
	}
	return decimal.NewFromInt32(int32(p.N)).Div(decimal.NewFromInt32(int32(p.D)))
}

type Signer struct {
	Key    string `json:"key"`
	Weight uint32 `json:"weight"`
}

func getSignerIds(signers []xdr.Signer) *[]byte {
	if len(signers) == 0 {
		return nil
	}

	var signerIds []Signer
	for _, s := range signers {
		signerIds = append(signerIds, Signer{Key: s.Key.Address(), Weight: uint32(s.Weight)})
	}
	data, _ := json.Marshal(signerIds)

	return &data
}
