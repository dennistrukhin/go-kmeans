package kmeans

import "golang.org/x/exp/constraints"

func min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		panic("no args passed to min func")
	}

	min := args[0]
	for i := 1; i < len(args); i++ {
		if args[i] < min {
			min = args[i]
		}
	}

	return min
}
