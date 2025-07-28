package pack

import (
	"slices"
)

func CalculatePacks(sizes []int, order int) map[int]int {
	res := make(map[int]int)
	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

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

		withNextSize := sum / nextSize
		prop := withNextSize * nextSize

		if prop == sum && withNextSize < qty {
			delKey = sizes[i]
			res[nextSize] = withNextSize
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
