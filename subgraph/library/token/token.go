package token

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/monitoring"
	"github.com/Tsisar/solana-indexer/subgraph/events"

	"github.com/Tsisar/solana-indexer/storage/model/subgraph"

	"gorm.io/gorm"
)

func UpsertUnderlyingToken(ctx context.Context, db *gorm.DB, event events.VaultInitEvent) (*subgraph.Token, error) {
	token := subgraph.Token{
		ID: event.UnderlyingToken.Mint.String(),
	}
	if _, err := token.Load(ctx, db); err != nil {
		return nil, fmt.Errorf("[UpsertUnderlyingToken] failed to load token: %w", err)
	}

	token.Decimals = event.UnderlyingToken.Decimals
	token.Name = event.UnderlyingToken.Metadata.Name
	token.Symbol = event.UnderlyingToken.Metadata.Symbol

	if err := token.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("[UpsertUnderlyingToken] failed to save token: %w", err)
	}
	monitoring.Token(token)

	return &token, nil
}

func UpsertShareToken(ctx context.Context, db *gorm.DB, event events.VaultInitEvent) (*subgraph.Token, error) {
	token := subgraph.Token{
		ID: event.ShareToken.Mint.String(),
	}
	if _, err := token.Load(ctx, db); err != nil {
		return nil, fmt.Errorf("[UpsertUnderlyingToken] failed to load token: %w", err)
	}

	token.Decimals = event.ShareToken.Decimals
	token.Name = event.ShareToken.Metadata.Name
	token.Symbol = event.ShareToken.Metadata.Symbol

	if err := token.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("[UpsertUnderlyingToken] failed to save token: %w", err)
	}
	monitoring.Token(token)

	return &token, nil
}
