package model

import "fmt"
import "math"

// row-major order matrix
type matrix4x4 [4][4]float64

type Transform struct {
	// store both matrix and its inverse
	m, mInv matrix4x4
}

func NewTransform(m matrix4x4) Transform {
	return Transform{
		m:    m,
		mInv: m.inverse(),
	}
}

func (t Transform) Point(p Vector) Vector {
	x, y, z := p.X, p.Y, p.Z
	pp := Vector{
		t.m[0][0]*x + t.m[0][1]*y + t.m[0][2]*z + t.m[0][3],
		t.m[1][0]*x + t.m[1][1]*y + t.m[1][2]*z + t.m[1][3],
		t.m[2][0]*x + t.m[2][1]*y + t.m[2][2]*z + t.m[2][3],
	}
	wp := t.m[3][0]*x + t.m[3][1]*y + t.m[3][2]*z + t.m[3][3]
	if wp == 1 || wp == 0 {
		return pp
	}
	return pp.Times(1 / wp)
}

func (t Transform) Vector(v Vector) Vector {
	x, y, z := v.X, v.Y, v.Z
	return Vector{
		t.m[0][0]*x + t.m[0][1]*y + t.m[0][2]*z,
		t.m[1][0]*x + t.m[1][1]*y + t.m[1][2]*z,
		t.m[2][0]*x + t.m[2][1]*y + t.m[2][2]*z,
	}
}

func (t Transform) Ray(r Ray) Ray {
	//TODO: floating-point rounding errors
	return NewRay(t.Point(r.Origin), t.Vector(r.Direction))
}

// For any invertible n-by-n matrices A and B, (AB)−1 = B−1A−1
// where -1 = inverse.
func (t1 Transform) Mul(t2 Transform) Transform {
	return Transform{
		m:    t1.m.multiply(t2.m),
		mInv: t2.mInv.multiply(t1.mInv),
	}
}

func (t Transform) Inverse() Transform {
	return Transform{
		m:    t.mInv,
		mInv: t.m,
	}
}

func (t Transform) Transpose() Transform {
	return Transform{
		m:    t.m.transpose(),
		mInv: t.mInv.transpose(),
	}
}

func Translate(delta Vector) Transform {
	return Transform{
		m: matrix4x4{
			{1, 0, 0, delta.X},
			{0, 1, 0, delta.Y},
			{0, 0, 1, delta.Z},
			{0, 0, 0, 1},
		},
		mInv: matrix4x4{
			{1, 0, 0, -delta.X},
			{0, 1, 0, -delta.Y},
			{0, 0, 1, -delta.Z},
			{0, 0, 0, 1},
		},
	}
}

func ScaleUniform(x float64) Transform {
	return Scale(x, x, x)
}

//TODO: source of divisionByZero NaNs
// error handling? test!
func Scale(x, y, z float64) Transform {
	if x == 0 || y == 0 || z == 0 {
		fmt.Println("SCALING BY ZERO")
	}
	return Transform{
		m: matrix4x4{
			{x, 0, 0, 0},
			{0, y, 0, 0},
			{0, 0, z, 0},
			{0, 0, 0, 1},
		},
		mInv: matrix4x4{
			{1 / x, 0, 0, 0},
			{0, 1 / y, 0, 0},
			{0, 0, 1 / z, 0},
			{0, 0, 0, 1},
		},
	}
}

//TODO: check: theta in radians
//TODO: I either don't understand these or
// they are horribly broken. fix/test
func RotateX(theta float64) Transform {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	m := matrix4x4{
		{1, 0, 0, 0},
		{0, cosTheta, -sinTheta, 0},
		{0, sinTheta, cosTheta, 0},
		{0, 0, 0, 1},
	}
	return Transform{m, m.transpose()}
}

func RotateY(theta float64) Transform {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	m := matrix4x4{
		{cosTheta, 0, sinTheta, 0},
		{0, 1, 0, 0},
		{-sinTheta, 0, cosTheta, 0},
		{0, 0, 0, 1},
	}
	return Transform{m, m.transpose()}
}

