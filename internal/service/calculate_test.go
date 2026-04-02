package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculatePacks(t *testing.T) {
	defaultPacks := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name      string
		order     int
		packSizes []int
		wantPacks map[int]int
		wantTotal int
	}{
		{
			name:      "order 1 → 1x250",
			order:     1,
			packSizes: defaultPacks,
			wantPacks: map[int]int{250: 1},
			wantTotal: 250,
		},
		{
			name:      "order 250 → 1x250",
			order:     250,
			packSizes: defaultPacks,
			wantPacks: map[int]int{250: 1},
			wantTotal: 250,
		},
		{
			name:      "order 251 → 1x500",
			order:     251,
			packSizes: defaultPacks,
			wantPacks: map[int]int{500: 1},
			wantTotal: 500,
		},
		{
			name:      "order 501 → 1x500 + 1x250",
			order:     501,
			packSizes: defaultPacks,
			wantPacks: map[int]int{500: 1, 250: 1},
			wantTotal: 750,
		},
		{
			name:      "order 12001 → 2x5000 + 1x2000 + 1x250",
			order:     12001,
			packSizes: defaultPacks,
			wantPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
			wantTotal: 12250,
		},
		{
			name:      "order 0 → empty",
			order:     0,
			packSizes: defaultPacks,
			wantPacks: map[int]int{},
			wantTotal: 0,
		},
		{
			name:      "primes: order 500000 → exact fit",
			order:     500000,
			packSizes: []int{23, 31, 53},
			wantPacks: map[int]int{23: 2, 31: 7, 53: 9429},
			wantTotal: 500000,
		},
		{
			name:      "duplicate pack sizes are deduplicated",
			order:     501,
			packSizes: []int{250, 250, 500, 500, 1000},
			wantPacks: map[int]int{500: 1, 250: 1},
			wantTotal: 750,
		},
		{
			name:      "overshipping when exact match impossible",
			order:     7,
			packSizes: []int{3, 5},
			wantPacks: map[int]int{3: 1, 5: 1},
			wantTotal: 8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculatePacks(tc.order, tc.packSizes)
			require.NoError(t, err)

			// Check total items.
			total := 0
			for size, qty := range result {
				total += size * qty
			}
			assert.Equal(t, tc.wantTotal, total, "total items mismatch")
			assert.Equal(t, tc.wantPacks, result, "packs mismatch")
		})
	}
}

func TestCalculatePacksErrors(t *testing.T) {
	_, err := calculatePacks(100, nil)
	assert.Error(t, err, "expected error for nil pack sizes")

	_, err = calculatePacks(100, []int{})
	assert.Error(t, err, "expected error for empty pack sizes")

	_, err = calculatePacks(100, []int{0, -5})
	assert.Error(t, err, "expected error for no valid pack sizes")
}

func BenchmarkCalculatePacks(b *testing.B) {
	packSizes := []int{17, 29, 47}
	for b.Loop() {
		calculatePacks(750000, packSizes)
	}
}