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
			bigmode := roundingModeToBig(mode)

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					sum := declhs.AddWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					bigctx := apd.Context{
						Precision:   38,
						MaxExponent: 6145,
						MinExponent: -6176,
						Rounding:    bigmode,
					}

					bigctx.Add(biglhs, biglhs, bigrhs)

					if !decimalsEqual(sum, biglhs, bigmode) {
						t.Errorf("%v.AddWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
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
			bigmode := roundingModeToBig(mode)

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					prd := declhs.MulWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					bigctx := apd.Context{
						Precision:   38,
						MaxExponent: 6145,
						MinExponent: -6176,
						Rounding:    bigmode,
					}

					bigctx.Mul(biglhs, biglhs, bigrhs)

					if !decimalsEqual(prd, biglhs, bigmode) {
						t.Errorf("%v.MulWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, prd, biglhs)
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
			bigmode := roundingModeToBig(mode)

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					sum := declhs.QuoWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					bigctx := apd.Context{
						Precision:   38,
						MaxExponent: 6145,
						MinExponent: -6176,
						Rounding:    bigmode,
					}

					bigctx.Quo(biglhs, biglhs, bigrhs)

					if !decimalsEqual(sum, biglhs, bigmode) {
						t.Errorf("%v.QuoWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
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
			bigmode := roundingModeToBig(mode)

			for _, lhs := range decimalValues {
				for _, rhs := range decimalValues {
					declhs := lhs.Decimal()
					decrhs := rhs.Decimal()
					sum := declhs.SubWithMode(decrhs, mode)

					lhs.Big(biglhs)
					rhs.Big(bigrhs)

					bigctx := apd.Context{
						Precision:   38,
						MaxExponent: 6145,
						MinExponent: -6176,
						Rounding:    bigmode,
					}

					bigctx.Sub(biglhs, biglhs, bigrhs)

					if !decimalsEqual(sum, biglhs, bigmode) {
						t.Errorf("%v.SubWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
					}
				}
			}
		})
	}
}