func RotateZ(theta float64) Transform {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	m := matrix4x4{
		{cosTheta, -sinTheta, 0, 0},
		{sinTheta, cosTheta, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return Transform{m, m.transpose()}
}

func Rotate(theta float64, axis Vector) Transform {
	v := axis.Normalize()
	s := math.Sin(theta)
	c := math.Cos(theta)
	minC := 1 - c
	m := matrix4x4{
		{minC*v.X*v.X + c, minC*v.X*v.Y + v.Z*s, minC*v.Z*v.X - v.Y*s, 0},
		{minC*v.X*v.Y - v.Z*s, minC*v.Y*v.Y + c, minC*v.Y*v.Z + v.X*s, 0},
		{minC*v.Z*v.X + v.Y*s, minC*v.Y*v.Z - v.X*s, minC*v.Z*v.Z + c, 0},
		{0, 0, 0, 1},
	}
	return NewTransform(m)
}

func (m1 matrix4x4) multiply(m2 matrix4x4) matrix4x4 {
	r := matrix4x4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[i][j] = m1[i][0]*m2[0][j] +
				m1[i][1]*m2[1][j] +
				m1[i][2]*m2[2][j] +
				m1[i][3]*m2[3][j]
		}
	}
	return r
}

func (m matrix4x4) transpose() matrix4x4 {
	return matrix4x4{
		{m[0][0], m[1][0], m[2][0], m[3][0]},
		{m[0][1], m[1][1], m[2][1], m[3][1]},
		{m[0][2], m[1][2], m[2][2], m[3][2]},
		{m[0][3], m[1][3], m[2][3], m[3][3]},
	}
}

func (m matrix4x4) determinant() float64 {
	return (m[0][0]*m[1][1]*m[2][2]*m[3][3] - m[0][0]*m[1][1]*m[2][3]*m[3][2] +
		m[0][0]*m[1][2]*m[2][3]*m[3][1] - m[0][0]*m[1][2]*m[2][1]*m[3][3] +
		m[0][0]*m[1][3]*m[2][1]*m[3][2] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
		m[0][1]*m[1][2]*m[2][3]*m[3][0] + m[0][1]*m[1][2]*m[2][0]*m[3][3] -
		m[0][1]*m[1][3]*m[2][0]*m[3][2] + m[0][1]*m[1][3]*m[2][2]*m[3][0] -
		m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][1]*m[1][0]*m[2][3]*m[3][2] +
		m[0][2]*m[1][3]*m[2][0]*m[3][1] - m[0][2]*m[1][3]*m[2][1]*m[3][0] +
		m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][2]*m[1][0]*m[2][3]*m[3][1] +
		m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][2]*m[1][1]*m[2][0]*m[3][3] -
		m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][3]*m[1][0]*m[2][2]*m[3][1] -
		m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][3]*m[1][1]*m[2][0]*m[3][2] -
		m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][3]*m[1][2]*m[2][1]*m[3][0])
}

