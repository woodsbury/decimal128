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

	if len(data) == 0 {
		return nil
	}

	// Pass an empty NaN payload because NaN values are not valid in JSON, so we should never
	// encounter a NaN.
	dec, err := parse(data, 0)
	if err != nil {
		// If there was a parse error, check if the first byte of the input looks like any known JSON values and
		// return a more informative JSON unmarshal error. Otherwise, return the original parse error.
		switch data[0] {
		case 'f', 't':
			return &json.UnmarshalTypeError{
				Value: "bool",
				Type:  reflect.TypeOf(Decimal{}),
			}
		case '"':
			return &json.UnmarshalTypeError{
				Value: "string",
				Type:  reflect.TypeOf(Decimal{}),
			}
		case '{':
			return &json.UnmarshalTypeError{
				Value: "JSON object",
				Type:  reflect.TypeOf(Decimal{}),
			}
		case '[':
			return &json.UnmarshalTypeError{
				Value: "JSON array",
				Type:  reflect.TypeOf(Decimal{}),
			}
		default:
			return err
		}
	}

	*d = dec
	return nil
}
