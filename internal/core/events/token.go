package events

import (
	"github.com/gagliardetto/solana-go"
)

type TokenMetaData struct {
	Name   string `borsh:"name"`
	Symbol string `borsh:"symbol"`
}

type TokenData struct {
	Mint     solana.PublicKey `borsh:"mint"`
	Account  solana.PublicKey `borsh:"account"`
	Decimals uint8            `borsh:"decimals"`
	Metadata TokenMetaData    `borsh:"metadata"`
}
