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

		//sumSampleColor := model.NewColor(0, 0, 0)
		//for i := 0; i < 1000; i++ {
		color := tracer.GetRayColor(ray, w.scene, 0)
		//sumSampleColor = sumSampleColor.Add(color)
		//}
		//color := sumSampleColor.Times(1.0 / 1000.0)

		w.out <- answer{q.x, q.y, color}
	}
}

func RenderNaive(scene *model.Scene, numWorkers int) Film {
	return render(scene, numWorkers, model.NewWhittedRayTracer, 1)
}

func RenderWithPathTracer(scene *model.Scene, numWorkers, numSamples int) Film {
	return render(scene, numWorkers, model.NewPathTracer, numSamples)
}

func render(scene *model.Scene, numWorkers int, newTracerFunc func() model.Tracer, numSamples int) Film {
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
				for s := 0; s < numSamples; s++ {
					ch <- question{x, y}
				}
			}
		}
		close(ch)
	}()

	numPixelSamples := h * w * numSamples
	for {
		if numPixelSamples == 0 {
			break
		}
		a := <-ans
		img.Add(a.x, a.y, a.color)
		numPixelSamples--
	}

	img.DivideBySamples(numSamples)

	return img
}
