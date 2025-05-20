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

// GormDataType tells GORM to use numeric type
func (BigInt) GormDataType() string {
	return "numeric"
}

// Value returns a string representation for storing in DB
func (b *BigInt) Value() (driver.Value, error) {
	if b == nil || b.Int == nil {
		return nil, nil
	}
	return b.String(), nil
}

// Scan reads the DB value into BigInt
func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b.Int = nil
		return nil
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("unsupported type for BigInt: %T", value)
	}
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse BigInt from string: %s", s)
	}
	b.Int = i
	return nil
}

// MarshalJSON serializes as quoted string
func (b BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", b.String())), nil
}

// UnmarshalJSON parses from quoted string or number
func (b *BigInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return b.Scan(s)
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		b.Int = big.NewInt(i)
		return nil
	}

	return fmt.Errorf("failed to unmarshal BigInt from JSON: %s", string(data))
}

// String returns string representation
func (b BigInt) String() string {
	if b.Int == nil {
		return "nil"
	}
	return b.Int.String()
}

// Equals returns true if both BigInts are equal
func (b BigInt) Equals(other BigInt) bool {
	if b.Int == nil || other.Int == nil {
		return false
	}
	return b.Cmp(other.Int) == 0
}

// Math

func (b *BigInt) Plus(other *BigInt) *BigInt {
	if b == nil || b.Int == nil || other == nil || other.Int == nil {
		return &BigInt{Int: nil}
	}
	return &BigInt{Int: new(big.Int).Add(b.Int, other.Int)}
}

func (b *BigInt) Sub(other *BigInt) *BigInt {
	if b == nil || b.Int == nil || other == nil || other.Int == nil {
		return &BigInt{Int: nil}
	}
	return &BigInt{Int: new(big.Int).Sub(b.Int, other.Int)}
}

func (b *BigInt) Mul(other *BigInt) *BigInt {
	if b == nil || b.Int == nil || other == nil || other.Int == nil {
		return &BigInt{Int: nil}
	}
	return &BigInt{Int: new(big.Int).Mul(b.Int, other.Int)}
}

// Conversion

func (b *BigInt) ToBigDecimal() *BigDecimal {
	if b == nil || b.Int == nil {
		return &BigDecimal{Float: nil}
	}
	return &BigDecimal{Float: new(big.Float).SetPrec(BigDecimalPrecision).SetInt(b.Int)}
}

// Zero sets value to 0
func (b *BigInt) Zero() {
	if b.Int == nil {
		b.Int = big.NewInt(0)
	} else {
		b.Int.SetInt64(0)
	}
}

// Constructors

func ZeroBigInt() BigInt {
	return BigInt{Int: big.NewInt(0)}
}

func BigIntFromString(s string) (*BigInt, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid BigInt string: %s", s)
	}
	return &BigInt{Int: i}, nil
}

func MustBigIntFromString(s string) BigInt {
	b, err := BigIntFromString(s)
	if err != nil {
		panic(err)
	}
	return *b
}

func BigIntFromUint64(v uint64) *BigInt {
	return &BigInt{Int: new(big.Int).SetUint64(v)}
}

func BigIntFromInt64(v int64) *BigInt {
	return &BigInt{Int: big.NewInt(v)}
}
