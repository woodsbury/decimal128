package decimal128

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// MustParse is like [Parse] but panics if the provided string cannot be parsed,
// instead of returning an error.
func MustParse(s string) Decimal {
	d, err := parse(s, payloadOpMustParse)
	if err != nil {
		panic("decimal128.MustParse(" + strconv.Quote(s) + "): invalid syntax")
	}

	return d
}

// Parse parses a Decimal value from the string provided. Parse accepts decimal
// floating point syntax. An underscore character '_' may appear between digits
// as a separator. Parse also recognises the string "NaN", and the (possibly
// signed) strings "Inf" and "Infinity", as their respective special floating
// point values. It ignores case when matching.
//
// If s is not syntactically well-formed, Parse returns an error that can be
// compared to [strconv.ErrSyntax] via [errors.Is].
//
// If the value is too precise to fit in a Decimal the result is rounded using
// the [DefaultRoundingMode]. If the value is greater than the largest possible
// Decimal value, Parse returns ±Inf and an error that can be compared to
// [strconv.ErrRange] via [errors.Is].
func Parse(s string) (Decimal, error) {
	return parse(s, payloadOpParse)
}

// Scan implements the [fmt.Scanner] interface. It supports the verbs 'e', 'E',
// 'f', 'F', 'g', 'G', and 'v'.
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
			return &scanSyntaxError{}
		}

		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'F' && r != 'f' {
			return &scanSyntaxError{}
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
			return &scanSyntaxError{}
		}

		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		if r != 'N' && r != 'n' {
			return &scanSyntaxError{}
		}

		*d = nan(payloadOpScan, 0, 0)
		return nil
	}

	f.UnreadRune()

	var sig64 uint64
	var nfrac int16
	var trunc int8
	caneof := false
	cansep := false
	cansgn := false
	eneg := false
	sawdig := false
	sawdot := false
	saweof := false
	sawexp := false

