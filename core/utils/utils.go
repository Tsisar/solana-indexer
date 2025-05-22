package utils

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"time"
)

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
