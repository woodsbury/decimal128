package decimal128

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	nanText    = []byte("NaN")
	padNaNText = []byte(" NaN")
	posNaNText = []byte("+NaN")
	negInfText = []byte("-Inf")
	padInfText = []byte(" Inf")
	posInfText = []byte("+Inf")

	digitPairs = [...][2]byte{
		{'0', '0'}, {'0', '1'}, {'0', '2'}, {'0', '3'}, {'0', '4'},
		{'0', '5'}, {'0', '6'}, {'0', '7'}, {'0', '8'}, {'0', '9'},
		{'1', '0'}, {'1', '1'}, {'1', '2'}, {'1', '3'}, {'1', '4'},
		{'1', '5'}, {'1', '6'}, {'1', '7'}, {'1', '8'}, {'1', '9'},
		{'2', '0'}, {'2', '1'}, {'2', '2'}, {'2', '3'}, {'2', '4'},
		{'2', '5'}, {'2', '6'}, {'2', '7'}, {'2', '8'}, {'2', '9'},
		{'3', '0'}, {'3', '1'}, {'3', '2'}, {'3', '3'}, {'3', '4'},
		{'3', '5'}, {'3', '6'}, {'3', '7'}, {'3', '8'}, {'3', '9'},
		{'4', '0'}, {'4', '1'}, {'4', '2'}, {'4', '3'}, {'4', '4'},
		{'4', '5'}, {'4', '6'}, {'4', '7'}, {'4', '8'}, {'4', '9'},
		{'5', '0'}, {'5', '1'}, {'5', '2'}, {'5', '3'}, {'5', '4'},
		{'5', '5'}, {'5', '6'}, {'5', '7'}, {'5', '8'}, {'5', '9'},
		{'6', '0'}, {'6', '1'}, {'6', '2'}, {'6', '3'}, {'6', '4'},
		{'6', '5'}, {'6', '6'}, {'6', '7'}, {'6', '8'}, {'6', '9'},
		{'7', '0'}, {'7', '1'}, {'7', '2'}, {'7', '3'}, {'7', '4'},
		{'7', '5'}, {'7', '6'}, {'7', '7'}, {'7', '8'}, {'7', '9'},
		{'8', '0'}, {'8', '1'}, {'8', '2'}, {'8', '3'}, {'8', '4'},
		{'8', '5'}, {'8', '6'}, {'8', '7'}, {'8', '8'}, {'8', '9'},
		{'9', '0'}, {'9', '1'}, {'9', '2'}, {'9', '3'}, {'9', '4'},
		{'9', '5'}, {'9', '6'}, {'9', '7'}, {'9', '8'}, {'9', '9'},
	}
)

// Append appends the decimal to the buffer, using the specified format.
func (d Decimal) Append(buf []byte, f *FmtFormat) []byte {
	width, hasWidth := f.width, f.hasWidth
	verb := f.verb

	if d.isSpecial() {
		pad := 0
		padSign := false
		printSign := false
		if verb != 'v' {
			printSign = f.flag('+')
			padSign = f.flag(' ')

			if hasWidth {
				pad = width
			}
		}

		return d.fmtSpecial(buf, pad, printSign, padSign, f.flag('-'), false)
	}

	var digs digits
	d.digits(&digs)
	prec, hasPrec := f.prec, f.hasPrec

	switch verb {
	case 'e', 'E':
		if !hasPrec {
			prec = 6
		}

		pad := 0
		if hasWidth {
			pad = width
		}

		digs.round(prec + 1)
		return digs.fmtE(buf, prec, pad, f.flag('#'), f.flag('+'), f.flag(' '), true, f.flag('-'), f.flag('0'), byte(verb))
	case 'f', 'F':
		if !hasPrec {
			prec = 6
		}

		if digs.exp < 0 {
			digs.round(digs.ndig + digs.exp + prec)
		}

		pad := 0
		if hasWidth {
			pad = width
		}

		return digs.fmtF(buf, prec, pad, f.flag('#'), f.flag('+'), f.flag(' '), f.flag('-'), f.flag('0'))
	case 'g', 'G':
		var maxprec int
		if f.flag('#') {
			if !hasPrec {
				if digs.ndig < 6 {
					prec = 6
				} else {
					prec = digs.ndig
				}

				maxprec = 6
			} else {
				maxprec = prec
			}

			if hasWidth {
				if digs.ndig > prec {
					if width > digs.ndig {
						width -= digs.ndig
						prec = digs.ndig
					} else {
						prec += width
						width -= digs.ndig
					}
				}
			}

			digs.round(prec)
		} else {
			if hasPrec {
				digs.round(prec)
				maxprec = prec
			} else if digs.ndig != 0 {
				prec = digs.ndig
				maxprec = 6
			} else {
				maxprec = 6
			}
		}

		eprec := 0
		if digs.ndig != 0 {
			eprec = digs.ndig - 1
		}

		exp := digs.exp + eprec

		pad := 0
		if hasWidth {
			pad = width
		}

		if exp < -4 || exp >= maxprec {
			e := byte('e')
			if verb == 'G' {
				e = byte('E')
			}

			return digs.fmtE(buf, prec-1, pad, f.flag('#'), f.flag('+'), f.flag(' '), true, f.flag('-'), f.flag('0'), e)
		} else {
			if f.flag('#') {
				prec -= digs.exp
				if digs.ndig == 0 {
					prec--
				} else {
					prec -= digs.ndig
				}
			} else {
				prec = 0
				if digs.exp < 0 {
					prec = -digs.exp
				}
			}

			return digs.fmtF(buf, prec, pad, f.flag('#'), f.flag('+'), f.flag(' '), f.flag('-'), f.flag('0'))
		}
	case 'v':
		prec := 0
		if digs.ndig != 0 {
			prec = digs.ndig - 1
		}

		exp := digs.exp + prec

		if exp < -4 || exp >= 6 {
			return digs.fmtE(buf, prec, 0, false, false, false, true, false, false, 'e')
		} else {
			prec = 0
			if digs.exp < 0 {
				prec = -digs.exp
			}

			return digs.fmtF(buf, prec, 0, false, false, false, false, false)
		}
	default:
		return append(buf, []byte(fmt.Sprintf("%%!%c(decimal128.Decimal=%s)", verb, d.String()))...)
	}
}

