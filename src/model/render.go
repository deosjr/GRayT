package model

func Render(scene *Scene, numWorkers int) Image {

	h := int(scene.Camera.Image.height)
	w := int(scene.Camera.Image.width)

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
		scene.Camera.Image.Set(a.x, a.y, a.color)
		numPixels--
	}

	return scene.Camera.Image
}
