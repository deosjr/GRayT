package model

import "math"

type Dimension int

const (
	X = iota
	Y
	Z
)

type Vector struct {
	X, Y, Z float64
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

func (u Vector) Times(f float64) Vector {
	return Vector{
		X: f * u.X,
		Y: f * u.Y,
		Z: f * u.Z,
	}
}

func (u Vector) Dot(v Vector) float64 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

func (u Vector) Length() float64 {
	return math.Sqrt(u.Dot(u))
}

func (u Vector) Normalize() Vector {
	l := u.Length()
	// dealing with degenerate case
	if l == 0 {
		return Vector{0, 0, 0}
	}
	return u.Times(1.0 / u.Length())
}

func (u Vector) Cross(v Vector) Vector {
	return Vector{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

func (u Vector) Get(i Dimension) float64 {
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

func PointFromRay(r Ray, d float64) Vector {
	return r.Origin.Add(r.Direction.Times(d))
}
