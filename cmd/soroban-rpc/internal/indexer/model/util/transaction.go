package util

type Signature struct {
	Hint      string `json:"hint"`
	Signature string `json:"signature"`
}

type TypeItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Bonds struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type SignerKey struct {
	Type                 string `json:"type"`
	Ed25519              string `json:"ed25519"`
	PreAuthTx            string `json:"pre_auth_tx"`
	HashX                string `json:"hash_x"`
	Ed25519SignedPayload string `json:"ed25519_signed_payload"`
}

type Preconditions struct {
	TimeBonds       *Bonds       `json:"time_bonds"`
	LedgerBounds    *Bonds       `json:"ledger_bonds"`
	MinSeqNum       *int64       `json:"min_seq_num"`
	MinSeqAge       *int64       `json:"min_seq_age"`
	MinSeqLedgerGap *int32       `json:"min_seq_ledger_gap"`
	ExtraSigners    *[]SignerKey `json:"extra_signers"`
}

type FeeBumpInfo struct {
	Fee            int32   `json:"fee"`
	SourceAccount  *string `json:"source_account"`
	MuxedAccountId *int64  `json:"muxed_account_id"`
}
