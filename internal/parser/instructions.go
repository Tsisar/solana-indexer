package parser

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/events"
	"github.com/Tsisar/solana-indexer/internal/maping"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

var client = rpc.New(config.App.RPCEndpoint)

func parseTokenInstructions(ctx context.Context, db *storage.Gorm, sig string, tx *rpc.GetTransactionResult) error {
	parsedTx, err := tx.Transaction.GetTransaction()
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}
	msg := &parsedTx.Message

	if err := resolveAddressLookupsIfNeeded(ctx, msg); err != nil {
		return fmt.Errorf("resolve lookups for tx %s: %w", sig, err)
	}

	if tx.Meta == nil || tx.Meta.InnerInstructions == nil {
		log.Debugf("No inner instructions for transaction %s", sig)
		return nil
	}

	for _, inner := range tx.Meta.InnerInstructions {
		for i, innerInstr := range inner.Instructions {
			if err := processInnerInstruction(ctx, db, msg, sig, tx, inner.Index, i, &innerInstr); err != nil {
				log.Warnf("inner parse error: %v", err)
			}
		}
	}

	return nil
}

func processInnerInstruction(ctx context.Context, db *storage.Gorm, msg *solana.Message, sig string,
	tx *rpc.GetTransactionResult, instrIndex uint16, innerIndex int, instr *solana.CompiledInstruction,
) error {
	if len(instr.Data) == 0 || instr.Data[0] > token.Instruction_InitializeMint2 {
		//log.Debugf("Skipping unknown token instruction ID: %d", innerInstr.Data[0])
		return nil
	}

	programID, err := msg.Account(instr.ProgramIDIndex)
	if err != nil {
		log.Warnf("ProgramIDIndex out of range: %d", instr.ProgramIDIndex)
		return nil
	}
	if !programID.Equals(solana.TokenProgramID) {
		return nil
	}

	var accounts []*solana.AccountMeta
	for _, accIdx := range instr.Accounts {
		pubKey, err := msg.Account(accIdx)
		if err != nil {
			log.Errorf("Account index out of range: %d", accIdx)
			continue
		}

		writable := false
		writable, err = msg.IsWritable(pubKey)
		if err != nil {
			log.Errorf("Failed to check if account is writable: %v", err)
			continue
		}
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  pubKey,
			IsSigner:   msg.IsSigner(pubKey),
			IsWritable: writable,
		})
	}

	decoded, err := token.DecodeInstruction(accounts, instr.Data)
	if err != nil {
		log.Errorf("inner parse error: decode instruction: %v", err)
		return nil
	}

	var blockTime int64
	if tx.BlockTime != nil {
		blockTime = int64(*tx.BlockTime)
	}

	name, mapped := mapTokenInstruction(decoded.Impl)
	if name == "" {
		log.Debugf("Skipping unknown token instruction ID: %d", decoded.TypeID.Uint8())
		return nil
	}

	evRecord := storage.Event{
		TransactionSignature: sig,
		Slot:                 tx.Slot,
		BlockTime:            blockTime,
		LogIndex:             int(instrIndex+1)*100 + innerIndex,
		Name:                 name,
	}
	evRecord.JsonEv, _ = json.Marshal(mapped)

	if err := db.SaveEvent(ctx, evRecord); err != nil {
		return fmt.Errorf("save event %s: %w", evRecord.Name, err)
	}

	maping.Instruction(evRecord)
	return nil
}

func resolveAddressLookupsIfNeeded(ctx context.Context, msg *solana.Message) error {
	if msg.IsVersioned() && len(msg.AddressTableLookups) > 0 && !msg.IsResolved() {
		addressTables := make(map[solana.PublicKey]solana.PublicKeySlice)
		for _, lookup := range msg.AddressTableLookups {
			//log.Debugf("LUT address:", lookup.AccountKey.String())
			addresses, err := fetchAddressLookupTable(ctx, lookup.AccountKey)
			if err != nil {
				return fmt.Errorf("failed to fetch LUT %s: %v", lookup.AccountKey.String(), err)
			}
			addressTables[lookup.AccountKey] = addresses
		}
		if err := msg.SetAddressTables(addressTables); err != nil {
			return fmt.Errorf("failed to set address tables: %w", err)
		}
		if err := msg.ResolveLookups(); err != nil {
			return fmt.Errorf("failed to resolve lookups: %w", err)
		}
	}
	return nil
}

