package kmeans

import (
	"golang.org/x/exp/constraints"
)

func minIndex[T constraints.Ordered](args []T) int {
	l := len(args)
	if l == 0 {
		panic("Empty args")
	}
	index := 0
	min := args[0]
	for i := 1; i < l; i++ {
		if args[i] < min {
			min = args[i]
			index = i
		}
	}
	return index
}
