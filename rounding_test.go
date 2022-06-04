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
							t.Errorf("%v.Round(%d) = %v, want %v", val, dp, res, decval)
						}

						continue
					}

					if dp*-1 < int(val.exp-exponentBias) {
						if !decval.Equal(res) {
							t.Errorf("%v.Round(%d) = %v, want %v", val, dp, res, decval)
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
						t.Errorf("%v.Round(%d) = %v, want %v", val, dp, res, bigval)
					}
				}
			}
		})
	}
}
