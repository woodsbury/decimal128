package decimal128

import (
	"math"
	"math/big"
	"testing"
)

var floatValues = map[Decimal]float64{
	zero(false): 0.0,
	zero(true):  math.Copysign(0.0, -1.0),
	compose(false, uint128{1, 0}, exponentBias): 1.0,
	compose(true, uint128{1, 0}, exponentBias):  -1.0,
	compose(false, uint128{0, 1}, exponentBias): float64(uint64(1<<63)) * 2.0,
	compose(true, uint128{0, 1}, exponentBias):  math.Copysign(float64(uint64(1<<63))*2.0, -1.0),
	inf(false): math.Inf(1),
	inf(true):  math.Inf(-1),
	NaN():      math.NaN(),
}

func TestDecimalFloat(t *testing.T) {
	t.Parallel()

	bignum := new(big.Float)

	for val, num := range floatValues {
		if val.IsNaN() {
			continue
		}

		res := val.Float()

		bignum.SetFloat64(num)

		if res.Cmp(bignum) != 0 || res.Signbit() != bignum.Signbit() {
			t.Errorf("%v.Float() = %v, want %v", val, res, bignum)
		}
	}
}

func TestDecimalFloat64(t *testing.T) {
	t.Parallel()

	for val, num := range floatValues {
		res := val.Float64()

		if !(res == num || math.IsNaN(res) && math.IsNaN(num)) || math.Signbit(res) != math.Signbit(num) {
			t.Errorf("%v.Float64() = %v, want %v", val, res, num)
		}
	}
}
