package events

import (
	"github.com/Tsisar/solana-indexer/subgraph/types"
)

type TransferInstruction struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Authority string           `json:"authority"`
	Amount    types.BigDecimal `json:"amount"`
}

type MintToInstruction struct {
	To     string           `json:"to"`
	Mint   string           `json:"mint"`
	Amount types.BigDecimal `json:"amount"`
}

type BurnInstruction struct {
	From   string           `json:"from"`
	Mint   string           `json:"mint"`
	Amount types.BigDecimal `json:"amount"`
}

type InitializeMintInstruction struct {
	Mint            string       `json:"mint"`
	MintAuthority   string       `json:"mint_authority"`
	FreezeAuthority string       `json:"freeze_authority"`
	Decimals        types.BigInt `json:"decimals"`
}

type InitializeAccountInstruction struct {
	Account string `json:"account"`
	Mint    string `json:"mint"`
	Owner   string `json:"owner"`
}
