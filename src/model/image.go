package model

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Color struct {
	Vector
}

func NewColor(r, g, b uint8) Color {
	return Color{Vector{
		X: float64(r) / 255,
		Y: float64(g) / 255,
		Z: float64(b) / 255,
	}}
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

func (i Image) Save() {
	m := image.NewRGBA(image.Rect(0, 0, i.width, i.height))
	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			pc := i.pixels[y][x]
			c := color.RGBA{pc.R(), pc.G(), pc.B(), 255}
			m.Set(x, y, c)
		}
	}

	f, err := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, m)
}

// RGB vector in [0,1] -> to [0,255]
func (c Color) R() uint8 {
	return uint8(c.X * 255)
}

func (c Color) G() uint8 {
	return uint8(c.Y * 255)
}

func (c Color) B() uint8 {
	return uint8(c.Z * 255)
}
