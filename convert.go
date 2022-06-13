package decimal128

import (
	"math"
	"math/big"
	"math/bits"
)

// FromFloat converts f into a Decimal.
func FromFloat(f *big.Float) Decimal {
	if f.IsInf() {
		return inf(f.Signbit())
	}

	if f.Sign() == 0 {
		return zero(f.Signbit())
	}

	r, _ := f.Rat(nil)
	return FromRat(r)
}

// FromFloat32 converts f into a Decimal.
func FromFloat32(f float32) Decimal {
	if math.IsNaN(float64(f)) {
		return nan(payloadOpFromFloat32, 0, 0)
	}

	return FromFloat64(float64(f))
}

// FromFloat64 converts f into a Decimal.
func FromFloat64(f float64) Decimal {
	if math.IsNaN(f) {
		return nan(payloadOpFromFloat64, 0, 0)
	}

	if math.IsInf(f, 0) {
		return inf(math.Signbit(f))
	}

	if f == 0.0 {
		return zero(math.Signbit(f))
	}

	fbits := math.Float64bits(f)
	mant := fbits & 0x000f_ffff_ffff_ffff
	exp := int16(fbits >> 52 & 0x07ff)
	neg := fbits&0x8000_0000_0000_0000 != 0

	if exp == 0 {
		exp = -1022
	} else {
		mant |= 0x0010_0000_0000_0000
		exp -= 1023
	}

	shift := int(52 - exp)

	if shift == 0 {
		return compose(neg, uint128{mant, 0}, 0)
	}

	var sig256 uint256
	exp = exponentBias
	trunc := int8(0)

	if shift < 0 {
		shift *= -1
		zeros := bits.LeadingZeros64(mant)

		if zeros > shift {
			zeros = shift
		}

		mant <<= zeros
		shift -= zeros

		sig256 = uint256{mant, 0, 0, 0}

		if shift >= 192 {
			sig256 = sig256.lsh(192)
			shift -= 192
		} else if shift >= 128 {
			sig256 = sig256.lsh(128)
			shift -= 128
		}

		for shift != 0 {
			if sig256[3] != 0 {
				var rem uint64
				sig256, rem = sig256.div1e19()
				exp += 19

				if rem != 0 {
					trunc = 1
				}
			}

			if shift > 64 {
				sig256 = sig256.lsh(64)
				shift -= 64
			} else {
				sig256 = sig256.lsh(uint(shift))
				break
			}
		}
	} else {
		zeros := bits.TrailingZeros64(mant)

		if zeros > shift {
			zeros = shift
		}

		mant >>= zeros
		shift -= zeros

		sig := uint128{mant, 0}
		sig256 = sig.mul1e38()
		exp -= 38

		for shift != 0 {
			if sig256[2] == 0 && sig256[3] == 0 {
				sig = uint128{sig256[0], sig256[1]}
				sig256 = sig.mul1e38()
				exp -= 38
			}

			if shift > 128 {
				if sig[0] != 0 || sig[1] != 0 {
					trunc = 1
				}

				sig256 = sig256.rsh(128)
				shift -= 128
			} else {
				if shift < 64 {
					if sig[0]&(1<<shift-1) != 0 {
						trunc = 1
					}
				} else {
					if sig[0] != 0 && sig[1]&(1<<(shift-64)-1) != 0 {
						trunc = 1
					}
				}

				sig256 = sig256.rsh(uint(shift))
				break
			}
		}
	}

	sig, exp := DefaultRoundingMode.reduce256(neg, sig256, exp, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// FromInt converts i into a Decimal.
func FromInt(i *big.Int) Decimal {
	neg := false
	if sgn := i.Sign(); sgn == 0 {
		return zero(false)
	} else if sgn < 0 {
		neg = true
	}

	exp := int16(exponentBias)
	trunc := int8(0)

	if bl := i.BitLen(); bl > 128 {
		i = new(big.Int).Set(i)
		r := new(big.Int)

		if bl > 256 {
			e18 := big.NewInt(1_000_000_000_000_000_000)

			for bl > 256 {
				i.QuoRem(i, e18, r)
				exp += 18

				if exp > maxBiasedExponent {
					return inf(neg)
				}

				bl = i.BitLen()

				if r.Sign() != 0 {
					trunc = 1
				}
			}
		}

		ten := big.NewInt(10)

		for bl > 128 {
			i.QuoRem(i, ten, r)
			exp++

			if exp > maxBiasedExponent {
				return inf(neg)
			}

			bl = i.BitLen()

			if r.Sign() != 0 {
				trunc = 1
			}
		}
	}

	var sig uint128

	b := i.Bits()
	for i := len(b) - 1; i >= 0; i-- {
		sig = sig.lsh(bits.UintSize)
		sig = sig.or64(uint64(b[i]))
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp, trunc)

	if exp > maxBiasedExponent {
		return inf(neg)
	}

	return compose(neg, sig, exp)
}

// FromInt32 converts i into a Decimal.
func FromInt32(i int32) Decimal {
	return FromInt64(int64(i))
}

// FromInt64 converts i into a Decimal.
func FromInt64(i int64) Decimal {
	if i == 0 {
		return zero(false)
	}

	neg := false
	if i < 0 {
		neg = true
		i *= -1
	}

	return compose(neg, uint128{uint64(i), 0}, 0)
}

// FromRat converts r into a Decimal.
func FromRat(r *big.Rat) Decimal {
	num := r.Num()

	if num.Sign() == 0 {
		return zero(false)
	}

	denom := r.Denom()

	return FromInt(num).Quo(FromInt(denom))
}

// FromUint32 converts i into a Decimal.
func FromUint32(i uint32) Decimal {
	return FromUint64(uint64(i))
}

// FromUint64 converts i into a Decimal.
func FromUint64(i uint64) Decimal {
	if i == 0 {
		return zero(false)
	}

	return compose(false, uint128{i, 0}, 0)
}

// Float converts d into a big.Float. It panics if d is NaN.
func (d Decimal) Float() *big.Float {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Float()")
		}

		return new(big.Float).SetInf(d.isNeg())
	}

	sig, exp := d.decompose()

	f := new(big.Float).SetPrec(128)

	if sig[1] == 0 {
		f.SetUint64(sig[0])
	} else {
		bigsig := new(big.Int).SetUint64(sig[1])
		bigsig.Lsh(bigsig, 64).Or(bigsig, new(big.Int).SetUint64(sig[0]))

		f.SetInt(bigsig)
	}

	if d.isNeg() {
		f.Neg(f)
	}

	if exp == exponentBias {
		return f
	}

	exp -= exponentBias

	var bigexp *big.Int
	if exp > 0 {
		bigexp = big.NewInt(int64(exp))
	} else {
		bigexp = big.NewInt(int64(exp * -1))
	}

	bigexp.Exp(big.NewInt(10), bigexp, nil)

	if exp > 0 {
		f.Mul(f, new(big.Float).SetInt(bigexp))
	} else {
		f.Quo(f, new(big.Float).SetInt(bigexp))
	}

	return f
}

