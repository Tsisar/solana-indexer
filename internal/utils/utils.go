package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"github.com/gagliardetto/solana-go"
	"math/big"
	"time"
)

var msPerDay = big.NewInt(86_400_000) // 24*60*60*1000
var daysPerYear = big.NewFloat(365.0)
var attempts = 5 // TODO: make this configurable from application settings

// Retry executes the provided function up to `attempts` times until it succeeds.
// If all attempts fail, it returns an error.
// Generic version that works for any return type.
func Retry[T any](fn func() (T, error)) (T, error) {
	var zero T

	for i := 0; i < attempts; i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		time.Sleep(time.Second)
	}
	return zero, fmt.Errorf("[utils] retry failed after %d attempts", attempts)
}

// Ptr returns a pointer to the given value of any type.
func Ptr[T any](v T) *T {
	return &v
}

// BlockTime converts a *solana.UnixTimeSeconds to int64.
// Returns 0 if the input is nil.
func BlockTime(bt *solana.UnixTimeSeconds) int64 {
	if bt != nil {
		return int64(*bt)
	}
	return 0
}

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

// Val returns the dereferenced value of the given pointer,
// or the zero value of T if the pointer is nil.
func Val[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}

// ToScaledBigDecimal converts a BigInt, which represents a raw on-chain amount,
// to a BigDecimal that represents the human-readable value, by dividing by 10^decimals.
func ToScaledBigDecimal(val *types.BigInt, decimals *types.BigInt) *types.BigDecimal {
	if val == nil || val.Int == nil {
		return &types.BigDecimal{Float: nil}
	}

	// Default to 0 decimals if not provided, so no scaling happens.
	if decimals == nil || decimals.Int == nil || decimals.Int.Sign() == 0 {
		return val.ToBigDecimal()
	}

	divisorInt := new(big.Int).Exp(big.NewInt(10), decimals.Int, nil)
	divisorBD := (&types.BigInt{Int: divisorInt}).ToBigDecimal()

	if divisorBD == nil || divisorBD.Sign() == 0 {
		return val.ToBigDecimal()
	}

	return val.ToBigDecimal().SafeDiv(divisorBD)
}

func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
