package decimal128

import (
	"fmt"
	"math"
	"testing"
)

func TestDecimalPayload(t *testing.T) {
	t.Parallel()

	if s := NaN().Payload().String(); s != "NaN()" {
		t.Errorf("NaN().Payload() = %s, want NaN()", s)
	}

	var d Decimal
	if err := d.Compose(2, false, nil, 0); err != nil {
		t.Errorf("Decimal.Compose(2, false, nil, 0) = %v, want <nil>", err)
	} else if s := d.Payload().String(); s != "Compose()" {
		t.Errorf("Decimal.Compose(2, false, nil, 0).Payload() = %s, want Compose()", s)
	}

	d = FromFloat32(float32(math.NaN()))
	if s := d.Payload().String(); s != "FromFloat32()" {
		t.Errorf("FromFloat32(NaN).Payload() = %s, want FromFloat32()", s)
	}

	d = FromFloat64(math.NaN())
	if s := d.Payload().String(); s != "FromFloat64()" {
		t.Errorf("FromFloat64(NaN).Payload() = %s, want FromFloat64()", s)
	}

	d = MustParse("NaN")
	if s := d.Payload().String(); s != "MustParse()" {
		t.Errorf("MustParse(NaN).Payload() = %s, want MustParse()", s)
	}

	d, err := Parse("NaN")
	if err != nil {
		t.Errorf("Parse(NaN) = (%v, %v), want (NaN, <nil>)", d, err)
	} else if s := d.Payload().String(); s != "Parse()" {
		t.Errorf("Parse(NaN).Payload() = %s, want Parse()", s)
	}

	_, err = fmt.Sscan("NaN", &d)
	if err != nil {
		t.Errorf("fmt.Sscan(NaN) = (%v, %v), want (NaN, <nil>)", d, err)
	} else if s := d.Payload().String(); s != "Scan()" {
		t.Errorf("Scan(NaN).Payload() = %s, want Scan()", s)
	}

	err = d.UnmarshalText([]byte("NaN"))
	if err != nil {
		t.Errorf("Decimal.UnmarshalText(NaN) = (%v, %v), want (NaN, <nil>)", d, err)
	} else if s := d.Payload().String(); s != "UnmarshalText()" {
		t.Errorf("Decimal.UnmarshalText(NaN).Payload() = %s, want UnmarshalText()", s)
	}

	d = Log(inf(true))
	if s := d.Payload().String(); s != "Log(-Infinite)" {
		t.Errorf("Log(-Inf).Payload() = %s, want Log(-Infinite)", s)
	}

	d = Log(FromInt64(-1))
	if s := d.Payload().String(); s != "Log(-Finite)" {
		t.Errorf("Log(-1).Payload() = %s, want Log(-Finite)", s)
	}

	d = Log10(inf(true))
	if s := d.Payload().String(); s != "Log10(-Infinite)" {
		t.Errorf("Log10(-Inf).Payload() = %s, want Log10(-Infinite)", s)
	}

	d = Log10(FromInt64(-1))
	if s := d.Payload().String(); s != "Log10(-Finite)" {
		t.Errorf("Log10(-1).Payload() = %s, want Log10(-Finite)", s)
	}

	d = Log1p(inf(true))
	if s := d.Payload().String(); s != "Log1p(-Infinite)" {
		t.Errorf("Log1p(-Inf).Payload() = %s, want Log1p(-Infinite)", s)
	}

	d = Log1p(FromInt64(-2))
	if s := d.Payload().String(); s != "Log1p(-Finite)" {
		t.Errorf("Log1p(-2).Payload() = %s, want Log1p(-Finite)", s)
	}

	d = Log2(inf(true))
	if s := d.Payload().String(); s != "Log2(-Infinite)" {
		t.Errorf("Log2(-Inf).Payload() = %s, want Log2(-Infinite)", s)
	}

	d = Log2(FromInt64(-1))
	if s := d.Payload().String(); s != "Log2(-Finite)" {
		t.Errorf("Log2(-1).Payload() = %s, want Log2(-Finite)", s)
	}

	d = Sqrt(inf(true))
	if s := d.Payload().String(); s != "Sqrt(-Infinite)" {
		t.Errorf("Sqrt(-Inf).Payload() = %s, want Sqrt(-Infinite)", s)
	}

	d = inf(false).Add(inf(true))
	if s := d.Payload().String(); s != "Add(Infinite, -Infinite)" {
		t.Errorf("Inf.Add(-Inf).Payload() = %s, want Add(Infinite, -Infinite)", s)
	}

	d = inf(true).Add(inf(false))
	if s := d.Payload().String(); s != "Add(-Infinite, Infinite)" {
		t.Errorf("-Inf.Add(Inf).Payload() = %s, want Add(-Infinite, Infinite)", s)
	}

	d = zero(false).Mul(inf(false))
	if s := d.Payload().String(); s != "Mul(Zero, Infinite)" {
		t.Errorf("0.Mul(Inf).Payload() = %s, want Mul(Zero, Infinite)", s)
	}

	d = zero(true).Mul(inf(true))
	if s := d.Payload().String(); s != "Mul(-Zero, -Infinite)" {
		t.Errorf("-0.Mul(-Inf).Payload() = %s, want Mul(-Zero, -Infinite)", s)
	}

	d = inf(false).Mul(zero(false))
	if s := d.Payload().String(); s != "Mul(Infinite, Zero)" {
		t.Errorf("Inf.Mul(0).Payload() = %s, want Mul(Infinite, Zero)", s)
	}

	d = inf(true).Mul(zero(true))
	if s := d.Payload().String(); s != "Mul(-Infinite, -Zero)" {
		t.Errorf("-Inf.Mul(-0).Payload() = %s, want Mul(-Infinite, -Zero)", s)
	}

	d = FromInt64(-1).Pow(New(1, -1))
	if s := d.Payload().String(); s != "Pow(-Finite, Finite)" {
		t.Errorf("-1.Pow(0.1).Payload() = %s, want Pow(-Finite, Finite)", s)
	}

	d = FromInt64(-1).Pow(New(-1, -1))
	if s := d.Payload().String(); s != "Pow(-Finite, -Finite)" {
		t.Errorf("-1.Pow(-0.1).Payload() = %s, want Pow(-Finite, -Finite)", s)
	}

	d = inf(false).Quo(inf(true))
	if s := d.Payload().String(); s != "Quo(Infinite, -Infinite)" {
		t.Errorf("Inf.Quo(-Inf).Payload() = %s, want Quo(Infinite, -Infinite)", s)
	}

	d = inf(true).Quo(inf(false))
	if s := d.Payload().String(); s != "Quo(-Infinite, Infinite)" {
		t.Errorf("-Inf.Quo(Inf).Payload() = %s, want Quo(-Infinite, Infinite)", s)
	}

	d = zero(false).Quo(zero(true))
	if s := d.Payload().String(); s != "Quo(Zero, -Zero)" {
		t.Errorf("0.Quo(-0).Payload() = %s, want Quo(Zero, -Zero)", s)
	}

	d = zero(true).Quo(zero(false))
	if s := d.Payload().String(); s != "Quo(-Zero, Zero)" {
		t.Errorf("-0.Quo(0).Payload() = %s, want Quo(-Zero, Zero)", s)
	}

	d, _ = inf(false).QuoRem(inf(true))
	if s := d.Payload().String(); s != "QuoRem(Infinite, -Infinite)" {
		t.Errorf("Inf.QuoRem(-Inf).Payload() = %s, want QuoRem(Infinite, -Infinite)", s)
	}

	d, _ = inf(true).QuoRem(inf(false))
	if s := d.Payload().String(); s != "QuoRem(-Infinite, Infinite)" {
		t.Errorf("-Inf.QuoRem(Inf).Payload() = %s, want QuoRem(-Infinite, Infinite)", s)
	}

	_, d = inf(true).QuoRem(FromInt64(1))
	if s := d.Payload().String(); s != "QuoRem(-Infinite, Finite)" {
		t.Errorf("-Inf.QuoRem(Inf).Payload() = %s, want QuoRem(-Infinite, Finite)", s)
	}

	d, _ = zero(false).QuoRem(zero(true))
	if s := d.Payload().String(); s != "QuoRem(Zero, -Zero)" {
		t.Errorf("0.QuoRem(-0).Payload() = %s, want QuoRem(Zero, -Zero)", s)
	}

	d, _ = zero(true).QuoRem(zero(false))
	if s := d.Payload().String(); s != "QuoRem(-Zero, Zero)" {
		t.Errorf("-0.QuoRem(0).Payload() = %s, want QuoRem(-Zero, Zero)", s)
	}

	_, d = FromInt64(1).QuoRem(zero(false))
	if s := d.Payload().String(); s != "QuoRem(Finite, Zero)" {
		t.Errorf("-0.QuoRem(0).Payload() = %s, want QuoRem(Finite, Zero)", s)
	}

	_, d = FromInt64(1).QuoRem(zero(true))
	if s := d.Payload().String(); s != "QuoRem(Finite, -Zero)" {
		t.Errorf("-0.QuoRem(0).Payload() = %s, want QuoRem(Finite, -Zero)", s)
	}

	d = inf(false).Sub(inf(false))
	if s := d.Payload().String(); s != "Sub(Infinite, Infinite)" {
		t.Errorf("Inf.Sub(Inf).Payload() = %s, want Sub(Infinite, Infinite)", s)
	}

	d = inf(true).Sub(inf(true))
	if s := d.Payload().String(); s != "Sub(-Infinite, -Infinite)" {
		t.Errorf("-Inf.Sub(-Inf).Payload() = %s, want Sub(-Infinite, -Infinite)", s)
	}
}