// Float32 converts d into a float32.
func (d Decimal) Float32() float32 {
	return float32(d.Float64())
}

// Float64 converts d into a float64.
func (d Decimal) Float64() float64 {
	if d.isSpecial() {
		if d.IsNaN() {
			return math.NaN()
		}

		if d.isNeg() {
			return math.Inf(-1)
		}

		return math.Inf(1)
	}

	sig, exp := d.decompose()

	if sig == (uint128{}) {
		f := 0.0
		if d.isNeg() {
			f = math.Copysign(f, -1.0)
		}

		return f
	}

	expf := math.Pow10(int(exp - exponentBias))

	if sig[1] != 0 {
		shift := 64 - bits.LeadingZeros64(sig[1])
		sig = sig.rsh(uint(shift))
		expf *= math.Exp2(float64(shift))
	}

	sigf := float64(sig[0])

	if d.isNeg() {
		sigf = math.Copysign(sigf, -1.0)
	}

	return sigf * expf
}

// Int converts d into a big.Int, truncating towards zero. It panics if d is
// NaN or infinite.
func (d Decimal) Int() *big.Int {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Int()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Int()")
		}

		panic("Decimal(+Inf).Int()")
	}

	sig, exp := d.decompose()

	i := new(big.Int)

	if exp < -maxDigits {
		return i
	}

	if sig[1] == 0 {
		i.SetUint64(sig[0])
	} else {
		i.SetUint64(sig[1])
		i.Lsh(i, 64).Or(i, new(big.Int).SetUint64(sig[0]))
	}

	if d.isNeg() {
		i.Neg(i)
	}

	if exp == exponentBias {
		return i
	}

	exp -= exponentBias

	var bigexp *big.Int
	if exp > 0 {
		bigexp = big.NewInt(int64(exp))
	} else {
		bigexp = big.NewInt(int64(exp * -1))
	}

	bigexp.Exp(big.NewInt(10), bigexp, nil)

	if exp > 0 {
		i.Mul(i, bigexp)
	} else {
		i.Quo(i, bigexp)
	}

	return i
}

