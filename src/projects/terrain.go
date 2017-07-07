package projects

import (
	"math"

	"model"
)

// assumption: r is a rectangle
func ToPointGrid(r model.Quadrilateral, roughSize float64) [][]model.Vector {
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

func GridToTriangles(grid [][]model.Vector) []model.Object {
	ylen := len(grid)
	xlen := len(grid[0])
	triangles := []model.Object{}
	for y := 0; y < ylen-1; y++ {
		for x := 0; x < xlen-1; x++ {
			p1 := grid[y][x]
			p2 := grid[y][x+1]
			p3 := grid[y+1][x+1]
			p4 := grid[y+1][x]
			t1 := model.NewTriangle(p1, p4, p2, model.NewColor(255, 0, 0))
			t2 := model.NewTriangle(p4, p3, p2, model.NewColor(0, 0, 255))
			triangles = append(triangles, t1, t2)
		}
	}
	return triangles
}
