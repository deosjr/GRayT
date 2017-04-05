package model

var STANDARD_ALBEDO = 0.18

type Object interface {
	Intersect(Ray) (intersection Vector, ok bool, distance float64)
	SurfaceNormal(point Vector) Vector
}

type Light struct {
	Origin    Vector
	Color     Color
	Intensity float64
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

func (s *Scene) Add(o Object) {
	s.Objects = append(s.Objects, o)
}

func (s *Scene) AddLight(o Vector, c Color, i float64) {
	s.Lights = append(s.Lights, Light{o, c, i})
}
