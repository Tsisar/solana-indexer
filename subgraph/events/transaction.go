package events

import (
	"github.com/Tsisar/solana-indexer/subgraph/types"
)

type Transaction struct {
	Signature string
	Slot      types.BigInt
	Timestamp types.BigInt
}
