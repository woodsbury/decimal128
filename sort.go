package decimal128

import "sort"

// Sort sorts a slice of Decimals in increasing order. NaN values order before
// other values.
func Sort(x []Decimal) {
	sort.Sort(decimalSlice(x))
}

type decimalSlice []Decimal

func (x decimalSlice) Len() int {
	return len(x)
}

func (x decimalSlice) Less(i, j int) bool {
	if x[i].IsNaN() {
		return !x[j].IsNaN()
	}

	return x[i].Cmp(x[j]).Less()
}

func (x decimalSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
