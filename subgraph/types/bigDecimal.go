package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

const BigDecimalPrecision = 66 // bits, roughly 20 decimal places

type BigDecimal struct {
	*big.Float
}

// Ensure GORM saves as NUMERIC
func (BigDecimal) GormDataType() string {
	return "numeric(38,20)"
}

// Constructors
func NewBigDecimalFromFloat(f float64) BigDecimal {
	return BigDecimal{Float: big.NewFloat(f).SetPrec(BigDecimalPrecision)}
}

func NewBigDecimalFromString(s string) (BigDecimal, error) {
	f, _, err := big.ParseFloat(s, 10, BigDecimalPrecision, big.ToNearestEven)
	if err != nil {
		return BigDecimal{}, fmt.Errorf("parse BigDecimal from string: %w", err)
	}
	return BigDecimal{Float: f}, nil
}

func MustBigDecimalFromString(s string) BigDecimal {
	val, err := NewBigDecimalFromString(s)
	if err != nil {
		panic(err)
	}
	return val
}

func ZeroBigDecimal() BigDecimal {
	return BigDecimal{Float: big.NewFloat(0).SetPrec(BigDecimalPrecision)}
}

func (b *BigDecimal) Zero() {
	if b.Float == nil {
		b.Float = big.NewFloat(0).SetPrec(BigDecimalPrecision)
	} else {
		b.Float.SetFloat64(0).SetPrec(BigDecimalPrecision)
	}
}

func (b BigDecimal) String() string {
	if b.Float == nil {
		return "nil"
	}
	return b.Text('f', -1)
}

func (b BigDecimal) Equals(other BigDecimal) bool {
	if b.Float == nil || other.Float == nil {
		return false
	}
	return b.Cmp(other.Float) == 0
}

// driver.Valuer
func (b BigDecimal) Value() (driver.Value, error) {
	if b.Float == nil {
		return nil, nil
	}
	return b.Text('f', -1), nil
}

// sql.Scanner
func (b *BigDecimal) Scan(value interface{}) error {
	if value == nil {
		b.Float = nil
		return nil
	}
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("unsupported type for BigDecimal: %T", value)
	}
	f, _, err := big.ParseFloat(s, 10, BigDecimalPrecision, big.ToNearestEven)
	if err != nil {
		return fmt.Errorf("failed to parse big.Float: %w", err)
	}
	b.Float = f.SetPrec(BigDecimalPrecision)
	return nil
}

// JSON Marshal/Unmarshal
func (b BigDecimal) MarshalJSON() ([]byte, error) {
	if b.Float == nil {
		return []byte("null"), nil
	}
	return json.Marshal(b.Text('f', -1))
}

func (b *BigDecimal) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return b.Scan(s)
	}
	var f64 float64
	if err := json.Unmarshal(data, &f64); err == nil {
		b.Float = big.NewFloat(f64).SetPrec(BigDecimalPrecision)
		return nil
	}
	return errors.New("invalid JSON format for BigDecimal")
}

// Math
func (b *BigDecimal) Plus(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Add(b.Float, other.Float)}
}

func (b *BigDecimal) Sub(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil {
		return &BigDecimal{Float: nil}
	}
	if other == nil || other.Float == nil {
		return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Set(b.Float)}
	}
	return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Sub(b.Float, other.Float)}
}

func (b *BigDecimal) Mul(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Mul(b.Float, other.Float)}
}

func (b *BigDecimal) SafeDiv(other *BigDecimal) *BigDecimal {
	if b == nil || b.Float == nil || other == nil || other.Float == nil || other.Float.Sign() == 0 {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).Quo(b.Float, other.Float)}
}
