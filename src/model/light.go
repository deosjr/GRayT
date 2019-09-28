package model

import "math"

type Light interface {
	Intensity(distance float32) float32
	Color() Color
	GetLightSegment(p Vector) Vector
}

type light struct {
	color     Color
	intensity float32
}

func (l light) Color() Color {
	return l.color
}

type PointLight struct {
	light
	origin Vector
}

func NewPointLight(o Vector, c Color, i float32) PointLight {
	return PointLight{
		origin: o,
		light: light{
			color:     c,
			intensity: i,
		},
	}
}

func (l PointLight) Intensity(r float32) float32 {
	return l.intensity / (4 * math.Pi * r * r)
}

func (l PointLight) GetLightSegment(p Vector) Vector {
	return VectorFromTo(p, l.origin)
}

type DistantLight struct {
	light
	direction Vector
}

func NewDistantLight(d Vector, c Color, i float32) DistantLight {
	return DistantLight{
		// d is direction of light rays
		// we want to store direction towards light source
		// i.e. its reverse
		direction: d.Times(-1).Normalize(),
		light: light{
			color:     c,
			intensity: i,
		},
	}
}

func (l DistantLight) Intensity(r float32) float32 {
	return l.intensity
}

func (l DistantLight) GetLightSegment(p Vector) Vector {
	return l.direction.Times(MAX_RAY_DISTANCE)
}
