package decimal128

// Exp returns e**d, the base-e exponential of d.
func Exp(d Decimal) Decimal {
	if d.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if d.Signbit() {
			return zero(false)
		}

		return inf(false)
	}

	if d.IsZero() {
		return one(false)
	}

	dSig, dExp := d.decompose()
	l10 := dSig.log10()

	if int(dExp) > exponentBias+5-l10 {
		if d.Signbit() {
			return zero(false)
		}

		return inf(false)
	}

	res, trunc := decomposed128{
		sig: dSig,
		exp: dExp - exponentBias,
	}.epow(int16(l10), int8(0))

	if res.exp > maxBiasedExponent-exponentBias+39 {
		if d.Signbit() {
			return zero(false)
		}

		return inf(false)
	}

	if d.Signbit() {
		res, trunc = decomposed128{
			sig: uint128{1, 0},
			exp: 0,
		}.quo(res, trunc)
	}

	sig, exp := DefaultRoundingMode.reduce128(false, res.sig, res.exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		if d.Signbit() {
			return zero(false)
		}

		return inf(false)
	}

	return compose(false, sig, exp)
}

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

	neg, res, trunc := decomposed128{
		sig: dSig,
		exp: dExp - exponentBias,
	}.log()

	sig, exp := DefaultRoundingMode.reduce128(neg, res.sig, res.exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// Log10 returns the decimal logarithm of d.
func Log10(d Decimal) Decimal {
	if d.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if d.Signbit() {
			return nan(payloadOpLog10, payloadValNegInfinite, 0)
		}

		return inf(false)
	}

	if d.IsZero() {
		return inf(true)
	}

	if d.Signbit() {
		return nan(payloadOpLog10, payloadValNegFinite, 0)
	}

	dSig, dExp := d.decompose()

	neg, res, trunc := decomposed128{
		sig: dSig,
		exp: dExp - exponentBias,
	}.log()

	res, trunc = res.mul(invLn10, trunc)

	sig, exp := DefaultRoundingMode.reduce128(neg, res.sig, res.exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// Log2 returns the binary logarithm of d.
func Log2(d Decimal) Decimal {
	if d.isSpecial() {
		if d.IsNaN() {
			return d
		}

		if d.Signbit() {
			return nan(payloadOpLog2, payloadValNegInfinite, 0)
		}

		return inf(false)
	}

	if d.IsZero() {
		return inf(true)
	}

	if d.Signbit() {
		return nan(payloadOpLog2, payloadValNegFinite, 0)
	}

	dSig, dExp := d.decompose()

	neg, res, trunc := decomposed128{
		sig: dSig,
		exp: dExp - exponentBias,
	}.log()

	res, trunc = res.mul(invLn2, trunc)

	sig, exp := DefaultRoundingMode.reduce128(neg, res.sig, res.exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

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

	for i := 0; i < 8; i++ {
		tmp, trunc = nrm.quo(res, trunc)
		res, trunc = res.add(tmp, trunc)
		res, trunc = half.mul(res, trunc)
	}

	res.exp += dExp / 2
	sig, exp := DefaultRoundingMode.reduce128(false, res.sig, res.exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(false)
	}

	return compose(false, sig, exp)
}
