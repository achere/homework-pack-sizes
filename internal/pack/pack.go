package pack

import (
	"context"
	"fmt"
	"maps"
	"slices"
)

var ErrInvalidArg = fmt.Errorf("invalid arguments received")

type PackSizeRepo interface {
	GetPackSizes(context.Context) ([]int, error)
	StorePackSizes(context.Context, []int) error
}

type packSolution struct {
	packs      map[int]int
	totalItems int
	totalPacks int
	isValid    bool
}

// CalculatePacksWithRepo calculates the number of packs for a given order, fetching pack sizes from a repository,
// using the same logic as the CalculatePacks(). It respects the context passed as the first parameter.
// It returns the calculated packs, a sorted slice of available pack sizes, and any error encountered.
func CalculatePacksWithRepo(ctx context.Context, repo PackSizeRepo, order int) (map[int]int, []int, error) {
	sizes, err := repo.GetPackSizes(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get pack sizes: %w", err)
	}

	packs, err := CalculatePacks(sizes, order)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't calculate packs: %w", err)
	}

	slices.Sort(sizes)

	return packs, sizes, nil
}

// SavePackSizes saves a new set of pack sizes to the repository.
// It ensures that all provided pack sizes are positive integers.
func SavePackSizes(ctx context.Context, repo PackSizeRepo, sizes []int) error {
	for _, s := range sizes {
		if s <= 0 {
			return fmt.Errorf("%w: size amount is not positive: %d", ErrInvalidArg, s)
		}
	}

	return repo.StorePackSizes(ctx, sizes)
}

// CalculatePacks calculates the number of packs of different sizes to fulfill an order according to the following rules:
// 1. Only whole packs can be sent. Packs cannot be broken open.
// 2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
// 3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each order.
//
// Parameters:
//   - sizes: A slice of integers representing the available pack sizes.
//   - order: The total number of items ordered.
//
// Returns:
//   - A map where the keys are the pack sizes and the values are the number of packs of that size.
//   - An error if the order amount or any of the pack sizes are not positive.
func CalculatePacks(sizes []int, order int) (map[int]int, error) {
	if order <= 0 {
		return nil, fmt.Errorf("%w: order amount is not positive", ErrInvalidArg)
	}

	for _, s := range sizes {
		if s <= 0 {
			return nil, fmt.Errorf("%w: size amount is not positive: %d", ErrInvalidArg, s)
		}
	}

	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

	res, valid := calculatePacksDp(sizes, order)

	if valid {
		return res, nil
	}

	// If cannot find optimal solution, return the greedy solution
	return calculatePacksGreedy(sizes, order), nil
}

// calculatePacksDp uses dynamic programming to calculate optimal pack sizes
// Expects sizes to be in descending order
func calculatePacksDp(sizes []int, order int) (map[int]int, bool) {
	// Calculate best solution for each order amount from 1 till max allowing overflow by largest size
	// dp[i] represents the best solution for i items
	maxItems := order + sizes[len(sizes)-1]
	dp := make([]packSolution, maxItems+1)

	// Initialize: 0 items requires 0 packs (valid base case)
	dp[0] = packSolution{
		packs:      make(map[int]int),
		totalItems: 0,
		totalPacks: 0,
		isValid:    true,
	}

	for items := 1; items <= maxItems; items++ {
		dp[items] = packSolution{isValid: false}

		for _, packSize := range sizes {
			// If a valid solution that accomodates current packSize exists,
			// use it as current solution and increment packs and count packs of current of size
			if packSize <= items && dp[items-packSize].isValid {
				packs := make(map[int]int)
				maps.Copy(packs, dp[items-packSize].packs)
				packs[packSize]++

				newSolution := packSolution{
					packs:      packs,
					totalItems: items,
					totalPacks: dp[items-packSize].totalPacks + 1,
					isValid:    true,
				}

				// If current solution is better than existing, update to current
				if !dp[items].isValid || isBetterSolution(newSolution, dp[items]) {
					dp[items] = newSolution
				}
			}
		}
	}

	// Find the best solution that satisfies the order (>= order items)
	var bestSolution packSolution
	bestSolution.isValid = false

	for items := order; items <= maxItems; items++ {
		if !dp[items].isValid {
			continue
		}

		if !bestSolution.isValid || isBetterSolution(dp[items], bestSolution) {
			bestSolution = dp[items]
			break
		}
	}

	return bestSolution.packs, bestSolution.isValid
}

// isBetterSolution determines if solution A is better than solution B
func isBetterSolution(a, b packSolution) bool {
	if a.totalItems != b.totalItems {
		return a.totalItems < b.totalItems
	}
	return a.totalPacks < b.totalPacks
}

// calculatePacksGreedy implemets a greedy strategy for calculating packs while handling some edge cases.
// Expects sizes to be sorted in descending order
func calculatePacksGreedy(sizes []int, order int) map[int]int {
	res := make(map[int]int)

	// Naive packing
	sizeIdx := 0
	sum := 0
	for sum < order && sizeIdx < len(sizes) {
		currSize := sizes[sizeIdx]
		rest := order - sum
		if currSize > rest {
			sizeIdx++
			continue
		}

		n := rest / currSize

		res[currSize] = n
		sum += n * currSize
		sizeIdx++
	}

	// Add last pack if order is not fulfilled
	if sum < order {
		rest := order - sum
		for i := len(sizes) - 1; i >= 0; i-- {
			if sizes[i] > rest {
				if p, ok := res[sizes[i]]; ok {
					res[sizes[i]] = p + 1
				} else {
					res[sizes[i]] = 1
				}
				break
			}
		}
	}

	// Optimize package quantity with bigger sizes
	qty := 0
	sum = 0
	delKey := 0
	for i := len(sizes) - 1; i >= 0; i-- {
		pack, ok := res[sizes[i]]
		if !ok {
			continue
		}

		qty += pack
		sum += sizes[i] * pack

		if i-1 < 0 {
			break
		}
		nextSize := sizes[i-1]

		nextQty := sum / nextSize
		nextSum := nextQty * nextSize

		if nextSum == sum && nextQty < qty {
			delKey = sizes[i]
			res[nextSize] = nextQty
			break
		}
	}

	if delKey > 0 {
		for key := range res {
			if key <= delKey {
				delete(res, key)
			}
		}
	}

	return res
}
