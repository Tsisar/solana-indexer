package events

import "math/big"

type Transaction struct {
	Signature string
	Slot      *big.Int
	Timestamp *big.Int
}
