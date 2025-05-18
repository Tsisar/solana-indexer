package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
)

// BigDecimal wraps big.Float for use with GORM and JSON
type BigDecimal struct {
	*big.Float
}

// GormDataType tells GORM to use "numeric" in SQL
func (BigDecimal) GormDataType() string {
	return "numeric"
}

// Value implements driver.Valuer (for writing to DB)
func (b BigDecimal) Value() (driver.Value, error) {
	if b.Float == nil {
		return nil, nil
	}
	return b.Text('f', -1), nil // store as text (NUMERIC-compatible)
}

// Scan implements sql.Scanner (for reading from DB)
func (b *BigDecimal) Scan(value interface{}) error {
	if value == nil {
		b.Float = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return b.setString(string(v))
	case string:
		return b.setString(v)
	default:
		return fmt.Errorf("unsupported type for BigDecimal: %T", value)
	}
}

func (b *BigDecimal) setString(s string) error {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return fmt.Errorf("failed to parse big.Float: %w", err)
	}
	b.Float = f
	return nil
}

// MarshalJSON implements json.Marshaler
func (b BigDecimal) MarshalJSON() ([]byte, error) {
	if b.Float == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", b.Text('f', -1))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (b *BigDecimal) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return b.setString(s)
	}

	// fallback to float64
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		b.Float = big.NewFloat(f)
		return nil
	}

	return fmt.Errorf("failed to unmarshal BigDecimal from JSON: %s", string(data))
}

// Sub returns a new BigDecimal which is the result of b - other.
func (b *BigDecimal) Sub(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil {
		return &BigDecimal{Float: nil}
	}
	if other == nil || other.Float == nil {
		// b - nil = b
		return &BigDecimal{Float: new(big.Float).Set(b.Float)}
	}
	result := new(big.Float).Sub(b.Float, other.Float)
	return &BigDecimal{Float: result}
}

// SafeDiv divides b by c and returns a new BigDecimal.
// Returns nil.Float if c is nil or zero to avoid panic/div-by-zero.
func (b *BigDecimal) SafeDiv(c *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || c == nil || c.Float == nil {
		return &BigDecimal{Float: nil}
	}

	if c.Float.Sign() == 0 {
		return &BigDecimal{Float: nil}
	}

	return &BigDecimal{Float: new(big.Float).Quo(b.Float, c.Float)}
}

func (b *BigDecimal) Mul(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).Mul(b.Float, other.Float)}
}

// Plus returns a + other as a new BigDecimal.
// Returns nil.Float if any operand is nil.
func (b *BigDecimal) Plus(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{
		Float: new(big.Float).Add(b.Float, other.Float),
	}
}
