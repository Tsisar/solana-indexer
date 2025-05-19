package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"math/big"
)

var msPerDay = big.NewInt(86_400_000) // 24*60*60*1000
var daysPerYear = big.NewFloat(365.0)

// GenerateId creates a SHA-256 hash from all input strings concatenated in order.
// Returns the result as a hex-encoded string.
func GenerateId(parts ...string) string {
	hasher := sha256.New()
	for _, s := range parts {
		hasher.Write([]byte(s))
	}
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum)
}

// MillisToDays converts duration in milliseconds to number of days as *types.BigDecimal.
func MillisToDays(durationMs *types.BigInt) *types.BigDecimal {
	if durationMs == nil || durationMs.Int == nil {
		return &types.BigDecimal{Float: nil}
	}

	ms := new(big.Float).SetInt(durationMs.Int)
	day := new(big.Float).SetInt(msPerDay)

	result := new(big.Float).Quo(ms, day)
	return &types.BigDecimal{Float: result}
}

// DaysToYearFactor computes 365 / durationInDays (BigDecimal)
func DaysToYearFactor(durationInDays *types.BigDecimal) *types.BigDecimal {
	if durationInDays == nil || durationInDays.Float == nil {
		return &types.BigDecimal{Float: nil}
	}

	return &types.BigDecimal{
		Float: new(big.Float).Quo(daysPerYear, durationInDays.Float),
	}
}

func FormatBigDecimal(b *types.BigDecimal, prec int) string {
	if b == nil || b.Float == nil {
		return "nil"
	}
	return b.Float.Text('f', prec)
}

func FormatBigInt(b *types.BigInt) string {
	if b == nil || b.Int == nil {
		return "nil"
	}
	return b.Int.String()
}
