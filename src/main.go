package main

import (
	"math"
	"model"
)

var (
	WIDTH      uint = 1600
	HEIGHT     uint = 1200
	NUMWORKERS      = 10

	ex = model.Vector{1, 0, 0}
	ey = model.Vector{0, 1, 0}
	ez = model.Vector{0, 0, 1}
)

func main() {

	camera := model.NewCamera(WIDTH, HEIGHT)

	scene := model.NewScene(camera)
	scene.AddLight(model.Vector{0, 4, 0}, model.NewColor(255, 255, 255), 1000)
	scene.AddLight(model.Vector{-5, 5, 0}, model.NewColor(255, 255, 255), 1000)
	// background
	scene.Add(model.NewPlane(model.Vector{0, 0, -10}, ex, ey, model.NewColor(50, 200, 240)))
	// floor
	scene.Add(model.NewPlane(model.Vector{0, -2, 0}, ez, ex, model.NewColor(45, 200, 45)))

	// triangles
	r := Rectangle{
		model.Vector{-1, -1, -4},
		model.Vector{1, -1, -4},
		model.Vector{1, -1, -2},
		model.Vector{-1, -1, -2},
		model.NewColor(255, 0, 0)}

	//scene.Add(r.Tesselate()...)
	scene.Add(gridToTriangles(r.ToPointGrid(1))...)

	img := model.Render(scene, NUMWORKERS)
	img.Save("out.png")
}

// TRYOUT: working towards Perlin Noise landscape

//   P1           P2
//    . --------- .
//    |           |
//    |           |
//    |           |
//    |           |
//    . --------- .
//   P4           P3

type Rectangle struct {
	P1, P2, P3, P4 model.Vector
	Color          model.Color
}

func (r Rectangle) Tesselate() []model.Object {
	return []model.Object{
		model.Triangle{r.P1, r.P4, r.P2, r.Color},
		model.Triangle{r.P4, r.P3, r.P2, r.Color},
	}
}

func (r Rectangle) ToPointGrid(roughSize float64) [][]model.Vector {
	xlen := model.VectorFromTo(r.P1, r.P2).Length()
	ylen := model.VectorFromTo(r.P1, r.P4).Length()
	numDivisionsX := math.Ceil(xlen / roughSize)
	numDivisionsY := math.Ceil(ylen / roughSize)
	pointSizeX := xlen / numDivisionsX
	pointSizeY := ylen / numDivisionsY
	xVector := model.VectorFromTo(r.P1, r.P2).Normalize().Times(pointSizeX)
	yVector := model.VectorFromTo(r.P1, r.P4).Normalize().Times(pointSizeY)

	numPointsX := int(numDivisionsX) + 1
	numPointsY := int(numDivisionsY) + 1

	grid := make([][]model.Vector, numPointsY)
	for y := 0; y < numPointsY; y++ {
		row := make([]model.Vector, numPointsX)
		for x := 0; x < numPointsX; x++ {
			row[x] = r.P1.Add(xVector.Times(float64(x))).Add(yVector.Times(float64(y)))
		}
		grid[y] = row
	}

	return grid
}

func gridToTriangles(grid [][]model.Vector) []model.Object {
	ylen := len(grid)
	xlen := len(grid[0])
	triangles := []model.Object{}
	for y := 0; y < ylen-1; y++ {
		for x := 0; x < xlen-1; x++ {
			p1 := grid[y][x]
			p2 := grid[y][x+1]
			p3 := grid[y+1][x+1]
			p4 := grid[y+1][x]
			t1 := model.Triangle{p1, p4, p2, model.NewColor(255, 0, 0)}
			t2 := model.Triangle{p4, p3, p2, model.NewColor(0, 0, 255)}
			triangles = append(triangles, t1, t2)
		}
	}
	return triangles
}
