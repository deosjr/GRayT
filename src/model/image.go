package model

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Color struct {
	r, g, b float64
}

func NewColor(r, g, b uint8) Color {
	return Color{
		r: float64(r) / 255,
		g: float64(g) / 255,
		b: float64(b) / 255,
	}
}

type Image struct {
	pixels        [][]Color
	width, height int
}

func newImage(w, h uint) Image {
	pixels := [][]Color{}
	for y := 0; y < int(h); y++ {
		pixels = append(pixels, make([]Color, w))
	}
	return Image{
		pixels: pixels,
		width:  int(w),
		height: int(h),
	}
}

func (i Image) Set(x, y int, c Color) {
	i.pixels[y][x] = c
}

func (i Image) Save(filename string) {
	m := image.NewRGBA(image.Rect(0, 0, i.width, i.height))
	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			pc := i.pixels[y][x]
			c := color.RGBA{pc.R(), pc.G(), pc.B(), 255}
			m.Set(x, y, c)
		}
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, m)
}

func (c Color) Add(d Color) Color {
	return Color{
		r: c.r + d.r,
		g: c.g + d.g,
		b: c.b + d.b,
	}
}

func (c Color) Times(f float64) Color {
	return Color{
		r: f * c.r,
		g: f * c.g,
		b: f * c.b,
	}
}

func float64touint8(f float64) uint8 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 255
	}
	return uint8(f * 255)
}

// RGB vector in [0,1] -> to [0,255]
func (c Color) R() uint8 {
	return float64touint8(c.r)
}

func (c Color) G() uint8 {
	return float64touint8(c.g)
}

func (c Color) B() uint8 {
	return float64touint8(c.b)
}
