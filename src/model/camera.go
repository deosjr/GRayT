package model

import (
	"math"
)

type Camera struct {
	Ray
	FieldOfView   float64
	Width, Height float64
	// precomputed for pixelray() method
	ulhc                    Vector
	u, v                    Vector
	pixelWidth, pixelHeight float64
}

func NewCamera(w, h uint) Camera {
	r := NewRay(Vector{0, 0, 0}, Vector{0, 0, -1})
	wf := float64(w)
	hf := float64(h)
	fov := 0.5 * math.Pi
	ulhc, u, v, pw, ph := precompute(r, wf, hf, fov)
	return Camera{
		Ray:         r,
		FieldOfView: fov,
		Width:       wf,
		Height:      hf,
		ulhc:        ulhc,
		u:           u,
		v:           v,
		pixelWidth:  pw,
		pixelHeight: ph,
	}
}

func precompute(r Ray, xres, yres, fov float64) (Vector, Vector, Vector, float64, float64) {
	w := r.Direction
	up := Vector{0, 1, 0}
	u := w.Cross(up)
	v := u.Cross(w)

	tanx := math.Tan(fov / 2)
	tany := math.Tan((yres / xres) * fov / 2)

	pixelWidth := tanx / (xres / 2)
	pixelHeight := tany / (yres / 2)
	ULHC := r.Origin.Add(w).Sub(u.Times((xres / 2) * pixelWidth)).Add(v.Times((yres / 2) * pixelHeight))
	ULHCmid := ULHC.Add(u.Times(pixelWidth / 2)).Sub(v.Times(pixelHeight / 2))
	return ULHCmid, u, v, pixelWidth, pixelHeight
}

// 3d translation of 2d point on view
// assumes width >= height
// view size is 2x2 in world, from -1,1 to 1,-1
// viewing window lives in ex x ey plane (z = -1)

// METHOD 1
// viewport coordinate system (u,v,w)
func (c Camera) PixelRay(x, y int) Ray {
	xfactor := c.pixelWidth * float64(x)
	yfactor := c.pixelHeight * float64(y)
	pixel := c.ulhc.Add(c.u.Times(xfactor)).Sub(c.v.Times(yfactor))

	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}

// TODO: currently unused!
// METHOD 2
// NDC = percentage of total width/height on (0,0) (1,1) screen
// with (0,0) as lower left hand corner (LLHC)
// Screen = that scaled to (-1,1) (1,-1) screen with (-1,1) as ULHC
func (c Camera) pixelRay(x, y int) Ray {
	xNDC := (float64(x) + 0.5) / c.Width
	yNDC := (float64(y) + 0.5) / c.Height

	xScreen := 2*xNDC - 1
	yScreen := 1 - 2*yNDC

	tanx := math.Tan(c.FieldOfView / 2)
	tany := math.Tan((c.Height / c.Width) * c.FieldOfView / 2)

	px := xScreen * tanx
	py := yScreen * tany
	pixel := Vector{px, py, -1}

	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}
