package decimal128

var e = Decimal{0x4e90_6acc_b26a_bb56, 0x2ffe_8605_8a4b_f4de}

// E returns the mathematical constant e
func E() Decimal {
	return e
}

var pi = Decimal{0xbabe_5564_e6f3_9f8f, 0x2ffe_9ae4_7957_96a7}

// Pi returns the mathematical constant π
func Pi() Decimal {
	return pi
}
