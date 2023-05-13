package decimal128

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"unsafe"

	"github.com/cockroachdb/apd/v3"
)

var maxDecimal *apd.Decimal

func init() {
	lo := new(apd.BigInt)
	lo.SetUint64(math.MaxUint64)

	maxDecimal = new(apd.Decimal)
	maxDecimal.Exponent = maxUnbiasedExponent
	maxDecimal.Coeff.SetUint64(0x0002_7fff_ffff_ffff)
	maxDecimal.Coeff.Lsh(&maxDecimal.Coeff, 64).Or(&maxDecimal.Coeff, lo)
}

type testForm uint8

const (
	regularForm testForm = iota
	infForm
	nanForm
)

type testDec struct {
	form testForm
	neg  bool
	sig  uint128
	exp  int16
}

func (td testDec) Big(dec *apd.Decimal) *apd.Decimal {
	switch td.form {
	case regularForm:
		dec.Form = apd.Finite
		dec.Negative = td.neg
		dec.Exponent = int32(td.exp - exponentBias)

		if td.sig[1] == 0 {
			dec.Coeff.SetUint64(td.sig[0])
		} else {
			lo := new(apd.BigInt)
			lo.SetUint64(td.sig[0])

			dec.Coeff.SetUint64(td.sig[1])
			dec.Coeff.Lsh(&dec.Coeff, 64).Or(&dec.Coeff, lo)
		}

		dec.Reduce(dec)
		dec.Negative = td.neg
	case infForm:
		dec.Form = apd.Infinite
		dec.Negative = td.neg
	case nanForm:
		dec.Form = apd.NaN
	default:
		panic("unhandled test decimal form")
	}

	return dec
}

func (td testDec) Decimal() Decimal {
	switch td.form {
	case regularForm:
		return compose(td.neg, td.sig, td.exp)
	case infForm:
		return inf(td.neg)
	case nanForm:
		return NaN()
	default:
		panic("unhandled test decimal form")
	}
}

func (td testDec) Float64() (float64, bool) {
	if td.form == nanForm {
		return math.NaN(), true
	}

	if td.form == infForm {
		sign := 1
		if td.neg {
			sign = -1
		}

		return math.Inf(sign), true
	}

	if td.sig[1] != 0 {
		return 0.0, false
	}

	sig64 := td.sig[0]

	if sig64 == 0 {
		if td.neg {
			return math.Copysign(0.0, -1.0), true
		}

		return 0.0, true
	}

	if sig64 > math.MaxUint32 {
		return 0.0, false
	}

	if td.exp < exponentBias {
		return 0.0, false
	}

	switch td.exp {
	case exponentBias:
	case exponentBias + 1:
		sig64 *= 10
	case exponentBias + 2:
		sig64 *= 100
	case exponentBias + 3:
		sig64 *= 1000
	case exponentBias + 4:
		sig64 *= 10_000
	case exponentBias + 5:
		sig64 *= 100_000
	default:
		return 0.0, false
	}

	if sig64 > math.MaxUint32 {
		return 0.0, false
	}

	if td.neg {
		return math.Copysign(float64(sig64), -1.0), true
	}

	return float64(sig64), true
}

func (td testDec) String() string {
	switch td.form {
	case regularForm:
		sign := ""
		if td.neg {
			sign = "-"
		}

		return fmt.Sprintf("%s%ve%d", sign, td.sig, td.exp-exponentBias)
	case infForm:
		sign := ""
		if td.neg {
			sign = "-"
		}

		return fmt.Sprintf("%sinf", sign)
	case nanForm:
		return "nan"
	default:
		panic("unhandled test decimal form")
	}
}

type testDataReader struct {
	t     *testing.T
	files []string
	f     *os.File
	r     *bufio.Reader
}

