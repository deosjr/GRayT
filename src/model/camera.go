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

func NewCamera(w, h uint) *Camera {
	r := NewRay(Vector{0, 0, 0}, Vector{0, 0, -1})
	wf := float64(w)
	hf := float64(h)
	fov := 0.5 * math.Pi
	return &Camera{
		Ray:         r,
		FieldOfView: fov,
		Width:       wf,
		Height:      hf,
	}
}

// TODO: move and rotate are WIP (and never called right now)

func (c *Camera) Move(v Vector) {
	newOrigin := c.Ray.Origin.Add(v)
	c.Ray = NewRay(newOrigin, c.Ray.Direction)
}

func (c *Camera) Rotate() {
	newDirection := Vector{} // TODO: funky rotation stuff
	c.Ray = NewRay(c.Ray.Origin, newDirection)
}

func (c *Camera) Precompute() {
	w := c.Ray.Direction
	up := Vector{0, 1, 0}
	c.u = w.Cross(up)
	c.v = c.u.Cross(w)

	tanx := math.Tan(c.FieldOfView / 2)
	tany := math.Tan((c.Height / c.Width) * c.FieldOfView / 2)

	c.pixelWidth = tanx / (c.Width / 2)
	c.pixelHeight = tany / (c.Height / 2)
	ULHC := c.Ray.Origin.Add(w).Sub(c.u.Times((c.Width / 2) * c.pixelWidth)).Add(c.v.Times((c.Height / 2) * c.pixelHeight))
	c.ulhc = ULHC.Add(c.u.Times(c.pixelWidth / 2)).Sub(c.v.Times(c.pixelHeight / 2))
}

// 3d translation of 2d point on view
// assumes width >= height
// view size is 2x2 in world, from -1,1 to 1,-1
// viewing window lives in ex x ey plane (z = -1)
// viewport coordinate system (u,v,w)
func (c *Camera) PixelRay(x, y int) Ray {
	xfactor := c.pixelWidth * float64(x)
	yfactor := c.pixelHeight * float64(y)
	pixel := c.ulhc.Add(c.u.Times(xfactor)).Sub(c.v.Times(yfactor))

	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}
