package model

var (
	NUMWORKERS = 10
)

func Render(scene *Scene) {

	h := int(scene.Camera.Image.height)
	w := int(scene.Camera.Image.width)

	ch := make(chan question, NUMWORKERS)
	ans := make(chan answer, NUMWORKERS)

	for i := 0; i < NUMWORKERS; i++ {
		go worker(ch, ans)
	}

	go func() {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				ch <- question{scene, x, y}
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

	scene.Camera.Image.Save()
}
