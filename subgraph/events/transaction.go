package events

import (
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph/types"
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
		Slot:       *types.BigIntFromUint64(event.Slot),
		Timestamp:  *types.BigIntFromInt64(event.BlockTime),
		EventIndex: *types.BigIntFromInt64(int64(event.LogIndex)),
	}
}