ReadRunes64:
	for !sawexp && sig64 < 0x18ff_ffff_ffff_ffff {
		r, _, err = f.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if caneof {
					saweof = true
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

			sig64 = sig64*10 + uint64(r-'0')

			if sawdot {
				nfrac++
			}
		case r == '.':
			if sawdot {
				f.UnreadRune()
				saweof = true

				break ReadRunes64
			}

			caneof = true
			cansep = false
			cansgn = false
			sawdot = true
		case r == 'E' || r == 'e':
			if !sawdig {
				return &scanSyntaxError{}
			}

			caneof = false
			cansep = false
			cansgn = true
			sawexp = true
		case r == '_':
			if !cansep {
				return &scanSyntaxError{}
			}

			caneof = false
			cansep = false
			cansgn = false
		default:
			f.UnreadRune()
			saweof = true

			break ReadRunes64
		}
	}

	sig := uint128{sig64, 0}
	var exp int16
	maxexp := false

	if !saweof {
	ReadRunes:
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
					if exp > exponentBias/10+1 {
						maxexp = true
					}

					exp *= 10
					exp += int16(r - '0')
				} else {
					if sig[1] <= 0x18ff_ffff_ffff_ffff {
						if sig[1] <= 0x027f_ffff_ffff_ffff {
							var r2 rune
							r2, _, err = f.ReadRune()
							if err == nil {
								if r2 >= '0' && r2 <= '9' {
									sig = sig.mul64(100)
									sig = sig.add64(uint64(r-'0')*10 + uint64(r2-'0'))

									if sawdot {
										nfrac += 2
									}

									continue
								} else {
									if err = f.UnreadRune(); err != nil {
										return err
									}
								}
							} else if !errors.Is(err, io.EOF) {
								return err
							}
						}

						sig = sig.mul64(10)
						sig = sig.add64(uint64(r - '0'))

						if sawdot {
							nfrac++
						}
					} else {
						if r != '0' {
							trunc = 1
						}

						if !sawdot {
							if exp < exponentBias+39 {
								nfrac--
							}
						}
					}
				}
			case r == '.':
				if sawdot || sawexp {
					f.UnreadRune()
					saweof = true

					break ReadRunes
				}

				caneof = true
				cansep = false
				cansgn = false
				sawdot = true
			case r == 'E' || r == 'e':
				if !sawdig {
					return &scanSyntaxError{}
				}

				if sawexp {
					f.UnreadRune()
					saweof = true

					break ReadRunes
				}

				caneof = false
				cansep = false
				cansgn = true
				sawexp = true
			case r == '-':
				if !cansgn {
					return &scanSyntaxError{}
				}

				caneof = false
				cansep = false
				cansgn = false
				eneg = true
			case r == '_':
				if !cansep {
					return &scanSyntaxError{}
				}

				caneof = false
				cansep = false
				cansgn = false
			case r == '+':
				if !cansgn {
					return &scanSyntaxError{}
				}

				caneof = false
				cansep = false
				cansgn = false
			default:
				f.UnreadRune()
				saweof = true

				break ReadRunes
			}
		}
	}

	if !caneof {
		return &scanSyntaxError{}
	}

	// If the exponent value is larger than the maximum supported exponent,
	// there are two cases where the value is still valid:
	//  - the exponent is negative, where the logical value rounds to 0
	//  - the significand is zero, where the logical value is 0
	//
	// Otherwise, return a range error.
	if maxexp {
		if eneg {
			*d = zero(neg)
			return nil
		}

		if sig == (uint128{}) {
			*d = zero(neg)
			return nil
		}

		return &scanRangeError{}
	}

	if eneg {
		exp *= -1
	}

	exp -= nfrac

	if exp > maxUnbiasedExponent+39 {
		return &scanRangeError{}
	}

	if exp < minUnbiasedExponent-39 {
		*d = zero(neg)
		return nil
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return &scanRangeError{}
	}

	*d = compose(neg, sig, exp)
	return nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (d *Decimal) UnmarshalText(data []byte) error {
	tmp, err := parse(data, payloadOpUnmarshalText)
	if err != nil {
		return err
	}

	*d = tmp
	return nil
}

