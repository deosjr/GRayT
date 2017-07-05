package model

import "math"

var STANDARD_ALBEDO = 0.18

type Object interface {
	Intersect(Ray) (distance float64, ok bool)
	SurfaceNormal(point Vector) Vector
	GetColor() Color
}

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

func (l PointLight) Intensity(r float64) float64 {
	return l.intensity / (4 * math.Pi * r * r)
}

func (l PointLight) Color() Color {
	return l.color
}

func (l PointLight) Origin() Vector {
	return l.origin
}

type Scene struct {
	Objects []Object
	Lights  []Light
	Camera  Camera
}

func NewScene(camera Camera) *Scene {
	return &Scene{
		Objects: []Object{},
		Lights:  []Light{},
		Camera:  camera,
	}
}

func (s *Scene) Add(o ...Object) {
	s.Objects = append(s.Objects, o...)
}

func (s *Scene) AddLight(o Vector, c Color, i float64) {
	s.Lights = append(s.Lights, PointLight{o, c, i})
}
