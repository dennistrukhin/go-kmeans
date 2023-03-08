package main

import "golang.org/x/exp/constraints"

func min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		panic("no args passed to min func")
	}

	m := args[0]
	for i := 1; i < len(args); i++ {
		if args[i] < m {
			m = args[i]
		}
	}

	return m
}
