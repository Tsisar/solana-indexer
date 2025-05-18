package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	*big.Int
}

// GormDataType informs GORM of the data type
func (BigInt) GormDataType() string {
	return "numeric"
}

func (b *BigInt) Value() (driver.Value, error) {
	if b == nil || b.Int == nil {
		return nil, nil
	}
	return b.String(), nil
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b.Int = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return b.setString(v)
	case []byte:
		return b.setString(string(v))
	default:
		return fmt.Errorf("unsupported type for BigInt: %T", value)
	}
}

func (b *BigInt) setString(s string) error {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse types.BigInt from string: %s", s)
	}
	b.Int = i
	return nil
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", b.String())), nil
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return b.setString(s)
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		b.Int = big.NewInt(i)
		return nil
	}

	return fmt.Errorf("failed to unmarshal BigInt from JSON: %s", string(data))
}

// ToBigDecimal converts BigInt to BigDecimal (preserving integer value as float).
func (b *BigInt) ToBigDecimal() *BigDecimal {
	if b == nil || b.Int == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).SetInt(b.Int)}
}

func (b *BigInt) Sub(other *BigInt) *BigInt {
	if b == nil || b.Int == nil || other == nil || other.Int == nil {
		return &BigInt{Int: nil}
	}
	return &BigInt{Int: new(big.Int).Sub(b.Int, other.Int)}
}

// Plus returns b + other as a new BigInt.
// Returns nil.Int if either operand is nil.
func (b *BigInt) Plus(other *BigInt) *BigInt {
	if b == nil || b.Int == nil || other == nil || other.Int == nil {
		return &BigInt{Int: nil}
	}
	return &BigInt{
		Int: new(big.Int).Add(b.Int, other.Int),
	}
}

func BigIntFromString(s string) (*BigInt, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid BigInt string: %s", s)
	}
	return &BigInt{Int: i}, nil
}

func BigIntFromUint64(v uint64) *BigInt {
	return &BigInt{Int: new(big.Int).SetUint64(v)}
}

func BigIntFromInt64(v int64) *BigInt {
	return &BigInt{Int: big.NewInt(v)}
}
