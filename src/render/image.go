package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/deosjr/GRayT/src/model"

	"github.com/icza/mjpeg"
)

type Film struct {
	pixels        []model.Color
	width, height int
}

func newFilm(w, h int) Film {
	return Film{
		pixels: make([]model.Color, w*h),
		width:  w,
		height: h,
	}
}

// row major order
func (f Film) getArrayIndex(x, y int) int {
	return y*f.width + x
}

func (f Film) Add(x, y int, c model.Color) {
	index := f.getArrayIndex(x, y)
	current := f.pixels[index]
	f.pixels[index] = current.Add(c)
}

func (f Film) DivideBySamples(n int) {
	for i, c := range f.pixels {
		f.pixels[i] = c.Times(1.0 / float32(n))
	}
}

func (f Film) Set(x, y int, c model.Color) {
	f.pixels[f.getArrayIndex(x, y)] = c
}

func (f Film) SaveAsPNG(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	png.Encode(file, f.toImage())
}

func (f Film) SaveAsJPEG(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	jpeg.Encode(file, f.toImage(), nil)
}

func (f Film) toImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	for x := 0; x < f.width; x++ {
		for y := 0; y < f.height; y++ {
			pc := f.pixels[f.getArrayIndex(x, y)]
			c := color.RGBA{pc.R(), pc.G(), pc.B(), 255}
			img.Set(x, y, c)
		}
	}
	return img
}

var avi_fps int32 = 10

func NewAVI(filename string, w, h uint) mjpeg.AviWriter {
	aw, err := mjpeg.New(filename, int32(w), int32(h), avi_fps)
	if err != nil {
		fmt.Println(err)
	}
	return aw
}

func AddToAVI(aw mjpeg.AviWriter, f Film) {
	img := f.toImage()
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, nil); err != nil {
		fmt.Println(err)
		return
	}
	if err := aw.AddFrame(buf.Bytes()); err != nil {
		fmt.Println(err)
	}
}

func SaveAVI(aw mjpeg.AviWriter) {
	if err := aw.Close(); err != nil {
		fmt.Println(err)
	}
}

// Unfortunately golang gif support is only for limited palette gifs
// which I guess is in the standard but is very lossy for our purpose
type gifImages []*image.Paletted

// The successive delay times, one per frame, in 100ths of a second.
var gif_delay = 10

func NewGIF() gifImages {
	return []*image.Paletted{}
}

func (g gifImages) Add(f Film) gifImages {
	m := image.NewPaletted(image.Rect(0, 0, f.width, f.height), palette.Plan9)
	for x := 0; x < f.width; x++ {
		for y := 0; y < f.height; y++ {
			pc := f.pixels[f.getArrayIndex(x, y)]
			c := color.RGBA{pc.R(), pc.G(), pc.B(), 255}
			m.Set(x, y, c)
		}
	}
	return append(g, m)
}

func (g gifImages) Save(filename string) {
	delay := make([]int, len(g))
	disposal := make([]byte, len(g))
	for i := 0; i < len(g); i++ {
		delay[i] = gif_delay
		disposal[i] = gif.DisposalPrevious
	}
	gg := &gif.GIF{
		Image:     g,
		Delay:     delay,
		LoopCount: 0,
		Disposal:  disposal,
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	err = gif.EncodeAll(f, gg)
	if err != nil {
		fmt.Println(err)
	}
}
