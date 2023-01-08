package decimal128

// Log returns the natural logarithm of d.
func Log(d Decimal) Decimal {
	if d.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if d.Signbit() {
			return nan(payloadOpLog, payloadValNegInfinite, 0)
		}

		return inf(false)
	}

	if d.IsZero() {
		return inf(true)
	}

	if d.Signbit() {
		return nan(payloadOpLog, payloadValNegFinite, 0)
	}

	dSig, dExp := d.decompose()
	l10 := int16(dSig.log10())
	dExp = (dExp - exponentBias) + l10

	msd := dSig.msd()
	sig := dSig
	exp := -l10
	oneSig := uint128{1, 0}
	oneExp := int16(0)

	for sig[1] <= 0x0002_7fff_ffff_ffff {
		sig = sig.mul64(10_000)
		exp -= 4

		oneSig = oneSig.mul64(10_000)
		oneExp -= 4
	}

	for sig[1] <= 0x18ff_ffff_ffff_ffff {
		sig = sig.mul64(10)
		exp--

		oneSig = oneSig.mul64(10)
		oneExp--
	}

	var rem uint128
	var trunc int8
	switch msd {
	case 2:
		sig, rem = sig.div(uint128{2, 0})
	case 3:
		sig, rem = sig.div(uint128{3, 0})
	case 4:
		sig, rem = sig.div(uint128{4, 0})
	case 5:
		sig, rem = sig.div(uint128{5, 0})
	case 6:
		sig, rem = sig.div(uint128{6, 0})
	case 7:
		sig, rem = sig.div(uint128{7, 0})
	case 8:
		sig, rem = sig.div(uint128{8, 0})
	case 9:
		sig, rem = sig.div(uint128{9, 0})
	}

	if rem != (uint128{}) {
		trunc = 1
	}

	nrm := decomposed128{
		sig: sig,
		exp: exp,
	}

	one := decomposed128{
		sig: oneSig,
		exp: oneExp,
	}

	_, num, _ := nrm.sub(one, int8(0))
	den, _ := nrm.add(one, int8(0))
	frc, trunc := num.quo(den, int8(0))
	sqr, _ := frc.mul(frc, int8(0))

	res := frc

	for i := uint64(3); i <= 65; i += 2 {
		// res += frc^i / i
		frc, _ = frc.mul(sqr, int8(0))
		tmp, _ := frc.quo(decomposed128{
			sig: uint128{i, 0},
			exp: 0,
		}, int8(0))

		res, trunc = res.add(tmp, trunc)
	}

	ln10 := decomposed128{
		sig: uint128{0x09bb_c25b_3ca8_1898, 0xad3a_2d01_4ad4_7d7a},
		exp: -38,
	}

	dExpNeg := false
	if dExp < 0 {
		dExp *= -1
		dExpNeg = true
	}

	ln10, _ = ln10.mul(decomposed128{
		sig: uint128{uint64(dExp), 0},
		exp: 0,
	}, int8(0))

	res, trunc = res.mul(decomposed128{
		sig: uint128{2, 0},
		exp: 0,
	}, trunc)

	neg := false
	if dExpNeg {
		neg, res, trunc = res.sub(ln10, trunc)
	} else {
		res, trunc = res.add(ln10, trunc)
	}

	var lnMSD decomposed128
	switch msd {
	case 2:
		lnMSD = decomposed128{
			sig: uint128{0x43d4_c3f7_1489_9de8, 0x3425_8773_b151_f6b7},
			exp: -38,
		}
	case 3:
		lnMSD = decomposed128{
			sig: uint128{0xbdf9_7603_a9e9_a361, 0x52a6_80c7_3db5_105b},
			exp: -38,
		}
	case 4:
		lnMSD = decomposed128{
			sig: uint128{0x87a9_87ee_2913_3bcf, 0x684b_0ee7_62a3_ed6e},
			exp: -38,
		}
	case 5:
		lnMSD = decomposed128{
			sig: uint128{0xc5e6_fe64_281e_7ab0, 0x7914_a58d_9982_86c2},
			exp: -38,
		}
	case 6:
		lnMSD = decomposed128{
			sig: uint128{0x86cc_083a_ef07_0713, 0x01ce_39fa_be73_4148},
			exp: -38,
		}
	case 7:
		lnMSD = decomposed128{
			sig: uint128{0x6930_c961_2ceb_f1e4, 0x9264_ddc2_a8f8_bea3},
			exp: -38,
		}
	case 8:
		lnMSD = decomposed128{
			sig: uint128{0xcb7e_4be5_3d9c_d9b6, 0x9c70_965b_13f5_e425},
			exp: -38,
		}
	case 9:
		lnMSD = decomposed128{
			sig: uint128{0x7bf2_ec07_53d3_46c1, 0xa54d_018e_7b6a_20b7},
			exp: -38,
		}
	}

	if dExpNeg {
		_, res, trunc = res.sub(lnMSD, trunc)
	} else {
		res, trunc = res.add(lnMSD, trunc)
	}

	sig, exp = ToNearestEven.reduce128(neg, res.sig, res.exp+exponentBias, trunc)

	return compose(neg, sig, exp)
}