// Format implements the [fmt.Formatter] interface. It supports the verbs 'e',
// 'E', 'f', 'F', 'g', 'G', and 'v', along with the format flags '+', '-', '#',
// ' ', and '0' and custom width and precision values. Decimal values interpret
// the format value the same way float32 and float64 does.
func (d Decimal) Format(st fmt.State, verb rune) {
	var f FmtFormat
	fmtFormatFromFmtState(st, verb, &f)
	st.Write(d.Append(nil, &f))
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (d Decimal) MarshalText() ([]byte, error) {
	if d.isSpecial() {
		return d.fmtSpecial(nil, 0, false, false, false, true), nil
	}

	var digs digits
	d.digits(&digs)

	prec := 0
	if digs.ndig != 0 {
		prec = digs.ndig - 1
	}

	exp := digs.exp + prec

	if exp < -4 || exp >= 6 {
		return digs.fmtE(nil, prec, 0, false, false, false, true, false, false, 'e'), nil
	}

	prec = 0
	if digs.exp < 0 {
		prec = -digs.exp
	}

	return digs.fmtF(nil, prec, 0, false, false, false, false, false), nil
}

// String returns a string representation of the Decimal value.
func (d Decimal) String() string {
	if d.isSpecial() {
		return string(d.fmtSpecial(nil, 0, false, false, false, false))
	}

	var digs digits
	d.digits(&digs)

	prec := 0
	if digs.ndig != 0 {
		prec = digs.ndig - 1
	}

	exp := digs.exp + prec

	if exp < -4 || exp >= 6 {
		return string(digs.fmtE(nil, prec, 0, false, false, false, true, false, false, 'e'))
	}

	prec = 0
	if digs.exp < 0 {
		prec = -digs.exp
	}

	return string(digs.fmtF(nil, prec, 0, false, false, false, false, false))
}

func (d Decimal) digits(digs *digits) {
	*digs = digits{}
	digs.neg = d.Signbit()

	sig, exp := d.decompose()

	if sig != (uint128{}) {
		digs.exp = int(exp - exponentBias)

		n := 0
		for sig != (uint128{}) {
			var rem uint64
			sig, rem = sig.div100()

			if n == 0 && rem == 0 {
				digs.exp += 2
				continue
			}

			pair := digitPairs[rem]

			if n == 0 && pair[1] == '0' {
				digs.exp++
				digs.dig[n] = pair[0]
				n++
				continue
			}

			if pair[0] == '0' && sig == (uint128{}) {
				digs.dig[n] = pair[1]
				n++
			} else {
				digs.dig[n], digs.dig[n+1] = pair[1], pair[0]
				n += 2
			}
		}

		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			digs.dig[i], digs.dig[j] = digs.dig[j], digs.dig[i]
		}

		digs.ndig = n
	}
}

