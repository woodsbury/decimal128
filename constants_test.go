package decimal128

import "testing"

func TestE(t *testing.T) {
	t.Parallel()

	res := E()
	num := MustParse("2.71828182845904523536028747135266249775724709369995957496696763")

	if !res.Equal(num) {
		t.Errorf("E() = %v, want %v", res, num)
	}
}

func TestPhi(t *testing.T) {
	t.Parallel()

	res := Phi()
	num := MustParse("1.61803398874989484820458683436563811772030917980576286213544862")

	if !res.Equal(num) {
		t.Errorf("Phi() = %v, want %v", res, num)
	}
}

func TestPi(t *testing.T) {
	t.Parallel()

	res := Pi()
	num := MustParse("3.14159265358979323846264338327950288419716939937510582097494459")

	if !res.Equal(num) {
		t.Errorf("Pi() = %v, want %v", res, num)
	}
}