// TODO: determinant of 0 means singular (noninvertible) matrix
func (m matrix4x4) inverse() matrix4x4 {
	r := matrix4x4{}
	d := m.determinant()
	if d == 0 {
		fmt.Println("ZERO DETERMINANT")
	}
	r[0][0] = (m[1][2]*m[2][3]*m[3][1] - m[1][3]*m[2][2]*m[3][1] + m[1][3]*m[2][1]*m[3][2] - m[1][1]*m[2][3]*m[3][2] - m[1][2]*m[2][1]*m[3][3] + m[1][1]*m[2][2]*m[3][3]) / d
	r[0][1] = (m[0][3]*m[2][2]*m[3][1] - m[0][2]*m[2][3]*m[3][1] - m[0][3]*m[2][1]*m[3][2] + m[0][1]*m[2][3]*m[3][2] + m[0][2]*m[2][1]*m[3][3] - m[0][1]*m[2][2]*m[3][3]) / d
	r[0][2] = (m[0][2]*m[1][3]*m[3][1] - m[0][3]*m[1][2]*m[3][1] + m[0][3]*m[1][1]*m[3][2] - m[0][1]*m[1][3]*m[3][2] - m[0][2]*m[1][1]*m[3][3] + m[0][1]*m[1][2]*m[3][3]) / d
	r[0][3] = (m[0][3]*m[1][2]*m[2][1] - m[0][2]*m[1][3]*m[2][1] - m[0][3]*m[1][1]*m[2][2] + m[0][1]*m[1][3]*m[2][2] + m[0][2]*m[1][1]*m[2][3] - m[0][1]*m[1][2]*m[2][3]) / d
	r[1][0] = (m[1][3]*m[2][2]*m[3][0] - m[1][2]*m[2][3]*m[3][0] - m[1][3]*m[2][0]*m[3][2] + m[1][0]*m[2][3]*m[3][2] + m[1][2]*m[2][0]*m[3][3] - m[1][0]*m[2][2]*m[3][3]) / d
	r[1][1] = (m[0][2]*m[2][3]*m[3][0] - m[0][3]*m[2][2]*m[3][0] + m[0][3]*m[2][0]*m[3][2] - m[0][0]*m[2][3]*m[3][2] - m[0][2]*m[2][0]*m[3][3] + m[0][0]*m[2][2]*m[3][3]) / d
	r[1][2] = (m[0][3]*m[1][2]*m[3][0] - m[0][2]*m[1][3]*m[3][0] - m[0][3]*m[1][0]*m[3][2] + m[0][0]*m[1][3]*m[3][2] + m[0][2]*m[1][0]*m[3][3] - m[0][0]*m[1][2]*m[3][3]) / d
	r[1][3] = (m[0][2]*m[1][3]*m[2][0] - m[0][3]*m[1][2]*m[2][0] + m[0][3]*m[1][0]*m[2][2] - m[0][0]*m[1][3]*m[2][2] - m[0][2]*m[1][0]*m[2][3] + m[0][0]*m[1][2]*m[2][3]) / d
	r[2][0] = (m[1][1]*m[2][3]*m[3][0] - m[1][3]*m[2][1]*m[3][0] + m[1][3]*m[2][0]*m[3][1] - m[1][0]*m[2][3]*m[3][1] - m[1][1]*m[2][0]*m[3][3] + m[1][0]*m[2][1]*m[3][3]) / d
	r[2][1] = (m[0][3]*m[2][1]*m[3][0] - m[0][1]*m[2][3]*m[3][0] - m[0][3]*m[2][0]*m[3][1] + m[0][0]*m[2][3]*m[3][1] + m[0][1]*m[2][0]*m[3][3] - m[0][0]*m[2][1]*m[3][3]) / d
	r[2][2] = (m[0][1]*m[1][3]*m[3][0] - m[0][3]*m[1][1]*m[3][0] + m[0][3]*m[1][0]*m[3][1] - m[0][0]*m[1][3]*m[3][1] - m[0][1]*m[1][0]*m[3][3] + m[0][0]*m[1][1]*m[3][3]) / d
	r[2][3] = (m[0][3]*m[1][1]*m[2][0] - m[0][1]*m[1][3]*m[2][0] - m[0][3]*m[1][0]*m[2][1] + m[0][0]*m[1][3]*m[2][1] + m[0][1]*m[1][0]*m[2][3] - m[0][0]*m[1][1]*m[2][3]) / d
	r[3][0] = (m[1][2]*m[2][1]*m[3][0] - m[1][1]*m[2][2]*m[3][0] - m[1][2]*m[2][0]*m[3][1] + m[1][0]*m[2][2]*m[3][1] + m[1][1]*m[2][0]*m[3][2] - m[1][0]*m[2][1]*m[3][2]) / d
	r[3][1] = (m[0][1]*m[2][2]*m[3][0] - m[0][2]*m[2][1]*m[3][0] + m[0][2]*m[2][0]*m[3][1] - m[0][0]*m[2][2]*m[3][1] - m[0][1]*m[2][0]*m[3][2] + m[0][0]*m[2][1]*m[3][2]) / d
	r[3][2] = (m[0][2]*m[1][1]*m[3][0] - m[0][1]*m[1][2]*m[3][0] - m[0][2]*m[1][0]*m[3][1] + m[0][0]*m[1][2]*m[3][1] + m[0][1]*m[1][0]*m[3][2] - m[0][0]*m[1][1]*m[3][2]) / d
	r[3][3] = (m[0][1]*m[1][2]*m[2][0] - m[0][2]*m[1][1]*m[2][0] + m[0][2]*m[1][0]*m[2][1] - m[0][0]*m[1][2]*m[2][1] - m[0][1]*m[1][0]*m[2][2] + m[0][0]*m[1][1]*m[2][2]) / d
	return r
}
