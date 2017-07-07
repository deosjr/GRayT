package model

// Builds on quadrilateral definition:
// Let P1 - P4 be the top and
// let P5 - P8 be the bottom quadrilateral
// P1 corresponding to P5 etc

type Cuboid struct {
	P1, P2, P3, P4, P5, P6, P7, P8 Vector
	Color                          Color
}

func (c Cuboid) Tesselate() []Object {
	F1 := Quadrilateral{c.P1, c.P2, c.P3, c.P4, c.Color}.Tesselate()
	F2 := Quadrilateral{c.P2, c.P1, c.P5, c.P6, c.Color}.Tesselate()
	F3 := Quadrilateral{c.P3, c.P2, c.P6, c.P7, c.Color}.Tesselate()
	F4 := Quadrilateral{c.P4, c.P3, c.P7, c.P8, c.Color}.Tesselate()
	F5 := Quadrilateral{c.P1, c.P4, c.P8, c.P5, c.Color}.Tesselate()
	F6 := Quadrilateral{c.P6, c.P5, c.P8, c.P7, c.Color}.Tesselate()
	return append(F1, append(F2, append(F3, append(F4, append(F5, F6...)...)...)...)...)
}

//   P1           P2
//    . --------- .
//    |           |
//    |           |
//    |           |
//    |           |
//    . --------- .
//   P4           P3

type Quadrilateral struct {
	P1, P2, P3, P4 Vector
	Color          Color
}

func (r Quadrilateral) Tesselate() []Object {
	return []Object{
		NewTriangle(r.P1, r.P4, r.P2, r.Color),
		NewTriangle(r.P4, r.P3, r.P2, r.Color),
	}
}