// Int32 converts d into an int32, truncating towards zero. It panics if d
// cannot be represented by an int32.
func (d Decimal) Int32() int32 {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Int32()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Int32()")
		}

		panic("Decimal(+Inf).Int32()")
	}

	sig, exp := d.decompose()

	if exp < -maxDigits {
		return 0
	}

	for sig[1] != 0 && exp < 0 {
		sig, _ = sig.div10()
		exp++
	}

	for sig[1] == 0 && exp > 0 {
		sig = sig.mul64(10)
		exp--
	}

	if sig[1] != 0 || exp != 0 {
		if d.isNeg() {
			panic("Decimal(<MinInt32).Int32()")
		}

		panic("Decimal(>MaxInt32).Int32()")
	}

	neg := d.isNeg()

	if neg {
		if sig[0] > math.MinInt32*-1 {
			panic("Decimal(<MinInt32).Int32()")
		}
	} else {
		if sig[0] > math.MaxInt64 {
			panic("Decimal(>MaxInt32).Int32()")
		}
	}

	i := int32(sig[0])

	if neg {
		i *= -1
	}

	return i
}

// Int64 converts d into an int64, truncating towards zero. It panics if d
// cannot be represented by an int64.
func (d Decimal) Int64() int64 {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Int64()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Int64()")
		}

		panic("Decimal(+Inf).Int64()")
	}

	sig, exp := d.decompose()

	if exp < -maxDigits {
		return 0
	}

	for sig[1] != 0 && exp < 0 {
		sig, _ = sig.div10()
		exp++
	}

	for sig[1] == 0 && exp > 0 {
		sig = sig.mul64(10)
		exp--
	}

	if sig[1] != 0 || exp != 0 {
		if d.isNeg() {
			panic("Decimal(<MinInt64).Int64()")
		}

		panic("Decimal(>MaxInt64).Int64()")
	}

	neg := d.isNeg()

	if neg {
		if sig[0] > math.MinInt64*-1 {
			panic("Decimal(<MinInt64).Int64()")
		}
	} else {
		if sig[0] > math.MaxInt64 {
			panic("Decimal(>MaxInt64).Int64()")
		}
	}

	i := int64(sig[0])

	if neg {
		i *= -1
	}

	return i
}

// Rat converts d into a big.Rat. It panics if d is NaN or infinite.
func (d Decimal) Rat() *big.Rat {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Rat()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Rat()")
		}

		panic("Decimal(+Inf).Rat()")
	}

	sig, exp := d.decompose()

	r := new(big.Rat)

	if exp == exponentBias && sig[1] == 0 {
		r.SetUint64(sig[0])
	} else {
		bigsig := new(big.Int).SetUint64(sig[1])
		bigsig.Lsh(bigsig, 64).Or(bigsig, new(big.Int).SetUint64(sig[0]))

		exp -= exponentBias

		var bigexp *big.Int
		if exp > 0 {
			bigexp = big.NewInt(int64(exp))
		} else {
			bigexp = big.NewInt(int64(exp * -1))
		}

		bigexp.Exp(big.NewInt(10), bigexp, nil)

		if exp > 0 {
			bigsig.Mul(bigsig, bigexp)
			r.SetInt(bigsig)
		} else {
			r.SetFrac(bigsig, bigexp)
		}
	}

	if d.isNeg() {
		r.Neg(r)
	}

	return r
}

// Uint32 converts d into a uint32, truncating towards zero. It panics if d
// cannot be represented by a uint32.
func (d Decimal) Uint32() uint32 {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Uint32()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Uint32()")
		}

		panic("Decimal(+Inf).Uint32()")
	}

	if d.isNeg() {
		panic("Decimal(<0).Uint32()")
	}

	sig, exp := d.decompose()

	if exp < -maxDigits {
		return 0
	}

	for exp < 0 {
		sig, _ = sig.div10()
		exp++

		if sig == (uint128{}) {
			exp = 0
			break
		}
	}

	for sig[1] == 0 && exp > 0 {
		sig = sig.mul64(10)
		exp--
	}

	if sig[1] != 0 || exp != 0 {
		panic("Decimal(>MaxUint32).Uint32()")
	}

	if sig[0] > math.MaxUint32 {
		panic("Decimal(>MaxUint32).Uint32()")
	}

	return uint32(sig[0])
}

// Uint64 converts d into a uint64, truncating towards zero. It panics if d
// cannot be represented by a uint64.
func (d Decimal) Uint64() uint64 {
	if d.isSpecial() {
		if d.IsNaN() {
			panic("Decimal(NaN).Uint64()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Uint64()")
		}

		panic("Decimal(+Inf).Uint64()")
	}

	if d.isNeg() {
		panic("Decimal(<0).Uint64()")
	}

	sig, exp := d.decompose()

	if exp < -maxDigits {
		return 0
	}

	for exp < 0 {
		sig, _ = sig.div10()
		exp++

		if sig == (uint128{}) {
			exp = 0
			break
		}
	}

	for sig[1] == 0 && exp > 0 {
		sig = sig.mul64(10)
		exp--
	}

	if sig[1] != 0 || exp != 0 {
		panic("Decimal(>MaxUint64).Uint64()")
	}

	return sig[0]
}
