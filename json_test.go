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
		if num.isInf() || num.IsNaN() {
			continue
		}

		if strings.Contains(val, "_") {
			continue
		}

		var res Decimal
		err := res.UnmarshalJSON([]byte(val))
		if !res.Equal(num) || err != nil {
			t.Errorf("Decimal.UnmarshalJSON(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
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

func BenchmarkDecimalUnmarshalJSON(b *testing.B) {
	texts := [][]byte{
		[]byte("0.00000000000000000000000000000000000005"),
		[]byte("0.00000000000000000000000000000000000005e30"),
		[]byte("500000000000000000000000000000000000000e30"),
		[]byte("0.00000000000000000000000000000000000005e-6150"),
		[]byte("0.00000000000000000000000000000000000005e-999999"),
		[]byte("0e99999999"),
		[]byte("0"),
		[]byte("+1234567890"),
		[]byte("00123.45600e10"),
		[]byte("123.456e-10"),
	}

	for _, text := range texts {
		b.Run(string(text), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var res Decimal
				err := res.UnmarshalJSON(text)
				if err != nil {
					b.Errorf("failed to scan %q: %v", text, err)
				}
			}
		})
	}
}
