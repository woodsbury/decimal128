package decimal128

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

var textValues = map[string]Decimal{
	"0":  zero(false),
	"+0": zero(false),
	"-0": zero(true),
	"0.00000000000000000000000000000000000005e-99999999": zero(false),
	"0.00000000000000000000000000000000000005e-60150":    zero(false),
	"0.00000000000000000000000000000000000005e-6150":     zero(false),
	"1":              compose(false, uint128{1, 0}, exponentBias),
	"-1e1":           compose(true, uint128{10, 0}, exponentBias),
	"-1_23e1_0":      compose(true, uint128{123, 0}, exponentBias+10),
	"00123.45600e10": compose(false, uint128{123456, 0}, exponentBias+7),
	"123.456e-10":    compose(false, uint128{123456, 0}, exponentBias-13),
	"0.00000000000000000000000000000000000005e-30": compose(false, uint128{5, 0}, exponentBias-68),
	"0.00000000000000000000000000000000000005":     compose(false, uint128{5, 0}, exponentBias-38),
	"0.00000000000000000000000000000000000005e30":  compose(false, uint128{5, 0}, exponentBias-8),
	"500000000000000000000000000000000000000e-30":  compose(false, uint128{5, 0}, exponentBias+8),
	"500000000000000000000000000000000000000":      compose(false, uint128{5, 0}, exponentBias+38),
	"500000000000000000000000000000000000000e30":   compose(false, uint128{5, 0}, exponentBias+68),
	"inf":         inf(false),
	"+Inf":        inf(false),
	"-Infinity":   inf(true),
	"nan":         NaN(),
	"NaN":         NaN(),
	"0e99999999":  zero(false),
	"-0e99999999": zero(true),
}

func TestDecimalScan(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		if val == "-Infinity" {
			// Scan only accepts "Inf" for infinity
			continue
		}

		var res Decimal
		_, err := fmt.Sscan(val, &res)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("fmt.Sscan(%q) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}

		_, err = fmt.Sscanf(val, "%g", &res)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("fmt.Sscanf(%q, \"%%g\") = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}

		val += "\n"

		_, err = fmt.Sscanf(val, "%g\n", &res)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("fmt.Sscanf(%q, \"%%g\\n\") = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}

		_, err = fmt.Sscanln(val, &res)
		if !(res.Equal(num) || res.IsNaN() && num.IsNaN()) || res.Signbit() != num.Signbit() || err != nil {
			t.Errorf("fmt.Sscanln(%q) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}

	var res1, res2 Decimal
	n, err := fmt.Sscanf("1 2", "%g %g", &res1, &res2)
	if n != 2 || !res1.Equal(FromUint64(1)) || !res2.Equal(FromUint64(2)) || err != nil {
		t.Errorf("fmt.Sscanf(\"1 2\") = (%v, %v, %v), want (1, 2, <nil>)", res1, res2, err)
	}

	n, err = fmt.Sscanln("1 2\n", &res1, &res2)
	if n != 2 || !res1.Equal(FromUint64(1)) || !res2.Equal(FromUint64(2)) || err != nil {
		t.Errorf("fmt.Sscanln(\"1 2\\n\") = (%v, %v, %v), want (1, 2, <nil>)", res1, res2, err)
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

	res, err := Parse("1e999999999")
	if !res.IsInf(1) || !errors.Is(err, strconv.ErrRange) {
		t.Errorf("Parse(1e999999999) = (%v, %v), want (Inf, value out of range)", res, err)
	}

	res, err = Parse("-1e999999999")
	if !res.IsInf(-1) || !errors.Is(err, strconv.ErrRange) {
		t.Errorf("Parse(-1e999999999) = (%v, %v), want (-Inf, value out of range)", res, err)
	}
}

func FuzzParse(f *testing.F) {
	f.Add("123_456.789e10")
	f.Add("+Inf")
	f.Add("-Inf")
	f.Add("NaN")

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()

		Parse(s)
	})
}