// Sqrt returns the square root of d.
func Sqrt(d Decimal) Decimal {
	if d.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if d.Signbit() {
			return nan(payloadOpSqrt, payloadValNegInfinite, 0)
		}

		return d
	}

	if d.IsZero() {
		return d
	}

	if d.Signbit() {
		return nan(payloadOpSqrt, payloadValNegFinite, 0)
	}

	dSig, dExp := d.decompose()
	l10 := int16(dSig.log10())
	dExp = (dExp - exponentBias) + l10

	var add decomposed128
	var mul decomposed128
	var nrm decomposed128
	if dExp&1 == 0 {
		add = decomposed128{
			sig: uint128{259, 0},
			exp: -3,
		}

		mul = decomposed128{
			sig: uint128{819, 0},
			exp: -3,
		}

		nrm = decomposed128{
			sig: dSig,
			exp: -l10,
		}
	} else {
		add = decomposed128{
			sig: uint128{819, 0},
			exp: -4,
		}

		mul = decomposed128{
			sig: uint128{259, 0},
			exp: -2,
		}

		nrm = decomposed128{
			sig: dSig,
			exp: -l10 - 1,
		}

		dExp++
	}

	res, trunc := nrm.mul(mul, int8(0))
	res, trunc = res.add(add, trunc)

	var tmp decomposed128
	half := decomposed128{
		sig: uint128{5, 0},
		exp: -1,
	}

	for i := 0; i < 9; i++ {
		tmp, trunc = nrm.quo(res, trunc)
		res, trunc = res.add(tmp, trunc)
		res, trunc = half.mul(res, trunc)
	}

	res.exp += dExp / 2
	sig, exp := ToNearestEven.reduce128(false, res.sig, res.exp+exponentBias, trunc)

	return compose(false, sig, exp)
}

// Add adds d and o, rounded using the DefaultRoundingMode, and returns the
// result.
func (d Decimal) Add(o Decimal) Decimal {
	return d.AddWithMode(o, DefaultRoundingMode)
}

// AddWithMode adds d and o, rounding using the provided rounding mode, and
// returns the result.
func (d Decimal) AddWithMode(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if o.IsNaN() {
			return o
		}

		if d.isInf() {
			neg := d.Signbit()

			if o.isInf() && neg != o.Signbit() {
				lhs := payloadValPosInfinite
				rhs := payloadValNegInfinite
				if neg {
					lhs = payloadValNegInfinite
					rhs = payloadValPosInfinite
				}

				return nan(payloadOpAdd, lhs, rhs)
			}

			return inf(neg)
		}

		return inf(o.Signbit())
	}

	return d.add(o, mode, false)
}

// Mul multiplies d and o, rounding using the DefaultRoundingMode, and returns
// the result.
func (d Decimal) Mul(o Decimal) Decimal {
	return d.MulWithMode(o, DefaultRoundingMode)
}

