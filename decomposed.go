package decimal128

type decomposed128 struct {
	sig uint128
	exp int16
}

func (d decomposed128) add(o decomposed128, trunc int8) (decomposed128, int8) {
	exp := d.exp - o.exp

	if exp < 0 {
		for exp <= -19 && o.sig[1] == 0 {
			o.sig = o.sig.mul64(10_000_000_000_000_000_000)
			o.exp -= 19
			exp += 19
		}

		for exp <= -4 && o.sig[1] <= 0x0002_7fff_ffff_ffff {
			o.sig = o.sig.mul64(10_000)
			o.exp -= 4
			exp += 4
		}

		for exp < 0 && o.sig[1] <= 0x18ff_ffff_ffff_ffff {
			o.sig = o.sig.mul64(10)
			o.exp--
			exp++
		}

		if exp <= -3 {
			var rem uint64
			d.sig, rem = d.sig.div1000()
			if rem != 0 {
				trunc = 1
			}

			if d.sig == (uint128{}) {
				d.exp = o.exp
				exp = 0
			} else {
				d.exp += 3
				exp += 3
			}
		}

		for exp < 0 {
			var rem uint64
			d.sig, rem = d.sig.div10()
			if rem != 0 {
				trunc = 1
			}

			if d.sig == (uint128{}) {
				d.exp = o.exp
				exp = 0
				break
			}

			d.exp++
			exp++
		}
	} else if exp > 0 {
		if exp >= 19 && d.sig[1] == 0 {
			d.sig = d.sig.mul64(10_000_000_000_000_000_000)
			d.exp -= 19
			exp -= 19
		}

		for exp >= 4 && d.sig[1] <= 0x0002_7fff_ffff_ffff {
			d.sig = d.sig.mul64(10_000)
			d.exp -= 4
			exp -= 4
		}

		for exp > 0 && d.sig[1] <= 0x18ff_ffff_ffff_ffff {
			d.sig = d.sig.mul64(10)
			d.exp--
			exp--
		}

		if exp >= 3 {
			var rem uint64
			o.sig, rem = o.sig.div1000()
			if rem != 0 {
				trunc = -1
			}

			if o.sig == (uint128{}) {
				exp = 0
			} else {
				exp -= 3
			}
		}

		for exp > 0 {
			var rem uint64
			o.sig, rem = o.sig.div10()
			if rem != 0 {
				trunc = -1
			}

			if o.sig == (uint128{}) {
				exp = 0
				break
			}

			exp--
		}
	}

	sig192 := d.sig.add(o.sig)
	exp = d.exp

	for sig192[2] >= 0x0000_0000_0000_ffff {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = 1
		}
	}

	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10()
		exp++

		if rem != 0 {
			trunc = 1
		}
	}

	return decomposed128{
		sig: uint128{sig192[0], sig192[1]},
		exp: exp,
	}, trunc
}

func (d decomposed128) mul(o decomposed128, trunc int8) (decomposed128, int8) {
	sig256 := d.sig.mul(o.sig)
	exp := d.exp + o.exp

	for sig256[3] > 0 {
		var rem uint64
		sig256, rem = sig256.div1e19()
		exp += 19

		if rem != 0 {
			trunc = 1
		}
	}

	sig192 := uint192{sig256[0], sig256[1], sig256[2]}

	for sig192[2] >= 0x0000_0000_0000_ffff {
		var rem uint64
		sig192, rem = sig192.div10000()
		exp += 4

		if rem != 0 {
			trunc = 1
		}
	}

	for sig192[2] > 0 {
		var rem uint64
		sig192, rem = sig192.div10()
		exp++

		if rem != 0 {
			trunc = 1
		}
	}

	return decomposed128{
		sig: uint128{sig192[0], sig192[1]},
		exp: exp,
	}, trunc
}

