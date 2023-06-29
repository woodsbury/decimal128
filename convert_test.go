package decimal128

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestDecimalFloat(t *testing.T) {
	t.Parallel()

	values := []float64{
		math.Inf(-1),
		-5.0e100,
		-5.0e50,
		-5.0,
		-1.0,
		-5.0e-50,
		-5.0e-100,
		math.Copysign(0.0, -1.0),
		0.0,
		5.0e-100,
		5.0e-50,
		1.0,
		5.0,
		5.0e50,
		5.0e100,
		math.Inf(1),
	}

	bigval := new(big.Float)
	res := new(big.Float)

	for _, val := range values {
		bigval.SetFloat64(val)
		dec := FromFloat(bigval)
		dec.Float(res)

		valstr := fmt.Sprintf("%.17e", bigval)
		resstr := fmt.Sprintf("%.17e", res)

		if resstr != valstr {
			t.Errorf("%v.Float() = %v, want %v", val, res, bigval)
		}
	}
}

func TestDecimalFloat64(t *testing.T) {
	t.Parallel()

	values := []float64{
		math.Inf(-1),
		-5.0e300,
		-5.0e100,
		-5.0e50,
		-5.0,
		-1.0,
		-5.0e-50,
		-5.0e-100,
		-5.0e-300,
		math.Copysign(0.0, -1.0),
		0.0,
		5.0e-300,
		5.0e-100,
		5.0e-50,
		1.0,
		5.0,
		5.0e50,
		5.0e100,
		5.0e300,
		math.Inf(1),
		math.NaN(),
	}

	for _, val := range values {
		dec := FromFloat64(val)
		res := dec.Float64()

		valstr := fmt.Sprintf("%.17e", val)
		resstr := fmt.Sprintf("%.17e", res)

		if resstr != valstr {
			t.Errorf("%v.Float64() = %v, want %v", val, res, val)
		}
	}
}

func TestDecimalFloat32(t *testing.T) {
	t.Parallel()

	values := []float32{
		float32(math.Inf(-1)),
		-5.0e30,
		-5.0,
		-1.0,
		-5.0e-30,
		float32(math.Copysign(0.0, -1.0)),
		0.0,
		5.0e-30,
		1.0,
		5.0,
		5.0e30,
		float32(math.Inf(1)),
		float32(math.NaN()),
	}

	for _, val := range values {
		dec := FromFloat32(val)
		res := dec.Float32()

		valstr := fmt.Sprintf("%.9e", val)
		resstr := fmt.Sprintf("%.9e", res)

		if resstr != valstr {
			t.Errorf("%v.Float32() = %v, want %v", val, res, val)
		}
	}
}

func TestDecimalInt(t *testing.T) {
	t.Parallel()

	half := New(5, -1)

	values := [][]big.Word{
		{},
		{1},
		{0, 1},
		{0, 12118780946182832128, 10241496150359294728, 9566998059167698496, 12476541910035139367, 4681}, // 1e100
	}

	bigval := new(big.Int)
	res := new(big.Int)

	for _, val := range values {
		bigval.SetBits(val)
		dec := FromInt(bigval)
		dec.Int(res)

		if bigval.Cmp(res) != 0 {
			t.Errorf("%v.Int() = %v, want %v", dec, res, bigval)
		}

		dec = dec.Add(half)
		dec.Int(res)

		if bigval.Cmp(res) != 0 {
			t.Errorf("%v.Int() = %v, want %v", dec, res, bigval)
		}

		bigval.Neg(bigval)
		dec = FromInt(bigval)
		dec.Int(res)

		if bigval.Cmp(res) != 0 {
			t.Errorf("%v.Int() = %v, want %v", dec, res, bigval)
		}

		dec = dec.Sub(half)
		dec.Int(res)

		if bigval.Cmp(res) != 0 {
			t.Errorf("%v.Int() = %v, want %v", dec, res, bigval)
		}
	}
}

func TestDecimalInt32(t *testing.T) {
	t.Parallel()

	half := New(5, -1)

	values := []int32{
		math.MinInt32,
		-5,
		-1,
		0,
		1,
		5,
		math.MaxInt32,
	}

	for _, val := range values {
		dec := FromInt32(val)
		res, ok := dec.Int32()

		if res != val || !ok {
			t.Errorf("%v.Int32() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}

		if val >= 0 {
			dec = dec.Add(half)
		} else {
			dec = dec.Sub(half)
		}

		res, ok = dec.Int32()

		if res != val || !ok {
			t.Errorf("%v.Int32() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}
	}
}

func TestDecimalInt64(t *testing.T) {
	t.Parallel()

	half := New(5, -1)

	values := []int64{
		math.MinInt64,
		math.MinInt32,
		-5,
		-1,
		0,
		1,
		5,
		math.MaxInt32,
		math.MaxInt64,
	}

	for _, val := range values {
		dec := FromInt64(val)
		res, ok := dec.Int64()

		if res != val || !ok {
			t.Errorf("%v.Int64() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}

		if val >= 0 {
			dec = dec.Add(half)
		} else {
			dec = dec.Sub(half)
		}

		res, ok = dec.Int64()

		if res != val || !ok {
			t.Errorf("%v.Int64() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}
	}
}

func TestDecimalRat(t *testing.T) {
	t.Parallel()

	values := [][2][]big.Word{
		{
			{},
			{1},
		},
		{
			{1},
			{1},
		},
		{
			{0, 1},
			{1},
		},
		{
			{0, 12118780946182832128, 10241496150359294728, 9566998059167698496, 12476541910035139367, 4681}, // 1e100
			{1},
		},
		{
			{1},
			{0, 12118780946182832128, 10241496150359294728, 9566998059167698496, 12476541910035139367, 4681}, // 1e100
		},
	}

	bignum := new(big.Int)
	bigdenom := new(big.Int)
	bigval := new(big.Rat)
	res := new(big.Rat)

	for _, val := range values {
		bignum.SetBits(val[0])
		bigdenom.SetBits(val[1])
		bigval.SetFrac(bignum, bigdenom)
		dec := FromRat(bigval)
		dec.Rat(res)

		if res.Cmp(bigval) != 0 {
			t.Errorf("%v.Rat() = %v, want %v", dec, res, bigval)
		}

		bigval.Neg(bigval)
		dec = FromRat(bigval)
		dec.Rat(res)

		if res.Cmp(bigval) != 0 {
			t.Errorf("%v.Rat() = %v, want %v", dec, res, bigval)
		}
	}
}

func TestDecimalUint32(t *testing.T) {
	t.Parallel()

	half := New(5, -1)

	values := []uint32{
		0,
		1,
		5,
		math.MaxUint32,
	}

	for _, val := range values {
		dec := FromUint32(val)
		res, ok := dec.Uint32()

		if res != val || !ok {
			t.Errorf("%v.Uint32() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}

		dec = dec.Add(half)
		res, ok = dec.Uint32()

		if res != val || !ok {
			t.Errorf("%v.Uint32() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}
	}
}

func TestDecimalUint64(t *testing.T) {
	t.Parallel()

	half := New(5, -1)

	values := []uint64{
		0,
		1,
		5,
		math.MaxUint32,
		math.MaxUint64,
	}

	for _, val := range values {
		dec := FromUint64(val)
		res, ok := dec.Uint64()

		if res != val || !ok {
			t.Errorf("%v.Uint64() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}

		dec = dec.Add(half)
		res, ok = dec.Uint64()

		if res != val || !ok {
			t.Errorf("%v.Uint32() = (%d, %t), want (%d, true)", dec, res, ok, val)
		}
	}
}
