package decimal128

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
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

	if tf.prec != -1 {
		builder.WriteRune('.')
		if tf.prec != 0 {
			builder.WriteString(strconv.Itoa(tf.prec))
		}
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
	for _, width := range []int{0, 4, 15} {
		for _, prec := range []int{-1, 0, 4, 15} {
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

	var appres []byte

	for _, format := range formats {
		fmtstr := format.String()

		for _, val := range decimalValues {
			decval := val.Decimal()
			res := fmt.Sprintf(fmtstr, decval)

			if strings.Contains(res, "PANIC") {
				t.Errorf("fmt.Sprintf(%v, %v) = %s", format, val, res)
				continue
			}

			appres := decval.Append(appres[:0], fmtstr[1:])
			if res != string(appres) {
				t.Errorf("%v.Append(%s) = %s, want %s", val, fmtstr, appres, res)
			}

			fltval, ok := val.Float64()
			if !ok {
				// Skipping for now
				continue
			}

			fltres := fmt.Sprintf(fmtstr, fltval)

			if res != fltres {
				t.Errorf("fmt.Sprintf(%v, %v) = %s, want %s", format, val, res, fltres)
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

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := decval.String()

		fltval, ok := val.Float64()
		if !ok {
			// Skipping for now
			continue
		}

		fltres := fmt.Sprintf("%v", fltval)

		if res != fltres {
			t.Errorf("%v.String() = %s, want %s", val, res, fltres)
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

func TestFormat(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	var appres []byte

	for _, format := range []byte{'e', 'E', 'f', 'F', 'g', 'G'} {
		for _, prec := range []int{-1, 0, 4, 15} {
			for _, val := range decimalValues {
				decval := val.Decimal()
				res := Format(decval, format, prec)

				appres = Append(appres[:0], decval, format, prec)
				if res != string(appres) {
					t.Errorf("Append(%v, %c, %d) = %s, want %s", val, format, prec, appres, res)
				}

				fltval, ok := val.Float64()
				if !ok {
					// Skipping for now
					continue
				}

				fltres := strconv.FormatFloat(fltval, format, prec, 64)
				if res != fltres {
					t.Errorf("Format(%v, %c, %d) = %s, want %s", val, format, prec, res, fltres)
				}
			}
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		fmt  string
		want string
	}{
		{
			name: "large number f format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "f",
			want: "12345678901234.678901",
		},
		{
			name: "large number e format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "e",
			want: "1.234568e+13",
		},
		{
			name: "large number padding f format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "80.40f",
			want: "                         12345678901234.6789012345678901234500000000000000000000",
		},
		{
			name: "large number padding e format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "80.40e",
			want: "                                  1.2345678901234678901234567890123450000000e+13",
		},
		{
			name: "special",
			txt:  "nan",
			fmt:  "f",
			want: "NaN",
		},
		{
			name: "special with padding",
			txt:  "nan",
			fmt:  "80.40f",
			want: "                                                                             NaN",
		},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			// Ensure correctness of test case before benchmarking.
			v := MustParse(tc.txt)
			buf := []byte{}
			buf = v.Append(buf, tc.fmt)
			if string(buf) != tc.want {
				b.Fatalf("Unexpected formatted value. got %s, "+
					"want %s", buf, tc.want)
			}

			// Benchmark.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = buf[:0]
				buf = v.Append(buf, tc.fmt)
			}
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		fmt  string
		want string
	}{
		{
			name: "large number f format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "%f",
			want: "12345678901234.678901",
		},
		{
			name: "large number e format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "%e",
			want: "1.234568e+13",
		},
		{
			name: "large number padding f format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "%80.40f",
			want: "                         12345678901234.6789012345678901234500000000000000000000",
		},
		{
			name: "large number padding e format",
			txt:  "12345678901234.67890123456789012345",
			fmt:  "%80.40e",
			want: "                                  1.2345678901234678901234567890123450000000e+13",
		},
		{
			name: "special",
			txt:  "nan",
			fmt:  "%f",
			want: "NaN",
		},
		{
			name: "special with padding",
			txt:  "nan",
			fmt:  "%80.40f",
			want: "                                                                             NaN",
		},
	}

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

			// Benchmark.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				fmt.Fprintf(io.Discard, tc.fmt, vptr)
			}
		})
	}
}

func BenchmarkMarshalText(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		want string
	}{
		{
			name: "small number",
			txt:  "1234.1234",
			want: "1234.1234",
		},
		{
			name: "large number",
			txt:  "12345678901234.67890123456789012345",
			want: "1.234567890123467890123456789012345e+13",
		},
		{
			name: "special",
			txt:  "nan",
			want: "NaN",
		},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			// Ensure correctness of test case before benchmarking.
			v := MustParse(tc.txt)
			got, err := v.MarshalText()
			if err != nil || string(got) != tc.want {
				b.Fatalf("Unexpected formatted value. got '%s', "+
					"want '%s' with %v", got, tc.want, err)
			}

			// Benchmark.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				v.MarshalText()
			}
		})
	}
}

func BenchmarkParseFormat(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		fmt  string
		want string
	}{
		{
			name: "f format",
			fmt:  "%f",
		},
		{
			name: "e format",
			fmt:  "%e",
		},
		{
			name: "padding f format",
			fmt:  "%80.40f",
		},
		{
			name: "padding e format",
			fmt:  "%80.40e",
		},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			// Benchmark.
			var args formatArgs
			for i := 0; i < b.N; i++ {
				parseFormat(tc.fmt, &args)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	tests := []struct {
		name string
		txt  string
		want string
	}{
		{
			name: "small number",
			txt:  "1234.1234",
			want: "1234.1234",
		},
		{
			name: "large number",
			txt:  "12345678901234.67890123456789012345",
			want: "1.234567890123467890123456789012345e+13",
		},
		{
			name: "special",
			txt:  "nan",
			want: "NaN",
		},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			// Ensure correctness of test case before benchmarking.
			v := MustParse(tc.txt)
			got := v.String()
			if got != tc.want {
				b.Fatalf("Unexpected formatted value. got '%s', "+
					"want '%s'", got, tc.want)
			}

			// Benchmark.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = v.String()
			}
		})
	}
}
