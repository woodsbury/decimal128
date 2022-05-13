package decimal128

type CmpResult int8

const (
	cmpLess    CmpResult = -1
	cmpEqual   CmpResult = 0
	cmpGreater CmpResult = 1
	cmpNaN     CmpResult = 2
)

func (cr CmpResult) Equal() bool {
	return cr == cmpEqual
}

func (cr CmpResult) Greater() bool {
	return cr == cmpGreater
}

func (cr CmpResult) Less() bool {
	return cr == cmpLess
}

func (d Decimal) Cmp(o Decimal) CmpResult {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() || o.isNaN() {
			return cmpNaN
		}

		if d.isInf() {
			neg := d.isNeg()

			if o.isInf() && neg == o.isNeg() {
				return cmpEqual
			}

			if neg {
				return cmpLess
			}

			return cmpGreater
		}

		if o.isInf() {
			if o.isNeg() {
				return cmpGreater
			}

			return cmpLess
		}
	}

	if d == o {
		return cmpEqual
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if dSig == (uint128{}) {
		if oSig == (uint128{}) {
			return cmpEqual
		}

		if o.isNeg() {
			return cmpGreater
		}

		return cmpLess
	}

	if oSig == (uint128{}) {
		if d.isNeg() {
			return cmpLess
		}

		return cmpGreater
	}

	neg := d.isNeg()

	if neg != o.isNeg() {
		if neg {
			return cmpLess
		}

		return cmpGreater
	}

	exp := dExp - oExp
	trunc := false

	var res CmpResult
	if neg {
		res = cmpLess
	} else {
		res = cmpGreater
	}

	if exp < 0 {
		if exp <= -19 {
			if exp < -maxDigits {
				return res * -1
			}

			var rem uint64
			dSig, rem = dSig.div1e19()
			if dSig == (uint128{}) {
				return res * -1
			}

			if rem != 0 {
				trunc = true
			}

			exp += 19
		}

		exp *= -1
		dSig, oSig = oSig, dSig
		res *= -1
	} else if exp >= 19 {
		if exp > maxDigits {
			return res
		}

		var rem uint64
		oSig, rem = oSig.div1e19()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 19
	}

	for exp > 8 {
		var rem uint64
		oSig, rem = oSig.div1e8()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 8
	}

	for exp > 3 {
		var rem uint64
		oSig, rem = oSig.div1000()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 3
	}

	for exp > 0 {
		var rem uint64
		oSig, rem = oSig.div10()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp--
	}

	sres := dSig.cmp(oSig)
	if sres == 0 {
		if trunc {
			return res * -1
		}

		return cmpEqual
	}

	if res == cmpLess {
		return CmpResult(sres * -1)
	}

	return CmpResult(sres)
}

func (d Decimal) CmpAbs(o Decimal) CmpResult {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() || o.isNaN() {
			return cmpNaN
		}

		if d.isInf() {
			if o.isInf() {
				return cmpEqual
			}

			return cmpGreater
		}

		if o.isInf() {
			return cmpLess
		}
	}

	if d == o {
		return cmpEqual
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if dSig == (uint128{}) {
		if oSig == (uint128{}) {
			return cmpEqual
		}

		return cmpLess
	}

	if oSig == (uint128{}) {
		return cmpGreater
	}

	exp := dExp - oExp
	trunc := false
	res := cmpGreater

	if exp < 0 {
		if exp <= -19 {
			if exp < -maxDigits {
				return cmpLess
			}

			var rem uint64
			dSig, rem = dSig.div1e19()
			if dSig == (uint128{}) {
				return cmpLess
			}

			if rem != 0 {
				trunc = true
			}

			exp += 19
		}

		exp *= -1
		dSig, oSig = oSig, dSig
		res = cmpLess
	} else if exp >= 19 {
		if exp > maxDigits {
			return cmpGreater
		}

		var rem uint64
		oSig, rem = oSig.div1e19()
		if oSig == (uint128{}) {
			return cmpGreater
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 19
	}

	for exp > 8 {
		var rem uint64
		oSig, rem = oSig.div1e8()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 8
	}

	for exp > 3 {
		var rem uint64
		oSig, rem = oSig.div1000()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp -= 3
	}

	for exp > 0 {
		var rem uint64
		oSig, rem = oSig.div10()
		if oSig == (uint128{}) {
			return res
		}

		if rem != 0 {
			trunc = true
		}

		exp--
	}

	sres := dSig.cmp(oSig)
	if sres == 0 {
		if trunc {
			return res * -1
		}

		return cmpEqual
	}

	if res == cmpLess {
		return CmpResult(sres * -1)
	}

	return CmpResult(sres)
}

func (d Decimal) Equal(o Decimal) bool {
	if d.isSpecial() || o.isSpecial() {
		if d.isNaN() || o.isNaN() {
			return false
		}

		if d.isInf() {
			return o.isInf() && d.isNeg() == o.isNeg()
		}

		if o.isInf() {
			return false
		}
	}

	if d == o {
		return true
	}

	dSig, dExp := d.decompose()
	oSig, oExp := o.decompose()

	if dSig == (uint128{}) {
		return oSig == (uint128{})
	}

	if oSig == (uint128{}) {
		return false
	}

	if d.isNeg() != o.isNeg() {
		return false
	}

	exp := dExp - oExp

	if exp < 0 {
		if exp <= -19 {
			if exp < -maxDigits {
				return false
			}

			var rem uint64
			dSig, rem = dSig.div1e19()
			if rem != 0 {
				return false
			}

			exp += 19
		}

		exp *= -1
		dSig, oSig = oSig, dSig
	} else if exp >= 19 {
		if exp > maxDigits {
			return false
		}

		var rem uint64
		oSig, rem = oSig.div1e19()
		if rem != 0 {
			return false
		}

		exp -= 19
	}

	for exp > 8 {
		var rem uint64
		oSig, rem = oSig.div1e8()
		if rem != 0 {
			return false
		}

		exp -= 8
	}

	for exp > 3 {
		var rem uint64
		oSig, rem = oSig.div1000()
		if rem != 0 {
			return false
		}

		exp -= 3
	}

	for exp > 0 {
		var rem uint64
		oSig, rem = oSig.div10()
		if rem != 0 {
			return false
		}

		exp--
	}

	return dSig == oSig
}

func (d Decimal) IsZero() bool {
	if d == (Decimal{}) {
		return true
	}

	if d.hi&0x6000_0000_0000_0000 == 0x6000_0000_0000_0000 {
		return false
	} else {
		return d.lo == 0 && d.hi&0x0001_ffff_ffff_ffff == 0
	}
}
