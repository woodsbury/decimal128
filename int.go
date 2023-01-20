package decimal128

import "math/bits"

var uint128PowersOf10 = [...]uint128{
	{0x0000_0000_0000_0001, 0x0000_0000_0000_0000},
	{0x0000_0000_0000_000a, 0x0000_0000_0000_0000},
	{0x0000_0000_0000_0064, 0x0000_0000_0000_0000},
	{0x0000_0000_0000_03e8, 0x0000_0000_0000_0000},
	{0x0000_0000_0000_2710, 0x0000_0000_0000_0000},
	{0x0000_0000_0001_86a0, 0x0000_0000_0000_0000},
	{0x0000_0000_000f_4240, 0x0000_0000_0000_0000},
	{0x0000_0000_0098_9680, 0x0000_0000_0000_0000},
	{0x0000_0000_05f5_e100, 0x0000_0000_0000_0000},
	{0x0000_0000_3b9a_ca00, 0x0000_0000_0000_0000},
	{0x0000_0002_540b_e400, 0x0000_0000_0000_0000},
	{0x0000_0017_4876_e800, 0x0000_0000_0000_0000},
	{0x0000_00e8_d4a5_1000, 0x0000_0000_0000_0000},
	{0x0000_0918_4e72_a000, 0x0000_0000_0000_0000},
	{0x0000_5af3_107a_4000, 0x0000_0000_0000_0000},
	{0x0003_8d7e_a4c6_8000, 0x0000_0000_0000_0000},
	{0x0023_86f2_6fc1_0000, 0x0000_0000_0000_0000},
	{0x0163_4578_5d8a_0000, 0x0000_0000_0000_0000},
	{0x0de0_b6b3_a764_0000, 0x0000_0000_0000_0000},
	{0x8ac7_2304_89e8_0000, 0x0000_0000_0000_0000},
	{0x6bc7_5e2d_6310_0000, 0x0000_0000_0000_0005},
	{0x35c9_adc5_dea0_0000, 0x0000_0000_0000_0036},
	{0x19e0_c9ba_b240_0000, 0x0000_0000_0000_021e},
	{0x02c7_e14a_f680_0000, 0x0000_0000_0000_152d},
	{0x1bce_cced_a100_0000, 0x0000_0000_0000_d3c2},
	{0x1614_0148_4a00_0000, 0x0000_0000_0008_4595},
	{0xdcc8_0cd2_e400_0000, 0x0000_0000_0052_b7d2},
	{0x9fd0_803c_e800_0000, 0x0000_0000_033b_2e3c},
	{0x3e25_0261_1000_0000, 0x0000_0000_204f_ce5e},
	{0x6d72_17ca_a000_0000, 0x0000_0001_431e_0fae},
	{0x4674_edea_4000_0000, 0x0000_000c_9f2c_9cd0},
	{0xc091_4b26_8000_0000, 0x0000_007e_37be_2022},
	{0x85ac_ef81_0000_0000, 0x0000_04ee_2d6d_415b},
	{0x38c1_5b0a_0000_0000, 0x0000_314d_c644_8d93},
	{0x378d_8e64_0000_0000, 0x0001_ed09_bead_87c0},
	{0x2b87_8fe8_0000_0000, 0x0013_4261_72c7_4d82},
	{0xb34b_9f10_0000_0000, 0x00c0_97ce_7bc9_0715},
	{0x00f4_36a0_0000_0000, 0x0785_ee10_d5da_46d9},
	{0x098a_2240_0000_0000, 0x4b3b_4ca8_5a86_c47a},
}

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

func (n uint128) add64(o uint64) uint128 {
	r0, carry := bits.Add64(n[0], o, 0)
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
		r = r.add64(1)
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

func (n uint128) log10() int {
	var l2 int
	if n[1] != 0 {
		l2 = bits.Len64(n[1]) + 63
	} else if n[0] != 0 {
		l2 = bits.Len64(n[0]) - 1
	} else {
		return 0
	}

	l10 := (l2 + 1) * 1233 >> 12
	if n.cmp(uint128PowersOf10[l10]) < 0 {
		l10--
	}

	return l10
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

func (n uint128) msd2() int {
	for n[1] >= 10 {
		n, _ = n.div1e19()
	}

	if n[1] != 0 {
		n, _ = n.div10()
	}

	n64 := n[0]

	for n64 >= 10000 {
		n64 /= 1000
	}

	for n64 >= 100 {
		n64 /= 10
	}

	return int(n64)
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

func (n uint128) mul1e38() uint256 {
	const o0 = 687399551400673280
	const o1 = 5421010862427522170

	u1, r0 := bits.Mul64(n[0], o0)
	v1, v0 := bits.Mul64(n[1], o0)
	w1, w0 := bits.Mul64(n[0], o1)
	x1, x0 := bits.Mul64(n[1], o1)

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

func (n uint128) or64(o uint64) uint128 {
	return uint128{n[0] | o, n[1]}
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

func (n uint128) sub64(o uint64) uint128 {
	r0, borrow := bits.Sub64(n[0], o, 0)
	r1 := n[1] - borrow

	return uint128{r0, r1}
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

func (n uint256) lsh(o uint) uint256 {
	var r0, r1, r2, r3 uint64
	if o > 192 {
		r3 = n[0] << (o - 192)
	} else if o > 128 {
		r2 = n[0] << (o - 128)
		r3 = n[1]<<(o-128) | n[0]>>(192-o)
	} else if o > 64 {
		r1 = n[0] << (o - 64)
		r2 = n[1]<<(o-64) | n[0]>>(128-o)
		r3 = n[2]<<(o-64) | n[1]>>(128-o)
	} else {
		r0 = n[0] << o
		r1 = n[1]<<o | n[0]>>(64-o)
		r2 = n[2]<<o | n[1]>>(64-o)
		r3 = n[3]<<o | n[2]>>(64-o)
	}

	return uint256{r0, r1, r2, r3}
}

func (n uint256) mul64(o uint64) uint256 {
	u1, r0 := bits.Mul64(n[0], o)
	v1, v0 := bits.Mul64(n[1], o)
	w1, w0 := bits.Mul64(n[2], o)
	x0 := n[3] * o

	r1, carry := bits.Add64(u1, v0, 0)
	r2, carry := bits.Add64(v1, w0, carry)
	r3, _ := bits.Add64(w1, x0, carry)

	return uint256{r0, r1, r2, r3}
}

func (n uint256) rsh(o uint) uint256 {
	var r0, r1, r2, r3 uint64
	if o > 192 {
		r0 = n[3] >> (o - 192)
	} else if o > 128 {
		r0 = n[2]>>(o-128) | n[3]<<(192-o)
		r1 = n[3] >> (o - 128)
	} else if o > 64 {
		r0 = n[1]>>(o-64) | n[2]<<(128-o)
		r1 = n[2]>>(o-64) | n[3]<<(128-o)
		r2 = n[3] >> (o - 64)
	} else {
		r0 = n[0]>>o | n[1]<<(64-o)
		r1 = n[1]>>o | n[2]<<(64-o)
		r2 = n[2]>>o | n[3]<<(64-o)
		r3 = n[3] >> o
	}

	return uint256{r0, r1, r2, r3}
}
