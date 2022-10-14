package decimal128

import (
	"fmt"
	"testing"
)

var textValues = map[string]Decimal{
	"0":  zero(false),
	"+0": zero(false),
	"-0": zero(true),
	"0.00000000000000000000000000000000000005e-6150": zero(false),
	"1":              compose(false, uint128{1, 0}, exponentBias),
	"-1e1":           compose(true, uint128{10, 0}, exponentBias),
	"-1_23e1_0":      compose(true, uint128{123, 0}, exponentBias+10),
	"00123.45600e10": compose(false, uint128{123456, 0}, exponentBias+7),
	"123.456e-10":    compose(false, uint128{123456, 0}, exponentBias-13),
	"0.00000000000000000000000000000000000005":     compose(false, uint128{5, 0}, exponentBias-38),
	"500000000000000000000000000000000000000":      compose(false, uint128{5, 0}, exponentBias+38),
	"500000000000000000000000000000000000000e6150": inf(false),
	"inf":       inf(false),
	"+Inf":      inf(false),
	"-Infinity": inf(true),
	"nan":       NaN(),
	"NaN":       NaN(),
}

func TestDecimalScan(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		var res Decimal
		_, err := fmt.Sscan(val, &res)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("fmt.Sscan(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}
}

func TestDecimalUnmarshalText(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		var res Decimal
		err := res.UnmarshalText([]byte(val))
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("Decimal.UnmarshalText(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		res, err := Parse(val)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("Parse(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}
}

func FuzzParse(f *testing.F) {
	f.Add("123_456.789e10")
	f.Add("+Inf")
	f.Add("-Inf")
	f.Add("NaN")

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()

		dec, err := Parse(s)
		if err != nil {
			return
		}

		_ = dec.String()
	})
}
