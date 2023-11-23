package decimal128

import (
	"math"
	"math/big"
	"sync"
	"testing"
)

var (
	uintValuesOnce sync.Once
	uint64Values   []uint64
	uint128Values  []uint128
	uint192Values  []uint192
	uint256Values  []uint256
	uint384Values  []uint384
)

func initUintValues() {
	uintValuesOnce.Do(func() {
		uint64Values = []uint64{
			0,
			1,
			2,
			3,
			4,
			5,
			6,
			7,
			8,
			9,
			10,
			2500,
			math.MaxUint32,
			0x0001_ffff_ffff_ffff,
			0x0002_7fff_ffff_ffff,
			math.MaxUint64,
		}

		n := len(uint64Values)

		uint128Values = make([]uint128, 0, n*n*4)
		uint192Values = make([]uint192, 0, n*n*n)
		uint256Values = make([]uint256, 0, n*n*n)
		uint384Values = make([]uint384, 0, n*n*n)

		for _, l0 := range uint64Values {
			for _, l1 := range uint64Values {
				val := uint128{l0, l1}
				uint128Values = append(uint128Values, val)

				if l0 > 10 || l1 > 0 {
					val, _ = val.sub(uint128{1, 0})
					uint128Values = append(uint128Values, val)
				}

				val = val.mul64(100_000_000)
				uint128Values = append(uint128Values, val)

				val, _ = val.sub(uint128{1, 0})
				uint128Values = append(uint128Values, val)

				for _, l2 := range uint64Values {
					uint192Values = append(uint192Values, uint192{l0, l1, l2})
					uint256Values = append(uint256Values, uint256{l0, l1, l2, l1})
					uint384Values = append(uint384Values, uint384{l0, l1, l2, l2, l1, l0})
				}
			}
		}
	})
}

func uint128ToBig(v uint128, r *big.Int) *big.Int {
	u := new(big.Int)
	r.SetUint64(v[1])
	u.SetUint64(v[0])
	r.Lsh(r, 64).Add(r, u)

	return r
}

func uint192ToBig(v uint192, r *big.Int) *big.Int {
	u := new(big.Int)
	r.SetUint64(v[2])
	u.SetUint64(v[1])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[0])
	r.Lsh(r, 64).Add(r, u)

	return r
}

func uint256ToBig(v uint256, r *big.Int) *big.Int {
	u := new(big.Int)
	r.SetUint64(v[3])
	u.SetUint64(v[2])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[1])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[0])
	r.Lsh(r, 64).Add(r, u)

	return r
}

func uint384ToBig(v uint384, r *big.Int) *big.Int {
	u := new(big.Int)
	r.SetUint64(v[5])
	u.SetUint64(v[4])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[3])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[2])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[1])
	r.Lsh(r, 64).Add(r, u)
	u.SetUint64(v[0])
	r.Lsh(r, 64).Add(r, u)

	return r
}

func TestUint128Add(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpsum := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint128Values {
			sum := lhs.add(rhs)

			uint128ToBig(lhs, biglhs)
			uint128ToBig(rhs, bigrhs)
			bigsum := biglhs.Add(biglhs, bigrhs)

			if uint192ToBig(sum, tmpsum).Cmp(bigsum) != 0 {
				t.Errorf("%v.add(%v) = %v, want %v", lhs, rhs, sum, bigsum)
			}
		}
	}
}

func TestUint128Add64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpsum := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint64Values {
			sum := lhs.add64(rhs)

			uint128ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigsum := biglhs.Add(biglhs, bigrhs)
			if bigsum.BitLen() > 128 {
				bigsum.SetBytes(bigsum.Bytes()[1:])
			}

			if uint128ToBig(sum, tmpsum).Cmp(bigsum) != 0 {
				t.Errorf("%v.add64(%v) = %v, want %v", lhs, rhs, sum, bigsum)
			}
		}
	}
}

