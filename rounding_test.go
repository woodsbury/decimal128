package decimal128

import "github.com/cockroachdb/apd/v3"

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
