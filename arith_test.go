package decimal128

import (
	"bytes"
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
	f.Fuzz(func(t *testing.T, xneg bool, xlo, xhi uint64, xexp uint16, yneg bool, ylo, yhi uint64, yexp uint16) {
		t.Parallel()

		if xhi > 0x0002_7fff_ffff_ffff || yhi > 0x0002_7fff_ffff_ffff {
			t.Skip()
		}

		if xexp > maxBiasedExponent || yexp > maxBiasedExponent {
			t.Skip()
		}

		x := compose(xneg, uint128{xlo, xhi}, int16(xexp))
		y := compose(yneg, uint128{ylo, yhi}, int16(yexp))

		bigx := x.Float(nil)
		bigy := y.Float(nil)
		bigres := new(big.Float)
		var bigbuf []byte

		res := x.Add(y)
		var buf []byte

		if !res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Add(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp < maxUnbiasedExponent && bigexp > minUnbiasedExponent {
				bigbuf = bigres.Append(bigbuf[:0], 'e', 3)
				idx = bytes.IndexByte(bigbuf, 'e')
				bigbuf[idx-1] = '0'

				if string(buf) != string(bigbuf) {
					t.Fail()
				}
			}
		}

		res = x.Mul(y)

		if !res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Mul(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp < maxUnbiasedExponent && bigexp > minUnbiasedExponent {
				bigbuf = bigres.Append(bigbuf[:0], 'e', 3)
				idx = bytes.IndexByte(bigbuf, 'e')
				bigbuf[idx-1] = '0'

				if string(buf) != string(bigbuf) {
					t.Fail()
				}
			}
		}

		res = x.Quo(y)

		if !res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Quo(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp < maxUnbiasedExponent && bigexp > minUnbiasedExponent {
				bigbuf = bigres.Append(bigbuf[:0], 'e', 3)
				idx = bytes.IndexByte(bigbuf, 'e')
				bigbuf[idx-1] = '0'

				if string(buf) != string(bigbuf) {
					t.Fail()
				}
			}
		}

		res = x.Sub(y)

		if !res.isSpecial() {
			buf = res.Append(buf[:0], ".3e")
			idx := bytes.IndexByte(buf, 'e')
			buf[idx-1] = '0'

			bigres.Sub(bigx, bigy)
			if bigexp := bigres.MantExp(nil); bigexp < maxUnbiasedExponent && bigexp > minUnbiasedExponent {
				bigbuf = bigres.Append(bigbuf[:0], 'e', 3)
				idx = bytes.IndexByte(bigbuf, 'e')
				bigbuf[idx-1] = '0'

				if string(buf) != string(bigbuf) {
					t.Fail()
				}
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
