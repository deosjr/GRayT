package model

import "math/rand"

type Scene struct {
	Objects []Object
	// TODO: consolidate these two
	Lights   []Light
	Emitters []Triangle

	Camera                Camera
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
	// TODO: loop over all objects and add emitters to emmiters list
	s.AccelerationStructure = NewBVH(s.Objects, SplitSurfaceAreaHeuristic)
}

func (s *Scene) randomEmitter(random *rand.Rand) Triangle {
	if len(s.Emitters) == 0 {
		panic("no light in scene!")
	}
	return s.Emitters[random.Intn(len(s.Emitters))]
}

func SetBackgroundColor(c Color) {
	BACKGROUND_COLOR = c
}
