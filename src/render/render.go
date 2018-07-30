package render

import "model"

// TODO: optimizations
//   - backface culling? only for opaque objects?
// - scaling: communicate over the wire
//   - memory: use protobuff ?

type question struct {
	x, y int
}

type answer struct {
	x, y  int
	color model.Color
}

func worker(scene *model.Scene, ch chan question, ans chan answer) {
	for q := range ch {
		ans <- answer{q.x, q.y, scene.GetColor(q.x, q.y)}
	}
}

func Render(scene *model.Scene, numWorkers int) Film {
	w, h := scene.Camera.Width(), scene.Camera.Height()
	img := newFilm(w, h)

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
