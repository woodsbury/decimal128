package decimal128

import (
	"bytes"
	"fmt"
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

func FuzzOperations(f *testing.F) {
	const minBinaryExponent = -20517

	f.Fuzz(func(t *testing.T, xlo, xhi, ylo, yhi uint64) {
		t.Parallel()

		x := Decimal{xlo, xhi}
		y := Decimal{ylo, yhi}

		var bigx *big.Float
		var bigy *big.Float
		var bigres *big.Float
		if !x.isSpecial() && !y.isSpecial() {
			bigx = x.Float(nil)
			bigy = y.Float(nil)
			bigres = new(big.Float)
		}

		var buf []byte
		var cmpbuf []byte

		res := x.Add(y)

		if bigres == nil || res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			if idx != -1 {
				buf[idx-1] = '0'
			}

			fltres := x.Float64() + y.Float64()
			cmpbuf = fmt.Appendf(cmpbuf[:0], "%.3e", fltres)
			idx = bytes.IndexByte(cmpbuf, 'e')
			if idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Add(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
				cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
				idx = bytes.IndexByte(cmpbuf, 'e')
				cmpbuf[idx-1] = '0'

				if string(buf) != string(cmpbuf) {
					t.Fail()
				}
			} else if string(buf) != "0.000e+00" {
				t.Fail()
			}
		}

		res = x.Mul(y)

		if bigres == nil || res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			if idx != -1 {
				buf[idx-1] = '0'
			}

			fltres := x.Float64() * y.Float64()
			cmpbuf = fmt.Appendf(cmpbuf[:0], "%.3e", fltres)
			idx = bytes.IndexByte(cmpbuf, 'e')
			if idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Mul(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
				cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
				idx = bytes.IndexByte(cmpbuf, 'e')
				cmpbuf[idx-1] = '0'

				if string(buf) != string(cmpbuf) {
					t.Fail()
				}
			} else if string(buf) != "0.000e+00" {
				t.Fail()
			}
		}

		res = x.Quo(y)

		if bigres == nil || res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			if idx != -1 {
				buf[idx-1] = '0'
			}

			fltx := x.Float64()
			if res.isInf() && !x.IsZero() && y.IsZero() && fltx == 0 {
				if x.Signbit() {
					fltx = -1
				} else {
					fltx = 1
				}
			}

			fltres := fltx / y.Float64()
			cmpbuf = fmt.Appendf(cmpbuf[:0], "%.3e", fltres)
			idx = bytes.IndexByte(cmpbuf, 'e')
			if idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Quo(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
				cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
				idx = bytes.IndexByte(cmpbuf, 'e')
				cmpbuf[idx-1] = '0'

				if string(buf) != string(cmpbuf) {
					t.Fail()
				}
			} else if string(buf) != "0.000e+00" {
				t.Fail()
			}
		}

		res = x.Sub(y)

		if res.isSpecial() || bigres == nil {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			if idx != -1 {
				buf[idx-1] = '0'
			}

			fltres := x.Float64() - y.Float64()
			cmpbuf = fmt.Appendf(cmpbuf[:0], "%.3e", fltres)
			idx = bytes.IndexByte(cmpbuf, 'e')
			if idx != -1 {
				cmpbuf[idx-1] = '0'
			}

			if string(buf) != string(cmpbuf) {
				t.Fail()
			}
		} else {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Sub(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp > minBinaryExponent {
				cmpbuf = bigres.Append(cmpbuf[:0], 'e', 3)
				idx = bytes.IndexByte(cmpbuf, 'e')
				cmpbuf[idx-1] = '0'

				if string(buf) != string(cmpbuf) {
					t.Fail()
				}
			} else if string(buf) != "0.000e+00" {
				t.Fail()
			}
		}
	})
}

func BenchmarkOperations(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, 0, len(decimalValues))
	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		if val.sig == (uint128{}) {
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
