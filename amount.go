package offchain

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
)

const FLOAT_PRECISION = 6

// Balance is a big integer amount as blockchain expects it for tx.
type Balance big.Int

// Amount is a decimal amount as a human expects it for readability.
type Amount decimal.Decimal

func (amount Balance) Bytes() []byte {
	bigInt := big.Int(amount)
	return bigInt.Bytes()
}

func (amount Balance) String() string {
	bigInt := big.Int(amount)
	return bigInt.String()
}

// Int converts an AmountBlockchain into *bit.Int
func (amount Balance) Int() *big.Int {
	bigInt := big.Int(amount)
	return &bigInt
}

func (amount Balance) Sign() int {
	bigInt := big.Int(amount)
	return bigInt.Sign()
}

// Uint64 converts an AmountBlockchain into uint64
func (amount Balance) Uint64() uint64 {
	bigInt := big.Int(amount)
	return bigInt.Uint64()
}

// UnmaskFloat64 converts an AmountBlockchain into float64 given the number of decimals
func (amount Balance) UnmaskFloat64() float64 {
	bigInt := big.Int(amount)
	bigFloat := new(big.Float).SetInt(&bigInt)
	exponent := new(big.Float).SetFloat64(math.Pow10(FLOAT_PRECISION))
	bigFloat = bigFloat.Quo(bigFloat, exponent)
	f64, _ := bigFloat.Float64()
	return f64
}

// Use the underlying big.Int.Cmp()
func (amount *Balance) Cmp(other *Balance) int {
	return amount.Int().Cmp(other.Int())
}

// Use the underlying big.Int.Add()
func (amount *Balance) Add(x *Balance) Balance {
	sum := new(big.Int)
	sum.Set((*big.Int)(amount))
	return Balance(*sum.Add(sum, x.Int()))
}

// Use the underlying big.Int.Sub()
func (amount *Balance) Sub(x *Balance) Balance {
	diff := new(big.Int)
	diff.Set((*big.Int)(amount))
	return Balance(*diff.Sub(diff, x.Int()))
}

// Use the underlying big.Int.Mul()
func (amount *Balance) Mul(x *Balance) Balance {
	prod := new(big.Int)
	prod.Set((*big.Int)(amount))
	return Balance(*prod.Mul(prod, x.Int()))
}

// Use the underlying big.Int.Div()
func (amount *Balance) Div(x *Balance) Balance {
	quot := new(big.Int)
	quot.Set((*big.Int)(amount))
	return Balance(*quot.Div(quot, x.Int()))
}

func (amount *Balance) Abs() Balance {
	abs := new(big.Int)
	abs.Set((*big.Int)(amount))
	return Balance(*abs.Abs(abs))
}

var zero = big.NewInt(0)

func (amount *Balance) IsZero() bool {
	return amount.Int().Cmp(zero) == 0
}

func (amount *Balance) ToHuman(decimals int32) Amount {
	dec := decimal.NewFromBigInt(amount.Int(), -decimals)
	return Amount(dec)
}

func MultiplyByFloat(amount Balance, multiplier float64) Balance {
	if amount.Uint64() == 0 {
		return amount
	}
	// We are computing (100000 * multiplier * amount) / 100000
	precision := uint64(1000000)
	multBig := NewAmountBlockchainFromUint64(uint64(float64(precision) * multiplier))
	divBig := NewAmountBlockchainFromUint64(precision)
	product := multBig.Mul(&amount)
	result := product.Div(&divBig)
	return result
}

// NewAmountBlockchainFromUint64 creates a new AmountBlockchain from a uint64
func NewAmountBlockchainFromInt64(i64 int64) (Balance, bool) {
	if i64 < 0 {
		return NewAmountBlockchainFromUint64(0), false
	}
	return NewAmountBlockchainFromUint64(uint64(i64)), true
}

// NewAmountBlockchainFromUint64 creates a new AmountBlockchain from a uint64
func NewAmountBlockchainFromUint64(u64 uint64) Balance {
	bigInt := new(big.Int).SetUint64(u64)
	return Balance(*bigInt)
}

// NewAmountBlockchainToMaskFloat64 creates a new AmountBlockchain as a float64 times 10^FLOAT_PRECISION
func NewAmountBlockchainToMaskFloat64(f64 float64) Balance {
	bigFloat := new(big.Float).SetFloat64(f64)
	exponent := new(big.Float).SetFloat64(math.Pow10(FLOAT_PRECISION))
	bigFloat = bigFloat.Mul(bigFloat, exponent)
	var bigInt big.Int
	bigFloat.Int(&bigInt)
	return Balance(bigInt)
}

// NewAmountBlockchainFromStr creates a new AmountBlockchain from a string
func NewAmountBlockchainFromStr(str string) Balance {
	var ok bool
	var bigInt *big.Int
	bigInt, ok = new(big.Int).SetString(str, 0)
	if !ok {
		return NewAmountBlockchainFromUint64(0)
	}
	return Balance(*bigInt)
}

// NewAmountFromString creates a new AmountHumanReadable from a string
func NewAmountFromString(str string) (Amount, error) {
	decimal, err := decimal.NewFromString(str)
	return Amount(decimal), err
}

// NewAmountHumanReadableFromFloat creates a new AmountHumanReadable from a float
func NewAmountHumanReadableFromFloat(float float64) Amount {
	return Amount(decimal.NewFromFloat(float))
}

func (amount Amount) Decimal() decimal.Decimal {
	return decimal.Decimal(amount)
}

func (amount Amount) ToBlockchain(decimals int32) Balance {
	factor := decimal.NewFromInt32(10).Pow(decimal.NewFromInt32(decimals))
	raised := ((decimal.Decimal)(amount)).Mul(factor)
	return Balance(*raised.BigInt())
}

func (amount Amount) String() string {
	return decimal.Decimal(amount).String()
}

func (amount Amount) Div(x Amount) Amount {
	return Amount(decimal.Decimal(amount).Div(decimal.Decimal(x)))
}

var _ json.Marshaler = Amount{}
var _ json.Unmarshaler = &Amount{}
var _ yaml.Unmarshaler = &Amount{}
var _ yaml.Marshaler = Amount{}
var _ yaml.IsZeroer = Amount{}

func (b Amount) MarshalYAML() (interface{}, error) {
	return b.String(), nil
}

func (b Amount) IsZero() bool {
	return decimal.Decimal(b).IsZero()
}

func (b *Amount) UnmarshalYAML(node *yaml.Node) error {
	value := node.Value
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "\"")
	value = strings.TrimSuffix(value, "\"")
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return fmt.Errorf("invalid decimal amount: %v", err)
	}
	*b = Amount(dec)
	return nil
}

func (b Amount) MarshalJSON() ([]byte, error) {
	return []byte("\"" + b.String() + "\""), nil
}

func (b *Amount) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	str := strings.Trim(string(p), "\"")
	decimal, err := decimal.NewFromString(str)
	if err != nil {
		return err
	}
	*b = Amount(decimal)
	return nil
}

var _ json.Marshaler = Balance{}
var _ json.Unmarshaler = &Balance{}

func (b Balance) MarshalJSON() ([]byte, error) {
	return []byte("\"" + b.String() + "\""), nil
}

func (b *Balance) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	str := strings.Trim(string(p), "\"")
	var z big.Int
	_, ok := z.SetString(str, 0)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	*b = Balance(z)
	return nil
}

type StringOrInt string

func (s StringOrInt) AsString() string {
	return string(s)
}

func (s StringOrInt) AsInt() (uint64, bool) {
	asInt, err := strconv.ParseUint(string(s), 0, 64)
	return asInt, err == nil
}
