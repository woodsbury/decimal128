package decimal128

import (
	"math"
	"testing"
)

func TestCbrt(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("cbrt(%v) = %v\n", &val, &res) {
		root := Cbrt(val)

		if !resultEqual(root, res) {
			t.Errorf("Cbrt(%v) = %v, want %v", val, root, res)
		}
	}
}

func TestExp(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("exp(%v) = %v\n", &val, &res) {
		exp := Exp(val)

		if !resultEqual(exp, res) {
			t.Errorf("Exp(%v) = %v, want %v", val, exp, res)
		}
	}
}

func TestExp10(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("exp10(%v) = %v\n", &val, &res) {
		exp := Exp10(val)

		if !resultEqual(exp, res) {
			t.Errorf("Exp10(%v) = %v, want %v", val, exp, res)
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

func TestExpm1(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("expm1(%v) = %v\n", &val, &res) {
		exp := Expm1(val)

		if !resultEqual(exp, res) {
			t.Errorf("Expm1(%v) = %v, want %v", val, exp, res)
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

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("log10(%v) = %v\n", &val, &res) {
		log := Log10(val)

		if !resultEqual(log, res) {
			t.Errorf("Log10(%v) = %v, want %v", val, log, res)
		}
	}
}

func TestLog1p(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var res Decimal

	for r.scan("log1p(%v) = %v\n", &val, &res) {
		log := Log1p(val)

		if !resultEqual(log, res) {
			t.Errorf("Log1p(%v) = %v, want %v", val, log, res)
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

func FuzzExp(f *testing.F) {
	values := []float64{
		0.0,
		math.Copysign(0.0, -1.0),
		0.5,
		1.0,
		5.0,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
	}

	for _, v := range values {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		t.Parallel()

		decv := FromFloat64(v)

		decres := Exp(decv).Float64()
		fltres := math.Exp(v)
		eps := math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp := math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Exp10(decv).Float64()
		fltres = math.Pow(10.0, v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Exp2(decv).Float64()
		fltres = math.Exp2(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Expm1(decv).Float64()
		fltres = math.Expm1(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}
	})
}

func FuzzLog(f *testing.F) {
	values := []float64{
		0.0,
		math.Copysign(0.0, -1.0),
		0.5,
		1.0,
		5.0,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
	}

	for _, v := range values {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		t.Parallel()

		decv := FromFloat64(v)

		decres := Log(decv).Float64()
		fltres := math.Log(v)
		eps := math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp := math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Log10(decv).Float64()
		fltres = math.Log10(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Log1p(decv).Float64()
		fltres = math.Log1p(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Log2(decv).Float64()
		fltres = math.Log2(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Log(decv, Log2(decv))
			t.Logf("%.24e", v)
			t.Log(v, fltres, decres, eps)
			t.Fail()
		}
	})
}

func FuzzRoot(f *testing.F) {
	values := []float64{
		0.0,
		math.Copysign(0.0, -1.0),
		0.5,
		1.0,
		5.0,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
	}

	for _, v := range values {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		t.Parallel()

		decv := FromFloat64(v)

		decres := Cbrt(decv).Float64()
		fltres := math.Cbrt(v)
		eps := math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp := math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}

		decres = Sqrt(decv).Float64()
		fltres = math.Sqrt(v)
		eps = math.Nextafter(math.Abs(math.Max(decres, fltres)), math.Inf(1)) - math.Abs(math.Max(decres, fltres))
		eps, epsexp = math.Frexp(eps)
		eps = math.Ldexp(eps, epsexp+50)
		if math.Abs(decres-fltres) > eps {
			t.Fail()
		}
	})
}

func BenchmarkExp(b *testing.B) {
	initDecimalValues()

	decvals := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		decvals[i] = val.Decimal()
	}

	b.Run("Exp", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Exp(decval)
			}
		}
	})

	b.Run("Exp10", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Exp10(decval)
			}
		}
	})

	b.Run("Exp2", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Exp2(decval)
			}
		}
	})

	b.Run("Expm1", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Expm1(decval)
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
		for b.Loop() {
			for _, decval := range decvals {
				Log(decval)
			}
		}
	})

	b.Run("Log10", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Log10(decval)
			}
		}
	})

	b.Run("Log1p", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Log1p(decval)
			}
		}
	})

	b.Run("Log2", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Log2(decval)
			}
		}
	})
}

func BenchmarkRoot(b *testing.B) {
	initDecimalValues()

	decvals := make([]Decimal, len(decimalValues))
	for i, val := range decimalValues {
		decvals[i] = val.Decimal()
	}

	b.Run("Cbrt", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Cbrt(decval)
			}
		}
	})

	b.Run("Sqrt", func(b *testing.B) {
		for b.Loop() {
			for _, decval := range decvals {
				Sqrt(decval)
			}
		}
	})
}
