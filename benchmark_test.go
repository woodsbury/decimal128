package decimal128

import (
	"testing"
)

func BenchmarkParse_short(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse("333333333.333333")
	}
}

func BenchmarkParse_long(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse("1234567890.0123456789")
	}
}

func BenchmarkParseDecimal_short(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseDecimal("333333333.333333")
	}
}

func BenchmarkParseDecimal_long(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseDecimal("1234567890.0123456789")
	}
}
