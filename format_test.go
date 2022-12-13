package decimal128

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

type testFmt struct {
	verb  rune
	flags []rune
	width int
	prec  int
}

func (tf testFmt) String() string {
	builder := strings.Builder{}
	builder.WriteRune('%')

	for _, flag := range tf.flags {
		builder.WriteRune(flag)
	}

	if tf.width != 0 {
		builder.WriteString(strconv.Itoa(tf.width))
	}

	if tf.prec != 0 {
		builder.WriteRune('.')
		builder.WriteString(strconv.Itoa(tf.prec))
	}

	builder.WriteRune(tf.verb)

	return builder.String()
}

func TestDecimalFormat(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	flags := []rune{' ', '#', '+', '-', '0'}

	var flagSets [][]rune
	for i := 0b00000; i <= 0b11111; i++ {
		var flagSet []rune
		for j := 0; j < 5; j++ {
			if i&(1<<j) != 0 {
				flagSet = append(flagSet, flags[j])
			}
		}

		flagSets = append(flagSets, flagSet)
	}

	var formats []testFmt
	for _, width := range []int{0, 15} {
		for _, prec := range []int{0, 15} {
			for _, flagSet := range flagSets {
				formats = append(formats,
					testFmt{'e', flagSet, width, prec},
					testFmt{'E', flagSet, width, prec},
					testFmt{'f', flagSet, width, prec},
					testFmt{'F', flagSet, width, prec},
					testFmt{'g', flagSet, width, prec},
					testFmt{'G', flagSet, width, prec},
				)
			}
		}
	}

	for _, format := range formats {
		for _, val := range decimalValues {
			decval := val.Decimal()
			fmtstr := format.String()
			res := fmt.Sprintf(fmtstr, decval)

			if strings.Contains(res, "PANIC") {
				t.Errorf("fmt.Sprintf(%v, %v) = %s", format, val, res)
				continue
			}

			var bigres string
			if val.form == regularForm {
				if fltval, ok := val.Float64(); ok {
					bigres = fmt.Sprintf(fmtstr, fltval)
				} else {
					// Skipping for now
					continue
				}
			} else if val.form == infForm {
				if val.neg {
					bigres = fmt.Sprintf(fmtstr, math.Inf(-1))
				} else {
					bigres = fmt.Sprintf(fmtstr, math.Inf(1))
				}
			} else {
				bigres = fmt.Sprintf(fmtstr, math.NaN())
			}

			if res != bigres {
				t.Errorf("fmt.Sprintf(%v, %v) = %s, want %s", format, val, res, bigres)
			}
		}
	}
}

func TestDecimalMarshalText(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		res, err := decval.MarshalText()

		if err != nil {
			t.Errorf("%v.MarshalText() = (%s, %v), want (%s, <nil>)", val, res, err, res)
		}

		var resval Decimal
		err = resval.UnmarshalText(res)

		if !(resval.Equal(decval) || resval.IsNaN() && decval.IsNaN()) || err != nil {
			t.Errorf("Decimal.UnmarshalText(%s) = (%v, %v), want (%v, <nil>)", res, resval, err, decval)
		}
	}
}

func TestDecimalString(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigval := new(apd.Decimal)

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := decval.String()

		var bigres string
		if val.form == regularForm {
			if fltval, ok := val.Float64(); ok {
				bigres = fmt.Sprintf("%v", fltval)
			} else {
				val.Big(bigval)

				prec := bigval.NumDigits() - 1
				exp := int64(bigval.Exponent) + prec

				if exp < -4 || exp > 5 {
					bigres = fmt.Sprintf("%e", bigval)

					if idx := strings.IndexRune(bigres, 'e'); idx != -1 && idx == len(bigres)-3 {
						idx += 2
						bigres = bigres[:idx] + "0" + bigres[idx:]
					}
				} else {
					bigres = fmt.Sprintf("%f", bigval)
				}
			}
		} else if val.form == infForm {
			if val.neg {
				bigres = fmt.Sprintf("%v", math.Inf(-1))
			} else {
				bigres = fmt.Sprintf("%v", math.Inf(1))
			}
		} else {
			bigres = fmt.Sprintf("%v", math.NaN())
		}

		if res != bigres {
			t.Errorf("%v.String() = %s, want %s", val, res, bigres)
		}

		fmtres := fmt.Sprintf("%v", decval)

		if fmtres != res {
			t.Errorf("fmt.Sprintf(%v) = %s, want %s", val, fmtres, res)
		}

		mshres, err := decval.MarshalText()

		if string(mshres) != res || err != nil {
			t.Errorf("%v.MarshalText() = (%s, %v), want (%s, <nil>)", val, mshres, err, res)
		}
	}
}

type devNull struct{}

func (d devNull) Write(b []byte) (int, error) {
	return len(b), nil
}

// BenchmarkFormat benchmarks formatting using the various format specifiers.
func BenchmarkFormat(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		fmt  string
		want string
	}{{
		name: "large number f format",
		txt:  "12345678901234.67890123456789012345",
		fmt:  "%f",
		want: "12345678901234.678901",
	}, {
		name: "large number e format",
		txt:  "12345678901234.67890123456789012345",
		fmt:  "%e",
		want: "1.234568e+13",
	}, {
		name: "large number padding f format",
		txt:  "12345678901234.67890123456789012345",
		fmt:  "%80.40f",
		want: "                         12345678901234.6789012345678901234500000000000000000000",
	}, {
		name: "large number padding e format",
		txt:  "12345678901234.67890123456789012345",
		fmt:  "%80.40e",
		want: "                                  1.2345678901234678901234567890123450000000e+13",
	}, {
		name: "special",
		txt:  "nan",
		fmt:  "%f",
		want: "NaN",
	}, {
		name: "special with padding",
		txt:  "nan",
		fmt:  "%80.40f",
		want: "                                                                             NaN",
	}}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			// Ensure correctness of test case before benchmarking.
			v := MustParse(tc.txt)
			vptr := &v // Avoid allocating during Fprintf() call
			got := fmt.Sprintf(tc.fmt, vptr)
			if got != tc.want {
				b.Fatalf("Unexpected formatted value. got %s, "+
					"want %s", got, tc.want)
			}
			w := devNull{}

			// Benchmark.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				fmt.Fprintf(w, tc.fmt, vptr)
			}
		})
	}
}