func TestUint128Cmp(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint128Values {
			lhs := lhs
			rhs := rhs

			res := lhs.cmp(rhs)

			uint128ToBig(lhs, biglhs)
			uint128ToBig(rhs, bigrhs)
			bigres := biglhs.Cmp(bigrhs)

			if res != bigres {
				t.Errorf("%v.cmp(%v) = %d, want %d", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint128Div(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint128Values {
			if rhs[0]|rhs[1] == 0 {
				continue
			}

			quo, rem := lhs.div(rhs)

			uint128ToBig(lhs, biglhs)
			uint128ToBig(rhs, bigrhs)
			bigquo, bigrem := biglhs.QuoRem(biglhs, bigrhs, bigrem)

			if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || uint128ToBig(rem, tmprem).Cmp(bigrem) != 0 {
				t.Errorf("%v.div(%v) = (%v, %v), want (%v, %v)", lhs, rhs, quo, rem, bigquo, bigrem)
			}
		}
	}
}

func TestUint128Div10(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div10()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Div100(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(100)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div100()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div100() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Div1000(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(1000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div1000()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1000() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Div10000(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div10000()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10000() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Div1e8(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(100_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div1e8()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e8() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Div1e19(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := new(big.Int).SetUint64(10_000_000_000_000_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint128Values {
		quo, rem := val.div1e19()

		uint128ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint128ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e19() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint128Log10(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint128Values {
		res := val.log10()

		uint128ToBig(val, bigval)
		bigres := len(bigval.Text(10)) - 1

		if res != bigres {
			t.Errorf("%v.log10() = %v, want %v", val, res, bigres)
		}
	}
}

func TestUint128Lsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint128Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.lsh(rhs)

			uint128ToBig(lhs, biglhs)
			bigres := biglhs.Lsh(biglhs, rhs)
			if bigres.BitLen() > 128 {
				b := bigres.Bytes()
				bigres.SetBytes(b[len(b)-16:])
			}

			if uint128ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.lsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint128Mul(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpprd := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint128Values {
			prd := lhs.mul(rhs)

			uint128ToBig(lhs, biglhs)
			uint128ToBig(rhs, bigrhs)
			bigprd := biglhs.Mul(biglhs, bigrhs)

			if uint256ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
				t.Errorf("%v.mul(%v) = %v, want %v", lhs, rhs, prd, bigprd)
			}
		}
	}
}

func TestUint128Mul1e38(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := new(big.Int).Exp(new(big.Int).SetUint64(10), new(big.Int).SetUint64(38), nil)

	bigval := new(big.Int)
	tmpprd := new(big.Int)

	for _, val := range uint128Values {
		prd := val.mul1e38()

		uint128ToBig(val, bigval)
		bigprd := bigval.Mul(bigval, c)

		if uint256ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
			t.Errorf("%v.mul1e38() = %v, want %v", val, prd, bigprd)
		}
	}
}

func TestUint128Mul64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpprd := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint64Values {
			prd := lhs.mul64(rhs)

			uint128ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigprd := biglhs.Mul(biglhs, bigrhs)
			if bigprd.BitLen() > 128 {
				b := bigprd.Bytes()
				bigprd.SetBytes(b[len(b)-16:])
			}

			if uint128ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
				t.Errorf("%v.mul64(%d) = %v, want %v", lhs, rhs, prd, bigprd)
			}
		}
	}
}

func TestUint128Or64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint64Values {
			res := lhs.or64(rhs)

			uint128ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigres := biglhs.Or(biglhs, bigrhs)

			if uint128ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.or64(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint128Rsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint128Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.rsh(rhs)

			uint128ToBig(lhs, biglhs)
			bigres := biglhs.Rsh(biglhs, rhs)

			if uint128ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.rsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint128String(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint128Values {
		res := val.String()

		uint128ToBig(val, bigval)
		bigres := bigval.String()

		if res != bigres {
			t.Errorf("%v.String() = %s, want %s", val, res, bigres)
		}
	}
}

func TestUint128Sub(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpdif := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint128Values {
			dif, brw := lhs.sub(rhs)

			uint128ToBig(lhs, biglhs)
			uint128ToBig(rhs, bigrhs)
			bigdif := biglhs.Sub(biglhs, bigrhs)
			bigbrw := uint(0)
			if bigdif.Sign() == -1 {
				bigbrw = 1

				b := make([]byte, 16)
				c := bigdif.Bytes()
				copy(b[16-len(c):], c)

				for i := range b {
					b[i] = ^b[i]
				}

				bigdif.SetBytes(b)
				bigdif.Add(bigdif, one)
			}

			if uint128ToBig(dif, tmpdif).Cmp(bigdif) != 0 || brw != bigbrw {
				t.Errorf("%v.sub(%v) = (%v, %v), want (%v, %v)", lhs, rhs, dif, brw, bigdif, bigbrw)
			}
		}
	}
}

func TestUint128Sub64(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpdif := new(big.Int)

	for _, lhs := range uint128Values {
		for _, rhs := range uint64Values {
			dif := lhs.sub64(rhs)

			uint128ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigdif := biglhs.Sub(biglhs, bigrhs)
			if bigdif.Sign() == -1 {
				b := make([]byte, 16)
				c := bigdif.Bytes()
				copy(b[16-len(c):], c)

				for i := range b {
					b[i] = ^b[i]
				}

				bigdif.SetBytes(b)
				bigdif.Add(bigdif, one)
			}

			if uint128ToBig(dif, tmpdif).Cmp(bigdif) != 0 {
				t.Errorf("%v.sub64(%v) = %v, want %v", lhs, rhs, dif, bigdif)
			}
		}
	}
}

func TestUint128Twos(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	bigval := new(big.Int)
	tmpres := new(big.Int)

	for _, val := range uint128Values {
		res := val.twos()

		uint128ToBig(val, bigval)
		b := make([]byte, 16)
		c := bigval.Bytes()
		copy(b[16-len(c):], c)

		for i := range b {
			b[i] = ^b[i]
		}

		bigval.SetBytes(b)
		bigres := bigval.Add(bigval, one)
		if bigres.BitLen() > 128 {
			b = bigres.Bytes()
			bigres.SetBytes(b[len(b)-16:])
		}

		if uint128ToBig(res, tmpres).Cmp(bigres) != 0 {
			t.Errorf("%v.twos() = %v, want %v", val, res, bigres)
		}
	}
}

func TestUint192Add(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpsum := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint192Values {
			sum := lhs.add(rhs)

			uint192ToBig(lhs, biglhs)
			uint192ToBig(rhs, bigrhs)
			bigsum := biglhs.Add(biglhs, bigrhs)

			if uint256ToBig(sum, tmpsum).Cmp(bigsum) != 0 {
				t.Errorf("%v.add(%v) = %v, want %v", lhs, rhs, sum, bigsum)
			}
		}
	}
}

func TestUint192Add64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpsum := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint64Values {
			sum := lhs.add64(rhs)

			uint192ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigsum := biglhs.Add(biglhs, bigrhs)
			if bigsum.BitLen() > 192 {
				bigsum.SetBytes(bigsum.Bytes()[1:])
			}

			if uint192ToBig(sum, tmpsum).Cmp(bigsum) != 0 {
				t.Errorf("%v.add64(%v) = %v, want %v", lhs, rhs, sum, bigsum)
			}
		}
	}
}

func TestUint192Cmp(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint192Values {
			lhs := lhs
			rhs := rhs

			res := lhs.cmp(rhs)

			uint192ToBig(lhs, biglhs)
			uint192ToBig(rhs, bigrhs)
			bigres := biglhs.Cmp(bigrhs)

			if res != bigres {
				t.Errorf("%v.cmp(%v) = %d, want %d", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint192Div(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint192Values {
			if rhs[0]|rhs[1]|rhs[2] == 0 {
				continue
			}

			quo, rem := lhs.div(rhs)

			uint192ToBig(lhs, biglhs)
			uint192ToBig(rhs, bigrhs)
			bigquo, bigrem := biglhs.QuoRem(biglhs, bigrhs, bigrem)

			if uint192ToBig(quo, tmpquo).Cmp(bigquo) != 0 || uint192ToBig(rem, tmprem).Cmp(bigrem) != 0 {
				t.Errorf("%v.div(%v) = (%v, %v), want (%v, %v)", lhs, rhs, quo, rem, bigquo, bigrem)
			}
		}
	}
}

func TestUint192Div10(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint192Values {
		quo, rem := val.div10()

		uint192ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint192ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint192Div10000(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint192Values {
		quo, rem := val.div10000()

		uint192ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint192ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10000() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint192Div1e8(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(100_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint192Values {
		quo, rem := val.div1e8()

		uint192ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint192ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e8() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint192Div1e19(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := new(big.Int).SetUint64(10_000_000_000_000_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint192Values {
		quo, rem := val.div1e19()

		uint192ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint192ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e8() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint192Log10(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint192Values {
		res := val.log10()

		uint192ToBig(val, bigval)
		bigres := len(bigval.Text(10)) - 1

		if res != bigres {
			t.Errorf("%v.log10() = %v, want %v", val, res, bigres)
		}
	}
}

func TestUint192Lsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint192Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.lsh(rhs)

			uint192ToBig(lhs, biglhs)
			bigres := biglhs.Lsh(biglhs, rhs)
			if bigres.BitLen() > 192 {
				b := bigres.Bytes()
				bigres.SetBytes(b[len(b)-24:])
			}

			if uint192ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.lsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint192Mul(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpprd := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint192Values {
			prd := lhs.mul(rhs)

			uint192ToBig(lhs, biglhs)
			uint192ToBig(rhs, bigrhs)
			bigprd := biglhs.Mul(biglhs, bigrhs)

			if uint384ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
				t.Errorf("%v.mul(%v) = %v, want %v", lhs, rhs, prd, bigprd)
			}
		}
	}
}

func TestUint192Mul64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpprd := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint64Values {
			prd := lhs.mul64(rhs)

			uint192ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigprd := biglhs.Mul(biglhs, bigrhs)
			if bigprd.BitLen() > 192 {
				b := bigprd.Bytes()
				bigprd.SetBytes(b[len(b)-24:])
			}

			if uint192ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
				t.Errorf("%v.mul64(%v) = %v, want %v", lhs, rhs, prd, bigprd)
			}
		}
	}
}

func TestUint192Msd2(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint192Values {
		res := val.msd2()

		uint192ToBig(val, bigval)
		bigtxt := bigval.Text(10)
		bigres := int(bigtxt[0]) - '0'
		if len(bigtxt) > 1 {
			bigres *= 10
			bigres += int(bigtxt[1]) - '0'
		}

		if res != bigres {
			t.Errorf("%v.msd2() = %v, want %v", val, res, bigres)
		}
	}
}

func TestUint192Pow2(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)
	tmpprd := new(big.Int)

	for _, val := range uint192Values {
		prd := val.pow2()

		uint192ToBig(val, bigval)
		bigprd := bigval.Mul(bigval, bigval)

		if uint384ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
			t.Errorf("%v.pow2() = %v, want %v", val, prd, bigprd)
		}
	}
}

func TestUint192Rsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint192Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.rsh(rhs)

			uint192ToBig(lhs, biglhs)
			bigres := biglhs.Rsh(biglhs, rhs)

			if uint192ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.rsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint192String(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint192Values {
		res := val.String()

		uint192ToBig(val, bigval)
		bigres := bigval.String()

		if res != bigres {
			t.Errorf("%v.String() = %s, want %s", val, res, bigres)
		}
	}
}

func TestUint192Sub(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpdif := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint192Values {
			dif, brw := lhs.sub(rhs)

			uint192ToBig(lhs, biglhs)
			uint192ToBig(rhs, bigrhs)
			bigdif := biglhs.Sub(biglhs, bigrhs)
			bigbrw := uint(0)
			if bigdif.Sign() == -1 {
				bigbrw = 1

				b := make([]byte, 24)
				c := bigdif.Bytes()
				copy(b[24-len(c):], c)

				for i := range b {
					b[i] = ^b[i]
				}

				bigdif.SetBytes(b)
				bigdif.Add(bigdif, one)
			}

			if uint192ToBig(dif, tmpdif).Cmp(bigdif) != 0 || brw != bigbrw {
				t.Errorf("%v.sub(%v) = (%v, %v), want (%v, %v)", lhs, rhs, dif, brw, bigdif, bigbrw)
			}
		}
	}
}

func TestUint192Sub64(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpdif := new(big.Int)

	for _, lhs := range uint192Values {
		for _, rhs := range uint64Values {
			dif := lhs.sub64(rhs)

			uint192ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigdif := biglhs.Sub(biglhs, bigrhs)
			if bigdif.Sign() == -1 {
				b := make([]byte, 24)
				c := bigdif.Bytes()
				copy(b[24-len(c):], c)

				for i := range b {
					b[i] = ^b[i]
				}

				bigdif.SetBytes(b)
				bigdif.Add(bigdif, one)
			}

			if uint192ToBig(dif, tmpdif).Cmp(bigdif) != 0 {
				t.Errorf("%v.sub64(%v) = %v, want %v", lhs, rhs, dif, bigdif)
			}
		}
	}
}

func TestUint192Twos(t *testing.T) {
	t.Parallel()

	initUintValues()

	one := big.NewInt(1)

	bigval := new(big.Int)
	tmpres := new(big.Int)

	for _, val := range uint192Values {
		res := val.twos()

		uint192ToBig(val, bigval)
		b := make([]byte, 24)
		c := bigval.Bytes()
		copy(b[24-len(c):], c)

		for i := range b {
			b[i] = ^b[i]
		}

		bigval.SetBytes(b)
		bigres := bigval.Add(bigval, one)
		if bigres.BitLen() > 192 {
			b = bigres.Bytes()
			bigres.SetBytes(b[len(b)-24:])
		}

		if uint192ToBig(res, tmpres).Cmp(bigres) != 0 {
			t.Errorf("%v.twos() = %v, want %v", val, res, bigres)
		}
	}
}

func TestUint256Div10(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint256Values {
		quo, rem := val.div10()

		uint256ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint256ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint256Div10000(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint256Values {
		quo, rem := val.div10000()

		uint256ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint256ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10000() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint256Div1e8(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(100_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint256Values {
		quo, rem := val.div1e8()

		uint256ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint256ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e8() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint256Div1e19(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := new(big.Int).SetUint64(10_000_000_000_000_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint256Values {
		quo, rem := val.div1e19()

		uint256ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint256ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e19() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint256Lsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint256Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.lsh(rhs)

			uint256ToBig(lhs, biglhs)
			bigres := biglhs.Lsh(biglhs, rhs)
			if bigres.BitLen() > 256 {
				b := bigres.Bytes()
				bigres.SetBytes(b[len(b)-32:])
			}

			if uint256ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.lsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint256Mul64(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	bigrhs := new(big.Int)
	tmpprd := new(big.Int)

	for _, lhs := range uint256Values {
		for _, rhs := range uint64Values {
			prd := lhs.mul64(rhs)

			uint256ToBig(lhs, biglhs)
			bigrhs.SetUint64(rhs)
			bigprd := biglhs.Mul(biglhs, bigrhs)
			if bigprd.BitLen() > 256 {
				b := bigprd.Bytes()
				bigprd.SetBytes(b[len(b)-32:])
			}

			if uint256ToBig(prd, tmpprd).Cmp(bigprd) != 0 {
				t.Errorf("%v.mul64(%d) = %v, want %v", lhs, rhs, prd, bigprd)
			}
		}
	}
}

func TestUint256Rsh(t *testing.T) {
	t.Parallel()

	initUintValues()

	biglhs := new(big.Int)
	tmpres := new(big.Int)

	for _, lhs := range uint256Values {
		for rhs := uint(0); rhs < 200; rhs += 10 {
			res := lhs.rsh(rhs)

			uint256ToBig(lhs, biglhs)
			bigres := biglhs.Rsh(biglhs, rhs)

			if uint256ToBig(res, tmpres).Cmp(bigres) != 0 {
				t.Errorf("%v.rsh(%d) = %v, want %v", lhs, rhs, res, bigres)
			}
		}
	}
}

func TestUint256String(t *testing.T) {
	t.Parallel()

	initUintValues()

	bigval := new(big.Int)

	for _, val := range uint256Values {
		res := val.String()

		uint256ToBig(val, bigval)
		bigres := bigval.String()

		if res != bigres {
			t.Errorf("%v.String() = %s, want %s", val, res, bigres)
		}
	}
}

func TestUint384Div10(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := big.NewInt(10)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint384Values {
		quo, rem := val.div10()

		uint384ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint384ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div10() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}

func TestUint384Div1e19(t *testing.T) {
	t.Parallel()

	initUintValues()

	c := new(big.Int).SetUint64(10_000_000_000_000_000_000)

	bigval := new(big.Int)
	bigrem := new(big.Int)
	tmpquo := new(big.Int)
	tmprem := new(big.Int)

	for _, val := range uint384Values {
		quo, rem := val.div1e19()

		uint384ToBig(val, bigval)
		bigquo, bigrem := bigval.QuoRem(bigval, c, bigrem)

		if uint384ToBig(quo, tmpquo).Cmp(bigquo) != 0 || tmprem.SetUint64(rem).Cmp(bigrem) != 0 {
			t.Errorf("%v.div1e19() = (%v, %v), want (%v, %v)", val, quo, rem, bigquo, bigrem)
		}
	}
}
