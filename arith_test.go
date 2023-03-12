package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func TestDecimalAdd(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			biglhs := new(apd.Decimal)
			bigrhs := new(apd.Decimal)
			bigsum := new(apd.Decimal)
			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    roundingModeToBig(mode),
			}

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					sum := declhs.AddWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					// apd is very slow at adding or subtracting two values
					// when their exponents differ by too much, so help it out
					// by adjusting numbers that are very far apart.
					if biglhs.Form == apd.Finite && biglhs.Coeff.BitLen() != 0 && bigrhs.Form == apd.Finite && bigrhs.Coeff.BitLen() != 0 {
						lhsexp := biglhs.Exponent + int32(biglhs.NumDigits())
						rhsexp := bigrhs.Exponent + int32(bigrhs.NumDigits())
						dexp := lhsexp - rhsexp

						if dexp > 40 {
							bigrhs.Coeff.SetInt64(1)
							bigrhs.Exponent = lhsexp - 40
						} else if dexp < -40 {
							biglhs.Coeff.SetInt64(1)
							biglhs.Exponent = rhsexp - 40
						}
					}

					bigctx.Add(bigsum, biglhs, bigrhs)

					if !decimalsEqual(sum, bigsum, bigctx.Rounding) {
						t.Errorf("%v.AddWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, bigsum)
					}
				}
			}
		})
	}
}

func TestDecimalMul(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			biglhs := new(apd.Decimal)
			bigrhs := new(apd.Decimal)
			bigprd := new(apd.Decimal)
			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    roundingModeToBig(mode),
			}

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					prd := declhs.MulWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)
					bigctx.Mul(bigprd, biglhs, bigrhs)

					if !decimalsEqual(prd, bigprd, bigctx.Rounding) {
						t.Errorf("%v.MulWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, prd, bigprd)
					}
				}
			}
		})
	}
}

func TestDecimalQuo(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			biglhs := new(apd.Decimal)
			bigrhs := new(apd.Decimal)
			bigquo := new(apd.Decimal)
			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    roundingModeToBig(mode),
			}

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					quo := declhs.QuoWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)
					bigctx.Quo(bigquo, biglhs, bigrhs)

					if !decimalsEqual(quo, bigquo, bigctx.Rounding) {
						t.Errorf("%v.QuoWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, quo, bigquo)
					}
				}
			}
		})
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

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			biglhs := new(apd.Decimal)
			bigrhs := new(apd.Decimal)
			bigdif := new(apd.Decimal)
			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    roundingModeToBig(mode),
			}

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					dif := declhs.SubWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					// apd is very slow at adding or subtracting two values
					// when their exponents differ by too much, so help it out
					// by adjusting numbers that are very far apart.
					if biglhs.Form == apd.Finite && biglhs.Coeff.BitLen() != 0 && bigrhs.Form == apd.Finite && bigrhs.Coeff.BitLen() != 0 {
						lhsexp := biglhs.Exponent + int32(biglhs.NumDigits())
						rhsexp := bigrhs.Exponent + int32(bigrhs.NumDigits())
						dexp := lhsexp - rhsexp

						if dexp > 40 {
							bigrhs.Coeff.SetInt64(1)
							bigrhs.Exponent = lhsexp - 40
						} else if dexp < -40 {
							biglhs.Coeff.SetInt64(1)
							biglhs.Exponent = rhsexp - 40
						}
					}

					bigctx.Sub(bigdif, biglhs, bigrhs)

					if !decimalsEqual(dif, bigdif, bigctx.Rounding) {
						t.Errorf("%v.SubWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, dif, bigdif)
					}
				}
			}
		})
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
