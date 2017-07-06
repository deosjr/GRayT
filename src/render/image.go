package render

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"model"
)

type Image struct {
	pixels        [][]model.Color
	width, height int
}

func newImage(w, h int) Image {
	pixels := [][]model.Color{}
	for y := 0; y < h; y++ {
		pixels = append(pixels, make([]model.Color, w))
	}
	return Image{
		pixels: pixels,
		width:  w,
		height: h,
	}
}

func (i Image) Set(x, y int, c model.Color) {
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
