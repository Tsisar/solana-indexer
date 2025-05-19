package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
)

// BigDecimalPrecision defines the bit precision for all BigDecimal operations.
const BigDecimalPrecision = 256

// BigDecimal wraps big.Float for use with GORM and JSON
// All operations use a mantissa of BigDecimalPrecision bits.
type BigDecimal struct {
	*big.Float
}

// GormDataType tells GORM to use "numeric" in SQL
func (BigDecimal) GormDataType() string {
	return "numeric(38,20)"
}

// Value implements driver.Valuer (for writing to DB)
func (b BigDecimal) Value() (driver.Value, error) {
	if b.Float == nil {
		return nil, nil
	}
	// use full precision text representation
	return b.Text('f', -1), nil
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
	// parse with defined precision
	f, _, err := big.ParseFloat(s, 10, BigDecimalPrecision, big.ToNearestEven)
	if err != nil {
		return fmt.Errorf("failed to parse big.Float: %w", err)
	}
	// ensure precision
	b.Float = f.SetPrec(BigDecimalPrecision)
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
	// try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return b.setString(s)
	}
	// fallback to float64
	var f64 float64
	if err := json.Unmarshal(data, &f64); err == nil {
		f := new(big.Float).SetPrec(BigDecimalPrecision).SetFloat64(f64)
		b.Float = f
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
		// b - nil = b copy
		return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Set(b.Float)}
	}
	result := new(big.Float).SetPrec(BigDecimalPrecision).Sub(b.Float, other.Float)
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
	result := new(big.Float).SetPrec(BigDecimalPrecision).Quo(b.Float, c.Float)
	return &BigDecimal{Float: result}
}

// Mul returns the product b * other as a new BigDecimal.
func (b *BigDecimal) Mul(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	result := new(big.Float).SetPrec(BigDecimalPrecision).Mul(b.Float, other.Float)
	return &BigDecimal{Float: result}
}

// Plus returns a + other as a new BigDecimal.
func (b *BigDecimal) Plus(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	result := new(big.Float).SetPrec(BigDecimalPrecision).Add(b.Float, other.Float)
	return &BigDecimal{Float: result}
}

func (b *BigDecimal) Zero() {
	if b.Float == nil {
		b.Float = big.NewFloat(0).SetPrec(BigDecimalPrecision)
	} else {
		b.Float.SetPrec(BigDecimalPrecision).SetFloat64(0)
	}
}

func ZeroBigDecimal() BigDecimal {
	return BigDecimal{Float: big.NewFloat(0).SetPrec(BigDecimalPrecision)}
}
