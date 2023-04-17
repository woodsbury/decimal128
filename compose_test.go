package decimal128

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"math/big"
	"testing"
)

type sqlConn struct{}

func (c sqlConn) Begin() (driver.Tx, error) {
	return nil, errors.New("unsupported")
}

func (c sqlConn) Close() error {
	return nil
}

func (c sqlConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unsupported")
}

func (c sqlConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var num sqlNumeric
	if err := num.Compose(args[0].Value.(Decimal).Decompose(nil)); err != nil {
		return nil, err
	}

	return &sqlRows{num: num}, nil
}

type sqlConnector struct{}

func (c sqlConnector) Connect(context.Context) (driver.Conn, error) {
	return sqlConn{}, nil
}

func (c sqlConnector) Driver() driver.Driver {
	return nil
}

type sqlNumeric struct {
	form byte
	neg  bool
	sig  []byte
	exp  int32
}

func (n *sqlNumeric) Compose(form byte, neg bool, sig []byte, exp int32) error {
	*n = sqlNumeric{
		form: form,
		neg:  neg,
		sig:  sig,
		exp:  exp,
	}

	return nil
}

func (n sqlNumeric) Decompose(buf []byte) (byte, bool, []byte, int32) {
	return n.form, n.neg, n.sig, n.exp
}

type sqlRows struct {
	eof bool
	num sqlNumeric
}

func (r *sqlRows) Close() error {
	return nil
}

func (r *sqlRows) Columns() []string {
	return []string{"num"}
}

func (r *sqlRows) Next(dest []driver.Value) error {
	if r.eof {
		return io.EOF
	}

	dest[0] = r.num
	r.eof = true
	return nil
}

func TestDecimalCompose(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	bigsig := new(big.Int)

	for _, val := range decimalValues {
		var form byte
		var neg bool
		var sig []byte
		var exp int32

		switch val.form {
		case regularForm:
			form = 0
			neg = val.neg

			if val.sig != (uint128{}) {
				sig = uint128ToBig(val.sig, bigsig).Bytes()
				exp = int32(val.exp) - exponentBias
			}
		case infForm:
			form = 1
			neg = val.neg
		case nanForm:
			form = 2
		}

		var dec Decimal
		err := dec.Compose(form, neg, sig, exp)

		if !resultEqual(dec, val.Decimal()) || err != nil {
			t.Errorf("Decimal.Compose(%d, %t, %x, %d) = (%v, %v), want (%v, <nil>)", form, neg, sig, exp, dec, err, val)
		}

		decform, decneg, decsig, decexp := dec.Decompose(nil)

		if decform != form || decneg != neg || !bytes.Equal(decsig, sig) || decexp != exp {
			t.Errorf("%v.Decompose() = (%d, %t, %x, %d), want (%d, %t, %x, %d)", dec, decform, decneg, decsig, decexp, form, neg, sig, exp)
		}
	}
}

func TestDecimalComposeSQL(t *testing.T) {
	t.Parallel()

	initDecimalValues()

	for _, val := range decimalValues {
		decval := val.Decimal()

		db := sql.OpenDB(sqlConnector{})

		rows, err := db.Query("select ?", decval)
		if err != nil {
			t.Fatalf("sql.Query() = %v, want <nil>", err)
		}

		if !rows.Next() {
			t.Fatalf("sql.Rows.Next() = false, want true")
		}

		var resval Decimal
		err = rows.Scan(&resval)

		if !resultEqual(resval, decval) || err != nil {
			t.Errorf("sql.Rows.Scan() = (%v, %v), want (%v, <nil>)", resval, err, resval)
		}
	}
}

func FuzzDecimalCompose(f *testing.F) {
	f.Add(byte(0), false, []byte{0, 1}, int32(0))
	f.Add(byte(1), false, []byte{0, 1}, int32(0))
	f.Add(byte(2), false, []byte{0, 1}, int32(0))

	f.Fuzz(func(t *testing.T, form byte, neg bool, sig []byte, exp int32) {
		t.Parallel()

		form %= 3

		var dec Decimal
		if err := dec.Compose(form, neg, sig, exp); err != nil {
			return
		}

		decform, decneg, _, _ := dec.Decompose(nil)

		if decform != form || (decform != 2 && decneg != neg) {
			t.Fail()
		}
	})
}
