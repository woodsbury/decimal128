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
				if val.sig == (uint128{}) {
					if val.neg {
						bigres = fmt.Sprintf(fmtstr, math.Copysign(0.0, -1.0))
					} else {
						bigres = fmt.Sprintf(fmtstr, 0.0)
					}
				} else if val.exp == exponentBias && val.sig[1] == 0 && val.sig[0] < math.MaxUint32 {
					if val.neg {
						bigres = fmt.Sprintf(fmtstr, math.Copysign(float64(val.sig[0]), -1.0))
					} else {
						bigres = fmt.Sprintf(fmtstr, float64(val.sig[0]))
					}
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
			if val.sig == (uint128{}) {
				if val.neg {
					bigres = fmt.Sprintf("%v", math.Copysign(0.0, -1.0))
				} else {
					bigres = fmt.Sprintf("%v", 0.0)
				}
			} else if val.exp == exponentBias && val.sig[1] == 0 && val.sig[0] < math.MaxUint32 {
				if val.neg {
					bigres = fmt.Sprintf("%v", math.Copysign(float64(val.sig[0]), -1.0))
				} else {
					bigres = fmt.Sprintf("%v", float64(val.sig[0]))
				}
			} else {
				val.Big(bigval)

				prec := bigval.NumDigits() - 1
				exp := int64(bigval.Exponent) + prec

				if exp < -4 || exp > prec {
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
