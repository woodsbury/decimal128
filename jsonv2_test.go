//go:build goexperiment.jsonv2

package decimal128

import (
	"bytes"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"strings"
	"testing"
)

func TestDecimalMarshalJSONTo(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		if val.form != regularForm {
			continue
		}

		decval := val.Decimal()
		res, err := json.Marshal(decval)

		if err != nil {
			t.Errorf("%v.MarshalJSONTo() = (%s, %v), want (%s, <nil>", val, res, err, res)
		}

		var resval Decimal
		err = json.Unmarshal(res, &resval)

		if !resval.Equal(decval) || err != nil {
			t.Errorf("Decimal.UnmarshalJSONFrom(%s) = (%v, %v), want (%v, <nil>)", res, resval, err, decval)
		}

		if fltval, ok := val.Float64(); ok {
			fltres, err := json.Marshal(fltval)

			if err != nil {
				t.Errorf("json.Marshal(%v) = (%s, %v), want (%s, <nil>)", fltval, fltres, err, res)
			}

			if string(fltres) != string(res) {
				t.Errorf("%v.MarshalJSONTo() = (%s, <nil>), want (%s, <nil>)", val, res, fltres)
			}
		}
	}

	type S struct {
		D Decimal `json:",string"`
	}

	s := S{D: New(12345, -2)}
	res, err := json.Marshal(s)
	if string(res) != `{"D":"123.45"}` || err != nil {
		t.Errorf("json.Marshal(%v) = (%s, %v), want ({\"D\":\"123.45\"}, <nil>)", s, res, err)
	}
}

func TestDecimalUnmarshalJSONFrom(t *testing.T) {
	t.Parallel()

	for val, num := range textValues {
		var res Decimal
		err := json.Unmarshal([]byte(val), &res)

		if num.isInf() || num.IsNaN() || strings.Contains(val, "_") || strings.HasPrefix(val, "00") || strings.HasPrefix(val, "+") {
			if err == nil {
				t.Errorf("Decimal.UnarshalJSONFrom(%s) = (0, <nil>), want (%v, cannot unmarshal)", val, res)
			}
		} else if !res.Equal(num) || err != nil {
			t.Errorf("Decimal.UnmarshalJSONFrom(%s) = (%v, %v), want (%v, <nil>)", val, res, err, num)
		}
	}

	num := New(123, -1)
	res := num
	err := json.Unmarshal([]byte("null"), &res)
	if !res.IsZero() || err != nil {
		t.Errorf("Decimal.UnmarshalJSONFrom(null) = (%v, %v), want (%v, <nil>)", res, err, num)
	}

	res = num
	err = json.Unmarshal([]byte("null"), &res, jsonv1.MergeWithLegacySemantics(true))
	if !res.Equal(num) || err != nil {
		t.Errorf("Decimal.UnmarshalJSONFrom(null, MergeWithLegacySemantics) = (%v, %v), want (%v, <nil>)", res, err, num)
	}

	err = json.Unmarshal([]byte("[]"), &res)
	jerr := new(json.SemanticError)
	if !errors.As(err, &jerr) || jerr.JSONKind != '[' {
		t.Errorf("Decimal.UnmarshalJSONFrom([]) = %v, want json.SemanticError", err)
	}

	err = json.Unmarshal([]byte("{}"), &res)
	if !errors.As(err, &jerr) || jerr.JSONKind != '{' {
		t.Errorf("Decimal.UnmarshalJSONFrom({}) = %v, want json.SemanticError", err)
	}

	err = json.Unmarshal([]byte("true"), &res)
	if !errors.As(err, &jerr) || jerr.JSONKind != 't' {
		t.Errorf("Decimal.UnmarshalJSONFrom(true) = %v, want json.SemanticError", err)
	}

	err = json.Unmarshal([]byte("\"1.23\""), &res)
	if !errors.As(err, &jerr) || jerr.JSONKind != '"' {
		t.Errorf("Decimal.UnmarshalJSONFrom(\"\") = %v, want json.SemanticError", err)
	}

	err = json.Unmarshal([]byte("\"1.23\""), &res, json.StringifyNumbers(true))
	if !res.Equal(New(123, -2)) || err != nil {
		t.Errorf("Decimal.UnmarshalJSONFrom(\"1.23\", StringifyNumbers) = (%v, %v), want (1.23, <nil>)", res, err)
	}

	type S struct {
		D Decimal `json:",string"`
	}

	var s S
	err = json.Unmarshal([]byte(`{"D":"123.45"}`), &s)
	if !s.D.Equal(New(12345, -2)) || err != nil {
		t.Errorf("json.Unmarshal({\"D\":\"123.45\"}) = (%v, %v), want (123.45, <nil>)", s.D, err)
	}

	err = json.Unmarshal([]byte(`{"D":123.45}`), &s)
	if !s.D.Equal(New(12345, -2)) || err != nil {
		t.Errorf("json.Unmarshal({\"D\":123.45}) = (%v, %v), want (123.45, <nil>)", s.D, err)
	}
}

func FuzzDecimalUnmarshalJSONFrom(f *testing.F) {
	f.Add([]byte("123456.789e10"))

	f.Fuzz(func(t *testing.T, data []byte) {
		t.Parallel()

		var dec Decimal
		dec.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(data)))
	})
}
