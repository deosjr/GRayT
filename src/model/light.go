package model

import "math"

type Light interface {
	Intensity(distance float64) float64
	Color() Color
	Origin() Vector
}

type PointLight struct {
	origin    Vector
	color     Color
	intensity float64
}

func NewPointLight(o Vector, c Color, i float64) PointLight {
	return PointLight{
		origin:    o,
		color:     c,
		intensity: i,
	}
}

func (l PointLight) Intensity(r float64) float64 {
	return l.intensity / (4 * math.Pi * r * r)
}

func (l PointLight) Color() Color {
	return l.color
}

func (l PointLight) Origin() Vector {
	return l.origin
}
