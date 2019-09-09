package model

type Scene struct {
	Objects []Object
	Lights  []Light
	Camera  Camera

	AccelerationStructure AccelerationStructure
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

func (s *Scene) AddLights(l ...Light) {
	s.Lights = append(s.Lights, l...)
}

func (s *Scene) Precompute() {
	s.AccelerationStructure = NewBVH(s.Objects, SplitMiddle)
}

func SetBackgroundColor(c Color) {
	BACKGROUND_COLOR = c
}
