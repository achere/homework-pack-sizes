package pack

import (
	"fmt"
	"slices"
)

var ErrInvalidArg = fmt.Errorf("invalid arguments received")

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

	res := make(map[int]int)
	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

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

	// Optimize package quantity
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

	return res, nil
}

