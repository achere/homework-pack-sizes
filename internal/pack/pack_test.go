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
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Order%d_Expected:%v", test.order, test.packs), func(t *testing.T) {
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
			_, err := pack.CalculatePacks(test.sizes, test.order)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, pack.ErrInvalidArg)
			}
		})
	}
}
