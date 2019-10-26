package render

import (
	"github.com/deosjr/GRayT/src/model"
)

// TODO: optimizations
//   - backface culling? only for opaque objects?
// - scaling: communicate over the wire
//   - memory: use protobuff ?

type worker struct {
	in  chan question
	out chan answer
}

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color model.Color
}

func (w worker) work(params Params) {
	tracer := getTracer(params.TracerType)
	if params.TracerType == model.WhittedStyle {
		params.NumSamples = 1
	}

	random := tracer.Random()
	for q := range w.in {
		x, y := float32(q.x), float32(q.y)
		var xvar, yvar float32 = 0.5, 0.5
		ray := params.Scene.Camera.PixelRay(x+xvar, y+yvar)
		color := model.NewColor(0, 0, 0)
		for i := 0; i < params.NumSamples; i++ {
			// anti-aliasing: first sample is exact middle of pixel
			// rest is randomly sampled
			if params.AntiAliasing && i != 0 {
				xvar, yvar = random.Float32(), random.Float32()
				ray = params.Scene.Camera.PixelRay(x+xvar, y+yvar)
			}
			sampleColor := tracer.GetRayColor(ray, params.Scene, 0)
			color = color.Add(sampleColor)
		}
		color = color.Times(1.0 / float32(params.NumSamples))
		w.out <- answer{q.x, q.y, color}
	}
}

type Params struct {
	Scene        *model.Scene
	NumWorkers   int
	NumSamples   int
	TracerType   model.TracerType
	AntiAliasing bool
}

func getTracer(tt model.TracerType) model.Tracer {
	switch tt {
	case model.WhittedStyle:
		return model.NewWhittedRayTracer()
	case model.Path:
		return model.NewPathTracer()
	case model.PathNextEventEstimate:
		return model.NewPathTracerNEE()
	}
	return nil
}

func Render(params Params) Film {
	w, h := params.Scene.Camera.Width(), params.Scene.Camera.Height()
	img := newFilm(w, h)

	inputChannel := make(chan question, params.NumWorkers)
	outputChannel := make(chan answer, params.NumWorkers)

	for i := 0; i < params.NumWorkers; i++ {
		worker := worker{
			in:  inputChannel,
			out: outputChannel,
		}
		go worker.work(params)
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
