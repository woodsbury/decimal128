package decimal128

import (
	"fmt"
	"testing"
)

func TestDecimalMarshalBinary(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		res, err := decval.MarshalBinary()

		fmtval := fmt.Sprintf("%.16x%.16x", decval.hi, decval.lo)
		fmtres := fmt.Sprintf("%x", res)

		if fmtres != fmtval || err != nil {
			t.Errorf("%v.MarshalBinary() = (%s, %v), want (%s, <nil>)", val, fmtres, err, fmtval)
		}

		var resval Decimal
		err = resval.UnmarshalBinary(res)

		if resval != decval || err != nil {
			t.Errorf("Decimal.UnmarshalBinary(%x) = (%v, %v), want (%v, <nil>)", res, resval, err, decval)
		}

		res, err = decval.AppendBinary(res[:0])

		fmtres = fmt.Sprintf("%x", res)

		if fmtres != fmtval || err != nil {
			t.Errorf("%v.AppendBinary() = (%s, %v), want (%s, <nil>)", val, fmtres, err, fmtval)
		}

		err = resval.UnmarshalBinary(res)

		if resval != decval || err != nil {
			t.Errorf("Decimal.UnmarshalBinary(%x) = (%v, %v), want (%v, <nil>)", res, resval, err, decval)
		}
	}
}

func BenchmarkDecimalAppendBinary(b *testing.B) {
	b.ReportAllocs()

	d := New(123456789, 10)
	var buf []byte

	for i := 0; i < b.N; i++ {
		buf, _ = d.AppendBinary(buf[:0])
	}
}

func BenchmarkDecimalMarshalBinary(b *testing.B) {
	b.ReportAllocs()

	d := New(123456789, 10)

	for i := 0; i < b.N; i++ {
		d.MarshalBinary()
	}
}
