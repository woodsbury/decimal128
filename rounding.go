package decimal128

import "fmt"

// Round rounds (or quantises) a Decimal value to the specified number of
// decimal places using the rounding mode provided.
//
// The value of dp affects how many digits after the decimal point the Decimal
// would have if it were printed in decimal notation (for example, by the '%f'
// verb in Format). It can be zero to round off all digits after the decimal
// point and return an integer, and can also be negative to round off digits
// before the decimal point.
//
// NaN and infinity values are left untouched.
func (d Decimal) Round(dp int, mode RoundingMode) Decimal {
	if d.isSpecial() {
		return d
	}

	sig, exp := d.decompose()

	if sig == (uint128{}) {
		return zero(d.isNeg())
	}

	dp = dp*-1 + exponentBias
	iexp := int(exp)

	if iexp >= dp {
		return d
	}

	if iexp < dp-maxDigits {
		return zero(d.isNeg())
	}

	var trunc int8
	var digit uint64

	for iexp < dp {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()

		if sig == (uint128{}) && digit == 0 {
			return zero(d.isNeg())
		}

		iexp++
	}

	neg := d.isNeg()
	sig, exp = mode.round(neg, sig, int16(iexp), trunc, digit)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// RoundingMode determines how a Decimal value is rounded when the result of an
// operation is greater than the format can hold.
type RoundingMode uint8

const (
	ToNearestEven RoundingMode = iota // == IEEE 754 roundTiesToEven
	ToNearestAway                     // == IEEE 754 roundTiesToAway
	ToZero                            // == IEEE 754 roundTowardZero
	AwayFromZero                      // no IEEE 754 equivalent
	ToNegativeInf                     // == IEEE 754 roundTowardNegative
	ToPositiveInf                     // == IEEE 754 roundTowardPostive
)

// String returns a string representation of the rounding mode.
func (rm RoundingMode) String() string {
	switch rm {
	case ToNearestEven:
		return "ToNearestEven"
	case ToNearestAway:
		return "ToNearestAway"
	case ToZero:
		return "ToZero"
	case AwayFromZero:
		return "AwayFromZero"
	case ToNegativeInf:
		return "ToNegativeInf"
	case ToPositiveInf:
		return "ToPositiveInf"
	default:
		return fmt.Sprintf("RoundingMode(%d)", uint8(rm))
	}
}

func (rm RoundingMode) reduce256(neg bool, sig256 uint256, exp int16) (uint128, int16) {
	var trunc int8

	for sig256[3] > 0 {
		var rem uint64
		sig256, rem = sig256.div1e19()
		exp += 19

		if rem != 0 {
			trunc = 1
		}
	}

	sig192 := uint192{sig256[0], sig256[1], sig256[2]}

	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = 1
		}
	}

	sig := uint128{sig192[0], sig192[1]}

	var digit uint64

	for sig[1] > 0x0002_7fff_ffff_ffff {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()
		exp++
	}

	for exp < minBiasedExponent {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()

		if sig == (uint128{}) && digit == 0 {
			trunc = 0
			digit = 0
			exp = 0
			break
		}

		exp++
	}

	for exp > maxBiasedExponent && sig[1] < 0x0002_7fff_ffff_ffff {
		tmp := sig.mul64(10)

		if tmp[1] <= 0x0002_7fff_ffff_ffff {
			sig = tmp
			exp--
		} else {
			break
		}
	}

	return rm.round(neg, sig, exp, trunc, digit)
}

func (rm RoundingMode) reduce192(neg bool, sig192 uint192, exp int16, trunc int8) (uint128, int16) {
	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = 1
		}
	}

	sig := uint128{sig192[0], sig192[1]}

	var digit uint64

	for sig[1] > 0x0002_7fff_ffff_ffff {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()
		exp++
	}

	for exp < minBiasedExponent {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()

		if sig == (uint128{}) && digit == 0 {
			trunc = 0
			digit = 0
			exp = 0
			break
		}

		exp++
	}

	for exp > maxBiasedExponent && sig[1] < 0x0002_7fff_ffff_ffff {
		tmp := sig.mul64(10)

		if tmp[1] <= 0x0002_7fff_ffff_ffff {
			sig = tmp
			exp--
		} else {
			break
		}
	}

	return rm.round(neg, sig, exp, trunc, digit)
}

func (rm RoundingMode) reduce128(neg bool, sig uint128, exp int16, trunc int8) (uint128, int16) {
	var digit uint64

	for sig[1] > 0x0002_7fff_ffff_ffff {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()
		exp++
	}

	for exp < minBiasedExponent {
		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()

		if sig == (uint128{}) && digit == 0 {
			trunc = 0
			digit = 0
			exp = 0
			break
		}

		exp++
	}

	for exp > maxBiasedExponent && sig[1] < 0x0002_7fff_ffff_ffff {
		tmp := sig.mul64(10)

		if tmp[1] <= 0x0002_7fff_ffff_ffff {
			sig = tmp
			exp--
		} else {
			break
		}
	}

	return rm.round(neg, sig, exp, trunc, digit)
}

func (rm RoundingMode) reduce64(neg bool, sig64 uint64, exp int16) (uint128, int16) {
	var trunc int8
	var digit uint64

	for exp < minBiasedExponent {
		if digit != 0 {
			trunc = 1
		}

		digit = sig64 % 10
		sig64 = sig64 / 10

		if sig64 == 0 && digit == 0 {
			trunc = 0
			digit = 0
			exp = 0
			break
		}

		exp++
	}

	sig := uint128{sig64, 0}

	for exp > maxBiasedExponent && sig[1] < 0x0002_7fff_ffff_ffff {
		tmp := sig.mul64(10)

		if tmp[1] <= 0x0002_7fff_ffff_ffff {
			sig = tmp
			exp--
		} else {
			break
		}
	}

	return rm.round(neg, sig, exp, trunc, digit)
}

func (rm RoundingMode) round(neg bool, sig uint128, exp int16, trunc int8, digit uint64) (uint128, int16) {
	for {
		var adjust int
		switch rm {
		case ToNearestEven:
			if trunc == 1 {
				if digit >= 5 {
					adjust = 1
				}
			} else if trunc == -1 {
				if digit > 5 {
					adjust = 1
				}
			} else {
				if digit > 5 {
					adjust = 1
				} else if digit == 5 {
					if sig[0]%2 != 0 {
						adjust = 1
					}
				}
			}
		case ToNearestAway:
			if digit >= 5 {
				adjust = 1
			}
		case ToZero:
			if trunc == -1 && digit == 0 {
				adjust = -1
			}
		case AwayFromZero:
			if trunc == 1 || digit != 0 {
				adjust = 1
			}
		case ToPositiveInf:
			if neg {
				if trunc == -1 && digit == 0 {
					adjust = -1
				}
			} else if trunc == 1 || digit != 0 {
				adjust = 1
			}
		case ToNegativeInf:
			if neg {
				if trunc == 1 || digit != 0 {
					adjust = 1
				}
			} else if trunc == -1 && digit == 0 {
				adjust = -1
			}
		}

		if adjust != 0 {
			var tsig uint128
			if adjust == 1 {
				tsig = sig.add64(1)
			} else {
				tsig = sig.sub64(1)
			}

			if tsig[1] > 0x0002_7fff_ffff_ffff {
				if digit != 0 {
					trunc = 1
				}

				sig, digit = sig.div10()
				exp++
				continue
			}

			sig = tsig
		}

		return sig, exp
	}
}

// DefaultRoundingMode is the rounding mode used by any methods where an
// alternate rounding mode isn't provided.
var DefaultRoundingMode RoundingMode = ToNearestEven
