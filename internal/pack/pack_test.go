package pack_test

import (
	"testing"

	"github.com/achere/homework-pack-sizes/internal/pack"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePacks(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		order int
		packs map[int]int
	}{
		{1, map[int]int{250: 1}},
		{250, map[int]int{250: 1}},
		{251, map[int]int{500: 1}},
		{501, map[int]int{500: 1, 250: 1}},
		{12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
	}

	for _, test := range tests {
		assert.Equal(t, test.packs, pack.CalculatePacks(sizes, test.order))
	}
}
