package decimal128

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func TestExp(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   49,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Exp(decval)

		val.Big(bigval)

		bigctx.Exp(bigres, bigval)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Exp(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestExp10(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   42,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	bigln10 := new(apd.Decimal)
	bigctx.Ln(bigln10, apd.New(10, 0))

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Exp10(decval)

		val.Big(bigval)

		bigctx.Mul(bigres, bigval, bigln10)
		bigctx.Exp(bigres, bigres)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Exp10(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestExp2(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("exp2(%v) = %v\n", &val, &res) {
		exp := Exp2(val)

		if !resultEqual(exp, res) {
			t.Errorf("Exp2(%v) = %v, want %v", val, exp, res)
		}
	}
}

func TestLog(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("log(%v) = %v\n", &val, &res) {
		log := Log(val)

		if !resultEqual(log, res) {
			t.Errorf("Log(%v) = %v, want %v", val, log, res)
		}
	}
}

func TestLog10(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)
	bigres := new(apd.Decimal)
	bigctx := apd.Context{
		Precision:   39,
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    apd.RoundHalfEven,
	}

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Log10(decval)

		val.Big(bigval)

		bigctx.Log10(bigres, bigval)

		if !decimalsEqual(res, bigres, bigctx.Rounding) {
			t.Errorf("Log10(%v) = %v, want %v", val, res, bigres)
		}
	}
}

func TestLog2(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("log2(%v) = %v\n", &val, &res) {
		log := Log2(val)

		if !resultEqual(log, res) {
			t.Errorf("Log2(%v) = %v, want %v", val, log, res)
		}
	}
}

func TestSqrt(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("sqrt(%v) = %v\n", &val, &res) {
		root := Sqrt(val)

		if !resultEqual(root, res) {
			t.Errorf("Sqrt(%v) = %v, want %v", val, root, res)
		}
	}
}

func BenchmarkExp(b *testing.B) {
	initDecimalValues()

	decvals := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		decvals[i] = val.Decimal()
	}

	b.Run("Exp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Exp(decval)
			}
		}
	})

	b.Run("Exp10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Exp10(decval)
			}
		}
	})

	b.Run("Exp2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Exp2(decval)
			}
		}
	})
}

func BenchmarkLog(b *testing.B) {
	initDecimalValues()

	decvals := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		decvals[i] = val.Decimal()
	}

	b.Run("Log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Log(decval)
			}
		}
	})

	b.Run("Log10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Log10(decval)
			}
		}
	})

	b.Run("Log2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, decval := range decvals {
				Log2(decval)
			}
		}
	})
}
