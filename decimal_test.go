package decimal128

import (
	"fmt"
	"math"
	"math/big"
	"sync"
	"testing"
	"unsafe"

	"github.com/cockroachdb/apd/v3"
)

var maxDecimal *apd.Decimal

func init() {
	lo := new(apd.BigInt)
	lo.SetUint64(math.MaxUint64)

	maxDecimal = new(apd.Decimal)
	maxDecimal.Exponent = maxBiasedExponent - exponentBias
	maxDecimal.Coeff.SetUint64(0x0002_7fff_ffff_ffff)
	maxDecimal.Coeff.Lsh(&maxDecimal.Coeff, 64).Add(&maxDecimal.Coeff, lo)
}

type testForm uint8

const (
	regularForm testForm = iota
	infForm
	nanForm
)

type testDec struct {
	form testForm
	neg  bool
	sig  uint128
	exp  int16
}

func (td testDec) Big(dec *apd.Decimal) *apd.Decimal {
	switch td.form {
	case regularForm:
		dec.Form = apd.Finite
		dec.Negative = td.neg
		dec.Exponent = int32(td.exp - exponentBias)

		if td.sig[1] == 0 {
			dec.Coeff.SetUint64(td.sig[0])
		} else {
			lo := new(apd.BigInt)
			lo.SetUint64(td.sig[0])

			dec.Coeff.SetUint64(td.sig[1])
			dec.Coeff.Lsh(&dec.Coeff, 64).Add(&dec.Coeff, lo)
		}

		dec.Reduce(dec)
		dec.Negative = td.neg
	case infForm:
		dec.Form = apd.Infinite
		dec.Negative = td.neg
	case nanForm:
		dec.Form = apd.NaN
	default:
		panic("unhandled test decimal form")
	}

	return dec
}

func (td testDec) Decimal() Decimal {
	switch td.form {
	case regularForm:
		return compose(td.neg, td.sig, td.exp)
	case infForm:
		return inf(td.neg)
	case nanForm:
		return nan()
	default:
		panic("unhandled test decimal form")
	}
}

