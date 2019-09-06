package render

import (
	"github.com/deosjr/GRayT/src/model"
)

// TODO: optimizations
//   - backface culling? only for opaque objects?
// - scaling: communicate over the wire
//   - memory: use protobuff ?

type worker struct {
	scene *model.Scene
	in    chan question
	out   chan answer
}

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color model.Color
}

func (w worker) work(tracer model.Tracer) {
	for q := range w.in {
		ray := w.scene.Camera.PixelRay(q.x, q.y)
		color := tracer.GetRayColor(ray, w.scene, 0)
		w.out <- answer{q.x, q.y, color}
	}
}

func RenderNaive(scene *model.Scene, numWorkers int) Film {
	return render(scene, numWorkers, model.NewNaiveRayTracer)
}

func RenderWithPathTracer(scene *model.Scene, numWorkers, numSamples int) Film {
	f := func() model.Tracer {
		return model.NewPathTracer(numSamples)
	}
	return render(scene, numWorkers, f)
}

func render(scene *model.Scene, numWorkers int, newTracerFunc func() model.Tracer) Film {
	w, h := scene.Camera.Width(), scene.Camera.Height()
	img := newFilm(w, h)

	ch := make(chan question, numWorkers)
	ans := make(chan answer, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := worker{
			scene: scene,
			in:    ch,
			out:   ans,
		}
		tracer := newTracerFunc()
		go worker.work(tracer)
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
