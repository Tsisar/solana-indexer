package events

import "github.com/gagliardetto/solana-go"

type TransferInstruction struct {
	From      *solana.PublicKey `json:"from"`
	To        *solana.PublicKey `json:"to"`
	Authority *solana.PublicKey `json:"authority"`
	Amount    *uint64           `json:"amount"`
}

type MintToInstruction struct {
	To     *solana.PublicKey `json:"to"`
	Mint   *solana.PublicKey `json:"mint"`
	Amount *uint64           `json:"amount"`
}

type BurnInstruction struct {
	From   *solana.PublicKey `json:"from"`
	Mint   *solana.PublicKey `json:"mint"`
	Amount *uint64           `json:"amount"`
}

type InitializeMintInstruction struct {
	Mint            *solana.PublicKey `json:"mint"`
	MintAuthority   *solana.PublicKey `json:"mint_authority"`
	FreezeAuthority *solana.PublicKey `json:"freeze_authority"`
	Decimals        *uint8            `json:"decimals"`
}

type InitializeAccountInstruction struct {
	Account *solana.PublicKey `json:"account"`
	Mint    *solana.PublicKey `json:"mint"`
	Owner   *solana.PublicKey `json:"owner"`
}