func (d decomposed128) quo(o decomposed128, trunc int8) (decomposed128, int8) {
	if d.sig == (uint128{}) {
		return decomposed128{
			sig: uint128{},
			exp: 0,
		}, trunc
	}

	if d.sig[1] == 0 {
		d.sig = d.sig.mul64(10_000_000_000_000_000_000)
		d.exp -= 19
	}

	for d.sig[1] <= 0x0002_7fff_ffff_ffff {
		d.sig = d.sig.mul64(10_000)
		d.exp -= 4
	}

	for d.sig[1] <= 0x18ff_ffff_ffff_ffff {
		d.sig = d.sig.mul64(10)
		d.exp--
	}

	for o.sig[1] >= 0x18ff_ffff_ffff_ffff {
		var rem uint64
		o.sig, rem = o.sig.div10()
		o.exp++

		if rem != 0 {
			trunc = 1
		}
	}

	sig, rem := d.sig.div(o.sig)
	exp := d.exp - o.exp

	for rem != (uint128{}) && sig[1] <= 0x18ff_ffff_ffff_ffff {
		for rem[1] <= 0x0002_7fff_ffff_ffff && sig[1] <= 0x0002_7fff_ffff_ffff {
			rem = rem.mul64(10_000)
			sig = sig.mul64(10_000)
			exp -= 4
		}

		for rem[1] <= 0x18ff_ffff_ffff_ffff && sig[1] <= 0x18ff_ffff_ffff_ffff {
			rem = rem.mul64(10)
			sig = sig.mul64(10)
			exp--
		}

		var tmp uint128
		tmp, rem = rem.div(o.sig)
		sig192 := sig.add(tmp)

		for sig192[2] != 0 {
			var rem192 uint64
			sig192, rem192 = sig192.div10()
			exp++

			if rem192 != 0 {
				trunc = 1
			}
		}

		sig = uint128{sig192[0], sig192[1]}
	}

	if rem != (uint128{}) {
		trunc = 1
	}

	return decomposed128{
		sig: sig,
		exp: exp,
	}, trunc
}

func (d decomposed128) sub(o decomposed128, trunc int8) (bool, decomposed128, int8) {
	exp := d.exp - o.exp

	if exp < 0 {
		for exp <= -19 && o.sig[1] == 0 {
			o.sig = o.sig.mul64(10_000_000_000_000_000_000)
			o.exp -= 19
			exp += 19
		}

		for exp <= -4 && o.sig[1] <= 0x0002_7fff_ffff_ffff {
			o.sig = o.sig.mul64(10_000)
			o.exp -= 4
			exp += 4
		}

		for exp < 0 && o.sig[1] <= 0x18ff_ffff_ffff_ffff {
			o.sig = o.sig.mul64(10)
			o.exp--
			exp++
		}

		if exp <= -3 {
			var rem uint64
			d.sig, rem = d.sig.div1000()
			if rem != 0 {
				trunc = 1
			}

			if d.sig == (uint128{}) {
				d.exp = o.exp
				exp = 0
			} else {
				d.exp += 3
				exp += 3
			}
		}

		for exp < 0 {
			var rem uint64
			d.sig, rem = d.sig.div10()
			if rem != 0 {
				trunc = 1
			}

			if d.sig == (uint128{}) {
				d.exp = o.exp
				exp = 0
				break
			}

			d.exp++
			exp++
		}
	} else if exp > 0 {
		if exp >= 19 && d.sig[1] == 0 {
			d.sig = d.sig.mul64(10_000_000_000_000_000_000)
			d.exp -= 19
			exp -= 19
		}

		for exp >= 4 && d.sig[1] <= 0x0002_7fff_ffff_ffff {
			d.sig = d.sig.mul64(10_000)
			d.exp -= 4
			exp -= 4
		}

		for exp > 0 && d.sig[1] <= 0x18ff_ffff_ffff_ffff {
			d.sig = d.sig.mul64(10)
			d.exp--
			exp--
		}

		if exp >= 3 {
			var rem uint64
			o.sig, rem = o.sig.div1000()
			if rem != 0 {
				trunc = -1
			}

			if o.sig == (uint128{}) {
				exp = 0
			} else {
				exp -= 3
			}
		}

		for exp > 0 {
			var rem uint64
			o.sig, rem = o.sig.div10()
			if rem != 0 {
				trunc = -1
			}

			if o.sig == (uint128{}) {
				exp = 0
				break
			}

			exp--
		}
	}

	neg := false
	sig, brw := d.sig.sub(o.sig)
	exp = d.exp

	if brw != 0 {
		sig = sig.twos()
		neg = true
		trunc *= -1
	}

	return neg, decomposed128{
		sig: sig,
		exp: exp,
	}, trunc
}
