package decimal128

import "math/bits"

type uint128 [2]uint64

func (n uint128) String() string {
	if n == (uint128{}) {
		return "0"
	}

	var buf [39]byte

	i := 39
	for n != (uint128{}) {
		var d uint64
		n, d = n.div10()
		i--
		buf[i] = '0' + byte(d)
	}

	return string(buf[i:])
}

func (n uint128) add(o uint128) uint192 {
	r0, carry := bits.Add64(n[0], o[0], 0)
	r1, r2 := bits.Add64(n[1], o[1], carry)

	return uint192{r0, r1, r2}
}

func (n uint128) add1() uint128 {
	r0, carry := bits.Add64(n[0], 1, 0)
	r1 := n[1] + carry

	return uint128{r0, r1}
}

func (n uint128) cmp(o uint128) int {
	if n[1] == o[1] {
		if n[0] == o[0] {
			return 0
		}

		if n[0] < o[0] {
			return -1
		}

		return 1
	}

	if n[1] < o[1] {
		return -1
	}

	return 1
}

func (n uint128) div(o uint128) (uint128, uint128) {
	if o[1] == 0 {
		var r0, r1, rem uint64
		if n[1] < o[0] {
			r0, rem = bits.Div64(n[1], n[0], o[0])
		} else {
			r1, rem = bits.Div64(0, n[1], o[0])
			r0, rem = bits.Div64(rem, n[0], o[0])
		}

		return uint128{r0, r1}, uint128{rem, 0}
	}

	i := uint(bits.LeadingZeros64(o[1]))
	u := o.lsh(i)
	v := n.rsh(1)
	r0, _ := bits.Div64(v[1], v[0], u[1])
	r0 >>= 63 - i
	if r0 != 0 {
		r0--
	}

	r := uint128{r0, 0}
	rem, _ := n.sub(o.mul64(r0))

	if rem.cmp(o) >= 0 {
		r = r.add1()
		rem, _ = rem.sub(o)
	}

	return r, rem
}

func (n uint128) div10() (uint128, uint64) {
	var r0, r1, rem uint64
	if n[1] < 10 {
		r0, rem = bits.Div64(n[1], n[0], 10)
	} else {
		r1, rem = bits.Div64(0, n[1], 10)
		r0, rem = bits.Div64(rem, n[0], 10)
	}

	return uint128{r0, r1}, rem
}

func (n uint128) div100() (uint128, uint64) {
	var r0, r1, rem uint64
	if n[1] < 10 {
		r0, rem = bits.Div64(n[1], n[0], 100)
	} else {
		r1, rem = bits.Div64(0, n[1], 100)
		r0, rem = bits.Div64(rem, n[0], 100)
	}

	return uint128{r0, r1}, rem
}

func (n uint128) div1000() (uint128, uint64) {
	var r0, r1, rem uint64
	if n[1] < 1000 {
		r0, rem = bits.Div64(n[1], n[0], 1000)
	} else {
		r1, rem = bits.Div64(0, n[1], 1000)
		r0, rem = bits.Div64(rem, n[0], 1000)
	}

	return uint128{r0, r1}, rem
}

func (n uint128) div1e8() (uint128, uint64) {
	var r0, r1, rem uint64
	if n[1] < 100_000_000 {
		r0, rem = bits.Div64(n[1], n[0], 100_000_000)
	} else {
		r1, rem = bits.Div64(0, n[1], 100_000_000)
		r0, rem = bits.Div64(rem, n[0], 100_000_000)
	}

	return uint128{r0, r1}, rem
}

func (n uint128) div1e19() (uint128, uint64) {
	var r0, r1, rem uint64
	if n[1] < 10_000_000_000_000_000_000 {
		r0, rem = bits.Div64(n[1], n[0], 10_000_000_000_000_000_000)
	} else {
		r1, rem = bits.Div64(0, n[1], 10_000_000_000_000_000_000)
		r0, rem = bits.Div64(rem, n[0], 10_000_000_000_000_000_000)
	}

	return uint128{r0, r1}, rem
}

