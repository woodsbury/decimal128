//go:build goexperiment.jsonv2

package decimal128

import (
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"strconv"
)

// MarshalJSONTo implements the [encoding/json/v2.MarshalerTo] interface.
func (d Decimal) MarshalJSONTo(enc *jsontext.Encoder) error {
	if d.isSpecial() {
		return fmt.Errorf("unsupported value: %v", d)
	}

	stringify, _ := json.GetOption(enc.Options(), json.StringifyNumbers)
	buf := enc.AvailableBuffer()
	if stringify {
		buf = append(buf, '"')
	}

	var digs digits
	d.digits(&digs)

	prec := 0
	if digs.ndig != 0 {
		prec = digs.ndig - 1
	}

	exp := digs.exp + prec

	if exp < -6 || exp >= 20 {
		buf = digs.fmtE(buf, prec, 0, false, false, false, false, false, false, 'e')
	} else {
		prec = 0
		if digs.exp < 0 {
			prec = -digs.exp
		}

		buf = digs.fmtF(buf, prec, 0, false, false, false, false, false)
	}

	if stringify {
		buf = append(buf, '"')
	}

	enc.WriteValue(buf)
	return nil
}

// UnmarshalJSONFrom implements the [encoding/json/v2.UnmarshalerFrom] interface.
func (d *Decimal) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	data, err := dec.ReadValue()
	if err != nil {
		return err
	}

	if string(data) == "null" {
		if legacy, _ := json.GetOption(dec.Options(), jsonv1.MergeWithLegacySemantics); legacy {
			return nil
		}

		*d = Decimal{}
		return nil
	}

	l := len(data)
	if l == 0 {
		if legacy, _ := json.GetOption(dec.Options(), jsonv1.MergeWithLegacySemantics); legacy {
			return nil
		}

		*d = Decimal{}
		return nil
	}

	var stringify bool
	if l >= 2 && data[0] == '"' && data[l-1] == '"' {
		stringify, _ = json.GetOption(dec.Options(), json.StringifyNumbers)
		if !stringify {
			return &json.SemanticError{
				JSONKind: '"',
			}
		}

		data = data[1 : l-1]
		l -= 2

		if l == 0 {
			if legacy, _ := json.GetOption(dec.Options(), jsonv1.MergeWithLegacySemantics); legacy {
				return nil
			}

			*d = Decimal{}
			return nil
		}
	}

	neg := false

	i := 0
	switch data[i] {
	case '+':
		i = 1
	case '-':
		neg = true
		i = 1
	}

	tmp, err := parseNumber(data[i:], neg, false)
	if err != nil {
		switch err.(type) {
		case parseNumberRangeError:
			return &json.SemanticError{
				JSONKind: '0',
				Err:      strconv.ErrRange,
			}
		case parseNumberSyntaxError:
			return &json.SemanticError{
				JSONKind: data.Kind(),
			}
		}
	}

	*d = tmp
	return nil
}
