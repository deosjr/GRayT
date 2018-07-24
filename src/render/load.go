package render

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"model"
)

// LoadObj assumes filename contains one triangle mesh object
func LoadObj(filename string, c model.Color) (model.Object, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return loadObj(scanner, c)
}

func loadObj(scanner *bufio.Scanner, c model.Color) (model.Object, error) {
	var vertices []model.Vector
	var triangles []model.Object
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key, values := fields[0], fields[1:]
		switch key {
		case "#":
			continue
		case "v":
			vertex, err := readVertex(values)
			if err != nil {
				return nil, err
			}
			vertices = append(vertices, vertex)
		case "f":
			face, err := readTriangle(values, vertices, c)
			if err != nil {
				return nil, err
			}
			triangles = append(triangles, face)
		default:
			fmt.Printf("Unexpected line: %s", line)
		}
	}
	return toObject(triangles)
}

func readVertex(coordinates []string) (model.Vector, error) {
	if len(coordinates) != 3 {
		return model.Vector{}, fmt.Errorf("Invalid coordinates: %v", coordinates)
	}
	p1, err := strconv.ParseFloat(coordinates[0], 64)
	if err != nil {
		return model.Vector{}, err
	}
	p2, err := strconv.ParseFloat(coordinates[1], 64)
	if err != nil {
		return model.Vector{}, err
	}
	p3, err := strconv.ParseFloat(coordinates[2], 64)
	if err != nil {
		return model.Vector{}, err
	}
	return model.Vector{p1, p2, p3}, nil
}

func readTriangle(indices []string, vertices []model.Vector, c model.Color) (model.Triangle, error) {
	if len(indices) != 3 {
		return model.Triangle{}, fmt.Errorf("Invalid indices: %v", indices)
	}
	i1, err := strconv.ParseInt(indices[0], 10, 64)
	if err != nil {
		return model.Triangle{}, err
	}

	numVertices := int64(len(vertices))
	if i1 < 1 || numVertices < i1 {
		return model.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i1, numVertices)
	}
	i2, err := strconv.ParseInt(indices[1], 10, 64)
	if err != nil {
		return model.Triangle{}, err
	}
	if i2 < 1 || numVertices < i2 {
		return model.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i2, numVertices)
	}
	i3, err := strconv.ParseInt(indices[2], 10, 64)
	if err != nil {
		return model.Triangle{}, err
	}
	if i3 < 1 || numVertices < i3 {
		return model.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i3, numVertices)
	}
	// TODO: coordinate handedness!
	return model.NewTriangle(vertices[i3-1], vertices[i2-1], vertices[i1-1], c), nil
}

func toObject(triangles []model.Object) (model.Object, error) {
	if len(triangles) == 0 {
		return nil, errors.New("Object list empty")
	}
	b := model.ObjectsBound(triangles)
	objectToOrigin := model.Translate(b.Centroid()).Inverse()

	for i, tobj := range triangles {
		t := tobj.(model.Triangle)
		triangles[i] = model.NewTriangle(
			objectToOrigin.Point(t.P0),
			objectToOrigin.Point(t.P1),
			objectToOrigin.Point(t.P2),
			t.Color)
	}
	return model.NewComplexObject(triangles), nil
}
