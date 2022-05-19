package decimal128

const (
	exponentBias      = 6176
	maxBiasedExponent = 12287
	minBiasedExponent = 0
	maxDigits         = 35
)

type Decimal struct {
	lo, hi uint64
}

func New(sig int64, exp int) Decimal {
	if sig == 0 {
		return zero(false)
	}

	neg := false
	if sig < 0 {
		neg = true
		sig *= -1
	}

	if exp < minBiasedExponent-exponentBias+19 {
		return zero(neg)
	}

	if exp > maxBiasedExponent-exponentBias+39 {
		return inf(neg)
	}

	sig128, exp16 := DefaultRoundingMode.reduce64(neg, uint64(sig), int16(exp+exponentBias))

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig128, exp16)
}

func compose(neg bool, sig uint128, exp int16) Decimal {
	var hi uint64
	if sig[1] > 0x0001_ffff_ffff_ffff {
		hi = 0x6000_0000_0000_0000 | uint64(exp)<<47 | sig[1]&0x7fff_ffff_ffff
	} else {
		hi = uint64(exp)<<49 | sig[1]
	}

	if neg {
		hi |= 0x8000_0000_0000_0000
	}

	return Decimal{sig[0], hi}
}

func inf(neg bool) Decimal {
	if neg {
		return Decimal{0, 0xf800_0000_0000_0000}
	}

	return Decimal{0, 0x7800_0000_0000_0000}
}

func nan() Decimal {
	return Decimal{0, 0x7c00_0000_0000_0000}
}

func zero(neg bool) Decimal {
	if neg {
		return Decimal{0, 0x8000_0000_0000_0000}
	}

	return Decimal{}
}

func (d Decimal) decompose() (uint128, int16) {
	var sig uint128
	var exp int16

	if d.hi&0x6000_0000_0000_0000 == 0x6000_0000_0000_0000 {
		sig = uint128{d.lo, d.hi&0x7fff_ffff_ffff | 0x0002_0000_0000_0000}
		exp = int16(d.hi & 0x1fff_8000_0000_0000 >> 47)
	} else {
		sig = uint128{d.lo, d.hi & 0x0001_ffff_ffff_ffff}
		exp = int16(d.hi & 0x7ffe_0000_0000_0000 >> 49)
	}

	return sig, exp
}

func (d Decimal) isInf() bool {
	return d.hi&0x7c00_0000_0000_0000 == 0x7800_0000_0000_0000
}

func (d Decimal) isNaN() bool {
	return d.hi&0x7c00_0000_0000_0000 == 0x7c00_0000_0000_0000
}

func (d Decimal) isNeg() bool {
	return d.hi&0x8000_0000_0000_0000 == 0x8000_0000_0000_0000
}

func (d Decimal) isSpecial() bool {
	return d.hi&0x7800_0000_0000_0000 == 0x7800_0000_0000_0000
}