func openTestData(t *testing.T) *testDataReader {
	dir := filepath.Join("testdata", t.Name())

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("error reading testdata directory %s: %v", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return &testDataReader{
		t:     t,
		files: files,
	}
}

func (tr *testDataReader) close() {
	if tr.f != nil {
		if err := tr.f.Close(); err != nil {
			tr.t.Fatalf("error closing file %s: %v", tr.f.Name(), err)
		}

		tr.f = nil
		tr.r = nil
	}

	tr.files = nil
}

func (tr *testDataReader) openNext() bool {
	if len(tr.files) == 0 {
		return false
	}

	if tr.f != nil {
		if err := tr.f.Close(); err != nil {
			tr.t.Fatalf("error closing file %s: %v", tr.f.Name(), err)
		}
	}

	var err error
	tr.f, err = os.Open(tr.files[0])
	if err != nil {
		tr.t.Fatalf("error opening file %s: %v", tr.files[0], err)
	}

	if tr.r == nil {
		tr.r = bufio.NewReader(tr.f)
	} else {
		tr.r.Reset(tr.f)
	}

	tr.files = tr.files[1:]

	return true
}

func (tr *testDataReader) scan(format string, args ...any) bool {
	for {
		if tr.f == nil {
			if tr.openNext() {
				continue
			}

			return false
		}

		n, err := fmt.Fscanf(tr.r, format, args...)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				if n != 0 {
					tr.t.Fatalf("error reading file %s: %v", tr.f.Name(), err)
				}

				if tr.openNext() {
					continue
				}

				return false
			}

			tr.t.Fatalf("error reading file %s: %v", tr.f.Name(), err)
		}

		return true
	}
}

type testDataResult struct {
	ToNearestEven Decimal
	ToNearestAway Decimal
	ToZero        Decimal
	AwayFromZero  Decimal
	ToNegativeInf Decimal
	ToPositiveInf Decimal
}

func (tr *testDataResult) Scan(f fmt.ScanState, verb rune) error {
	if verb != 'v' {
		return errors.New("bad verb '%" + string(verb) + "' for testDataResult")
	}

	tok, err := f.Token(true, nil)
	if err != nil {
		return err
	}

	index := bytes.IndexByte(tok, ';')

	var data []byte
	if index == -1 {
		data = tok
	} else {
		data = tok[:index]
	}

	var res Decimal
	if err := res.UnmarshalText(data); err != nil {
		return err
	}

	tr.ToNearestEven = res
	tr.ToNearestAway = res
	tr.ToZero = res
	tr.AwayFromZero = res
	tr.ToNegativeInf = res
	tr.ToPositiveInf = res

	for index != -1 {
		tok = tok[index+1:]

		sep := bytes.IndexByte(tok, ':')
		if sep == -1 {
			return errors.New("invalid value")
		}

		index = bytes.IndexByte(tok[sep+1:], ';')

		if index == -1 {
			data = tok[sep+1:]
		} else {
			data = tok[sep+1 : index]
		}

		if err := res.UnmarshalText(data); err != nil {
			return err
		}

		modes := tok[:sep]

		for sep != -1 {
			sep = bytes.IndexByte(modes, ',')

			var mode []byte
			if sep == -1 {
				mode = modes
			} else {
				mode = modes[:sep]
				modes = modes[sep+1:]
			}

			switch string(mode) {
			case "NE":
				tr.ToNearestEven = res
			case "NA":
				tr.ToNearestAway = res
			case "Z":
				tr.ToZero = res
			case "FZ":
				tr.AwayFromZero = res
			case "NI":
				tr.ToNegativeInf = res
			case "PI":
				tr.ToPositiveInf = res
			default:
				return errors.New("invalid value \"" + string(mode) + "\"")
			}
		}
	}

	return nil
}

func (tr *testDataResult) equal(val Decimal, mode RoundingMode) bool {
	return resultEqual(val, tr.result(mode))
}

