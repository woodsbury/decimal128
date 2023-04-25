package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func TestDecimalCmp(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	biglhs := new(apd.Decimal)
	bigrhs := new(apd.Decimal)

	for _, lhs := range decimalValues {
		for _, rhs := range decimalValues {
			declhs := lhs.Decimal()
			decrhs := rhs.Decimal()
			res := declhs.Cmp(decrhs)

			if lhs.form == nanForm || rhs.form == nanForm {
				if res.Equal() {
					t.Errorf("%v.Cmp(%v).Equal() = %t, want %t", lhs, rhs, true, false)
				}

				if res.Greater() {
					t.Errorf("%v.Cmp(%v).Greater() = %t, want %t", lhs, rhs, true, false)
				}

				if res.Less() {
					t.Errorf("%v.Cmp(%v).Less() = %t, want %t", lhs, rhs, true, false)
				}
			} else {
				lhs.Big(biglhs)
				rhs.Big(bigrhs)

				bigres := biglhs.Cmp(bigrhs)

				if (bigres == 0) != res.Equal() {
					t.Errorf("%v.Cmp(%v).Equal() == %t, want %t", lhs, rhs, res.Equal(), bigres == 0)
				}

				if (bigres > 0) != res.Greater() {
					t.Errorf("%v.Cmp(%v).Greater() = %t, want %t", lhs, rhs, res.Greater(), bigres > 0)
				}

				if (bigres < 0) != res.Less() {
					t.Errorf("%v.Cmp(%v).Less() = %t, want %t", lhs, rhs, res.Less(), bigres < 0)
				}
			}
		}
	}
}

func TestDecimalCmpAbs(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	biglhs := new(apd.Decimal)
	bigrhs := new(apd.Decimal)

	for _, lhs := range decimalValues {
		for _, rhs := range decimalValues {
			declhs := lhs.Decimal()
			decrhs := rhs.Decimal()
			res := declhs.CmpAbs(decrhs)

			if lhs.form == nanForm || rhs.form == nanForm {
				if res.Equal() {
					t.Errorf("%v.CmpAbs(%v).Equal() = %t, want %t", lhs, rhs, true, false)
				}

				if res.Greater() {
					t.Errorf("%v.CmpAbs(%v).Greater() = %t, want %t", lhs, rhs, true, false)
				}

				if res.Less() {
					t.Errorf("%v.CmpAbs(%v).Less() = %t, want %t", lhs, rhs, true, false)
				}
			} else {
				lhs.Big(biglhs)
				rhs.Big(bigrhs)

				biglhs.Abs(biglhs)
				bigrhs.Abs(bigrhs)

				bigres := biglhs.Cmp(bigrhs)

				if (bigres == 0) != res.Equal() {
					t.Errorf("%v.CmpAbs(%v).Equal() == %t, want %t", lhs, rhs, res.Equal(), bigres == 0)
				}

				if (bigres > 0) != res.Greater() {
					t.Errorf("%v.CmpAbs(%v).Greater() = %t, want %t", lhs, rhs, res.Greater(), bigres > 0)
				}

				if (bigres < 0) != res.Less() {
					t.Errorf("%v.CmpAbs(%v).Less() = %t, want %t", lhs, rhs, res.Less(), bigres < 0)
				}
			}
		}
	}
}

func TestDecimalEqual(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	biglhs := new(apd.Decimal)
	bigrhs := new(apd.Decimal)

	for _, lhs := range decimalValues {
		for _, rhs := range decimalValues {
			declhs := lhs.Decimal()
			decrhs := rhs.Decimal()
			res := declhs.Equal(decrhs)

			var bigres bool
			if lhs.form == nanForm || rhs.form == nanForm {
				bigres = false
			} else {
				lhs.Big(biglhs)
				rhs.Big(bigrhs)

				bigres = biglhs.Cmp(bigrhs) == 0
			}

			if res != bigres {
				t.Errorf("%v.Equal(%v) = %t, want %t", lhs, rhs, res, bigres)
			}
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
