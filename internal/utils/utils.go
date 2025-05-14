package utils

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"time"
)

var attempts = 5 //TODO: get if from config

func Retry[T any](fn func() (T, error)) (T, error) {
	var zero T

	for i := 0; i < attempts; i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		time.Sleep(time.Second)
	}
	return zero, fmt.Errorf("retry failed after %d attempts", attempts)
}

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// BlockTime converts a solana.UnixTimeSeconds to int64.
func BlockTime(bt *solana.UnixTimeSeconds) int64 {
	if bt != nil {
		return int64(*bt)
	}
	return 0
}
