package decimal128

import "fmt"

type RoundingMode uint8

const (
	ToNearestEven RoundingMode = iota
	ToNearestAway
	ToZero
	AwayFromZero
	ToNegativeInf
	ToPositiveInf
)

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
	var trunc bool

	for sig256[3] > 0 {
		var rem uint64
		sig256, rem = sig256.div1e19()
		exp += 19

		if rem != 0 {
			trunc = true
		}
	}

	sig192 := uint192{sig256[0], sig256[1], sig256[2]}

	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = true
		}
	}

	sig := uint128{sig192[0], sig192[1]}

	var digit uint64
	for sig[1] > 0x0002_7fff_ffff_ffff {
		if digit != 0 {
			trunc = true
		}

		sig, digit = sig.div10()
		exp++
	}

	for exp < minBiasedExponent {
		if sig == (uint128{}) {
			digit = 0
			exp = 0
			break
		}

		if digit != 0 {
			trunc = true
		}

		sig, digit = sig.div10()
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

	for {
		incr := false
		switch rm {
		case ToNearestEven:
			if trunc {
				if digit >= 5 {
					incr = true
				}
			} else {
				if digit > 5 {
					incr = true
				} else if digit == 5 {
					if sig[0]%2 != 0 {
						incr = true
					}
				}
			}
		case ToNearestAway:
			if digit >= 5 {
				incr = true
			}
		case AwayFromZero:
			if trunc || digit != 0 {
				incr = true
			}
		case ToPositiveInf:
			if (trunc || digit != 0) && !neg {
				incr = true
			}
		case ToNegativeInf:
			if (trunc || digit != 0) && neg {
				incr = true
			}
		}

		if incr {
			tsig := sig.add1()
			if tsig[1] > 0x0002_7fff_ffff_ffff {
				sig, digit = sig.div10()
				exp++
				trunc = true
				continue
			}

			sig = tsig
		}

		return sig, exp
	}
}

func (rm RoundingMode) reduce192(neg bool, sig192 uint192, exp int16, trunc bool) (uint128, int16) {
	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = true
		}
	}

	sig := uint128{sig192[0], sig192[1]}

	var digit uint64
	for sig[1] > 0x0002_7fff_ffff_ffff {
		if digit != 0 {
			trunc = true
		}

		sig, digit = sig.div10()
		exp++
	}

	for exp < minBiasedExponent {
		if sig == (uint128{}) {
			digit = 0
			exp = 0
			break
		}

		if digit != 0 {
			trunc = true
		}

		sig, digit = sig.div10()
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

	for {
		incr := false
		switch rm {
		case ToNearestEven:
			if trunc {
				if digit >= 5 {
					incr = true
				}
			} else {
				if digit > 5 {
					incr = true
				} else if digit == 5 {
					if sig[0]%2 != 0 {
						incr = true
					}
				}
			}
		case ToNearestAway:
			if digit >= 5 {
				incr = true
			}
		case AwayFromZero:
			if trunc || digit != 0 {
				incr = true
			}
		case ToPositiveInf:
			if (trunc || digit != 0) && !neg {
				incr = true
			}
		case ToNegativeInf:
			if (trunc || digit != 0) && neg {
				incr = true
			}
		}

		if incr {
			tsig := sig.add1()
			if tsig[1] > 0x0002_7fff_ffff_ffff {
				sig, digit = sig.div10()
				exp++
				trunc = true
				continue
			}

			sig = tsig
		}

		return sig, exp
	}
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
		if sig == (uint128{}) {
			digit = 0
			exp = 0
			break
		}

		if digit != 0 {
			trunc = 1
		}

		sig, digit = sig.div10()
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

	for {
		incr := false
		switch rm {
		case ToNearestEven:
			if trunc == 1 {
				if digit >= 5 {
					incr = true
				}
			} else if trunc == -1 {
				if digit > 5 {
					incr = true
				}
			} else {
				if digit > 5 {
					incr = true
				} else if digit == 5 {
					if sig[0]%2 != 0 {
						incr = true
					}
				}
			}
		case ToNearestAway:
			if digit >= 5 {
				incr = true
			}
		case AwayFromZero:
			if trunc == 1 || digit != 0 {
				incr = true
			}
		case ToPositiveInf:
			if (trunc == 1 || digit != 0) && !neg {
				incr = true
			}
		case ToNegativeInf:
			if (trunc == 1 || digit != 0) && neg {
				incr = true
			}
		}

		if incr {
			tsig := sig.add1()
			if tsig[1] > 0x0002_7fff_ffff_ffff {
				if trunc == 0 && digit != 0 {
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

var (
	DefaultRoundingMode RoundingMode = ToNearestEven
)
