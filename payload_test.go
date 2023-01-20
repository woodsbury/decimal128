package decimal128

import (
	"fmt"
	"testing"
)

func TestDecimalPayload(t *testing.T) {
	t.Parallel()

	if s := NaN().Payload().String(); s != "NaN()" {
		t.Errorf("NaN().Payload() = %s, want NaN()", s)
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
		t.Errorf("UnmarshalText(NaN).Payload() = %s, want UnmarshalText()", s)
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

	d = inf(false).Sub(inf(false))
	if s := d.Payload().String(); s != "Sub(Infinite, Infinite)" {
		t.Errorf("Inf.Sub(Inf).Payload() = %s, want Sub(Infinite, Infinite)", s)
	}

	d = inf(true).Sub(inf(true))
	if s := d.Payload().String(); s != "Sub(-Infinite, -Infinite)" {
		t.Errorf("-Inf.Sub(-Inf).Payload() = %s, want Sub(-Infinite, -Infinite)", s)
	}
}
