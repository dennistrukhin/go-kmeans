# K-Means for golang

[![Go Report Card](https://goreportcard.com/badge/github.com/dennistrukhin/go-kmeans)](https://goreportcard.com/report/github.com/dennistrukhin/go-kmeans)

Поддерживает любые типы данных. Для работы необходимо передать три функции:
- Metric - метрика, рассчитывающая расстояние между двумя точками
- Average - функция расчёта среднего значения для проивзольного количества точек
- Seeder - функция инициализации центроидов

Результатом является:

- Массив центроидов
- Массив принадлжености точек данных к центроидам

Массив на выходе имеет ту же размерность, что и массив на входе, но содержит 
вместо типа входных данных int с номером центроида

### Пример использования

Программа для кластеризации точек по цвету в jpeg-изображении

```go
package main

import (
	"fmt"
	"github.com/dennistrukhin/go-kmeans"
	"image"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"time"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m := func(x, y RGB) float64 {
		dR := int(y.R) - int(x.R)
		dG := int(y.G) - int(x.G)
		dB := int(y.B) - int(x.B)
		return math.Sqrt(float64(dR*dR + dG*dG + dB*dB))
	}
	s := func(_ int) RGB {
		return RGB{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
		}
	}
	a := func(args ...RGB) RGB {
		if len(args) == 0 {
			return RGB{}
		}

		r_, g_, b_ := 0, 0, 0
		for _, x := range args {
			r_ += int(x.R)
			g_ += int(x.G)
			b_ += int(x.B)
		}
		return RGB{
			R: uint8(r_ / len(args)),
			G: uint8(g_ / len(args)),
			B: uint8(b_ / len(args)),
		}
	}

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	file, err := os.Open("/Users/dennis/Downloads/1.jpeg")
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := make([]RGB, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[y*height+x] = RGB{
				R: uint8(r & 0xff),
				G: uint8(g & 0xff),
				B: uint8(b & 0xff),
			}
		}
	}

	k := kmeans.New[RGB](m, s, a)
	centroids, _ := k.Partition(pixels, 5)
	for _, c := range centroids {
		fmt.Printf("%v", c)
	}
}
```