func parse[D []byte | string](d D, op Payload) (Decimal, error) {
	if len(d) == 0 {
		return Decimal{}, &parseSyntaxError{string(d)}
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
		return Decimal{}, &parseSyntaxError{string(d)}
	} else if l == 3 {
		if (d[0] == 'I' || d[0] == 'i') && (d[1] == 'N' || d[1] == 'n') && (d[2] == 'F' || d[2] == 'f') {
			return inf(neg), nil
		}

		if (d[0] == 'N' || d[0] == 'n') && (d[1] == 'A' || d[1] == 'a') && (d[2] == 'N' || d[2] == 'n') {
			return nan(op, 0, 0), nil
		}
	} else if l == 8 {
		if (d[0] == 'I' || d[0] == 'i') && (d[1] == 'N' || d[1] == 'n') && (d[2] == 'F' || d[2] == 'f') && (d[3] == 'I' || d[3] == 'i') && (d[4] == 'N' || d[4] == 'n') && (d[5] == 'I' || d[5] == 'i') && (d[6] == 'T' || d[6] == 't') && (d[7] == 'Y' || d[7] == 'y') {
			return inf(neg), nil
		}
	}

	var sig64 uint64
	var nfrac int16
	var trunc int8
	caneof := false
	cansep := false
	cansgn := false
	eneg := false
	sawdig := false
	sawdot := false
	sawexp := false

	i := 0
	for ; !sawexp && sig64 <= 0x18ff_ffff_ffff_ffff && i < l; i++ {
		switch c := d[i]; true {
		case c >= '0' && c <= '9':
			caneof = true
			cansep = true
			cansgn = false
			sawdig = true

			sig64 = sig64*10 + uint64(c-'0')

			if sawdot {
				nfrac++
			}
		case c == '.':
			if sawdot {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = true
			cansep = false
			cansgn = false
			sawdot = true
		case c == 'E' || c == 'e':
			if !sawdig {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = true
			sawexp = true
		case c == '_':
			if !cansep {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
		default:
			return Decimal{}, &parseSyntaxError{string(d)}
		}
	}

	sig := uint128{sig64, 0}
	var exp int16
	maxexp := false

	for ; i < l; i++ {
		switch c := d[i]; true {
		case c >= '0' && c <= '9':
			caneof = true
			cansep = true
			cansgn = false
			sawdig = true

			if sawexp {
				if exp > exponentBias/10+1 {
					maxexp = true
				}

				exp *= 10
				exp += int16(c - '0')
			} else {
				if sig[1] <= 0x18ff_ffff_ffff_ffff {
					if sig[1] <= 0x027f_ffff_ffff_ffff && i < l-1 {
						c2 := d[i+1]
						if c2 >= '0' && c2 <= '9' {
							sig = sig.mul64(100)
							sig = sig.add64(uint64(c-'0')*10 + uint64(c2-'0'))

							if sawdot {
								nfrac += 2
							}

							i++
							continue
						}
					}

					sig = sig.mul64(10)
					sig = sig.add64(uint64(c - '0'))

					if sawdot {
						nfrac++
					}
				} else {
					if c != '0' {
						trunc = 1
					}

					if !sawdot {
						if exp < exponentBias+39 {
							nfrac--
						}
					}
				}
			}
		case c == '.':
			if sawdot || sawexp {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = true
			cansep = false
			cansgn = false
			sawdot = true
		case c == 'E' || c == 'e':
			if !sawdig || sawexp {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = true
			sawexp = true
		case c == '-':
			if !cansgn {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
			eneg = true
		case c == '_':
			if !cansep {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
		case c == '+':
			if !cansgn {
				return Decimal{}, &parseSyntaxError{string(d)}
			}

			caneof = false
			cansep = false
			cansgn = false
		default:
			return Decimal{}, &parseSyntaxError{string(d)}
		}
	}

	if !caneof {
		return Decimal{}, &parseSyntaxError{string(d)}
	}

	// If the exponent value is larger than the maximum supported exponent,
	// there are two cases where the value is still valid:
	//  - the exponent is negative, where the logical value rounds to 0
	//  - the significand is zero, where the logical value is 0
	//
	// Otherwise, return a range error.
	if maxexp {
		if eneg {
			return zero(neg), nil
		}

		if sig == (uint128{}) {
			return zero(neg), nil
		}

		return inf(neg), &parseRangeError{string(d)}
	}

	if eneg {
		exp *= -1
	}

	exp -= nfrac

	if exp > maxUnbiasedExponent+39 {
		return inf(neg), &parseRangeError{string(d)}
	}

	if exp < minUnbiasedExponent-39 {
		return zero(neg), nil
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return inf(neg), &parseRangeError{string(d)}
	}

	return compose(neg, sig, exp), nil
}

type parseRangeError struct {
	s string
}

func (err *parseRangeError) Error() string {
	return "parsing " + strconv.Quote(err.s) + ": value out of range"
}

func (err *parseRangeError) Is(target error) bool {
	return target == strconv.ErrRange
}

type parseSyntaxError struct {
	s string
}

func (err *parseSyntaxError) Error() string {
	return "parsing " + strconv.Quote(err.s) + ": invalid syntax"
}

func (err *parseSyntaxError) Is(target error) bool {
	return target == strconv.ErrSyntax
}

type scanRangeError struct{}

func (err *scanRangeError) Error() string {
	return "parsing decimal: value out of range"
}

func (err *scanRangeError) Is(target error) bool {
	return target == strconv.ErrRange
}

type scanSyntaxError struct{}

func (err *scanSyntaxError) Error() string {
	return "parsing decimal: invalid syntax"
}

func (err *scanSyntaxError) Is(target error) bool {
	return target == strconv.ErrSyntax
}
