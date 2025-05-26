package parser

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"sync"
	"time"
)

// lutCacheEntry represents a cached result of a fetched Address Lookup Table (LUT).
// It stores the decoded addresses, how many addresses there are, and when it was cached.
type lutCacheEntry struct {
	Addresses solana.PublicKeySlice
	Count     int
	Timestamp time.Time
}

var (
	lutCache     sync.Map   // Caches LUT results by PublicKey
	lutCacheLock sync.Mutex // Prevents concurrent LUT fetches for the same key
)

const lutTTL = 24 * time.Hour // Time-to-live for LUT cache entries

// maxUint8Slice returns the maximum value found in two uint8 slices.
func maxUint8Slice(a, b []uint8) uint8 {
	max := uint8(0)
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	for _, v := range b {
		if v > max {
			max = v
		}
	}
	return max
}

// getOrFetchLUT returns cached addresses for the given lookup table key,
// or fetches them from the chain if not cached, outdated, or incomplete.
//
// expectedMaxIndex is used to ensure that the cache contains enough addresses.
// This function double-checks the cache before and after acquiring a lock to
// avoid redundant RPC calls in concurrent contexts.
func getOrFetchLUT(ctx context.Context, key solana.PublicKey, expectedMaxIndex int) (solana.PublicKeySlice, error) {
	// First optimistic cache read
	if val, ok := lutCache.Load(key); ok {
		entry := val.(lutCacheEntry)

		// Check if the cache is fresh and sufficient
		if time.Since(entry.Timestamp) < lutTTL && expectedMaxIndex < entry.Count {
			return entry.Addresses, nil
		}
	}

	// Lock to prevent concurrent fetches for the same key
	lutCacheLock.Lock()
	defer lutCacheLock.Unlock()

	// Check again inside the lock to avoid race conditions
	if val, ok := lutCache.Load(key); ok {
		entry := val.(lutCacheEntry)
		if time.Since(entry.Timestamp) < lutTTL && expectedMaxIndex < entry.Count {
			return entry.Addresses, nil
		}
	}

	// Fetch from blockchain
	addresses, err := fetchAddressLookupTable(ctx, key)
	if err != nil {
		return nil, err
	}

	// Cache the result
	lutCache.Store(key, lutCacheEntry{
		Addresses: addresses,
		Count:     len(addresses),
		Timestamp: time.Now(),
	})

	return addresses, nil
}