// MulWithMode multiplies d and o, rounding using the provided rounding mode,
// and returns the result.
func (d Decimal) MulWithMode(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if o.IsNaN() {
			return o
		}

		if !d.isSpecial() {
			sig, _ := d.decompose()
			if sig == (uint128{}) {
				lhs := payloadValPosZero
				if d.Signbit() {
					lhs = payloadValNegZero
				}

				rhs := payloadValPosInfinite
				if o.Signbit() {
					rhs = payloadValNegInfinite
				}

				return nan(payloadOpMul, lhs, rhs)
			}
		} else if !o.isSpecial() {
			sig, _ := o.decompose()
			if sig == (uint128{}) {
				lhs := payloadValPosInfinite
				if d.Signbit() {
					lhs = payloadValNegInfinite
				}

				rhs := payloadValPosZero
				if o.Signbit() {
					rhs = payloadValNegZero
				}

				return nan(payloadOpMul, lhs, rhs)
			}
		}

		return inf(d.Signbit() != o.Signbit())
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	sig256 := dSig.mul(oSig)

	if sig256 == (uint256{}) {
		return zero(d.Signbit() != o.Signbit())
	}

	exp := (dExp - exponentBias) + (oExp - exponentBias) + exponentBias

	neg := d.Signbit() != o.Signbit()
	sig, exp := mode.reduce256(neg, sig256, exp, 0)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// Quo divides d by o, rounding using the DefaultRoundingMode, and returns the
// result.
func (d Decimal) Quo(o Decimal) Decimal {
	return d.QuoWithMode(o, DefaultRoundingMode)
}

// QuoWithMode divides d by o, rounding using the provided rounding mode, and
// returns the result.
func (d Decimal) QuoWithMode(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if o.IsNaN() {
			return o
		}

		if d.isInf() {
			if o.isInf() {
				lhs := payloadValPosInfinite
				if d.Signbit() {
					lhs = payloadValNegInfinite
				}

				rhs := payloadValPosInfinite
				if o.Signbit() {
					rhs = payloadValNegInfinite
				}

				return nan(payloadOpQuo, lhs, rhs)
			}

			return inf(d.Signbit() != o.Signbit())
		}

		if o.isInf() {
			return zero(d.Signbit() != o.Signbit())
		}
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if oSig == (uint128{}) {
		if dSig == (uint128{}) {
			lhs := payloadValPosZero
			if d.Signbit() {
				lhs = payloadValNegZero
			}

			rhs := payloadValPosZero
			if o.Signbit() {
				rhs = payloadValNegZero
			}

			return nan(payloadOpQuo, lhs, rhs)
		}

		return inf(d.Signbit() != o.Signbit())
	}

	if dSig == (uint128{}) {
		return zero(d.Signbit() != o.Signbit())
	}

	if dSig[1] == 0 {
		dSig = dSig.mul64(10_000_000_000_000_000_000)
		dExp -= 19
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
			exp++

			if rem192 != 0 {
				trunc = 1
			}
		}

		sig = uint128{sig192[0], sig192[1]}
	}

	if rem != (uint128{}) {
		trunc = 1
	}

	neg := d.Signbit() != o.Signbit()
	sig, exp = mode.reduce128(neg, sig, exp, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// QuoRem divides d by o, rounding using the DefaultRoundingMode, and returns
// the result as an integer quotient and a remainder.
func (d Decimal) QuoRem(o Decimal) (Decimal, Decimal) {
	return d.QuoRemWithMode(o, DefaultRoundingMode)
}

// QuoRem divides d by o, rounding using the provided rounding mode, and
// returns the result as an integer quotient and a remainder.
func (d Decimal) QuoRemWithMode(o Decimal, mode RoundingMode) (Decimal, Decimal) {
	if d.isSpecial() || o.isSpecial() {
		if d.IsNaN() {
			return d, d
		}

		if o.IsNaN() {
			return o, o
		}

		if d.isInf() {
			lhs := payloadValPosInfinite
			if d.Signbit() {
				lhs = payloadValNegInfinite
			}

			if o.isInf() {
				rhs := payloadValPosInfinite
				if o.Signbit() {
					rhs = payloadValNegInfinite
				}

				res := nan(payloadOpQuoRem, lhs, rhs)
				return res, res
			}

			rhs := payloadValPosFinite
			if o.IsZero() {
				if o.Signbit() {
					rhs = payloadValNegZero
				} else {
					rhs = payloadValPosZero
				}
			} else if o.Signbit() {
				rhs = payloadValNegFinite
			}

			return inf(d.Signbit() != o.Signbit()), nan(payloadOpQuoRem, lhs, rhs)
		}

		if o.isInf() {
			return zero(d.Signbit() != o.Signbit()), d
		}
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if oSig == (uint128{}) {
		rhs := payloadValPosZero
		if o.Signbit() {
			rhs = payloadValNegZero
		}

		if dSig == (uint128{}) {
			lhs := payloadValPosZero
			if d.Signbit() {
				lhs = payloadValNegZero
			}

			res := nan(payloadOpQuoRem, lhs, rhs)
			return res, res
		}

		lhs := payloadValPosFinite
		if d.Signbit() {
			lhs = payloadValNegFinite
		}

		return inf(d.Signbit() != o.Signbit()), nan(payloadOpQuoRem, lhs, rhs)
	}

	if dSig == (uint128{}) {
		return zero(d.Signbit() != o.Signbit()), zero(d.Signbit())
	}

	exp := (dExp - exponentBias) - (oExp - exponentBias)

	if exp < 0 {
		if exp <= -19 && oSig[1] == 0 {
			oSig = oSig.mul64(10_000_000_000_000_000_000)
			exp += 19
		}

		for exp <= -4 && oSig[1] <= 0x0002_7fff_ffff_ffff {
			oSig = oSig.mul64(10_000)
			exp += 4
		}

		for exp < 0 && oSig[1] <= 0x18ff_ffff_ffff_ffff {
			oSig = oSig.mul64(10)
			exp++
		}

		if exp < 0 || oSig.cmp(dSig) > 0 {
			return zero(d.Signbit() != o.Signbit()), d
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
	}

	sig, rem := dSig.div(oSig)
	trunc := int8(0)

	qexp := exp + exponentBias
	rexp := dExp

	for exp > 0 && rem != (uint128{}) && sig[1] <= 0x0002_7fff_ffff_ffff {
		for exp >= 4 && rem[1] <= 0x0002_7fff_ffff_ffff && sig[1] <= 0x0002_7fff_ffff_ffff {
			rem = rem.mul64(10_000)
			sig = sig.mul64(10_000)
			exp -= 4
			qexp -= 4
			rexp -= 4
		}

		for exp > 0 && rem[1] <= 0x18ff_ffff_ffff_ffff && sig[1] <= 0x18ff_ffff_ffff_ffff {
			rem = rem.mul64(10)
			sig = sig.mul64(10)
			exp--
			qexp--
			rexp--
		}

		var tmp uint128
		tmp, rem = rem.div(oSig)
		sig192 := sig.add(tmp)

		for sig192[2] != 0 {
			var rem192 uint64
			sig192, rem192 = sig192.div10()
			qexp++

			if rem192 != 0 {
				trunc = 1
			}
		}

		sig = uint128{sig192[0], sig192[1]}
	}

	for exp > 0 && rem != (uint128{}) {
		for exp >= 4 && rem[1] <= 0x0002_7fff_ffff_ffff {
			rem = rem.mul64(10_000)
			exp -= 4
			rexp -= 4
		}

		for exp > 0 && rem[1] <= 0x18ff_ffff_ffff_ffff {
			rem = rem.mul64(10)
			exp--
			rexp--
		}

		var tmp uint128
		tmp, rem = rem.div(oSig)

		if tmp != (uint128{}) {
			trunc = 1
		}
	}

	qneg := d.Signbit() != o.Signbit()
	qsig, qexp := mode.reduce128(qneg, sig, qexp, trunc)

	rneg := d.Signbit()
	rsig, rexp := mode.reduce128(rneg, rem, rexp, 0)

	quo := compose(qneg, qsig, qexp)

	if qexp > maxBiasedExponent {
		quo = inf(qneg)
	}

	if rexp > maxBiasedExponent {
		return quo, inf(rneg)
	}

	return quo, compose(rneg, rsig, rexp)
}

// Sub subtracts o from d, rounding using the DefaultRoundingMode, and returns
// the result.
func (d Decimal) Sub(o Decimal) Decimal {
	return d.SubWithMode(o, DefaultRoundingMode)
}

// SubWithMode subtracts o from d, rounding using the provided rounding mode,
// and returns the result.
func (d Decimal) SubWithMode(o Decimal, mode RoundingMode) Decimal {
	if d.isSpecial() || o.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if o.IsNaN() {
			return o
		}

		if d.isInf() {
			neg := d.Signbit()

			if o.isInf() && neg == o.Signbit() {
				lhs := payloadValPosInfinite
				rhs := payloadValPosInfinite
				if neg {
					lhs = payloadValNegInfinite
					rhs = payloadValNegInfinite
				}

				return nan(payloadOpSub, lhs, rhs)
			}

			return inf(neg)
		}

		return inf(!o.Signbit())
	}

	return d.add(o, mode, true)
}

func (d Decimal) add(o Decimal, mode RoundingMode, subtract bool) Decimal {
	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if dSig == (uint128{}) {
		if oSig == (uint128{}) {
			if subtract {
				return zero(d.Signbit() && !o.Signbit())
			} else {
				return zero(d.Signbit() && o.Signbit())
			}
		}

		if subtract {
			return compose(!o.Signbit(), oSig, oExp)
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

	dNeg := d.Signbit()
	oNeg := o.Signbit()
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
