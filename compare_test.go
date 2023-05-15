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

func (tcr testCmpResult) less() bool {
	return tcr == -1
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

		if clt, rlt := cmp.Less(), res.less(); clt != rlt {
			t.Errorf("%v.Cmp(%v).Less() = %t, want %t", lhs, rhs, clt, rlt)
		}

		if ceq, req := lhs.Equal(rhs), res.equal(); ceq != req {
			t.Errorf("%v.Equal(%v) = %t, want %t", lhs, rhs, ceq, req)
		}
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

		if clt, rlt := cmp.Less(), res.less(); clt != rlt {
			t.Errorf("%v.CmpAbs(%v).Less() = %t, want %t", lhs, rhs, clt, rlt)
		}

		if ceq, req := Abs(lhs).Equal(Abs(rhs)), res.equal(); ceq != req {
			t.Errorf("Abs(%v).Equal(Abs(%v)) = %t, want %t", lhs, rhs, ceq, req)
		}
	}
}

func TestDecimalIsZero(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		zero := decval.IsZero()

		var res bool
		if val.form == regularForm && val.sig == (uint128{}) {
			res = true
		}

		if zero != res {
			t.Errorf("%v.IsZero() = %t, want %t", val, zero, res)
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
				lhs.Cmp(rhs)
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