func (tr *testDataResult) result(mode RoundingMode) Decimal {
	switch mode {
	case ToNearestEven:
		return tr.ToNearestEven
	case ToNearestAway:
		return tr.ToNearestAway
	case ToZero:
		return tr.ToZero
	case AwayFromZero:
		return tr.AwayFromZero
	case ToNegativeInf:
		return tr.ToNegativeInf
	case ToPositiveInf:
		return tr.ToPositiveInf
	default:
		panic("invalid rounding mode " + mode.String())
	}
}

type testDataResultPair struct {
	first  testDataResult
	second testDataResult
	sep    byte
}

func (tr *testDataResultPair) Scan(f fmt.ScanState, verb rune) error {
	if verb != 'v' {
		return errors.New("bad verb '%" + string(verb) + "' for testDataResultPair")
	}

	tok, err := f.Token(true, nil)
	if err != nil {
		return err
	}

	index := bytes.IndexByte(tok, ';')

	var data []byte
	if index == -1 {
		data = tok
	} else {
		data = tok[:index]
	}

	pairSep := bytes.IndexByte(data, tr.sep)
	if pairSep == -1 {
		return errors.New("invalid value")
	}

	var firstRes Decimal
	if err := firstRes.UnmarshalText(data[:pairSep]); err != nil {
		return err
	}

	var secondRes Decimal
	if err := secondRes.UnmarshalText(data[pairSep+1:]); err != nil {
		return err
	}

	tr.first.ToNearestEven = firstRes
	tr.first.ToNearestAway = firstRes
	tr.first.ToZero = firstRes
	tr.first.AwayFromZero = firstRes
	tr.first.ToNegativeInf = firstRes
	tr.first.ToPositiveInf = firstRes

	tr.second.ToNearestEven = secondRes
	tr.second.ToNearestAway = secondRes
	tr.second.ToZero = secondRes
	tr.second.AwayFromZero = secondRes
	tr.second.ToNegativeInf = secondRes
	tr.second.ToPositiveInf = secondRes

	for index != -1 {
		tok = tok[index+1:]

		sep := bytes.IndexByte(tok, ':')
		if sep == -1 {
			return errors.New("invalid value")
		}

		index = bytes.IndexByte(tok[sep+1:], ';')

		if index == -1 {
			data = tok[sep+1:]
		} else {
			data = tok[sep+1 : index]
		}

		pairSep = bytes.IndexByte(data, tr.sep)
		if pairSep == -1 {
			return errors.New("invalid value")
		}

		if err := firstRes.UnmarshalText(data[:pairSep]); err != nil {
			return err
		}

		if err := secondRes.UnmarshalText(data[pairSep+1:]); err != nil {
			return err
		}

		modes := tok[:sep]

		for sep != -1 {
			sep = bytes.IndexByte(modes, ',')

			var mode []byte
			if sep == -1 {
				mode = modes
			} else {
				mode = modes[:sep]
				modes = modes[sep+1:]
			}

			switch string(mode) {
			case "NE":
				tr.first.ToNearestEven = firstRes
				tr.second.ToNearestEven = secondRes
			case "NA":
				tr.first.ToNearestAway = firstRes
				tr.second.ToNearestAway = secondRes
			case "Z":
				tr.first.ToZero = firstRes
				tr.second.ToZero = secondRes
			case "FZ":
				tr.first.AwayFromZero = firstRes
				tr.second.AwayFromZero = secondRes
			case "NI":
				tr.first.ToNegativeInf = firstRes
				tr.second.ToNegativeInf = secondRes
			case "PI":
				tr.first.ToPositiveInf = firstRes
				tr.second.ToPositiveInf = secondRes
			default:
				return errors.New("invalid value \"" + string(mode) + "\"")
			}
		}
	}

	return nil
}

func (tr *testDataResultPair) equal(first, second Decimal, mode RoundingMode) bool {
	return tr.first.equal(first, mode) && tr.second.equal(second, mode)
}

func resultEqual(x, y Decimal) bool {
	if x.IsNaN() && y.IsNaN() {
		return true
	}

	if x.Signbit() != y.Signbit() {
		return false
	}

	return x.Equal(y)
}

