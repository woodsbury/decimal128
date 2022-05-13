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
					sum := declhs.AddWithRounding(decrhs, mode)

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
						t.Errorf("%v.AddWithRounding(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
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
					prd := declhs.MulWithRounding(decrhs, mode)

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
						t.Errorf("%v.MulWithRounding(%v, %v) = %v, want %v", lhs, rhs, mode, prd, biglhs)
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
					sum := declhs.QuoWithRounding(decrhs, mode)

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
						t.Errorf("%v.QuoWithRounding(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
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
					sum := declhs.SubWithRounding(decrhs, mode)

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
						t.Errorf("%v.SubWithRounding(%v, %v) = %v, want %v", lhs, rhs, mode, sum, biglhs)
					}
				}
			}
		})
	}
}

func TestQuoE(t *testing.T) {
	lhs := testDec{regularForm, true, uint128{0x0001_7fff_ffff_ffff, 0x0002_7fff_ffff_ffff}, exponentBias - 6176}
	rhs := testDec{regularForm, false, uint128{0x0002_7fff_ffff_ffff, 0x0002_7fff_ffff_ffff}, exponentBias - 19}

	declhs := lhs.Decimal()
	decrhs := rhs.Decimal()
	declhs = declhs.Quo(decrhs)

	biglhs := lhs.Big(new(apd.Decimal))
	bigrhs := rhs.Big(new(apd.Decimal))

	t.Error(lhs, "+", rhs, "=")
	t.Error(declhs)

	bigctx := apd.Context{
		Precision:   38,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    roundingModeToBig(DefaultRoundingMode),
	}

	bigctx.Quo(biglhs, biglhs, bigrhs)
	t.Error(biglhs)

	t.Error(decimalsEqual(declhs, biglhs, roundingModeToBig(DefaultRoundingMode)))
}
