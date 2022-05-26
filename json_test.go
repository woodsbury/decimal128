package decimal128

import "testing"

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
	}
}