var (
	decimalValuesOnce sync.Once
	decimalValues     []testDec
)

func initDecimalValues() {
	decimalValuesOnce.Do(func() {
		initUintValues()

		var exponentValues []int16
		if testing.Short() {
			exponentValues = []int16{
				minBiasedExponent,
				exponentBias - 19,
				exponentBias,
				exponentBias + 19,
				maxBiasedExponent,
			}
		} else {
			exponentValues = []int16{
				minBiasedExponent,
				minBiasedExponent + exponentBias/2,
				exponentBias - 34,
				exponentBias - 19,
				exponentBias - 5,
				exponentBias - 1,
				exponentBias,
				exponentBias + 1,
				exponentBias + 5,
				exponentBias + 19,
				exponentBias + 34,
				maxBiasedExponent - exponentBias/2,
				maxBiasedExponent,
			}
		}

		for _, sighi := range uint64Values {
			if sighi > 0x0002_7fff_ffff_ffff {
				continue
			}

			for _, siglo := range uint64Values {
				sig := uint128{siglo, sighi}

				if sig == (uint128{}) {
					continue
				}

				for _, exp := range exponentValues {
					decimalValues = append(decimalValues,
						testDec{regularForm, false, sig, exp},
						testDec{regularForm, true, sig, exp},
					)
				}
			}
		}

		decimalValues = append(decimalValues, testDec{regularForm, false, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{regularForm, true, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{infForm, false, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{infForm, true, uint128{}, 0})
		decimalValues = append(decimalValues, testDec{nanForm, false, uint128{}, 0})
	})
}

func decimalToBig(v Decimal) *apd.Decimal {
	r := new(apd.Decimal)

	if v.isInf() {
		r.Form = apd.Infinite
		r.Negative = v.Signbit()
		return r
	}

	if v.IsNaN() {
		r.Form = apd.NaN
		return r
	}

	r.Negative = v.Signbit()

	sig, exp := v.decompose()
	r.Exponent = int32(exp - exponentBias)

	if sig[1] == 0 {
		r.Coeff.SetUint64(sig[0])
	} else {
		lo := new(apd.BigInt)
		lo.SetUint64(sig[0])

		r.Coeff.SetUint64(sig[1])
		r.Coeff.Lsh(&r.Coeff, 64).Or(&r.Coeff, lo)
	}

	r.Reduce(r)
	return r
}

func decimalsEqual(x Decimal, y *apd.Decimal, mode apd.Rounder) bool {
	if x.isSpecial() {
		if x.IsNaN() {
			return y.Form == apd.NaN || y.Form == apd.NaNSignaling
		}

		if x.isInf() {
			if y.Negative != x.Signbit() {
				return false
			}

			if y.Form == apd.Infinite {
				return true
			}

			neg := y.Negative
			y.Negative = false
			cmp := y.Cmp(maxDecimal)
			y.Negative = neg

			return cmp > 0
		}
	} else if y.Form != apd.Finite {
		return false
	}

	if x.Signbit() != y.Negative {
		// apd appears to always return -0 when rounding towards -infinity,
		// even if the operands are themselves zero.
		if x.IsZero() && y.Coeff.IsInt64() && y.Coeff.Int64() == 0 && mode == apd.RoundFloor {
			return true
		}

		return false
	}

	bigx := decimalToBig(x)

	bigctx := apd.Context{
		Precision:   uint32(bigx.NumDigits()),
		MaxExponent: 6145,
		MinExponent: -6176,
		Rounding:    mode,
	}

	bigctx.Round(y, y)

	// apd appears to return the wrong result during rounding in some scenarios
	// when the result underflows, returning +/-1e-6176 instead of 0.
	if x.IsZero() && y.Exponent <= -6176 && y.Coeff.IsInt64() && y.Coeff.Int64() == 1 {
		return true
	}

	// apd appears to return the wrong result during rounding in some scenarios
	// when the result is just under the allowed exponent range, returning 0
	// instead of +/-1e-6176.
	if bigx.Coeff.IsInt64() && bigx.Coeff.Int64() == 1 && bigx.Exponent == -6176 && y.Coeff.IsInt64() && y.Coeff.Int64() == 0 {
		return true
	}

	return bigx.Cmp(y) == 0
}

func TestAbs(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := Abs(decval)

		absval := val
		absval.neg = false
		absres := absval.Decimal()

		if !(res.Equal(absres) || res.IsNaN() && absres.IsNaN()) && res.Signbit() == absres.Signbit() {
			t.Errorf("Abs(%v) = %v, want %v", val, res, absres)
		}
	}
}

func TestCompose(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		dec := compose(val.neg, val.sig, val.exp)

		if dec.isSpecial() {
			t.Errorf("%v.isSpecial() = true, want false", val)
		}

		if dec.isInf() {
			t.Errorf("%v.isInf() = true, want false", val)
		}

		if dec.IsNaN() {
			t.Errorf("%v.IsNaN() = true, want false", val)
		}

		if dec.Signbit() != val.neg {
			t.Errorf("%v.Signbit() = %t, want %t", val, dec.Signbit(), val.neg)
		}

		sig, exp := dec.decompose()

		if sig != val.sig || exp != val.exp {
			t.Errorf("%v.decompose() = (%v, %d), want (%v, %d)", val, sig, exp, val.sig, val.exp)
		}
	}
}

func TestDecimalNeg(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()
		res := decval.Neg()

		negval := val
		negval.neg = !val.neg
		negres := negval.Decimal()

		if !(res.Equal(negres) || res.IsNaN() && negres.IsNaN()) && res.Signbit() == negres.Signbit() {
			t.Errorf("%v.Neg() = %v, want %v", val, res, negval.Decimal())
		}
	}
}

