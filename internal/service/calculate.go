package service

import (
	"fmt"
	"math"
	"sort"
)

// CalculateService handles optimal pack calculation for orders.
type CalculateService struct {
	name      string
	packStore PackStore
}

// NewCalculateService creates a new CalculateService with the given storage.
func NewCalculateService(packStore PackStore) *CalculateService {
	return &CalculateService{
		name:      "calculate service",
		packStore: packStore,
	}
}

// Calculate returns optimal pack combination for a given order.
// Priorities:
//  1. Minimum total items (>= order, whole packs only).
//  2. Minimum number of packs (for that total).
//
// Uses DP (coin change variant) with GCD optimization.
func (s *CalculateService) Calculate(order int) (map[int]int, error) {
	if order <= 0 {
		return map[int]int{}, nil
	}

	packSizes, err := s.packStore.GetPackSizes()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}

	if len(packSizes) == 0 {
		return nil, fmt.Errorf("no pack sizes configured")
	}

	return calculatePacks(order, packSizes)
}

// calculatePacks is the core DP algorithm, separated for testability.
func calculatePacks(order int, packSizes []int) (map[int]int, error) {
	// Remove duplicates and non-positive values, sort descending.
	sizes := uniquePositive(packSizes)
	if len(sizes) == 0 {
		return nil, fmt.Errorf("no valid pack sizes provided")
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// GCD optimization: divide all sizes and target by common divisor.
	g := sizes[0]
	for _, s := range sizes[1:] {
		g = gcd(g, s)
	}

	normSizes := make([]int, len(sizes))
	for i, s := range sizes {
		normSizes[i] = s / g
	}

	// Normalized target: ceil(order / g).
	normTarget := (order + g - 1) / g

	// Upper bound: normTarget + max normalized size.
	maxNorm := normSizes[0]
	upperBound := normTarget + maxNorm

	// dp[t] = minimum number of packs to reach exactly t normalized units.
	// used[t] = which normalized pack size was added last to reach t.
	const inf = math.MaxInt32
	dp := make([]int, upperBound+1)
	used := make([]int, upperBound+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0

	// Fill DP: unbounded knapsack / coin change.
	for t := 1; t <= upperBound; t++ {
		for _, ns := range normSizes {
			prev := t - ns
			if prev < 0 || dp[prev] == inf {
				continue
			}
			candidate := dp[prev] + 1
			// Strict less-than: keeps the first found (larger pack) on ties.
			if candidate < dp[t] {
				dp[t] = candidate
				used[t] = ns
			}
		}
	}

	// Find first reachable t >= normTarget.
	bestT := -1
	for t := normTarget; t <= upperBound; t++ {
		if dp[t] != inf {
			bestT = t
			break
		}
	}
	if bestT == -1 {
		return nil, fmt.Errorf("cannot fulfil order of %d with given pack sizes", order)
	}

	// Backtrack through used[] to recover pack combination.
	packs := make(map[int]int)
	for t := bestT; t > 0; {
		normPack := used[t]
		packs[normPack*g]++
		t -= normPack
	}

	return packs, nil
}

// gcd computes greatest common divisor using Euclidean algorithm.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	
	return a
}

// uniquePositive removes duplicates and non-positive values.
func uniquePositive(nums []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, n := range nums {
		if n > 0 && !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	
	return result
}