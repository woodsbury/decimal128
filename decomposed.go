package decimal128

var (
	ln10 = decomposed128{
		sig: uint128{0x09bb_c25b_3ca8_1898, 0xad3a_2d01_4ad4_7d7a},
		exp: -38,
	}

	invLn10 = decomposed128{
		sig: uint128{0x2817_8c5e_dd06_4c45, 0x20ac_351d_00b1_29f6},
		exp: -38,
	}

	invLn2 = decomposed128{
		sig: uint128{0x43d4_c3f7_1489_9de8, 0x3425_8773_b151_f6b7},
		exp: -38,
	}

	ln = [...]uint128{
		{0xc57f_c7fc_0b2b_0026, 0x072b_9b77_8c2c_1083}, // ln(1.1)
		{0x3be7_3b96_9654_c697, 0x0db7_62ad_5584_8050}, // ln(1.2)
		{0xa858_b1b9_15df_0b68, 0x13bc_f3b1_dd60_6152}, // ln(1.3)
		{0xa349_cafd_04cd_7733, 0x1950_3835_0f76_37e0}, // ln(1.4)
		{0x7a24_b20c_9560_0579, 0x1e80_f953_8c63_19a4}, // ln(1.5)
		{0x0597_4d81_157e_5f06, 0x235b_f0cd_7a73_5d63}, // ln(1.6)
		{0xa4a9_d782_e965_ec5e, 0x27eb_8743_f332_c960}, // ln(1.7)
		{0xb60b_eda3_2bb4_cc10, 0x2c38_5c00_e1e7_99f4}, // ln(1.8)
		{0x7675_df87_db0c_3b63, 0x3049_a7fc_43d2_47f3}, // ln(1.9)
		{0x43d4_c3f7_1489_9de8, 0x3425_8773_b151_f6b7}, // ln(2.0)
		{0x1d6e_7d09_9a2d_7cac, 0x37d1_3188_9bd9_5185}, // ln(2.1)
		{0x0954_8bf3_1fb4_9e0d, 0x3b51_22eb_3d7e_073b}, // ln(2.2)
		{0x9528_5cb8_ffb0_2a2c, 0x3ea9_3f07_7f26_6fb9}, // ln(2.3)
		{0x7fbb_ff8d_aade_647f, 0x41dc_ea21_06d6_7707}, // ln(2.4)
		{0x8212_3a6d_1394_dcc9, 0x44ef_1e19_e830_900b}, // ln(2.5)
		{0xec2d_75b0_2a68_a94f, 0x47e2_7b25_8eb2_5809}, // ln(2.6)
		{0x3030_9faf_c114_d18a, 0x4ab9_5554_6e4a_b399}, // ln(2.7)
		{0xe71e_8ef4_1957_151a, 0x4d75_bfa8_c0c8_2e97}, // ln(2.8)
		{0xcbef_c0e1_c298_304d, 0x5019_9539_4076_3df0}, // ln(2.9)
		{0xbdf9_7603_a9e9_a361, 0x52a6_80c7_3db5_105b}, // ln(3.0)
		{0x581c_cc5f_f1f3_f4bc, 0x551e_0316_24cb_bf58}, // ln(3.1)
		{0x496c_1178_2a07_fced, 0x5781_7841_2bc5_541a}, // ln(3.2)
		{0x8379_3dff_b514_a386, 0x59d2_1c3e_c9e1_20df}, // ln(3.3)
		{0xe87e_9b79_fdef_8a46, 0x5c11_0eb7_a484_c017}, // ln(3.4)
		{0x255c_056a_1862_53fc, 0x5e3f_564e_f7a6_c7ec}, // ln(3.5)
		{0xf9e0_b19a_403e_69f8, 0x605d_e374_9339_90ab}, // ln(3.6)
		{0x4994_cb3a_11f0_4fcd, 0x626d_92d3_e742_799a}, // ln(3.7)
		{0xba4a_a37e_ef95_d94b, 0x646f_2f6f_f524_3eaa}, // ln(3.8)
		{0x6652_27bc_bfc8_aec9, 0x6663_7479_1b15_71ae}, // ln(3.9)
		{0x87a9_87ee_2913_3bcf, 0x684b_0ee7_62a3_ed6e}, // ln(4.0)
		{0x7656_c6c8_5770_229d, 0x6a26_9ee2_2347_bcd5}, // ln(4.1)
		{0x6143_4100_aeb7_1a93, 0x6bf6_b8fc_4d2b_483c}, // ln(4.2)
		{0x7d30_8c5e_97fe_01ab, 0x6dbb_e74b_7b2b_765a}, // ln(4.3)
		{0x4d29_4fea_343e_3bf5, 0x6f76_aa5e_eecf_fdf2}, // ln(4.4)
		{0x381e_2810_3f49_a8d9, 0x7127_7a1a_ca18_2a00}, // ln(4.5)
		{0xd8fd_20b0_1439_c813, 0x72ce_c67b_3078_6670}, // ln(4.6)
		{0x6bec_4ab3_73e9_4271, 0x746c_f842_6b3d_2da5}, // ln(4.7)
		{0xc390_c384_bf68_0266, 0x7602_7194_b828_6dbe}, // ln(4.8)
		{0xc8a5_d067_1d2f_cb2f, 0x778f_8e84_071c_ffcc}, // ln(4.9)
		{0xc5e6_fe64_281e_7ab0, 0x7914_a58d_9982_86c2}, // ln(5.0)
		{0x62a3_4d86_934f_8fbf, 0x7a92_080b_30e7_d9bc}, // ln(5.1)
		{0x3002_39a7_3ef2_4737, 0x7c08_0299_4004_4ec1}, // ln(5.2)
		{0x91d8_b56a_1351_c2f3, 0x7d76_dd73_5fb9_29e6}, // ln(5.3)
		{0x7405_63a6_d59e_6f71, 0x7ede_dcc8_1f9c_aa50}, // ln(5.4)
		{0x8b66_c660_3349_7ad6, 0x8040_4105_25ae_9746}, // ln(5.5)
		{0x2af3_52eb_2de0_b302, 0x819b_471c_721a_254f}, // ln(5.6)
		{0x346f_558b_84f5_dec4, 0x82f0_28c3_8187_584f}, // ln(5.7)
		{0x0fc4_84d8_d721_ce34, 0x843f_1cac_f1c8_34a8}, // ln(5.8)
		{0x1a60_529a_bde3_dbce, 0x8588_56bd_3913_7bcc}, // ln(5.9)
		{0x86cc_083a_ef07_0713, 0x01ce_39fa_be73_4148}, // ln(6.0)
		{0x1c63_4881_c6a5_0faf, 0x880a_5ffb_17f2_a89b}, // ln(6.1)
		{0x9bf1_9057_067d_92a4, 0x8943_8a89_d61d_b60f}, // ln(6.2)
		{0xdb67_f30d_4417_200d, 0x8a77_b24f_d98e_61e0}, // ln(6.3)
		{0x8d40_d56f_3e91_9ad5, 0x8ba6_ffb4_dd17_4ad1}, // ln(6.4)
		{0x6e3f_b01d_3dfd_8618, 0x8cd1_993f_76e2_e815}, // ln(6.5)
		{0xc74e_01f6_c99e_416e, 0x8df7_a3b2_7b33_1796}, // ln(6.6)
		{0x16a6_26d1_beaf_26d8, 0x8f19_4228_2970_e2a9}, // ln(6.7)
		{0x2c53_5f71_1279_282d, 0x9036_962b_55d6_b6cf}, // ln(6.8)
		{0x5321_d2bc_a999_cd8c, 0x914f_bfce_bcdb_8015}, // ln(6.9)
		{0x6930_c961_2ceb_f1e4, 0x9264_ddc2_a8f8_bea3}, // ln(7.0)
		{0xf6c4_b737_9855_6f73, 0x9376_0d69_0f5d_8206}, // ln(7.1)
		{0x3db5_7591_54c8_07e0, 0x9483_6ae8_448b_8763}, // ln(7.2)
		{0x54d4_cfde_f90c_6b2a, 0x958d_113c_66ac_34d9}, // ln(7.3)
		{0x8d69_8f31_2679_edb5, 0x9693_1a47_9894_7051}, // ln(7.4)
		{0x400b_b070_bd7e_8029, 0x9795_9ee1_25e5_a067}, // ln(7.5)
		{0xfe1f_6776_041f_7732, 0x9894_b6e3_a676_3561}, // ln(7.6)
		{0x2eb0_915d_3816_f209, 0x9990_793a_3524_cf27}, // ln(7.7)
		{0xaa26_ebb3_d452_4cb0, 0x9a88_fbec_cc67_6865}, // ln(7.8)
		{0x698f_c977_956d_34e1, 0x9b7e_542b_d945_b932}, // ln(7.9)
		{0xcb7e_4be5_3d9c_d9b6, 0x9c70_965b_13f5_e425}, // ln(8.0)
		{0xee2a_15b3_6afe_74ea, 0x9d5f_d61b_abff_c3f4}, // ln(8.1)
		{0xba2b_8abf_6bf9_c085, 0x9e4c_2655_d499_b38c}, // ln(8.2)
		{0x7690_842a_ef44_069e, 0x9f35_9941_bcdd_cf09}, // ln(8.3)
		{0xa518_04f7_c340_b87b, 0xa01c_406f_fe7d_3ef3}, // ln(8.4)
		{0x6a90_d5e7_1184_670f, 0xa100_2cd1_8cb5_5023}, // ln(8.5)
		{0xc105_5055_ac87_9f93, 0xa1e1_6ebf_2c7d_6d11}, // ln(8.6)
		{0x89e9_36e5_6c81_d3ae, 0xa2c0_1600_7e2b_4e4c}, // ln(8.7)
		{0x90fe_13e1_48c7_d9dc, 0xa39c_31d2_a021_f4a9}, // ln(8.8)
		{0x4827_72c5_2db9_d2bc, 0xa475_d0ee_7186_7e7d}, // ln(8.9)
		{0x7bf2_ec07_53d3_46c1, 0xa54d_018e_7b6a_20b7}, // ln(9.0)
		{0x1189_7b1a_42ca_fd4c, 0xa621_d174_8659_1ff6}, // ln(9.1)
		{0x1cd1_e4a7_28c3_65fb, 0xa6f4_4dee_e1ca_5d28}, // ln(9.2)
		{0x1616_4263_9bdd_981d, 0xa7c4_83dd_6280_cfb4}, // ln(9.3)
		{0xafc1_0eaa_8872_e059, 0xa892_7fb6_1c8f_245c}, // ln(9.4)
		{0x3c5c_ddec_032a_b614, 0xa95e_4d89_dd54_ceb6}, // ln(9.5)
		{0x0765_877b_d3f1_a04e, 0xaa27_f908_697a_6476}, // ln(9.6)
		{0x3958_344e_deef_43dc, 0xaaef_8d84_82a8_6461}, // ln(9.7)
		{0x0c7a_945e_31b9_6916, 0xabb5_15f7_b86e_f684}, // ln(9.8)
		{0x4172_b403_5efe_46e7, 0xac78_9d06_0796_313b}, // ln(9.9)
	}
)

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

