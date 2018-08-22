package projects

import (
	"math"
	"math/rand"
	"time"

	perlin "github.com/aquilax/go-perlin"

	"github.com/deosjr/GRayT/src/model"
)

// NOTE: this folder will not be part of this project
// should move to its own project using this one as a library

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PerlinHeightMap(grid [][]model.Vector) [][]model.Vector {
	xSize, ySize := len(grid), len(grid[0])
	// alpha, beta, n iterations, random seed
	p := perlin.NewPerlin(2, 2, 3, rand.Int63())
	for y, row := range grid {
		for x, _ := range row {
			nx := float64(x)/float64(xSize) - 0.5
			ny := float64(y)/float64(ySize) - 0.5
			noise := 0.5 * p.Noise2D(nx, ny)
			noise += 0.7 * p.Noise2D(2*nx, 2*ny)
			noise += 0.25 * p.Noise2D(4*nx, 4*ny)
			noise += 0.15 * p.Noise2D(8*nx, 8*ny)
			// normalize
			noise = noise / (0.5 + 0.7 + 0.25 + 0.15)
			// map from [-1,1] to [0,1]
			noise = (noise + 1) / 2
			noise = math.Pow(noise, 3.75)
			//TODO: *2 -1 manually for camera position, remove once transformations are in
			grid[y][x].Y = noise*2 - 1
		}
	}
	return grid
}

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