func (td testDec) Float(flt *big.Float) *big.Float {
	if td.form == nanForm {
		panic("NaN big.Float")
	}

	if td.form == infForm {
		return flt.SetInf(td.neg)
	}

	flt.SetPrec(128)

	if td.sig[1] == 0 {
		flt.SetUint64(td.sig[0])
	} else {
		bigsig := new(big.Int).SetUint64(td.sig[1])
		bigsig.Lsh(bigsig, 64).Add(bigsig, new(big.Int).SetUint64(td.sig[0]))

		flt.SetInt(bigsig)
	}

	if td.neg {
		flt.Neg(flt)
	}

	if td.exp == exponentBias {
		return flt
	}

	exp := td.exp - exponentBias

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

func (td testDec) String() string {
	switch td.form {
	case regularForm:
		sign := ""
		if td.neg {
			sign = "-"
		}

		return fmt.Sprintf("%s%ve%d", sign, td.sig, td.exp-exponentBias)
	case infForm:
		sign := ""
		if td.neg {
			sign = "-"
		}

		return fmt.Sprintf("%sinf", sign)
	case nanForm:
		return "nan"
	default:
		panic("unhandled test decimal form")
	}
}

var (
	decimalValuesOnce sync.Once
	decimalValues     []testDec
)

func initDecimalValues() {
	decimalValuesOnce.Do(func() {
		initUintValues()

		var exponentValues []int16
		if testing.Short() {
			exponentValues = []int16{
				minBiasedExponent,
				exponentBias - 19,
				exponentBias,
				exponentBias + 19,
				maxBiasedExponent,
			}
		} else {
			exponentValues = []int16{
				minBiasedExponent,
				minBiasedExponent + exponentBias/2,
				exponentBias - 34,
				exponentBias - 19,
				exponentBias,
				exponentBias + 19,
				exponentBias + 34,
				maxBiasedExponent - exponentBias/2,
				maxBiasedExponent,
			}
		}

		for _, sighi := range uint64Values {
			if sighi > 0x0002_7fff_ffff_ffff {
				continue
			}

			for _, siglo := range uint64Values {
				sig := uint128{siglo, sighi}

				for _, exp := range exponentValues {
					decimalValues = append(decimalValues,
						testDec{regularForm, false, sig, exp},
						testDec{regularForm, true, sig, exp},
					)
				}
			}
		}

		decimalValues = append(decimalValues, testDec{infForm, false, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{infForm, true, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{nanForm, false, uint128{}, 0})
	})
}

func decimalToBig(v Decimal) *apd.Decimal {
	r := new(apd.Decimal)

	if v.isInf() {
		r.Form = apd.Infinite
		r.Negative = v.isNeg()
		return r
	}

	if v.isNaN() {
		r.Form = apd.NaN
		return r
	}

	r.Negative = v.isNeg()

	sig, exp := v.decompose()
	r.Exponent = int32(exp - exponentBias)

	if sig[1] == 0 {
		r.Coeff.SetUint64(sig[0])
	} else {
		lo := new(apd.BigInt)
		lo.SetUint64(sig[0])

		r.Coeff.SetUint64(sig[1])
		r.Coeff.Lsh(&r.Coeff, 64).Add(&r.Coeff, lo)
	}

	r.Reduce(r)
	return r
}

func decimalsEqual(x Decimal, y *apd.Decimal, mode apd.Rounder) bool {
	if x.isSpecial() {
		if x.isNaN() {
			return y.Form == apd.NaN || y.Form == apd.NaNSignaling
		}

		if x.isInf() {
			if y.Negative != x.isNeg() {
				return false
			}

			if y.Form == apd.Infinite {
				return true
			}

			neg := y.Negative
			y.Negative = false
			cmp := y.Cmp(maxDecimal)
			y.Negative = neg

			return cmp > 0
		}
	} else if y.Form != apd.Finite {
		return false
	}

	if x.isNeg() != y.Negative {
		// apd appears to always return -0 when rounding towards -infinity,
		// even if the operands are themselves zero.
		if x.IsZero() && y.Coeff.IsInt64() && y.Coeff.Int64() == 0 && mode == apd.RoundFloor {
			return true
		}

		return false
	}

	bigx := decimalToBig(x)

	bigctx := apd.Context{
		Precision:   uint32(bigx.NumDigits()),
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    mode,
	}

	bigctx.Round(y, y)

	// apd appears to return the wrong result during rounding in some scenarios
	// when the result underflows, returning +/-1e-6176 instead of 0.
	if x.IsZero() && y.Exponent <= -6176 && y.Coeff.IsInt64() && y.Coeff.Int64() == 1 {
		return true
	}

	return bigx.Cmp(y) == 0
}

func TestCompose(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		dec := compose(val.neg, val.sig, val.exp)

		if dec.isSpecial() {
			t.Errorf("%v.isSpecial() = true, want false", val)
		}

		if dec.isInf() {
			t.Errorf("%v.isInf() = true, want false", val)
		}

		if dec.isNaN() {
			t.Errorf("%v.isNaN() = true, want false", val)
		}

		if dec.isNeg() != val.neg {
			t.Errorf("%v.isNeg() = %t, want %t", val, dec.isNeg(), val.neg)
		}

		sig, exp := dec.decompose()

		if sig != val.sig || exp != val.exp {
			t.Errorf("%v.decompose() = (%v, %d), want (%v, %d)", val, sig, exp, val.sig, val.exp)
		}
	}
}

func TestInf(t *testing.T) {
	t.Parallel()

	dec := inf(false)

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if !dec.isInf() {
		t.Errorf("%v.isInf() = false, want true", dec)
	}

	if dec.isNaN() {
		t.Errorf("%v.isNaN() = true, want false", dec)
	}

	if dec.isNeg() {
		t.Errorf("%v.isNeg() = true, want false", dec)
	}

	dec = inf(true)

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if !dec.isInf() {
		t.Errorf("%v.isInf() = false, want true", dec)
	}

	if dec.isNaN() {
		t.Errorf("%v.isNaN() = true, want false", dec)
	}

	if !dec.isNeg() {
		t.Errorf("%v.isNeg() = false, want true", dec)
	}
}

func TestNaN(t *testing.T) {
	t.Parallel()

	dec := nan()

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if dec.isInf() {
		t.Errorf("%v.isInf() = true, want false", dec)
	}

	if !dec.isNaN() {
		t.Errorf("%v.isNaN() = false, want true", dec)
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	res := unsafe.Sizeof(Decimal{})

	if res != 16 {
		t.Errorf("unsafe.Sizeof(Decimal{}) = %d, want 16", res)
	}
}

func FuzzDecimal(f *testing.F) {
	f.Add(uint64(0), uint64(0))
	f.Add(uint64(math.MaxUint64), uint64(math.MaxUint64))

	f.Fuzz(func(t *testing.T, hi, lo uint64) {
		t.Parallel()

		dec := Decimal{hi, lo}

		if dec.isSpecial() {
			if dec.isInf() == dec.isNaN() {
				t.Fail()
			}
		} else {
			if dec.isInf() || dec.isNaN() {
				t.Fail()
			}

			sig, exp := dec.decompose()
			res := compose(dec.isNeg(), sig, exp)

			if res != dec {
				t.Fail()
			}
		}
	})
}
