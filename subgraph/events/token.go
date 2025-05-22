package events

import (
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"github.com/gagliardetto/solana-go"
)

type TokenMetaData struct {
	Name   string
	Symbol string
}

type TokenData struct {
	Mint     solana.PublicKey
	Account  solana.PublicKey
	Decimals types.BigInt
	Metadata TokenMetaData
}
