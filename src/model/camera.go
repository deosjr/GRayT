package model

import "math"

type Camera struct {
	Ray
	FieldOfView float64
	Image       Image
}

func NewCamera(w, h uint) Camera {
	img := newImage(w, h)
	r := NewRay(Vector{0, 0, 0}, Vector{0, 0, -1})
	return Camera{
		Ray:         r,
		FieldOfView: 0.5 * math.Pi,
		Image:       img,
	}
}

// 3d translation of 2d point on view
// assumes width >= height
// view size is 2x2 in world, from -1,1 to 1,-1
func (c Camera) PixelRay(x, y int) Ray {
	w := float64(c.Image.width)
	h := float64(c.Image.height)
	aspectRatio := w / h
	tan := math.Tan(c.FieldOfView / 2)

	xNDC := (float64(x) + 0.5) / w
	yNDC := (float64(y) + 0.5) / h

	xScreen := 2*xNDC - 1
	yScreen := 1 - 2*yNDC

	px := xScreen * aspectRatio * tan
	py := yScreen * tan

	pixel := Vector{px, py, -1}
	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}
