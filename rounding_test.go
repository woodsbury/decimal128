package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

var roundingModes = []RoundingMode{
	ToNearestEven,
	ToNearestAway,
	ToZero,
	AwayFromZero,
	ToNegativeInf,
	ToPositiveInf,
}

func roundingModeToBig(mode RoundingMode) apd.Rounder {
	switch mode {
	case ToNearestEven:
		return apd.RoundHalfEven
	case ToNearestAway:
		return apd.RoundHalfUp
	case ToZero:
		return apd.RoundDown
	case AwayFromZero:
		return apd.RoundUp
	case ToNegativeInf:
		return apd.RoundFloor
	case ToPositiveInf:
		return apd.RoundCeiling
	default:
		panic("rounding mode not handled")
	}
}

func TestDecimalCeil(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)

	for dp := -2; dp <= 2; dp += 2 {
		for _, val := range decimalValues {
			decval := val.Decimal()
			res := decval.Ceil(dp)

			if decval.isSpecial() {
				if !(decval.Equal(res) || decval.IsNaN() && res.IsNaN()) {
					t.Errorf("%v.Ceil(%d) = %v, want %v", val, dp, res, decval)
				}

				continue
			}

			if dp*-1 < int(val.exp-exponentBias) {
				if !decval.Equal(res) {
					t.Errorf("%v.Ceil(%d) = %v, want %v", val, dp, res, decval)
				}

				continue
			}

			val.Big(bigval)

			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    apd.RoundCeiling,
			}

			bigval.Exponent += int32(dp)
			bigctx.Ceil(bigres, bigval)
			bigres.Exponent -= int32(dp)

			if !decimalsEqual(res, bigres, apd.RoundCeiling) {
				t.Errorf("%v.Ceil(%d) = %v, want %v", val, dp, res, bigres)
			}
		}
	}
}

func TestDecimalFloor(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)

	for dp := -2; dp <= 2; dp += 2 {
		for _, val := range decimalValues {
			decval := val.Decimal()
			res := decval.Floor(dp)

			if decval.isSpecial() {
				if !(decval.Equal(res) || decval.IsNaN() && res.IsNaN()) {
					t.Errorf("%v.Floor(%d) = %v, want %v", val, dp, res, decval)
				}

				continue
			}

			if dp*-1 < int(val.exp-exponentBias) {
				if !decval.Equal(res) {
					t.Errorf("%v.Floor(%d) = %v, want %v", val, dp, res, decval)
				}

				continue
			}

			val.Big(bigval)

			bigctx := apd.Context{
				Precision:   38,
				MaxExponent: 6145,
				MinExponent: -6176,
				Rounding:    apd.RoundFloor,
			}

			bigval.Exponent += int32(dp)
			bigctx.Floor(bigres, bigval)
			bigres.Exponent -= int32(dp)

			if !decimalsEqual(res, bigres, apd.RoundFloor) {
				t.Errorf("%v.Floor(%d) = %v, want %v", val, dp, res, bigres)
			}
		}
	}
}

func TestDecimalRound(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, mode := range roundingModes {
		mode := mode

		t.Run(mode.String(), func(t *testing.T) {
			t.Parallel()

			bigval := new(apd.Decimal)
			bigmode := roundingModeToBig(mode)

			for dp := -2; dp <= 2; dp += 2 {
				for _, val := range decimalValues {
					decval := val.Decimal()
					res := decval.Round(dp, mode)

					if decval.isSpecial() {
						if !(decval.Equal(res) || decval.IsNaN() && res.IsNaN()) {
							t.Errorf("%v.Round(%d, %v) = %v, want %v", val, dp, mode, res, decval)
						}

						continue
					}

					if dp*-1 < int(val.exp-exponentBias) {
						if !decval.Equal(res) {
							t.Errorf("%v.Round(%d, %v) = %v, want %v", val, dp, mode, res, decval)
						}

						continue
					}

					val.Big(bigval)

					bigctx := apd.Context{
						Precision:   38,
						MaxExponent: 6145,
						MinExponent: -6176,
						Rounding:    bigmode,
					}

					bigctx.Quantize(bigval, bigval, int32(dp*-1))

					if !decimalsEqual(res, bigval, bigmode) {
						t.Errorf("%v.Round(%d, %v) = %v, want %v", val, dp, mode, res, bigval)
					}
				}
			}
		})
	}
}
