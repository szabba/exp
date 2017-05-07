package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
)

var (
	fname string

	width, height int

	p float64
	N int
)

func main() {
	flag.StringVar(&fname, "fname", "caves.png", "name of the output file")

	flag.IntVar(&width, "width", 500, "width of the map")
	flag.IntVar(&height, "height", 500, "height of the map")

	flag.Float64Var(&p, "p", 0.5, "probability of a cell initially containing a floor")
	flag.IntVar(&N, "N", 1, "number of smoothing iterations")

	flag.Parse()

	g := newGrid(width, height)
	g.step(func(x, y int) int {
		if rand.Float64() < p {
			return Floor
		}
		return Wall
	})

	for i := 0; i < N; i++ {
		g.step(boxSmooth(g, 1, 1))
		g.step(boxSmooth(g, 3, 3))
	}

	f, err := os.Create("caves.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	png.Encode(f, draw(g, colorFor))
}

func boxSmooth(g *grid, Δx, Δy int) rule {
	area := float64((1 + 2*Δx) * (1 + 2*Δy))
	return func(x, y int) int {

		floors := 0.0
		for i := -Δx; i < Δx+1; i++ {
			for j := -Δy; j < Δy+1; j++ {
				if g.get(x+i, y+j) == Floor {
					floors++
				}
			}
		}

		if floors > p*area {
			return Wall
		}
		if floors < p*area {
			return Floor
		}
		return g.get(x, y)
	}
}

const (
	Floor = iota
	Wall
)

func colorFor(v int) color.Color {
	switch v {
	case Floor:
		return color.White
	case Wall:
		return color.Black
	default:
		return color.Transparent
	}
}

func draw(g *grid, f func(int) color.Color) image.Image {
	img := image.NewNRGBA64(image.Rect(0, 0, g.width, g.height))
	g.eachLocation(func(x, y int) {
		img.Set(x, y, f(g.get(x, y)))
	})
	return img
}

type grid struct {
	width, height int
	current, next []int
}

func newGrid(width, height int) *grid {
	return &grid{
		width, height,
		make([]int, width*height),
		make([]int, width*height),
	}
}

func (g *grid) get(x, y int) int {
	for x < 0 {
		x += g.width
	}
	for x >= g.width {
		x -= g.width
	}
	for y < 0 {
		y += g.height
	}
	for y >= g.height {
		y -= g.height
	}
	return g.current[y*g.width+x]
}

type rule func(x, y int) int

func (g *grid) step(r rule) {
	g.eachLocation(func(x, y int) {
		g.setNext(x, y, r(x, y))
	})
	g.swap()
}

func (g *grid) eachLocation(f func(x, y int)) {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			f(x, y)
		}
	}
}

func (g *grid) swap() {
	g.current, g.next = g.next, g.current
}

func (g *grid) setNext(x, y, v int) {
	g.next[y*g.width+x] = v
}
