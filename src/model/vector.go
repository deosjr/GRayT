package model

import (
	"math"

	"github.com/deosjr/GRayT/src/simd"
)

type Dimension int

const (
	X = iota
	Y
	Z
)

var Dimensions = []Dimension{X, Y, Z}

type Vector struct {
	X, Y, Z float32
}

func VectorFromTo(u, v Vector) Vector {
	return v.Sub(u)
}

func (u Vector) Add(v Vector) Vector {
	return Vector{
		X: u.X + v.X,
		Y: u.Y + v.Y,
		Z: u.Z + v.Z,
	}
}

func (u Vector) Sub(v Vector) Vector {
	return Vector{
		X: u.X - v.X,
		Y: u.Y - v.Y,
		Z: u.Z - v.Z,
	}
}

func (u Vector) Times(f float32) Vector {
	return Vector{
		X: f * u.X,
		Y: f * u.Y,
		Z: f * u.Z,
	}
}

func (u Vector) Dot(v Vector) float32 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

func (u Vector) Length() float32 {
	return float32(math.Sqrt(float64(u.Dot(u))))
}

func (u Vector) Normalize() Vector {
	l := u.Length()
	// dealing with degenerate case
	if l == 0 {
		return Vector{0, 0, 0}
	}
	return u.Times(1.0 / l)
}

func (u Vector) Cross(v Vector) Vector {
	return Vector{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

func VectorMin(u, v Vector) Vector {
	uf := [4]float32{u.X, u.Y, u.Z, 0}
	vf := [4]float32{v.X, v.Y, v.Z, 0}
	f := simd.Min(uf, vf)
	return Vector{f[0], f[1], f[2]}
}

func VectorMax(u, v Vector) Vector {
	uf := [4]float32{u.X, u.Y, u.Z, 0}
	vf := [4]float32{v.X, v.Y, v.Z, 0}
	f := simd.Max(uf, vf)
	return Vector{f[0], f[1], f[2]}
}

func (u Vector) Get(i Dimension) float32 {
	switch i {
	case X:
		return u.X
	case Y:
		return u.Y
	}
	// case Z:
	return u.Z
}

type Ray struct {
	Origin    Vector
	Direction Vector
}

func NewRay(o, d Vector) Ray {
	return Ray{
		Origin:    o,
		Direction: d.Normalize(),
	}
}

func PointFromRay(r Ray, d float32) Vector {
	return r.Origin.Add(r.Direction.Times(d))
}

func (r Ray) toSimd() ([4]float32, [4]float32, [4]float32, [4]float32, [4]float32, [4]float32) {
	rox := [4]float32{r.Origin.X, r.Origin.X, r.Origin.X, r.Origin.X}
	roy := [4]float32{r.Origin.Y, r.Origin.Y, r.Origin.Y, r.Origin.Y}
	roz := [4]float32{r.Origin.Z, r.Origin.Z, r.Origin.Z, r.Origin.Z}
	rdx := [4]float32{r.Direction.X, r.Direction.X, r.Direction.X, r.Direction.X}
	rdy := [4]float32{r.Direction.Y, r.Direction.Y, r.Direction.Y, r.Direction.Y}
	rdz := [4]float32{r.Direction.Z, r.Direction.Z, r.Direction.Z, r.Direction.Z}
	return rox, roy, roz, rdx, rdy, rdz
}
