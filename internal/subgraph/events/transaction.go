package events

import (
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
)

type Transaction struct {
	Signature  string
	Slot       types.BigInt
	Timestamp  types.BigInt
	EventIndex types.BigInt
}

func NewTransaction(event core.Event) Transaction {
	return Transaction{
		Signature:  event.TransactionSignature,
		Slot:       *types.NewBigIntFromUint64(event.Slot),
		Timestamp:  *types.NewBigIntFromInt64(event.BlockTime),
		EventIndex: *types.NewBigIntFromInt64(int64(event.LogIndex)),
	}
}
