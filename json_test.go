package decimal128

import (
	"encoding/json"
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