func (d decomposed128) log() (bool, decomposed128, int8) {
	l10 := int16(d.sig.log10())
	d.exp = (d.exp - exponentBias) + l10

	msd := d.sig.msd2()
	sig := d.sig
	exp := -l10
	oneSig := uint128{1, 0}
	oneExp := int16(0)

	for sig[1] <= 0x0002_7fff_ffff_ffff {
		sig = sig.mul64(10_000)
		exp -= 4

		oneSig = oneSig.mul64(10_000)
		oneExp -= 4
	}

	for sig[1] <= 0x18ff_ffff_ffff_ffff {
		sig = sig.mul64(10)
		exp--

		oneSig = oneSig.mul64(10)
		oneExp--
	}

	nrm := decomposed128{
		sig: sig,
		exp: exp,
	}

	if msd < 10 {
		msd *= 10
	}

	var trunc int8
	if msd > 10 {
		nrm, trunc = nrm.quo(decomposed128{
			sig: uint128{uint64(msd), 0},
			exp: -1,
		}, 0)
	}

	one := decomposed128{
		sig: oneSig,
		exp: oneExp,
	}

	_, num, _ := nrm.sub(one, int8(0))
	den, _ := nrm.add(one, int8(0))
	frc, trunc := num.quo(den, trunc)
	sqr, _ := frc.mul(frc, int8(0))

	res := frc

	for i := uint64(3); i <= 25; i += 2 {
		// res += frc^i / i
		frc, _ = frc.mul(sqr, int8(0))
		tmp, _ := frc.quo(decomposed128{
			sig: uint128{i, 0},
			exp: 0,
		}, int8(0))

		res, trunc = res.add(tmp, trunc)
	}

	expNeg := false
	if d.exp < 0 {
		d.exp *= -1
		expNeg = true
	}

	lnExp, _ := ln10.mul(decomposed128{
		sig: uint128{uint64(d.exp), 0},
		exp: 0,
	}, int8(0))

	res, trunc = res.mul(decomposed128{
		sig: uint128{2, 0},
		exp: 0,
	}, trunc)

	neg := false
	if expNeg {
		neg, res, trunc = res.sub(lnExp, trunc)
	} else {
		res, trunc = res.add(lnExp, trunc)
	}

	if msd > 10 {
		lnMSD := decomposed128{
			sig: ln[msd-11],
			exp: -38,
		}

		if expNeg {
			_, res, trunc = res.sub(lnMSD, trunc)
		} else {
			res, trunc = res.add(lnMSD, trunc)
		}
	}

	return neg, res, trunc
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
