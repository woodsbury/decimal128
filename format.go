package decimal128

import (
	"fmt"
	"unsafe"
)

var (
	nanText    = []byte("NaN")
	padNaNText = []byte(" NaN")
	posNaNText = []byte("+NaN")
	negInfText = []byte("-Inf")
	padInfText = []byte(" Inf")
	posInfText = []byte("+Inf")
	spaceText  = []byte{' ', ' ', ' '}

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

// Append formats the Decimal according to the provided format specifier and
// appends the result to the provided byte slice, returning the updated byte
// slice. The format specifier can be any value supported by [Decimal.Format], without
// the leading %.
func (d Decimal) Append(buf []byte, format string) []byte {
	var args formatArgs
	parseFormat(format, &args)

	if args.verb == 0 {
		return append(buf, "%!(NOVERB)"...)
	}

	if d.isSpecial() {
		pad := 0
		padSign := false
		printSign := false
		if args.verb != 'v' {
			printSign = args.printSign
			padSign = args.padSign

			width, hasWidth := args.width()
			if hasWidth {
				pad = width
			}
		}

		return d.appendSpecial(buf, pad, printSign, padSign, args.padRight)
	}

	return d.format(buf, &args)
}

// Format implements the [fmt.Formatter] interface. It supports the verbs 'e',
// 'E', 'f', 'F', 'g', 'G', and 'v', along with the format flags '+', '-', '#',
// ' ', and '0' and custom width and precision values. Decimal values interpret
// the format value the same way float32 and float64 does.
func (d Decimal) Format(f fmt.State, verb rune) {
	if d.isSpecial() {
		pad := 0
		padSign := false
		printSign := false
		if verb != 'v' {
			printSign = f.Flag('+')
			padSign = f.Flag(' ')

			width, hasWidth := f.Width()
			if hasWidth {
				pad = width
			}
		}

		d.writeSpecial(f, pad, printSign, padSign, f.Flag('-'))
		return
	}

	prec, hasPrec := f.Precision()
	if !hasPrec {
		prec = -1
	}

	width, hasWidth := f.Width()
	if !hasWidth {
		width = -1
	}

	args := formatArgs{
		forceDP:   f.Flag('#'),
		printSign: f.Flag('+'),
		padSign:   f.Flag(' '),
		padRight:  f.Flag('-'),
		padZero:   f.Flag('0'),
		verb:      byte(verb),
		prec:      prec,
		wid:       width,
	}

	f.Write(d.format(nil, &args))
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (d Decimal) MarshalText() ([]byte, error) {
	if d.isSpecial() {
		return d.appendSpecial(nil, 0, false, false, false), nil
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
	var buf []byte
	if d.isSpecial() {
		buf = d.appendSpecial(buf, 0, false, false, false)
	} else {
		var digs digits
		d.digits(&digs)

		prec := 0
		if digs.ndig != 0 {
			prec = digs.ndig - 1
		}

		exp := digs.exp + prec

		if exp < -4 || exp >= 6 {
			buf = digs.fmtE(buf, prec, 0, false, false, false, true, false, false, 'e')
		} else {
			prec = 0
			if digs.exp < 0 {
				prec = -digs.exp
			}

			buf = digs.fmtF(buf, prec, 0, false, false, false, false, false)
		}
	}

	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

func (d Decimal) appendSpecial(buf []byte, pad int, printSign, padSign, padRight bool) []byte {
	var value []byte
	if d.IsNaN() {
		if printSign {
			value = posNaNText
		} else if padSign {
			value = padNaNText
		} else {
			value = nanText
		}
	} else {
		if d.Signbit() {
			value = negInfText
		} else {
			if padSign && !printSign {
				value = padInfText
			} else {
				value = posInfText
			}
		}
	}

	if cap(buf) == 0 {
		sizeHint := len(value)
		if pad > sizeHint {
			sizeHint = pad
		}

		buf = make([]byte, 0, sizeHint)
	}

	n := len(value)
	if p := pad - n; p > 0 {
		if padRight {
			buf = append(buf, value...)

			for i := n; i < pad; i++ {
				buf = append(buf, ' ')
			}
		} else {
			for i := 0; i < p; i++ {
				buf = append(buf, ' ')
			}

			buf = append(buf, value...)
		}
	} else {
		buf = append(buf, value...)
	}

	return buf
}

func (d Decimal) digits(digs *digits) {
	*digs = digits{}
	digs.neg = d.Signbit()

	sig, exp := d.decompose()

	if sig[0]|sig[1] != 0 {
		digs.exp = int(exp - exponentBias)

		n := 0
		for sig[1] != 0 {
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

			if pair[0] == '0' && sig[0]|sig[1] == 0 {
				digs.dig[n] = pair[1]
				n++
			} else {
				digs.dig[n], digs.dig[n+1] = pair[1], pair[0]
				n += 2
			}
		}

		sig64 := sig[0]

		for sig64 != 0 {
			rem := sig64 % 100
			sig64 /= 100

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

			if pair[0] == '0' && sig64 == 0 {
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

func (d Decimal) format(buf []byte, args *formatArgs) []byte {
	var digs digits
	d.digits(&digs)

	prec, hasPrec := args.precision()
	width, hasWidth := args.width()

	switch args.verb {
	case 'e', 'E':
		if !hasPrec {
			prec = 6
		}

		pad := 0
		if hasWidth {
			pad = width
		}

		digs.round(prec + 1)
		return digs.fmtE(buf, prec, pad, args.forceDP, args.printSign, args.padSign, true, args.padRight, args.padZero, args.verb)
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

		return digs.fmtF(buf, prec, pad, args.forceDP, args.printSign, args.padSign, args.padRight, args.padZero)
	case 'g', 'G':
		var maxprec int
		if args.forceDP {
			if !hasPrec {
				if digs.ndig < 6 {
					prec = 6
				} else {
					prec = digs.ndig
				}

				maxprec = 6
			} else {
				if prec == 0 {
					prec = 1
				}

				maxprec = prec
			}

			digs.round(prec)
		} else {
			if hasPrec {
				if prec == 0 {
					prec = 1
				}

				digs.round(prec)
				maxprec = prec
				prec = digs.ndig
			} else if digs.ndig != 0 {
				maxprec = 6
				prec = digs.ndig
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
			if args.verb == 'G' {
				e = byte('E')
			}

			return digs.fmtE(buf, prec-1, pad, args.forceDP, args.printSign, args.padSign, true, args.padRight, args.padZero, e)
		} else {
			if args.forceDP {
				prec -= digs.exp
				if digs.ndig == 0 {
					prec--
				} else {
					prec -= digs.ndig
				}
			} else {
				prec = 0
				if digs.exp < 0 {
					prec -= digs.exp
				}
			}

			return digs.fmtF(buf, prec, pad, args.forceDP, args.printSign, args.padSign, args.padRight, args.padZero)
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
		}

		return digs.fmtF(buf, prec, 0, false, false, false, false, false)
	default:
		return fmt.Appendf(buf, "%%!%c(decimal128.Decimal=%s)", args.verb, d.String())
	}
}

func (d Decimal) writeSpecial(f fmt.State, pad int, printSign, padSign, padRight bool) {
	var value []byte
	if d.IsNaN() {
		if printSign {
			value = posNaNText
		} else if padSign {
			value = padNaNText
		} else {
			value = nanText
		}
	} else {
		if d.Signbit() {
			value = negInfText
		} else {
			if padSign && !printSign {
				value = padInfText
			} else {
				value = posInfText
			}
		}
	}

	n := len(value)
	if p := pad - n; p > 0 {
		if padRight {
			f.Write(value)

			i := n
			for ; i < pad-2; i += 3 {
				f.Write(spaceText)
			}

			if i < pad-1 {
				f.Write(spaceText[:2])
			} else if i < pad {
				f.Write(spaceText[:1])
			}
		} else {
			i := 0
			for ; i < p-2; i += 3 {
				f.Write(spaceText)
			}

			if i < p-1 {
				f.Write(spaceText[:2])
			} else if i < p {
				f.Write(spaceText[:1])
			}

			f.Write(value)
		}
	} else {
		f.Write(value)
	}
}

type digits struct {
	neg  bool
	dig  [39]byte
	exp  int
	ndig int
}

func (d *digits) fmtE(buf []byte, prec, pad int, forceDP, printSign, padSign, padExp, padRight, padZero bool, e byte) []byte {
	if cap(buf) == 0 {
		// Attempt to pre-size buffer to avoid multiple allocations. This might
		// overshoot the actual needed size. Calculation is:
		// sign + decimal point + 'e+/-' + exponent + zero + digits
		sizeHint := 1 + 1 + 2 + 4 + 1 + d.ndig
		if pad > sizeHint {
			sizeHint = pad
		}

		buf = make([]byte, 0, sizeHint)
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
	if cap(buf) == 0 {
		// Attempt to pre-size buffer to avoid multiple allocations. This might
		// overshoot the actual needed size. Calculation is:
		// sign + decimal point + zero + digits
		sizeHint := 1 + 1 + 1 + d.ndig
		if pad > sizeHint {
			sizeHint = pad
		}

		buf = make([]byte, 0, sizeHint)
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
		if len(buf) < pad {
			if cap(buf) < pad {
				tmp := make([]byte, pad)
				copy(tmp, buf)
				buf = tmp
			} else {
				buf = buf[:pad]
			}
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

type formatArgs struct {
	forceDP   bool
	printSign bool
	padSign   bool
	padRight  bool
	padZero   bool
	verb      byte
	prec      int
	wid       int
}

func parseFormat(format string, args *formatArgs) {
	*args = formatArgs{
		prec: -1,
		wid:  -1,
	}

	var c byte
	i := 0
	end := len(format)

parseFlags:
	for ; i < end; i++ {
		c = format[i]
		switch c {
		case ' ':
			args.padSign = true
		case '#':
			args.forceDP = true
		case '+':
			args.printSign = true
		case '-':
			args.padRight = true
			args.padZero = false
		case '0':
			args.padZero = !args.padRight
		default:
			break parseFlags
		}
	}

	if i >= end {
		return
	}

	if c >= '1' && c <= '9' {
		args.wid = int(c - '0')
		i++

		for ; i < end; i++ {
			c = format[i]
			if c < '0' || c > '9' {
				break
			}

			if args.wid < 1e5 {
				args.wid = args.wid*10 + int(c-'0')
			} else {
				args.wid = -1
			}
		}

		if i >= end {
			return
		}
	}

	if c == '.' {
		i++
		if i >= end {
			args.prec = 0
			return
		}

		c = format[i]
		if c < '0' || c > '9' {
			args.prec = 0
		} else {
			args.prec = int(c - '0')
			i++

			for ; i < end; i++ {
				c = format[i]
				if c < '0' || c > '9' {
					break
				}

				if args.prec < 1e5 {
					args.prec = args.prec*10 + int(c-'0')
				} else {
					args.prec = -1
				}
			}

			if i >= end {
				return
			}
		}
	}

	if i != end-1 {
		return
	}

	args.verb = c
}

func (args formatArgs) precision() (int, bool) {
	if args.prec < 0 {
		return 0, false
	}

	return args.prec, true
}

func (args formatArgs) width() (int, bool) {
	if args.wid < 0 {
		return 0, false
	}

	return args.wid, true
}
