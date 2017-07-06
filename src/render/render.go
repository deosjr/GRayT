package render

import "model"

func Render(scene *Scene, numWorkers int) Image {

	w, h := int(scene.Camera.Width), int(scene.Camera.Height)
	img := newImage(w, h)

	ch := make(chan question, numWorkers)
	ans := make(chan answer, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(scene, ch, ans)
	}

	go func() {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				ch <- question{x, y}
			}
		}
		close(ch)
	}()

	numPixels := h * w
	for {
		if numPixels == 0 {
			break
		}
		a := <-ans
		img.Set(a.x, a.y, a.color)
		numPixels--
	}

	return img
}

type Scene struct {
	Objects []model.Object
	Lights  []model.Light
	Camera  model.Camera
}

func NewScene(camera model.Camera) *Scene {
	return &Scene{
		Objects: []model.Object{},
		Lights:  []model.Light{},
		Camera:  camera,
	}
}

func (s *Scene) Add(o ...model.Object) {
	s.Objects = append(s.Objects, o...)
}

func (s *Scene) AddLights(l ...model.Light) {
	s.Lights = append(s.Lights, l...)
}