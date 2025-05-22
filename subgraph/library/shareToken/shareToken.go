package shareToken

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"github.com/Tsisar/solana-indexer/subgraph/utils"
	"gorm.io/gorm"
)

func InitializeAccount(ctx context.Context, db *gorm.DB, ev events.InitializeAccountInstruction) error {
	tokenAccount := subgraph.TokenAccount{ID: ev.Account}
	if _, err := tokenAccount.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load token account: %w", err)
	}
	tokenAccount.Mint = ev.Mint
	tokenAccount.Owner = ev.Owner

	if err := tokenAccount.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save token account: %w", err)
	}
	return nil
}

func Mint(ctx context.Context, db *gorm.DB, ev events.MintToInstruction, transaction events.Transaction) error {
	// Get the list of share token mints for all vaults
	mints, err := subgraph.GetShareTokenMints(ctx, db)
	if err != nil {
		return fmt.Errorf("[shareToken] failed to get share token mints: %w", err)
	}

	if !utils.Contains(mints, ev.Mint) {
		log.Debugf("[shareToken] mint not in tokens list: %s", ev.Mint)
		return nil
	}

	shareToken := subgraph.ShareToken{ID: ev.To}
	if _, err := shareToken.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load share token: %w", err)
	}
	shareToken.TotalMinted = utils.Val(shareToken.TotalMinted.Plus(&ev.Amount))
	shareToken.CurrentPrice = getCurrentPrice(ctx, db, ev.Mint)
	if err := shareToken.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save share token: %w", err)
	}

	tokenMint := subgraph.TokenMint{ID: transaction.Signature}
	if _, err := tokenMint.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load token mint: %w", err)
	}
	tokenMint.Amount = ev.Amount
	tokenMint.ToID = ev.To
	tokenMint.MintID = ev.Mint

	if err := tokenMint.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save token mint: %w", err)
	}

	return nil
}

func Burn(ctx context.Context, db *gorm.DB, ev events.BurnInstruction, transaction events.Transaction) error {
	// Get the list of share token mints for all vaults
	mints, err := subgraph.GetShareTokenMints(ctx, db)
	if err != nil {
		return fmt.Errorf("[shareToken] failed to get share token mints: %w", err)
	}

	if !utils.Contains(mints, ev.Mint) {
		log.Debugf("[shareToken] mint not in tokens list: %s", ev.Mint)
		return nil
	}

	shareToken := subgraph.ShareToken{ID: ev.From}
	if _, err := shareToken.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load share token: %w", err)
	}
	shareToken.TotalBurnt = utils.Val(shareToken.TotalBurnt.Plus(&ev.Amount))
	shareToken.CurrentPrice = getCurrentPrice(ctx, db, ev.Mint)
	if err := shareToken.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save share token: %w", err)
	}

	tokenBurn := subgraph.TokenBurn{ID: transaction.Signature}
	if _, err := tokenBurn.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load token burn: %w", err)
	}
	tokenBurn.Amount = ev.Amount
	tokenBurn.FromID = ev.From
	tokenBurn.MintID = ev.Mint

	if err := tokenBurn.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save token burn: %w", err)
	}
	return nil
}

func Transfer(ctx context.Context, db *gorm.DB, ev events.TransferInstruction, transaction events.Transaction) error {
	mint, err := getMint(ctx, db, ev.From, ev.To)
	if err != nil {
		return fmt.Errorf("[shareToken] failed to get mint: %w", err)
	}
	log.Debugf("[shareToken] mint: %s", mint)

	// Get the list of share token mints for all vaults
	mints, err := subgraph.GetShareTokenMints(ctx, db)
	if err != nil {
		return fmt.Errorf("[shareToken] failed to get share token mints: %w", err)
	}

	if !utils.Contains(mints, mint) {
		log.Debugf("[shareToken] mint not in tokens list: %s", mint)
		return nil
	}

	//TODO: Fix this part (this is todo from previous code)
	shareTokenIn := subgraph.ShareToken{ID: ev.To}
	if _, err := shareTokenIn.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load share token: %w", err)
	}
	shareTokenIn.TotalTransferIn = utils.Val(shareTokenIn.TotalTransferIn.Plus(&ev.Amount))
	shareTokenIn.CurrentPrice = getCurrentPrice(ctx, db, mint)
	if err := shareTokenIn.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save share token: %w", err)
	}

	shareTokenOut := subgraph.ShareToken{ID: ev.From}
	if _, err := shareTokenOut.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load share token: %w", err)
	}
	shareTokenOut.TotalTransferOut = utils.Val(shareTokenOut.TotalTransferOut.Plus(&ev.Amount))
	shareTokenOut.CurrentPrice = getCurrentPrice(ctx, db, mint)
	if err := shareTokenOut.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save share token: %w", err)
	}

	shareTokenTransfer := subgraph.ShareTokenTransfer{ID: transaction.Signature}
	if _, err := shareTokenTransfer.Load(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to load share token transfer: %w", err)
	}

	shareTokenTransfer.Amount = ev.Amount
	shareTokenTransfer.FromID = ev.From
	shareTokenTransfer.ToID = ev.To
	shareTokenTransfer.AuthorityID = ev.Authority
	shareTokenTransfer.MintID = mint

	if err := shareTokenTransfer.Save(ctx, db); err != nil {
		return fmt.Errorf("[shareToken] failed to save share token transfer: %w", err)
	}
	return nil
}

func getCurrentPrice(ctx context.Context, db *gorm.DB, tokenId string) types.BigInt {
	log.Debugf("[shareToken] get current price fo token: %s", tokenId)
	token := subgraph.Token{ID: tokenId}
	ok, err := token.Load(ctx, db)
	if err != nil || !ok {
		log.Warnf("[shareToken] failed to load token: %v", err)
		return types.ZeroBigInt()
	}
	return token.CurrentPrice
}

func getMint(ctx context.Context, db *gorm.DB, from, to string) (string, error) {
	tokenAccountFrom := subgraph.TokenAccount{ID: from}
	ok, err := tokenAccountFrom.Load(ctx, db)
	if err != nil {
		return "", fmt.Errorf("[shareToken] failed to load token account: %v", err)
	}
	if ok {
		return tokenAccountFrom.Mint, nil
	} else {
		log.Warnf("[shareToken] mint token account not found: %s", from)
	}

	tokenAccountTo := subgraph.TokenAccount{ID: to}
	ok, err = tokenAccountTo.Load(ctx, db)
	if err != nil {
		return "", fmt.Errorf("[shareToken] failed to load token account: %v", err)
	}
	if ok {
		return tokenAccountTo.Mint, nil
	} else {
		log.Warnf("[shareToken] mint token account not found: %s", to)
	}

	return "", fmt.Errorf("[shareToken] mint token account not found: %s", to)
}
