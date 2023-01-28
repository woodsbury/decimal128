package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func TestLog(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   39,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Log(decval)

		val.Big(bigval)

		bigctx.Ln(bigres, bigval)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Log(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestLog10(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   39,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Log10(decval)

		val.Big(bigval)

		bigctx.Log10(bigres, bigval)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Log10(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestLog2(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   39,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	bigln2 := new(apd.Decimal)
	bigctx.Ln(bigln2, apd.New(2, 0))

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Log2(decval)

		val.Big(bigval)

		bigctx.Ln(bigres, bigval)
		bigctx.Quo(bigres, bigres, bigln2)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Log2(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestSqrt(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   39,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Sqrt(decval)

		val.Big(bigval)

		bigctx.Sqrt(bigres, bigval)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Sqrt(%v) = %v, want %v", val, res, bigres)
		}
	}
}
