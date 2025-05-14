package events

import (
	"crypto/sha256"
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

var Discriminators = make(map[[8]byte]string)

func init() {
	for name := range Registry {
		hash := sha256.Sum256([]byte("event:" + name))
		var disc [8]byte
		copy(disc[:], hash[:8])
		Discriminators[disc] = name
	}
}
