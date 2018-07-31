package model

import "math"

type Light interface {
	Intensity(distance float64) float64
	Color() Color
	VectorFromPoint(p Vector) Vector
}

type light struct {
	color     Color
	intensity float64
}

func (l light) Color() Color {
	return l.color
}

type PointLight struct {
	light
	origin Vector
}

func NewPointLight(o Vector, c Color, i float64) PointLight {
	return PointLight{
		origin: o,
		light: light{
			color:     c,
			intensity: i,
		},
	}
}

func (l PointLight) Intensity(r float64) float64 {
	return l.intensity / (4 * math.Pi * r * r)
}

func (l PointLight) VectorFromPoint(p Vector) Vector {
	return VectorFromTo(p, l.origin)
}

type DistantLight struct {
	light
	direction Vector
}

func NewDistantLight(d Vector, c Color, i float64) DistantLight {
	return DistantLight{
		direction: d,
		light: light{
			color:     c,
			intensity: i,
		},
	}
}

func (l DistantLight) Intensity(r float64) float64 {
	return l.intensity
}

func (l DistantLight) VectorFromPoint(p Vector) Vector {
	distantSource := p.Add(l.direction)
	return VectorFromTo(p, distantSource)
}
