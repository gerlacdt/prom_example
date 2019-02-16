package http

import (
	"testing"
)

func TestRandn(t *testing.T) {
	random := NewRandom()
	range2 := Range{min: 100, max: 200}
	for i := 0; i < 1000; i++ {
		result := random.randn(range2)
		if result < range2.min || result > range2.max {
			t.Errorf("Randn retuned value is out of range, value: %d", result)
		}
	}
}