func (d Decimal) fmtSpecial(buf []byte, pad int, printSign, padSign, padRight, copyBuf bool) []byte {

	if d.IsNaN() {
		if printSign {
			buf = posNaNText
		} else if padSign {
			buf = padNaNText
		} else {
			buf = nanText
		}
	} else {
		if d.Signbit() {
			buf = negInfText
		} else {
			if padSign && !printSign {
				buf = padInfText
			} else {
				buf = posInfText
			}
		}
	}

	if p := pad - len(buf); p > 0 {
		tmp := make([]byte, pad)

		if padRight {
			n := copy(tmp, buf)

			for i := n; i < pad; i++ {
				tmp[i] = ' '
			}
		} else {
			for i := 0; i < p; i++ {
				tmp[i] = ' '
			}

			copy(tmp[p:], buf)
		}

		buf = tmp
		copyBuf = false
	}

	if copyBuf {
		tmp := make([]byte, len(buf))
		copy(tmp, buf)
		buf = tmp
	}

	return buf
}

type digits struct {
	neg  bool
	dig  [39]byte
	exp  int
	ndig int
}

func (d *digits) fmtE(buf []byte, prec, pad int, forceDP, printSign, padSign, padExp, padRight, padZero bool, e byte) []byte {

	// Attempt to pre-size buffer to avoid multiple allocations. Currently,
	// this might overshoot the actual needed size. Calculation is:
	// sign + decimal separator + exponent + signficant digits + padding.
	sizeHint := 1 + 1 + 5 + d.ndig
	if pad > sizeHint {
		sizeHint = pad
	}
	buf = buf[len(buf):]
	if cap(buf) < sizeHint {
		buf = append(buf, make([]byte, sizeHint)...)[:0]
	}

	if d.neg {
		buf = append(buf, '-')
	} else if printSign {
		buf = append(buf, '+')
	} else if padSign {
		buf = append(buf, ' ')
	}

	if d.ndig == 0 {
		buf = append(buf, '0')
	} else {
		buf = append(buf, d.dig[0])
	}

	if prec > 0 {
		buf = append(buf, '.')

		i := 0
		if d.ndig > 1 {
			buf = append(buf, d.dig[1:d.ndig]...)
			i = d.ndig - 1
		}

		for ; i < prec; i++ {
			buf = append(buf, '0')
		}
	} else if forceDP {
		buf = append(buf, '.')
	}

	buf = append(buf, e)

	exp := d.exp
	if d.ndig > 1 {
		exp += d.ndig - 1
	}

	if exp < 0 {
		exp = -exp
		buf = append(buf, '-')
	} else {
		buf = append(buf, '+')
	}

	if exp < 10 {
		if padExp {
			buf = append(buf, '0', '0'+byte(exp))
		} else {
			buf = append(buf, '0'+byte(exp))
		}
	} else if exp < 100 {
		buf = append(buf, '0'+byte(exp/10), '0'+byte(exp%10))
	} else if exp < 1000 {
		buf = append(buf, '0'+byte(exp/100), '0'+byte(exp/10%10), '0'+byte(exp%10))
	} else {
		buf = append(buf, '0'+byte(exp/1000), '0'+byte(exp/100%10), '0'+byte(exp/10%10), '0'+byte(exp%10))
	}

	buf = d.pad(buf, pad, printSign, padSign, padRight, padZero)
	return buf
}

func (d *digits) fmtF(buf []byte, prec, pad int, forceDP, printSign, padSign, padRight, padZero bool) []byte {

	// Attempt to pre-size buffer to avoid multiple allocations. Currently,
	// this might overshoot the actual needed size. Calculation is:
	// sign + decimal separator + signficant digits + padding.
	sizeHint := 1 + 1 + d.ndig
	if pad > sizeHint {
		sizeHint = pad
	}
	buf = buf[len(buf):]
	if cap(buf) < sizeHint {
		buf = append(buf, make([]byte, sizeHint)...)[:0]
	}

	if d.neg {
		buf = append(buf, '-')
	} else if printSign {
		buf = append(buf, '+')
	} else if padSign {
		buf = append(buf, ' ')
	}

	dp := 0
	if d.ndig == 0 {
		buf = append(buf, '0')
	} else {
		dp = d.ndig + d.exp

		if dp > 0 {
			if d.ndig > dp {
				buf = append(buf, d.dig[:dp]...)
			} else {
				buf = append(buf, d.dig[:d.ndig]...)

				for i := d.ndig; i < dp; i++ {
					buf = append(buf, '0')
				}
			}
		} else {
			buf = append(buf, '0')
		}
	}

	if prec > 0 {
		buf = append(buf, '.')

		for ; dp < 0; dp++ {
			prec--
			buf = append(buf, '0')
		}

		i := 0
		if d.ndig > dp {
			buf = append(buf, d.dig[dp:d.ndig]...)
			i = d.ndig - dp
		}

		for ; i < prec; i++ {
			buf = append(buf, '0')
		}
	} else if forceDP {
		buf = append(buf, '.')
	}

	buf = d.pad(buf, pad, printSign, padSign, padRight, padZero)
	return buf
}

