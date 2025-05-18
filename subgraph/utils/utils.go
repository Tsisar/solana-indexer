package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	coremodel "github.com/Tsisar/solana-indexer/storage/model/core"
	"math/big"
)

const msPerDay = 86_400_000

func GenerateId(event coremodel.Event) string {
	hasher := sha256.New()
	hasher.Write([]byte(event.TransactionSignature))
	hasher.Write([]byte(event.JsonEv))
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum)
}

func MillisToDays(durationMs int64) float64 {
	return float64(durationMs) / float64(msPerDay)
}

func ParseBigIntFromString(s string) (*big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid big.Int string: %s", s)
	}
	return n, nil
}
