package account

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/storage/model/subgraph"
	"gorm.io/gorm"
)

func UpdateAccount(ctx context.Context, db *gorm.DB, authorityId, tokenAccountId, shareAccountId string) error {
	authorityAccount := subgraph.Account{ID: authorityId}
	if _, err := authorityAccount.Load(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to load authority account: %w", err)
	}
	if err := authorityAccount.Save(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to save authority account: %w", err)
	}

	tokenAccount := subgraph.TokenWallet{ID: tokenAccountId}
	if _, err := tokenAccount.Load(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to load token account: %w", err)
	}
	tokenAccount.AuthorityID = authorityId
	if err := tokenAccount.Save(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to save token account: %w", err)
	}

	shareAccount := subgraph.TokenWallet{ID: shareAccountId}
	if _, err := shareAccount.Load(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to load share account: %w", err)
	}
	shareAccount.AuthorityID = authorityId
	if err := shareAccount.Save(ctx, db); err != nil {
		return fmt.Errorf("[account] failed to save share account: %w", err)
	}

	return nil
}
