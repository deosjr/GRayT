package render

import (
	"github.com/deosjr/GRayT/src/model"
)

// TODO: optimizations
//   - backface culling? only for opaque objects?
// - scaling: communicate over the wire
//   - memory: use protobuff ?

type worker struct {
	scene   *model.Scene
	in      chan question
	out     chan answer
	samples int
}

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color model.Color
}

// TODO: pass as parameter
const antiAliasing = true

func (w worker) work(tracer model.Tracer) {
	random := tracer.Random()
	for q := range w.in {
		x, y := float32(q.x), float32(q.y)
		var xvar, yvar float32 = 0.5, 0.5
		ray := w.scene.Camera.PixelRay(x+xvar, y+yvar)
		color := model.NewColor(0, 0, 0)
		for i := 0; i < w.samples; i++ {
			// anti-aliasing: first sample is exact middle of pixel
			// rest is randomly sampled
			if antiAliasing && i != 0 {
				xvar, yvar = random.Float32(), random.Float32()
				ray = w.scene.Camera.PixelRay(x+xvar, y+yvar)
			}
			sampleColor := tracer.GetRayColor(ray, w.scene, 0)
			color = color.Add(sampleColor)
		}
		color = color.Times(1.0 / float32(w.samples))
		w.out <- answer{q.x, q.y, color}
	}
}

func RenderNaive(scene *model.Scene, numWorkers int) Film {
	return render(scene, numWorkers, model.NewWhittedRayTracer, 1)
}

func RenderWithPathTracer(scene *model.Scene, numWorkers, numSamples int) Film {
	return render(scene, numWorkers, model.NewPathTracerNEE, numSamples)
}

func render(scene *model.Scene, numWorkers int, newTracerFunc func() model.Tracer, numSamples int) Film {
	w, h := scene.Camera.Width(), scene.Camera.Height()
	img := newFilm(w, h)

	inputChannel := make(chan question, numWorkers)
	outputChannel := make(chan answer, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := worker{
			scene:   scene,
			in:      inputChannel,
			out:     outputChannel,
			samples: numSamples,
		}
		tracer := newTracerFunc()
		go worker.work(tracer)
	}

	go func() {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				inputChannel <- question{x, y}
			}
		}
		close(inputChannel)
	}()

	numPixelSamples := h * w
	for {
		if numPixelSamples == 0 {
			break
		}
		a := <-outputChannel
		img.Set(a.x, a.y, a.color)
		numPixelSamples--
	}

	return img
}