// pad adds padding to the passed buf.
func (d *digits) pad(buf []byte, pad int, printSign, padSign, padRight, padZero bool) []byte {
	p := pad - len(buf)
	if p <= 0 {
		// No need for padding.
		return buf
	}

	padChar := byte(' ')
	if padZero {
		padChar = byte('0')
	}

	if padRight {
		for i := 0; i < p; i++ {
			buf = append(buf, padChar)
		}
	} else {
		// Determine where to keep the sign.
		i := 0
		if padZero && (d.neg || printSign || padSign) {
			i = 1
			p++
		}

		// Grow buf until it fits the nb + padding.
		for len(buf) < pad {
			buf = append(buf, 0)
		}

		// Move the existing number to the end of the buffer.
		copy(buf[p:], buf[i:])

		// Fill left-padding chars.
		for ; i < p; i++ {
			buf[i] = padChar
		}
	}

	return buf
}

func (d *digits) round(prec int) {
	if d.ndig <= prec {
		return
	}

	if prec < 0 {
		d.exp += d.ndig
		d.ndig = 0
		return
	}

	up := false
	if d.ndig > 1 && d.ndig == prec+1 && d.dig[prec] == '5' {
		up = (d.dig[prec-1]-'0')%2 != 0
	} else {
		up = d.dig[prec] >= '5'
	}

	if up {
		i := prec - 1
		for i >= 0 && d.dig[i] == '9' {
			i--
		}

		if i == -1 {
			d.dig[0] = '1'
			d.exp += d.ndig
			d.ndig = 1
		} else {
			d.dig[i]++
			prec = i + 1
			d.exp += d.ndig - prec
			d.ndig = prec
		}
	} else {
		i := prec - 1
		for i >= 0 && d.dig[i] == '0' {
			i--
		}

		prec = i + 1
		d.exp += d.ndig - prec
		d.ndig = prec
	}
}

// FmtFormat stores a formatting layout for Decimal values.
type FmtFormat struct {
	width    int
	hasWidth bool
	prec     int
	hasPrec  bool

	verb rune

	minus bool
	plus  bool
	sharp bool
	space bool
	zero  bool
}

func (f *FmtFormat) flag(b int) bool {
	switch b {
	case '-':
		return f.minus
	case '+':
		return f.plus
	case '#':
		return f.sharp
	case ' ':
		return f.space
	case '0':
		return f.zero
	}
	return false
}

// ParseFmtFormat parses the string as a decimal formatter.
//
// TODO: only a subset of formatting strings are supported.
func ParseFmtFormat(f string) (*FmtFormat, error) {
	if len(f) < 2 {
		return nil, fmt.Errorf("format string has len < 2")

	}
	if f[0] != '%' {
		return nil, fmt.Errorf("format string does not start with '%%'")
	}

	verb := rune(f[len(f)-1])
	switch verb {
	case 'f', 'g', 'e':
	default:
		return nil, fmt.Errorf("unsupported format verb '%s'", string(verb))
	}

	var width, prec int
	var hasWidth, hasPrec bool
	dp := strings.Index(f, ".")
	if dp > -1 && dp >= len(f)-2 {
		return nil, fmt.Errorf("format has decimal point without precision")

	}
	if dp == -1 {
		hasWidth = len(f) > 2
		dp = len(f) - 1
	} else {
		hasWidth = dp > 1
		hasPrec = true
	}

	if hasWidth {
		sWid := f[1:dp]
		w, err := strconv.ParseInt(sWid, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid width %q: %v", sWid, err)

		}
		width = int(w)
	}
	if hasPrec {
		sPrec := f[dp+1 : len(f)-1]
		p, err := strconv.ParseInt(sPrec, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid precision %q: %v", sPrec, err)

		}
		prec = int(p)
	}

	return &FmtFormat{
		prec:     prec,
		hasPrec:  hasPrec,
		width:    width,
		hasWidth: hasWidth,
		verb:     verb,
	}, nil
}

func fmtFormatFromFmtState(st fmt.State, verb rune, f *FmtFormat) {
	wid, hasWid := st.Width()
	prec, hasPrec := st.Precision()
	*f = FmtFormat{
		prec:     prec,
		hasPrec:  hasPrec,
		width:    wid,
		hasWidth: hasWid,
		verb:     verb,
		minus:    st.Flag('-'),
		plus:     st.Flag('+'),
		sharp:    st.Flag('#'),
		space:    st.Flag(' '),
		zero:     st.Flag('0'),
	}
}

func mustParseFmtFormat(f string) *FmtFormat {
	af, err := ParseFmtFormat(f)
	if err != nil {
		panic(err)

	}
	return af
}
