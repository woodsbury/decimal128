package decimal128

import (
	"encoding/json"
	"reflect"
)

func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.isSpecial() {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(d),
			Str:   d.String(),
		}
	}

	panic("not implemented")
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	panic("not implemented")
}
