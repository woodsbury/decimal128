package decimal128

func (d Decimal) Add(o Decimal) Decimal {
	return d.AddWithRounding(o, DefaultRoundingMode)
}

func (d Decimal) AddWithRounding(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() {
			return d
		}

		if o.isNaN() {
			return o
		}

		if d.isInf() {
			neg := d.isNeg()

			if o.isInf() && neg != o.isNeg() {
				return nan()
			}

			return inf(neg)
		}

		return inf(o.isNeg())
	}

	return d.add(o, mode, false)
}

func (d Decimal) Mul(o Decimal) Decimal {
	return d.MulWithRounding(o, DefaultRoundingMode)
}

func (d Decimal) MulWithRounding(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() {
			return d
		}

		if o.isNaN() {
			return o
		}

		if !d.isSpecial() {
			sig, _ := d.decompose()
			if sig == (uint128{}) {
				return nan()
			}
		} else if !o.isSpecial() {
			sig, _ := o.decompose()
			if sig == (uint128{}) {
				return nan()
			}
		}

		return inf(d.isNeg() != o.isNeg())
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	sig256 := dSig.mul(oSig)

	if sig256 == (uint256{}) {
		return zero(d.isNeg() != o.isNeg())
	}

	exp := (dExp - exponentBias) + (oExp - exponentBias) + exponentBias

	neg := d.isNeg() != o.isNeg()
	sig, exp := mode.reduce256(neg, sig256, exp)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

func (d Decimal) Quo(o Decimal) Decimal {
	return d.QuoWithRounding(o, DefaultRoundingMode)
}

func (d Decimal) QuoWithRounding(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() {
			return d
		}

		if o.isNaN() {
			return o
		}

		if d.isInf() {
			if o.isInf() {
				return nan()
			}

			return inf(d.isNeg() != o.isNeg())
		}

		if o.isInf() {
			return zero(d.isNeg() != o.isNeg())
		}
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if oSig == (uint128{}) {
		if dSig == (uint128{}) {
			return nan()
		}

		return inf(d.isNeg() != o.isNeg())
	}

	if dSig == (uint128{}) {
		return zero(d.isNeg() != o.isNeg())
	}

	for dSig[1] <= 0x0002_7fff_ffff_ffff {
		dSig = dSig.mul64(10_000)
		dExp -= 4
	}

	for dSig[1] <= 0x18ff_ffff_ffff_ffff {
		dSig = dSig.mul64(10)
		dExp--
	}

	sig, rem := dSig.div(oSig)
	exp := (dExp - exponentBias) - (oExp - exponentBias) + exponentBias
	trunc := int8(0)

	for rem != (uint128{}) && sig[1] <= 0x0002_7fff_ffff_ffff {
		for rem[1] <= 0x0002_7fff_ffff_ffff && sig[1] <= 0x0002_7fff_ffff_ffff {
			rem = rem.mul64(10_000)
			sig = sig.mul64(10_000)
			exp -= 4
		}

		for rem[1] <= 0x18ff_ffff_ffff_ffff && sig[1] <= 0x18ff_ffff_ffff_ffff {
			rem = rem.mul64(10)
			sig = sig.mul64(10)
			exp--
		}

		var tmp uint128
		tmp, rem = rem.div(oSig)
		sig192 := sig.add(tmp)

		for sig192[2] != 0 {
			var rem192 uint64
			sig192, rem192 = sig192.div10()
			exp--

			if rem192 != 0 {
				trunc = 1
			}
		}

		sig = uint128{sig192[0], sig192[1]}
	}

	if rem != (uint128{}) {
		trunc = 1
	}

	neg := d.isNeg() != o.isNeg()
	sig, exp = mode.reduce128(neg, sig, exp, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

func (d Decimal) Sub(o Decimal) Decimal {
	return d.SubWithRounding(o, DefaultRoundingMode)
}

func (d Decimal) SubWithRounding(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() {
			return d
		}

		if o.isNaN() {
			return o
		}

		if d.isInf() {
			neg := d.isNeg()

			if o.isInf() && neg == o.isNeg() {
				return nan()
			}

			return inf(neg)
		}

		return inf(!o.isNeg())
	}

	return d.add(o, mode, true)
}

func (d Decimal) add(o Decimal, mode RoundingMode, subtract bool) Decimal {
	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if dSig == (uint128{}) {
		if oSig == (uint128{}) {
			if subtract {
				return zero(d.isNeg() && !o.isNeg())
			} else {
				return zero(d.isNeg() && o.isNeg())
			}
		}

		if subtract {
			return compose(!o.isNeg(), oSig, oExp)
		}

		return o
	}

	if oSig == (uint128{}) {
		return d
	}

	exp := dExp - oExp
	trunc := int8(0)

	if exp < 0 {
		if exp <= -19 && oSig[1] == 0 {
			oSig = oSig.mul64(10_000_000_000_000_000_000)
			oExp -= 19
			exp += 19
		}

		for exp <= -4 && oSig[1] <= 0x0002_7fff_ffff_ffff {
			oSig = oSig.mul64(10_000)
			oExp -= 4
			exp += 4
		}

		for exp < 0 && oSig[1] <= 0x18ff_ffff_ffff_ffff {
			oSig = oSig.mul64(10)
			oExp--
			exp++
		}

		if exp < -maxDigits {
			if dSig != (uint128{}) {
				dSig = uint128{}
				trunc = 1
			}

			dExp = oExp
			exp = 0
		}

		if exp <= -3 {
			var rem uint64
			dSig, rem = dSig.div1000()
			if rem != 0 {
				trunc = 1
			}

			if dSig == (uint128{}) {
				dExp = oExp
				exp = 0
			} else {
				dExp += 3
				exp += 3
			}
		}

		for exp < 0 {
			var rem uint64
			dSig, rem = dSig.div10()
			if rem != 0 {
				trunc = 1
			}

			if dSig == (uint128{}) {
				dExp = oExp
				exp = 0
				break
			}

			dExp++
			exp++
		}
	} else if exp > 0 {
		if exp >= 19 && dSig[1] == 0 {
			dSig = dSig.mul64(10_000_000_000_000_000_000)
			dExp -= 19
			exp -= 19
		}

		for exp >= 4 && dSig[1] <= 0x0002_7fff_ffff_ffff {
			dSig = dSig.mul64(10_000)
			dExp -= 4
			exp -= 4
		}

		for exp > 0 && dSig[1] <= 0x18ff_ffff_ffff_ffff {
			dSig = dSig.mul64(10)
			dExp--
			exp--
		}

		if exp > maxDigits {
			if oSig != (uint128{}) {
				oSig = uint128{}
				trunc = -1
			}

			exp = 0
		}

		if exp >= 3 {
			var rem uint64
			oSig, rem = oSig.div1000()
			if rem != 0 {
				trunc = -1
			}

			if oSig == (uint128{}) {
				exp = 0
			} else {
				exp -= 3
			}
		}

		for exp > 0 {
			var rem uint64
			oSig, rem = oSig.div10()
			if rem != 0 {
				trunc = -1
			}

			if oSig == (uint128{}) {
				exp = 0
				break
			}

			exp--
		}
	}

	dNeg := d.isNeg()
	oNeg := o.isNeg()
	if subtract {
		oNeg = !oNeg
	}

	neg := dNeg

	var sig uint128
	if dNeg == oNeg {
		sig192 := dSig.add(oSig)

		if sig192 == (uint192{}) {
			return zero(mode == ToNegativeInf)
		}

		if trunc == -1 {
			trunc = 1
		}

		sig, exp = mode.reduce192(neg, sig192, dExp, trunc)
	} else {
		var brw uint
		sig, brw = dSig.sub(oSig)
		if brw != 0 {
			sig = sig.twos()
			neg = !neg
			trunc *= -1
		} else if sig == (uint128{}) {
			return zero(mode == ToNegativeInf)
		}

		sig, exp = mode.reduce128(neg, sig, dExp, trunc)
	}

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}
