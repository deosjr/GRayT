package model

type Light struct {
	Origin Vector
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

func (s *Scene) AddLight(x, y, z float64) {
	s.Lights = append(s.Lights, Light{Vector{x, y, z}})
}
