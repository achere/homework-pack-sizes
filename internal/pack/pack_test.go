package pack_test

import (
	"fmt"
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
		t.Run(fmt.Sprintf("Order%d_Expected:%v", test.order, test.packs), func(t *testing.T) {
			t.Parallel()
			expected, _ := pack.CalculatePacks(sizes, test.order)
			assert.Equal(t, test.packs, expected)
		})
	}
}

func TestCalculatePacksEdgeCases(t *testing.T) {
	tests := []struct {
		order int
		sizes []int
		packs map[int]int
	}{
		{74, []int{25, 100}, map[int]int{25: 3}},
		{76, []int{25, 100}, map[int]int{100: 1}},
		{24, []int{25, 100}, map[int]int{25: 1}},
		{99, []int{100, 80, 20}, map[int]int{100: 1}},
		{14, []int{5, 12}, map[int]int{5: 3}},
		{59, []int{5, 12, 22}, map[int]int{22: 2, 5: 3}},
		{58, []int{5, 12, 22}, map[int]int{22: 1, 12: 3}},
		{34, []int{25, 15, 8}, map[int]int{15: 2, 8: 1}},
		{500000, []int{23, 31, 53}, map[int]int{23: 2, 31: 7, 53: 9429}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Order%d_Expected:%v", test.order, test.packs), func(t *testing.T) {
			t.Parallel()
			expected, _ := pack.CalculatePacks(test.sizes, test.order)
			assert.Equal(t, test.packs, expected)
		})
	}
}

func TestCalculatePacksBadInput(t *testing.T) {
	tests := []struct {
		name  string
		order int
		sizes []int
	}{
		{"Negative order", -2, []int{5, 10}},
		{"Negative size", -2, []int{5, 10}},
		{"Zero order", 0, []int{5, 10}},
		{"Zero size", 2, []int{0, 10}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := pack.CalculatePacks(test.sizes, test.order)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, pack.ErrInvalidArg)
			}
		})
	}
}
