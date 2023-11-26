package decimal128

import (
	"bytes"
	"math"
	"math/big"
	"testing"
)

func TestDecimalAdd(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v + %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			sum := lhs.AddWithMode(rhs, mode)

			if !res.equal(sum, mode) {
				t.Errorf("%v.AddWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, sum, res.result(mode))
			}
		}
	}
}

func TestDecimalMul(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v * %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			prd := lhs.MulWithMode(rhs, mode)

			if !res.equal(prd, mode) {
				t.Errorf("%v.MulWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, prd, res.result(mode))
			}
		}
	}
}

func TestDecimalPow(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v ^ %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			pwr := lhs.PowWithMode(rhs, mode)

			if !res.equal(pwr, mode) {
				t.Errorf("%v.PowWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, pwr, res.result(mode))
			}
		}
	}
}

func TestDecimalQuo(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v / %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			quo := lhs.QuoWithMode(rhs, mode)

			if !res.equal(quo, mode) {
				t.Errorf("%v.QuoWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, quo, res.result(mode))
			}
		}
	}
}

func TestDecimalQuoRem(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResultPair

	res.sep = 'r'

	for r.scan("%v / %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			quo, rem := lhs.QuoRemWithMode(rhs, mode)

			if !res.equal(quo, rem, mode) {
				t.Errorf("%v.QuoRemWithMode(%v, %v) = (%v, %v), want (%v, %v)", lhs, rhs, mode, quo, rem, res.first.result(mode), res.second.result(mode))
			}
		}
	}
}

func TestDecimalSub(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testDataResult

	for r.scan("%v - %v = %v\n", &lhs, &rhs, &res) {
		for _, mode := range roundingModes {
			dif := lhs.SubWithMode(rhs, mode)

			if !res.equal(dif, mode) {
				t.Errorf("%v.SubWithMode(%v, %v) = %v, want %v", lhs, rhs, mode, dif, res.result(mode))
			}
		}
	}
}

func FuzzArith(f *testing.F) {
	values := []float64{
		0.0,
		math.Copysign(0.0, -1.0),
		0.5,
		1.0,
		5.0,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
	}

	for _, x := range values {
		for _, y := range values {
			f.Add(x, y)
		}
	}

	f.Fuzz(func(t *testing.T, x, y float64) {
		t.Parallel()

		decx := FromFloat64(x)
		decy := FromFloat64(y)

		decres := decx.Add(decy).Float64()
		fltres := x + y
		eps := math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp := math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = decx.Mul(decy).Float64()
		fltres = x * y
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = decx.Pow(decy).Float64()
		fltres = math.Pow(x, y)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = decx.Quo(decy).Float64()
		fltres = x / y
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		res, rem := decx.QuoRem(decy)
		decres, decrem := res.Float64(), rem.Float64()
		decres = decres*y + decrem
		fltres = x
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = decx.Sub(decy).Float64()
		fltres = x - y
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+20)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}
	})
}

func FuzzBigArith(f *testing.F) {
	const minBinaryExponent = -20517

	f.Fuzz(func(t *testing.T, xlo, xhi, ylo, yhi uint64) {
		t.Parallel()

		x := Decimal{xlo, xhi}
		y := Decimal{ylo, yhi}

		if x.IsNaN() || y.IsNaN() {
			t.Skip()
		}

		bigx := x.Float(nil)
		bigy := y.Float(nil)
		bigres := new(big.Float)

		var buf []byte
		var cmpbuf []byte

		res := x.Add(y)
		buf = res.Append(buf[:0], ".3e")
		if idx := bytes.IndexByte(buf, 'e'); idx != -1 {
			buf[idx-1] = '0'
		}

		bigres.Add(bigx, bigy)
		if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
			cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
			if idx := bytes.IndexByte(cmpbuf, 'e'); idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else if string(buf) != "0.000e+00" {
			t.Fail()
		}

		res = x.Mul(y)
		buf = res.Append(buf[:0], ".3e")
		if idx := bytes.IndexByte(buf, 'e'); idx != -1 {
			buf[idx-1] = '0'
		}

		bigres.Mul(bigx, bigy)
		if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
			cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
			if idx := bytes.IndexByte(cmpbuf, 'e'); idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else if string(buf) != "0.000e+00" {
			t.Fail()
		}

		if !y.IsZero() && !y.isInf() {
			res = x.Quo(y)
			buf = res.Append(buf[:0], ".3e")
			if idx := bytes.IndexByte(buf, 'e'); idx != -1 {
				buf[idx-1] = '0'
			}

			bigres.Quo(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
				cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
				if idx := bytes.IndexByte(cmpbuf, 'e'); idx != -1 {
					cmpbuf[idx-1] = '0'
				}

				if string(buf) != string(cmpbuf) {
					t.Fail()
				}
			} else if string(buf) != "0.000e+00" {
				t.Fail()
			}
		}

		res = x.Sub(y)
		buf = res.Append(buf[:0], ".3e")
		if idx := bytes.IndexByte(buf, 'e'); idx != -1 {
			buf[idx-1] = '0'
		}

		bigres.Sub(bigx, bigy)
		if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
			cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
			if idx := bytes.IndexByte(cmpbuf, 'e'); idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else if string(buf) != "0.000e+00" {
			t.Fail()
		}
	})
}

func BenchmarkArith(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, 0, len(decimalValues))
	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		if val.sig[0]|val.sig[1] == 0 {
			continue
		}

		values = append(values, val.Decimal())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, lhs := range values {
			for _, rhs := range values {
				lhs.Add(rhs)
				lhs.Mul(rhs)
				lhs.Quo(rhs)
				lhs.Sub(rhs)
			}
		}
	}
}
