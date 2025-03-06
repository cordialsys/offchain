package offchain_test

import (
	"testing"

	oc "github.com/cordialsys/offchain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestNewAmountBlockchainFromUint64(t *testing.T) {
	require := require.New(t)
	amount := oc.NewAmountBlockchainFromUint64(123)
	require.NotNil(amount)
	require.Equal(amount.Uint64(), uint64(123))
	require.Equal(amount.String(), "123")
}

func TestNewAmountBlockchainFromFloat64(t *testing.T) {
	require := require.New(t)
	amount := oc.NewAmountBlockchainToMaskFloat64(1.23)
	require.NotNil(amount)
	require.Equal(amount.Uint64(), uint64(1230000))
	require.Equal(amount.String(), "1230000")

	amountFloat := amount.UnmaskFloat64()
	require.Equal(amountFloat, 1.23)
}

func TestAmountHumanReadable(t *testing.T) {
	require := require.New(t)
	amountDec, _ := decimal.NewFromString("10.3")
	amount := oc.Amount(amountDec)
	require.NotNil(amount)
	require.Equal(amount.String(), "10.3")
}

func TestNewAmountHumanReadableFromStr(t *testing.T) {
	require := require.New(t)
	amount, err := oc.NewAmountFromString("10.3")
	require.NoError(err)
	require.NotNil(amount)
	require.Equal(amount.String(), "10.3")

	amount, err = oc.NewAmountFromString("0")
	require.NoError(err)
	require.NotNil(amount)
	require.Equal(amount.String(), "0")

	amount, err = oc.NewAmountFromString("")
	require.Error(err)
	require.NotNil(amount)
	require.Equal(amount.String(), "0")

	amount, err = oc.NewAmountFromString("invalid")
	require.Error(err)
	require.NotNil(amount)
	require.Equal(amount.String(), "0")
}

func TestNewBlockchainAmountStr(t *testing.T) {
	require := require.New(t)
	amount := oc.NewAmountBlockchainFromStr("10")
	require.EqualValues(amount.Uint64(), 10)

	amount = oc.NewAmountBlockchainFromStr("10.1")
	require.EqualValues(amount.Uint64(), 0)

	amount = oc.NewAmountBlockchainFromStr("0x10")
	require.EqualValues(amount.Uint64(), 16)
}

func TestAdd(t *testing.T) {
	require := require.New(t)
	sum := oc.NewAmountBlockchainFromStr("100000000000000000000")
	toAdd := oc.NewAmountBlockchainFromStr("100")

	iterations := 100
	for i := 0; i < iterations; i++ {
		sum = sum.Add(&toAdd)
		require.Equal(toAdd.String(), "100")
	}

	require.Equal(sum.String(), "100000000000000010000")

	// switch
	sum = oc.NewAmountBlockchainFromStr("100")
	toAdd = oc.NewAmountBlockchainFromStr("100000000000000000000")

	for i := 0; i < iterations; i++ {
		sum = sum.Add(&toAdd)
		require.Equal(toAdd.String(), "100000000000000000000")
	}

	require.Equal(sum.String(), "10000000000000000000100")
}

func TestSub(t *testing.T) {
	require := require.New(t)
	sum := oc.NewAmountBlockchainFromUint64(10000000)
	toDiff := oc.NewAmountBlockchainFromUint64(100)

	iterations := 100
	for i := 0; i < iterations; i++ {
		sum = sum.Sub(&toDiff)
	}

	require.EqualValues(sum.Uint64(), 10000000-int(toDiff.Uint64())*iterations)
	require.EqualValues(toDiff.Uint64(), 100)
}

func TestMult(t *testing.T) {
	require := require.New(t)
	sum := oc.NewAmountBlockchainFromUint64(1)
	toMul := oc.NewAmountBlockchainFromUint64(100)

	iterations := 3
	for i := 0; i < iterations; i++ {
		sum = sum.Mul(&toMul)
	}

	require.EqualValues(sum.Uint64(), 1000000)
	require.EqualValues(toMul.Uint64(), 100)
}

func TestDiv(t *testing.T) {
	require := require.New(t)
	quot := oc.NewAmountBlockchainFromUint64(1000000000)
	toDiv := oc.NewAmountBlockchainFromUint64(100)

	iterations := 3
	for i := 0; i < iterations; i++ {
		quot = quot.Div(&toDiv)
	}

	require.EqualValues(quot.Uint64(), 1000)
	require.EqualValues(toDiv.Uint64(), 100)
}

func TestCombined(t *testing.T) {
	// Here we chain a bunch of operations together that span many magnitudes
	// to trigger undlying slice reallocation of big integers.  The expectation
	// is that big integers are properly cloned, and do not mutate each other.
	require := require.New(t)

	a := oc.NewAmountBlockchainFromStr("10")
	b := oc.NewAmountBlockchainFromStr("10000")
	c := oc.NewAmountBlockchainFromStr("1000000000")
	d := oc.NewAmountBlockchainFromStr("1000000000000000")
	e := oc.NewAmountBlockchainFromStr("5000000000000000000000000")
	f := oc.NewAmountBlockchainFromStr("100")

	// (a+b) = 10010
	a_plus_b := a.Add(&b)

	// c - (a + b) = 999989990
	c_minus_a_plus_b := c.Sub(&a_plus_b)

	// d * (c - (a + b)) = 999989990000000000000000
	d_times_c_minus_a_plus_b := d.Mul(&c_minus_a_plus_b)

	// e - d * (c - (a + b)) = 4000010010000000000000000
	e_minus_d_times_c_minus_a_plus_b := e.Sub(&d_times_c_minus_a_plus_b)

	// (e - d * (c - (a + b)))//f = 40000100100000000000000
	e_minus_d_times_c_minus_a_plus_b_div_f := e_minus_d_times_c_minus_a_plus_b.Div(&f)

	// (e - d * (c - (a + b)))  +  (e - d * (c - (a + b)))//f = 4040010110100000000000000
	sum_last_two := e_minus_d_times_c_minus_a_plus_b.Add(&e_minus_d_times_c_minus_a_plus_b_div_f)

	// Confirm all variables are not mutated
	require.Equal(a.String(), "10")
	require.Equal(b.String(), "10000")
	require.Equal(c.String(), "1000000000")
	require.Equal(d.String(), "1000000000000000")
	require.Equal(e.String(), "5000000000000000000000000")

	// All results + intermediates are correct
	require.Equal(a_plus_b.String(), "10010")
	require.Equal(c_minus_a_plus_b.String(), "999989990")
	require.Equal(d_times_c_minus_a_plus_b.String(), "999989990000000000000000")
	require.Equal(e_minus_d_times_c_minus_a_plus_b.String(), "4000010010000000000000000")
	require.Equal(e_minus_d_times_c_minus_a_plus_b_div_f.String(), "40000100100000000000000")
	require.Equal(sum_last_two.String(), "4040010110100000000000000")
}
