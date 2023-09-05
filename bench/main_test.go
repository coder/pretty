package bench

import (
	"testing"
)

func BenchmarkPretty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pretty()
	}
}
