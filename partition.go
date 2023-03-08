package kmeans

import (
	"runtime"
	"sync"
)

type Metric[T interface{}] func(x, y T) float64
type Seeder[T interface{}] func(n int) T
type Average[T interface{}] func(args ...T) T

type KMeans[T interface{}] struct {
	metric  Metric[T]
	seeder  Seeder[T]
	average Average[T]
	epsilon float64
	maxGens int
}

func New[T interface{}](m Metric[T], s Seeder[T], a Average[T]) *KMeans[T] {
	return &KMeans[T]{
		metric:  m,
		seeder:  s,
		average: a,
		epsilon: 1.0,
		maxGens: 500,
	}
}

func (km *KMeans[T]) WithEpsilon(e float64) {
	km.epsilon = e
}

func (km *KMeans[T]) WithMaxGens(n int) {
	km.maxGens = n
}

func (km *KMeans[T]) Partition(data []T, num int) ([]T, *[]int) {
	dl := len(data)
	centroids := make([]T, num)
	mapping := make([]int, dl)
	for i := 0; i < dl; i++ {
		mapping[i] = -1
	}
	for i := 0; i < num; i++ {
		centroids[i] = km.seeder(i)
	}

	gen := 0
	numProc := runtime.GOMAXPROCS(0)
	batchSize := dl/numProc + 1

	for {
		gen++
		var wgd sync.WaitGroup
		for p := 0; p < numProc; p++ {
			start := p * batchSize
			end := min((p+1)*batchSize, dl)
			wgd.Add(1)
			go func(wg *sync.WaitGroup, m *[]int, start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					d := make([]float64, num)
					for j, c := range centroids {
						d[j] = km.metric(data[i], c)
					}
					(*m)[i] = minIndex(d)
				}
			}(&wgd, &mapping, start, end)
		}
		wgd.Wait()

		var wgc sync.WaitGroup
		converged := 0
		for p := 0; p < num; p++ {
			wgc.Add(1)
			go func(wg *sync.WaitGroup, n int) {
				defer wg.Done()
				points := make([]T, 1)
				for j := 0; j < dl; j++ {
					if mapping[j] == n {
						points = append(points, data[j])
					}
				}
				avg := km.average(points...)
				step := km.metric(avg, centroids[n])
				if step <= km.epsilon {
					converged++
				}
				centroids[n] = avg
			}(&wgc, p)
		}
		wgc.Wait()

		if converged == num {
			break
		}
		if gen >= km.maxGens {
			break
		}
	}

	return centroids, &mapping
}
