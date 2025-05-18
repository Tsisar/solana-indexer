package events

import (
	"github.com/gagliardetto/solana-go"
	"math/big"
)

type TokenMetaData struct {
	Name   string
	Symbol string
}

type TokenData struct {
	Mint     solana.PublicKey
	Account  solana.PublicKey
	Decimals big.Int
	Metadata TokenMetaData
}
