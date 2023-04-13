package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func TestDecimalAdd(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v + %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			sum := lhs.AddWithMode(rhs, mode)

			if !res.equal(sum, mode) {
				t.Errorf("%v.AddWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, res.result(mode))
			}
		}
	}
}

func TestDecimalMul(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v * %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			prd := lhs.MulWithMode(rhs, mode)

			if !res.equal(prd, mode) {
				t.Errorf("%v.MulWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, prd, res.result(mode))
			}
		}
	}
}

func TestDecimalQuo(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v / %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			quo := lhs.QuoWithMode(rhs, mode)

			if !res.equal(quo, mode) {
				t.Errorf("%v.QuoWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, quo, res.result(mode))
			}
		}
	}
}

func TestDecimalQuoRem(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			biglhs := new(apd.Decimal)
			bigrhs := new(apd.Decimal)
			bigquo := new(apd.Decimal)
			bigrem := new(apd.Decimal)
			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    roundingModeToBig(mode),
			}

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					if dexp := lhs.exp - rhs.exp; dexp < -128 || dexp > 128 {
						// apd is very slow at finding the integer quotient or
						// remainder of two values when their exponents differ
						// by too much, skip these for now.
						continue
					}

					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					quo, rem := declhs.QuoRemWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)
					bigctx.Precision = 12325
					bigctx.QuoInteger(bigquo, biglhs, bigrhs)
					bigctx.Rem(bigrem, biglhs, bigrhs)
					bigctx.Precision = 38
					bigctx.Round(bigquo, bigquo)
					bigctx.Round(bigrem, bigrem)

					if !decimalsEqual(quo, bigquo, bigctx.Rounding) || !decimalsEqual(rem, bigrem, bigctx.Rounding) {
						t.Errorf("%v.QuoRemWithMode(%v, %v) = (%v, %v), want (%v, %v)", lhs, rhs, mode, quo, rem, bigquo, bigrem)
					}
				}
			}
		})
	}
}

func TestDecimalSub(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v - %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			dif := lhs.SubWithMode(rhs, mode)

			if !res.equal(dif, mode) {
				t.Errorf("%v.SubWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, dif, res.result(mode))
			}
		}
	}
}

func BenchmarkOperations(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, 0, len(decimalValues))
	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		if val.sig == (uint128{}) {
			continue
		}

		values = append(values, val.Decimal())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, lhs := range values {
			for _, rhs := range values {
				lhs.Add(rhs)
				lhs.Mul(rhs)
				lhs.Quo(rhs)
				lhs.Sub(rhs)
			}
		}
	}
}
