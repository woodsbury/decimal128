package decimal128

import "testing"

var roundingModes = []RoundingMode{
	ToNearestEven,
	ToNearestAway,
	ToZero,
	AwayFromZero,
	ToNegativeInf,
	ToPositiveInf,
}

func TestDecimalCeil(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var dp int
	var res Decimal

	for r.scan("ceil(%v, %v) = %v\n", &val, &dp, &res) {
		rnd := val.Ceil(dp)

		if !resultEqual(rnd, res) {
			t.Errorf("%v.Ceil(%d) = %v, want %v", val, dp, rnd, res)
		}
	}
}

func TestDecimalFloor(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var dp int
	var res Decimal

	for r.scan("floor(%v, %v) = %v\n", &val, &dp, &res) {
		rnd := val.Floor(dp)

		if !resultEqual(rnd, res) {
			t.Errorf("%v.Floor(%d) = %v, want %v", val, dp, rnd, res)
		}
	}
}

func TestDecimalRound(t *testing.T) {
	t.Parallel()

	r := openTestData(t)
	defer r.close()

	var val Decimal
	var dp int
	var res testDataResult

	for r.scan("round(%v, %v) = %v\n", &val, &dp, &res) {
		for _, mode := range roundingModes {
			rnd := val.Round(dp, mode)

			if !res.equal(rnd, mode) {
				t.Errorf("%v.Round(%d, %v) = %v, want %v", val, dp, mode, rnd, res.result(mode))
			}
		}
	}
}

func BenchmarkReduce128(b *testing.B) {
	initUintValues()

	exponents := []int16{exponentBias / 2, exponentBias, exponentBias * 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, mode := range roundingModes {
			for _, val := range uint128Values {
				for _, exp := range exponents {
					mode.reduce128(false, val, exp, 0)
				}
			}
		}
	}
}

func BenchmarkReduce192(b *testing.B) {
	initUintValues()

	exponents := []int16{exponentBias / 2, exponentBias, exponentBias * 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, mode := range roundingModes {
			for _, val := range uint192Values {
				for _, exp := range exponents {
					mode.reduce192(false, val, exp, 0)
				}
			}
		}
	}
}

func BenchmarkReduce256(b *testing.B) {
	initUintValues()

	exponents := []int16{exponentBias / 2, exponentBias, exponentBias * 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, mode := range roundingModes {
			for _, val := range uint256Values {
				for _, exp := range exponents {
					mode.reduce256(false, val, exp, 0)
				}
			}
		}
	}
}
