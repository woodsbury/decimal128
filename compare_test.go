package decimal128

import (
	"errors"
	"fmt"
	"testing"
)

type testCmpResult int8

func (tcr *testCmpResult) Scan(f fmt.ScanState, verb rune) error {
	if verb != 'v' {
		return errors.New("bad verb '%" + string(verb) + "' for testCmpResult")
	}

	tok, err := f.Token(true, nil)
	if err != nil {
		return err
	}

	if len(tok) != 1 {
		return errors.New("invalid value")
	}

	switch tok[0] {
	case '!':
		*tcr = -2
	case '<':
		*tcr = -1
	case '=':
		*tcr = 0
	case '>':
		*tcr = 1
	default:
		return errors.New("invalid value")
	}

	return nil
}

func (tcr testCmpResult) equal() bool {
	return tcr == 0
}

func (tcr testCmpResult) greater() bool {
	return tcr == 1
}

func (tcr testCmpResult) greaterOrEqual() bool {
	return tcr == 1 || tcr == 0
}

func (tcr testCmpResult) less() bool {
	return tcr == -1
}

func (tcr testCmpResult) lessOrEqual() bool {
	return tcr == -1 || tcr == 0
}

func TestDecimalCmp(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testCmpResult

	for r.scan("%v cmp %v = %v\n", &lhs, &rhs, &res) {
		cmp := lhs.Cmp(rhs)

		if ceq, req := cmp.Equal(), res.equal(); ceq != req {
			t.Errorf("%v.Cmp(%v).Equal() = %t, want %t", lhs, rhs, ceq, req)
		}

		if cgt, rgt := cmp.Greater(), res.greater(); cgt != rgt {
			t.Errorf("%v.Cmp(%v).Greater() = %t, want %t", lhs, rhs, cgt, rgt)
		}

		if cge, rge := cmp.GreaterOrEqual(), res.greaterOrEqual(); cge != rge {
			t.Errorf("%v.Cmp(%v).GreaterOrEqual() = %t, want %t", lhs, rhs, cge, rge)
		}

		if clt, rlt := cmp.Less(), res.less(); clt != rlt {
			t.Errorf("%v.Cmp(%v).Less() = %t, want %t", lhs, rhs, clt, rlt)
		}

		if cle, rle := cmp.LessOrEqual(), res.lessOrEqual(); cle != rle {
			t.Errorf("%v.Cmp(%v).LessOrEqual() = %t, want %t", lhs, rhs, cle, rle)
		}

		if ceq, req := lhs.Equal(rhs), res.equal(); ceq != req {
			t.Errorf("%v.Equal(%v) = %t, want %t", lhs, rhs, ceq, req)
		}
	}

	x := New(1, 0)
	ysig := int64(1_000_000_000)
	yexp := -9

	for yexp < 0 {
		y := New(ysig, yexp)

		cmp := x.Cmp(y)

		if !cmp.Equal() {
			t.Errorf("1.Cmp(%de%d).Equal() = false, want true", ysig, yexp)
		}

		if !x.Equal(y) {
			t.Errorf("1.Equal(%de%d) = false, want true", ysig, yexp)
		}

		cmp = y.Cmp(x)

		if !cmp.Equal() {
			t.Errorf("%de%d.Cmp(1).Equal() = false, want true", ysig, yexp)
		}

		if !y.Equal(x) {
			t.Errorf("%de%d.Equal(1) = false, want true", ysig, yexp)
		}

		ysig /= 10
		yexp++
	}
}

func TestDecimalCmpAbs(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res testCmpResult

	for r.scan("%v cmpabs %v = %v\n", &lhs, &rhs, &res) {
		cmp := lhs.CmpAbs(rhs)

		if ceq, req := cmp.Equal(), res.equal(); ceq != req {
			t.Errorf("%v.CmpAbs(%v).Equal() = %t, want %t", lhs, rhs, ceq, req)
		}

		if cgt, rgt := cmp.Greater(), res.greater(); cgt != rgt {
			t.Errorf("%v.CmpAbs(%v).Greater() = %t, want %t", lhs, rhs, cgt, rgt)
		}

		if cge, rge := cmp.GreaterOrEqual(), res.greaterOrEqual(); cge != rge {
			t.Errorf("%v.CmpAbs(%v).GreaterOrEqual() = %t, want %t", lhs, rhs, cge, rge)
		}

		if clt, rlt := cmp.Less(), res.less(); clt != rlt {
			t.Errorf("%v.CmpAbs(%v).Less() = %t, want %t", lhs, rhs, clt, rlt)
		}

		if cle, rle := cmp.LessOrEqual(), res.lessOrEqual(); cle != rle {
			t.Errorf("%v.CmpAbs(%v).LessOrEqual() = %t, want %t", lhs, rhs, cle, rle)
		}

		if ceq, req := Abs(lhs).Equal(Abs(rhs)), res.equal(); ceq != req {
			t.Errorf("Abs(%v).Equal(Abs(%v)) = %t, want %t", lhs, rhs, ceq, req)
		}
	}

	x := New(1, 0)
	ysig := int64(1_000_000_000)
	yexp := -9

	for yexp < 0 {
		y := New(ysig, yexp)

		cmp := x.CmpAbs(y)

		if !cmp.Equal() {
			t.Errorf("1.CmpAbs(%de%d).Equal() = false, want true", ysig, yexp)
		}

		cmp = y.CmpAbs(x)

		if !cmp.Equal() {
			t.Errorf("%de%d.CmpAbs(1).Equal() = false, want true", ysig, yexp)
		}

		ysig /= 10
		yexp++
	}
}

func TestDecimalIsZero(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		zero := decval.IsZero()

		var res bool
		if val.form == regularForm && val.sig[0]|val.sig[1] == 0 {
			res = true
		}

		if zero != res {
			t.Errorf("%v.IsZero() = %t, want %t", val, zero, res)
		}
	}
}

func TestMax(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res Decimal

	for r.scan("max(%v, %v) = %v\n", &lhs, &rhs, &res) {
		max := Max(lhs, rhs)

		if !resultEqual(max, res) {
			t.Errorf("Max(%v, %v) = %v, want %v", lhs, rhs, max, res)
		}
	}
}

func TestMin(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var lhs Decimal
	var rhs Decimal
	var res Decimal

	for r.scan("min(%v, %v) = %v\n", &lhs, &rhs, &res) {
		min := Min(lhs, rhs)

		if !resultEqual(min, res) {
			t.Errorf("Min(%v, %v) = %v, want %v", lhs, rhs, min, res)
		}
	}
}

func BenchmarkDecimalCmp(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		values[i] = val.Decimal()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, lhs := range values {
			for _, rhs := range values {
				lhs.Cmp(rhs)
			}
		}
	}
}

func BenchmarkDecimalCmpAbs(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		values[i] = val.Decimal()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, lhs := range values {
			for _, rhs := range values {
				lhs.CmpAbs(rhs)
			}
		}
	}
}

func BenchmarkDecimalEqual(b *testing.B) {
	initDecimalValues()

	values := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		values[i] = val.Decimal()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, lhs := range values {
			for _, rhs := range values {
				lhs.Equal(rhs)
			}
		}
	}
}
