package model

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Image struct {
	pixels        [][]uint8
	width, height int
}

func newImage(w, h uint) Image {
	pixels := [][]uint8{}
	for y := 0; y < int(h); y++ {
		pixels = append(pixels, make([]uint8, w))
	}
	return Image{
		pixels: pixels,
		width:  int(w),
		height: int(h),
	}
}

func (i Image) Set(x, y int, c uint8) {
	i.pixels[y][x] = c
}

func (i Image) Save() {
	m := image.NewRGBA(image.Rect(0, 0, i.width, i.height))
	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			c := color.RGBA{255, 255, 255, i.pixels[y][x]}
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