func TestInf(t *testing.T) {
	t.Parallel()

	dec := Inf(1)

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if !dec.isInf() {
		t.Errorf("%v.isInf() = false, want true", dec)
	}

	if dec.IsNaN() {
		t.Errorf("%v.IsNaN() = true, want false", dec)
	}

	if dec.Signbit() {
		t.Errorf("%v.Signbit() = true, want false", dec)
	}

	dec = Inf(-1)

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if !dec.isInf() {
		t.Errorf("%v.isInf() = false, want true", dec)
	}

	if dec.IsNaN() {
		t.Errorf("%v.IsNaN() = true, want false", dec)
	}

	if !dec.Signbit() {
		t.Errorf("%v.Signbit() = false, want true", dec)
	}
}

func TestNaN(t *testing.T) {
	t.Parallel()

	dec := NaN()

	if !dec.isSpecial() {
		t.Errorf("%v.isSpecial() = false, want true", dec)
	}

	if dec.isInf() {
		t.Errorf("%v.isInf() = true, want false", dec)
	}

	if !dec.IsNaN() {
		t.Errorf("%v.IsNaN() = false, want true", dec)
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	res := unsafe.Sizeof(Decimal{})

	if res != 16 {
		t.Errorf("unsafe.Sizeof(Decimal{}) = %d, want 16", res)
	}
}

func FuzzDecimal(f *testing.F) {
	f.Add(uint64(0), uint64(0))
	f.Add(uint64(math.MaxUint64), uint64(math.MaxUint64))

	f.Fuzz(func(t *testing.T, hi, lo uint64) {
		t.Parallel()

		dec := Decimal{hi, lo}

		if dec.isSpecial() {
			if dec.isInf() == dec.IsNaN() {
				t.Fail()
			}
		} else {
			if dec.isInf() || dec.IsNaN() {
				t.Fail()
			}

			sig, exp := dec.decompose()
			res := compose(dec.Signbit(), sig, exp)

			if res != dec {
				t.Fail()
			}
		}
	})
}
