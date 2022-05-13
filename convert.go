package decimal128

import (
	"math"
	"math/big"
)

func (d Decimal) Float() *big.Float {
	if d.isSpecial() {
		if d.isNaN() {
			panic("Decimal(NaN).Float()")
		}

		return new(big.Float).SetInf(d.isNeg())
	}

	sig, exp := d.decompose()

	flt := new(big.Float).SetPrec(128)

	if sig[1] == 0 {
		flt.SetUint64(sig[0])
	} else {
		bigsig := new(big.Int).SetUint64(sig[1])
		bigsig.Lsh(bigsig, 64).Add(bigsig, new(big.Int).SetUint64(sig[0]))

		flt.SetInt(bigsig)
	}

	if d.isNeg() {
		flt.Neg(flt)
	}

	if exp == exponentBias {
		return flt
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
		flt.Mul(flt, new(big.Float).SetInt(bigexp))
	} else {
		flt.Quo(flt, new(big.Float).SetInt(bigexp))
	}

	return flt
}

func (d Decimal) Float32() float32 {
	return float32(d.Float64())
}

func (d Decimal) Float64() float64 {
	if d.isSpecial() {
		if d.isNaN() {
			return math.NaN()
		}

		if d.isNeg() {
			return math.Inf(-1)
		}

		return math.Inf(1)
	}

	sig, exp := d.decompose()

	if sig == (uint128{}) {
		flt := 0.0
		if d.isNeg() {
			flt = math.Copysign(flt, -1.0)
		}

		return flt
	}

	if exp == exponentBias && sig[1] == 0 {
		flt := float64(sig[0])
		if d.isNeg() {
			flt = math.Copysign(flt, -1.0)
		}

		return flt
	}

	panic("not implemented")
}

func (d Decimal) Rat() *big.Rat {
	if d.isSpecial() {
		if d.isNaN() {
			panic("Decimal(NaN).Rat()")
		}

		if d.isNeg() {
			panic("Decimal(-Inf).Rat()")
		}

		panic("Decimal(+Inf).Rat()")
	}

	sig, exp := d.decompose()

	rat := new(big.Rat)

	if exp == exponentBias && sig[1] == 0 {
		rat.SetUint64(sig[0])
	} else {
		bigsig := new(big.Int).SetUint64(sig[1])
		bigsig.Lsh(bigsig, 64).Add(bigsig, new(big.Int).SetUint64(sig[0]))

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
			rat.SetInt(bigsig)
		} else {
			rat.SetFrac(bigsig, bigexp)
		}
	}

	if d.isNeg() {
		rat.Neg(rat)
	}

	return rat
}
