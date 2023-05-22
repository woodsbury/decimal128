package decimal128

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDecimalMarshalJSON(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		decval := val.Decimal()
		res, err := decval.MarshalJSON()

		if err != nil {
			t.Errorf("%v.MarshalJSON() = (%s, %v), want (%s, <nil>", val, res, err, res)
		}

		var resval Decimal
		err = resval.UnmarshalJSON(res)

		if !resval.Equal(decval) || err != nil {
			t.Errorf("Decimal.UnmarshalJSON(%s) = (%v, %v), want (%v, <nil>)", res, resval, err, decval)
		}

		if fltval, ok := val.Float64(); ok {
			fltres, err := json.Marshal(fltval)

			if err != nil {
				t.Errorf("json.Marshal(%v) = (%s, %v), want (%s, <nil>)", fltval, fltres, err, res)
			}

			if string(fltres) != string(res) {
				t.Errorf("%v.MarshalJSON() = (%s, <nil>), want (%s, <nil>)", val, res, fltres)
			}
		}
	}
}

func TestDecimalUnmarshalJSON(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		var res Decimal
		err := res.UnmarshalJSON([]byte(val))

		if num.isInf() || num.IsNaN() || strings.Contains(val, "_") {
			if err == nil {
				t.Errorf("Decimal.UnarshalJSON(%s) = (0, <nil>), want (%v, cannot unmarshal)", val, res)
			}
		} else if !res.Equal(num) || err != nil {
			t.Errorf("Decimal.UnmarshalJSON(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}

	num := New(123, -1)
	res := num
	err := res.UnmarshalJSON([]byte("null"))
	if !res.Equal(num) || err != nil {
		t.Errorf("Decimal.UnmarshalJSON(null) = (%v, %v), want (%v, <nil>)", res, err, num)
	}

	res = num
	err = res.UnmarshalJSON(nil)
	if !res.Equal(num) || err != nil {
		t.Errorf("Decimal.UnmarshalJSON() = (%v, %v), want (%v, <nil>)", res, err, num)
	}
}

func FuzzDecimalUnmarshalJSON(f *testing.F) {
	f.Add([]byte("123456.789e10"))

	f.Fuzz(func(t *testing.T, data []byte) {
		t.Parallel()

		var dec Decimal
		if err := dec.UnmarshalJSON(data); err != nil {
			return
		}

		_, _ = dec.MarshalJSON()
	})
}
