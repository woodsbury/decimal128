package decimal128

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

func Parse(s string) (Decimal, error) {
	return parse(s)
}

func (d *Decimal) Scan(f fmt.ScanState, verb rune) error {
	switch verb {
	case 'e', 'E', 'f', 'F', 'g', 'G', 'v':
	default:
		return errors.New("bad verb '%" + string(verb) + "' for Decimal")
	}

	f.SkipSpace()
	r, _, err := f.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return io.ErrUnexpectedEOF
		}

		return err
	}

	neg := false

	if r == '-' {
		neg = true
	} else if r != '+' {
		f.UnreadRune()
	}

	r, _, err = f.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return io.ErrUnexpectedEOF
		}

		return err
	}

	if r == 'I' || r == 'i' {
		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'N' && r != 'n' {
			return &scanError{}
		}

		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'F' && r != 'f' {
			return &scanError{}
		}

		*d = inf(neg)
		return nil
	}

	if r == 'N' || r == 'n' {
		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'A' && r != 'a' {
			return &scanError{}
		}

		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'N' && r != 'n' {
			return &scanError{}
		}

		*d = nan()
		return nil
	}

	f.UnreadRune()

	var sig uint128
	var exp int16
	var nfrac int16
	var trunc int8
	caneof := false
	cansep := false
	cansgn := false
	eneg := false
	sawdig := false
	sawdot := false
	sawexp := false

	for {
		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if caneof {
					break
				}

				return io.ErrUnexpectedEOF
			}

			return err
		}

		switch true {
		case r >= '0' && r <= '9':
			caneof = true
			cansep = true
			cansgn = false
			sawdig = true

			if sawexp {
				if exp < exponentBias+39 {
					exp *= 10
					exp += int16(r - '0')
				}
			} else {
				if sig[1] <= 0x18ff_ffff_ffff_ffff {
					sig = sig.mul64(10)
					sig = sig.add64(uint64(r - '0'))

					if sawdot {
						nfrac++
					}
				} else if r != '0' {
					trunc = 1
				}
			}
		case r == '.':
			if sawdot || sawexp {
				return &scanError{}
			}

			caneof = true
			cansep = false
			cansgn = false
			sawdot = true
		case r == 'E' || r == 'e':
			if !sawdig || sawexp {
				return &scanError{}
			}

			caneof = false
			cansep = false
			cansgn = true
			sawexp = true
		case r == '+':
			if !cansgn {
				return &scanError{}
			}

			caneof = false
			cansep = false
			cansgn = false
		case r == '-':
			if !cansgn {
				return &scanError{}
			}

			caneof = false
			cansep = false
			cansgn = false
			eneg = true
		case r == '_':
			if !cansep {
				return &scanError{}
			}

			caneof = false
			cansep = false
			cansgn = false
		}
	}

	if !caneof {
		return &scanError{}
	}

	if eneg {
		exp *= -1
	}

	exp -= nfrac

	if exp > maxBiasedExponent-exponentBias+39 {
		*d = inf(neg)
		return nil
	}

	if exp < minBiasedExponent-exponentBias-39 {
		*d = zero(neg)
		return nil
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		*d = inf(neg)
		return nil
	}

	*d = compose(neg, sig, exp)
	return nil
}

func (d *Decimal) UnmarshalText(data []byte) error {
	tmp, err := parse(data)
	if err != nil {
		return err
	}

	*d = tmp
	return nil
}

func parse[D []byte | string](d D) (Decimal, error) {
	if len(d) == 0 {
		return Decimal{}, &parseError{string(d)}
	}

	neg := false

	if d[0] == '+' {
		d = d[1:]
	} else if d[0] == '-' {
		neg = true
		d = d[1:]
	}

	l := len(d)

	if l == 0 {
		return Decimal{}, &parseError{string(d)}
	} else if l == 3 {
		if (d[0] == 'I' || d[0] == 'i') && (d[1] == 'N' || d[1] == 'n') && (d[2] == 'F' || d[2] == 'f') {
			return inf(neg), nil
		}

		if (d[0] == 'N' || d[0] == 'n') && (d[1] == 'A' || d[1] == 'a') && (d[2] == 'N' || d[2] == 'n') {
			return nan(), nil
		}
	} else if l == 8 {
		if (d[0] == 'I' || d[0] == 'i') && (d[1] == 'N' || d[1] == 'n') && (d[2] == 'F' || d[2] == 'f') && (d[3] == 'I' || d[3] == 'i') && (d[4] == 'N' || d[4] == 'n') && (d[5] == 'I' || d[5] == 'i') && (d[6] == 'T' || d[6] == 't') && (d[7] == 'Y' || d[7] == 'y') {
			return inf(neg), nil
		}
	}

	var sig uint128
	var exp int16
	var nfrac int16
	var trunc int8
	caneof := false
	cansep := false
	cansgn := false
	eneg := false
	sawdig := false
	sawdot := false
	sawexp := false

	for i := 0; i < l; i++ {
		switch c := d[i]; true {
		case c >= '0' && c <= '9':
			caneof = true
			cansep = true
			cansgn = false
			sawdig = true

			if sawexp {
				if exp < exponentBias+39 {
					exp *= 10
					exp += int16(c - '0')
				}
			} else {
				if sig[1] <= 0x18ff_ffff_ffff_ffff {
					sig = sig.mul64(10)
					sig = sig.add64(uint64(c - '0'))

					if sawdot {
						nfrac++
					}
				} else if c != '0' {
					trunc = 1
				}
			}
		case c == '.':
			if sawdot || sawexp {
				return Decimal{}, &parseError{string(d)}
			}

			caneof = true
			cansep = false
			cansgn = false
			sawdot = true
		case c == 'E' || c == 'e':
			if !sawdig || sawexp {
				return Decimal{}, &parseError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = true
			sawexp = true
		case c == '+':
			if !cansgn {
				return Decimal{}, &parseError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
		case c == '-':
			if !cansgn {
				return Decimal{}, &parseError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
			eneg = true
		case c == '_':
			if !cansep {
				return Decimal{}, &parseError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
		default:
			return Decimal{}, &parseError{string(d)}
		}
	}

	if !caneof {
		return Decimal{}, &parseError{string(d)}
	}

	if eneg {
		exp *= -1
	}

	exp -= nfrac

	if exp > maxBiasedExponent-exponentBias+39 {
		return inf(neg), nil
	}

	if exp < minBiasedExponent-exponentBias-39 {
		return zero(neg), nil
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(neg), nil
	}

	return compose(neg, sig, exp), nil
}

type parseError struct {
	s string
}

func (err *parseError) Error() string {
	return "parsing " + strconv.Quote(err.s) + ": invalid syntax"
}

type scanError struct{}

func (err *scanError) Error() string {
	return "expected decimal"
}