func (n uint128) lsh(o uint) uint128 {
	var r0, r1 uint64
	if o > 64 {
		r1 = n[0] << (o - 64)
	} else {
		r0 = n[0] << o
		r1 = n[1]<<o | n[0]>>(64-o)
	}

	return uint128{r0, r1}
}

func (n uint128) mul(o uint128) uint256 {
	u1, r0 := bits.Mul64(n[0], o[0])
	v1, v0 := bits.Mul64(n[1], o[0])
	w1, w0 := bits.Mul64(n[0], o[1])
	x1, x0 := bits.Mul64(n[1], o[1])

	r1, carry := bits.Add64(u1, v0, 0)
	r2, y0 := bits.Add64(v1, w1, carry)
	r1, carry = bits.Add64(r1, w0, 0)
	r2, y1 := bits.Add64(r2, x0, carry)
	r3, _ := bits.Add64(x1, y0, y1)

	return uint256{r0, r1, r2, r3}
}

func (n uint128) mul64(o uint64) uint128 {
	r1, r0 := bits.Mul64(n[0], o)
	r1 += n[1] * o

	return uint128{r0, r1}
}

func (n uint128) not() uint128 {
	return uint128{^n[0], ^n[1]}
}

func (n uint128) rsh(o uint) uint128 {
	var r0, r1 uint64
	if o > 64 {
		r0 = n[1] >> (o - 64)
	} else {
		r0 = n[0]>>o | n[1]<<(64-o)
		r1 = n[1] >> o
	}

	return uint128{r0, r1}
}

func (n uint128) sub(o uint128) (uint128, uint) {
	r0, borrow := bits.Sub64(n[0], o[0], 0)
	r1, borrow := bits.Sub64(n[1], o[1], borrow)

	return uint128{r0, r1}, uint(borrow)
}

func (n uint128) twos() uint128 {
	r0, carry := bits.Add64(^n[0], 1, 0)
	r1 := ^n[1] + carry

	return uint128{r0, r1}
}

type uint192 [3]uint64

func (n uint192) String() string {
	if n == (uint192{}) {
		return "0"
	}

	var buf [58]byte

	i := 58
	for n != (uint192{}) {
		var d uint64
		n, d = n.div10()
		i--
		buf[i] = '0' + byte(d)
	}

	return string(buf[i:])
}

func (n uint192) div10() (uint192, uint64) {
	r2, rem := bits.Div64(0, n[2], 10)
	r1, rem := bits.Div64(rem, n[1], 10)
	r0, rem := bits.Div64(rem, n[0], 10)

	return uint192{r0, r1, r2}, rem
}

func (n uint192) div10000() (uint192, uint64) {
	r2, rem := bits.Div64(0, n[2], 10000)
	r1, rem := bits.Div64(rem, n[1], 10000)
	r0, rem := bits.Div64(rem, n[0], 10000)

	return uint192{r0, r1, r2}, rem
}

type uint256 [4]uint64

func (n uint256) String() string {
	if n == (uint256{}) {
		return "0"
	}

	var buf [78]byte

	i := 78
	for n != (uint256{}) {
		var d uint64
		n, d = n.div10()
		i--
		buf[i] = '0' + byte(d)
	}

	return string(buf[i:])
}

func (n uint256) div10() (uint256, uint64) {
	r3, rem := bits.Div64(0, n[3], 10)
	r2, rem := bits.Div64(rem, n[2], 10)
	r1, rem := bits.Div64(rem, n[1], 10)
	r0, rem := bits.Div64(rem, n[0], 10)

	return uint256{r0, r1, r2, r3}, rem
}

func (n uint256) div1e19() (uint256, uint64) {
	r3, rem := bits.Div64(0, n[3], 10_000_000_000_000_000_000)
	r2, rem := bits.Div64(rem, n[2], 10_000_000_000_000_000_000)
	r1, rem := bits.Div64(rem, n[1], 10_000_000_000_000_000_000)
	r0, rem := bits.Div64(rem, n[0], 10_000_000_000_000_000_000)

	return uint256{r0, r1, r2, r3}, rem
}
