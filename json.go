package decimal128

import (
	"encoding/json"
	"reflect"
)

// MarshalJSON implements the [encoding/json.Marshaler] interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.isSpecial() {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(d),
			Str:   d.String(),
		}
	}

	var digs digits
	d.digits(&digs)

	prec := 0
	if digs.ndig != 0 {
		prec = digs.ndig - 1
	}

	exp := digs.exp + prec

	if exp < -6 || exp >= 20 {
		return digs.fmtE(nil, prec, 0, false, false, false, false, false, false, 'e'), nil
	}

	prec = 0
	if digs.exp < 0 {
		prec = -digs.exp
	}

	return digs.fmtF(nil, prec, 0, false, false, false, false, false), nil
}

// UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	l := len(data)

	if l == 0 {
		return nil
	}

	neg := false

	i := 0
	if data[0] == '+' {
		i = 1
	} else if data[0] == '-' {
		neg = true
		i = 1
	}

	var sig64 uint64
	var nfrac int16
	var trunc int8
	caneof := false
	cansgn := false
	eneg := false
	sawdig := false
	sawdot := false
	sawexp := false

	for ; !sawexp && sig64 < 0x18ff_ffff_ffff_ffff && i < l; i++ {
		switch c := data[i]; true {
		case c >= '0' && c <= '9':
			caneof = true
			cansgn = false
			sawdig = true

			sig64 = sig64*10 + uint64(c-'0')

			if sawdot {
				nfrac++
			}
		case c == '.':
			if sawdot {
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = true
			cansgn = false
			sawdot = true
		case c == 'E' || c == 'e':
			if !sawdig {
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = true
			cansgn = true
			sawexp = true
		default:
			err := &json.UnmarshalTypeError{
				Type: reflect.TypeOf(Decimal{}),
			}

			switch data[0] {
			case 'f', 't':
				err.Value = "bool"
			case '"':
				err.Value = "string"
			default:
				err.Value = "number " + string(data)
			}

			return err
		}
	}

	sig := uint128{sig64, 0}
	var exp int16
	maxexp := false

	for ; i < l; i++ {
		switch c := data[i]; true {
		case c >= '0' && c <= '9':
			caneof = true
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
						c2 := data[i+1]
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
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = true
			cansgn = false
			sawdot = true
		case c == 'E' || c == 'e':
			if !sawdig || sawexp {
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = true
			cansgn = true
			sawexp = true
		case c == '-':
			if !cansgn {
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = false
			cansgn = false
			eneg = true
		case c == '+':
			if !cansgn {
				return &json.UnmarshalTypeError{
					Value: "number " + string(data),
					Type:  reflect.TypeOf(Decimal{}),
				}
			}

			caneof = false
			cansgn = false
		default:
			err := &json.UnmarshalTypeError{
				Type: reflect.TypeOf(Decimal{}),
			}

			switch data[0] {
			case 'f', 't':
				err.Value = "bool"
			case '"':
				err.Value = "string"
			default:
				err.Value = "number " + string(data)
			}

			return err
		}
	}

	if !caneof {
		return &json.UnmarshalTypeError{
			Value: "number " + string(data),
			Type:  reflect.TypeOf(Decimal{}),
		}
	}

	// If the exponent value is larger than the maximum supported exponent, there are two cases
	// where the value is still valid:
	//  - the exponent is negative, where the logical value rounds to 0
	//  - the significand is zero, where the logical value is 0
	//
	// Otherwise, return a range error.
	if maxexp {
		if eneg {
			*d = zero(neg)
			return nil
		}
		if sig[0] == 0 && sig[1] == 0 {
			*d = zero(neg)
			return nil
		}

		return &json.UnmarshalTypeError{
			Value: "number " + string(data),
			Type:  reflect.TypeOf(Decimal{}),
		}
	}

	if eneg {
		exp *= -1
	}

	exp -= nfrac

	if exp > maxUnbiasedExponent+39 {
		return &json.UnmarshalTypeError{
			Value: "number " + string(data),
			Type:  reflect.TypeOf(Decimal{}),
		}
	}

	if exp < minUnbiasedExponent-39 {
		*d = zero(neg)
		return nil
	}

	sig, exp = DefaultRoundingMode.reduce128(neg, sig, exp+exponentBias, trunc)

	if exp > maxBiasedExponent {
		return &json.UnmarshalTypeError{
			Value: "number " + string(data),
			Type:  reflect.TypeOf(Decimal{}),
		}
	}

	*d = compose(neg, sig, exp)
	return nil
}