func mapTokenInstruction(inst interface{}) (string, any) {
	switch i := inst.(type) {
	case *token.Transfer:
		return "TransferInstruction", events.TransferInstruction{
			From:      &i.GetSourceAccount().PublicKey,
			To:        &i.GetDestinationAccount().PublicKey,
			Authority: &i.GetOwnerAccount().PublicKey,
			Amount:    i.Amount,
		}
	case *token.TransferChecked:
		return "TransferCheckedInstruction", events.TransferCheckedInstruction{
			From:      &i.GetSourceAccount().PublicKey,
			To:        &i.GetDestinationAccount().PublicKey,
			Authority: &i.GetOwnerAccount().PublicKey,
			Mint:      &i.GetMintAccount().PublicKey,
			Amount:    i.Amount,
			Decimals:  i.Decimals,
		}
	case *token.MintTo:
		return "MintToInstruction", events.MintToInstruction{
			To:     &i.GetDestinationAccount().PublicKey,
			Mint:   &i.GetMintAccount().PublicKey,
			Amount: i.Amount,
		}
	case *token.MintToChecked:
		return "MintToCheckedInstruction", events.MintToCheckedInstruction{
			To:       &i.GetDestinationAccount().PublicKey,
			Mint:     &i.GetMintAccount().PublicKey,
			Amount:   i.Amount,
			Decimals: i.Decimals,
		}
	case *token.Burn:
		return "BurnInstruction", events.BurnInstruction{
			From:   &i.GetSourceAccount().PublicKey,
			Mint:   &i.GetMintAccount().PublicKey,
			Amount: i.Amount,
		}
	case *token.BurnChecked:
		return "BurnCheckedInstruction", events.BurnCheckedInstruction{
			From:     &i.GetSourceAccount().PublicKey,
			Mint:     &i.GetMintAccount().PublicKey,
			Amount:   i.Amount,
			Decimals: i.Decimals,
		}
	case *token.InitializeMint2:
		return "InitializeMint2Instruction", events.InitializeMint2Instruction{
			Mint:            &i.GetMintAccount().PublicKey,
			MintAuthority:   i.MintAuthority,
			FreezeAuthority: i.FreezeAuthority,
			Decimals:        i.Decimals,
		}
	case *token.InitializeAccount3:
		return "InitializeAccount3Instruction", events.InitializeAccount3Instruction{
			Mint:  &i.GetMintAccount().PublicKey,
			Owner: i.Owner,
		}
	default:
		return "", nil
	}
}

func fetchAddressLookupTable(ctx context.Context, address solana.PublicKey) (solana.PublicKeySlice, error) {
	resp, err := client.GetAccountInfoWithOpts(
		ctx,
		address,
		&rpc.GetAccountInfoOpts{
			Encoding:   "base64",
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get LUT account info: %w", err)
	}
	if resp == nil || resp.Value == nil || resp.Value.Data == nil {
		return nil, fmt.Errorf("empty LUT account data")
	}

	data := resp.Value.Data.GetBinary()
	if len(data) < 8 {
		return nil, fmt.Errorf("invalid LUT data (too short)")
	}

	numAddresses := binary.LittleEndian.Uint32(data[4:8])
	available := (len(data) - 8) / 32
	if available < int(numAddresses) {
		//log.Warnf("LUT %s claims %d addresses, but only %d available — truncating", address, numAddresses, available)
		numAddresses = uint32(available)
	}

	addresses := make([]solana.PublicKey, 0, numAddresses)
	offset := 8

	for i := uint32(0); i < numAddresses; i++ {
		var pub solana.PublicKey
		copy(pub[:], data[offset:offset+32])
		addresses = append(addresses, pub)
		offset += 32
	}

	return addresses, nil
}
