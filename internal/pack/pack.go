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

	return res
}
