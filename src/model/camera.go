package model

import (
	"math"
)

type Camera interface {
	PixelRay(x, y float32) Ray
	Width() int
	Height() int
	LookAt(from, to, up Vector)
}

type projectiveCamera struct {
	cameraToWorld  Transform
	cameraToScreen Transform
	rasterToCamera Transform
	screenToRaster Transform
	rasterToScreen Transform
	w, h           uint
}

func (c *projectiveCamera) Width() int {
	return int(c.w)
}
func (c *projectiveCamera) Height() int {
	return int(c.h)
}

func (c *projectiveCamera) LookAt(from, to, up Vector) {
	dir := VectorFromTo(from, to).Normalize()
	left := dir.Normalize().Cross(up).Normalize()
	newUp := left.Cross(dir)
	cameraToWorld := matrix4x4{
		{left.X, newUp.X, dir.X, from.X},
		{left.Y, newUp.Y, dir.Y, from.Y},
		{left.Z, newUp.Z, dir.Z, from.Z},
		{0, 0, 0, 1},
	}

	// TODO: pbrt returns the inverse: why?
	// it has to do with coordinate system handidness I think
	// if I invert my camera movements are backwards...
	c.cameraToWorld = Transform{
		m:    cameraToWorld,
		mInv: cameraToWorld.inverse(),
	}
}

func (c *projectiveCamera) cameraTransforms(w, h uint) {
	c.w, c.h = w, h
	aspectRatio := float32(w) / float32(h)

	// assumption: screenWindow is (-1,-1) to (1,1)
	// scale to aspectRatio
	var pMinX, pMaxX, pMinY, pMaxY float32
	if aspectRatio > 1.0 {
		pMinX = -aspectRatio
		pMaxX = aspectRatio
		pMinY = -1.0
		pMaxY = 1.0
	} else {
		pMinX = -1.0
		pMaxX = 1.0
		pMinY = -1.0 / aspectRatio
		pMaxY = 1.0 / aspectRatio
	}

	ulhc := Translate(Vector{-pMinX, -pMaxY, 0})
	ndc := Scale(1.0/(pMaxX-pMinX), 1.0/(pMinY-pMaxY), 1)
	resolution := Scale(float32(w), float32(h), 1)
	c.screenToRaster = resolution.Mul(ndc.Mul(ulhc))
	c.rasterToScreen = c.screenToRaster.Inverse()
	c.rasterToCamera = c.cameraToScreen.Inverse().Mul(c.rasterToScreen)
}

type OrthographicCamera struct {
	projectiveCamera
}

func orthographic(zNear, zFar float32) Transform {
	return Scale(1, 1, 1.0/(zFar-zNear)).Mul(Translate(Vector{0, 0, -zNear}))
}

func (c *OrthographicCamera) PixelRay(x, y float32) Ray {
	pCamera := c.rasterToCamera.Point(Vector{x, y, 0})
	r := NewRay(pCamera, Vector{0, 0, 1})
	return c.cameraToWorld.Ray(r)
}

func NewOrthographicCamera(w, h uint) *OrthographicCamera {
	c := &OrthographicCamera{
		projectiveCamera: projectiveCamera{
			cameraToScreen: orthographic(0, 1),
		},
	}
	c.cameraTransforms(w, h)
	return c
}

type PerspectiveCamera struct {
	projectiveCamera
}

func perspective(fov, n, f float32) Transform {
	persp := matrix4x4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, f / (f - n), -f * n / (f - n)},
		{0, 0, 1, 0},
	}
	invTanAng := 1.0 / float32(math.Tan(float64(fov/2.0)))
	return Scale(invTanAng, invTanAng, 1).Mul(NewTransform(persp))
}

func (c *PerspectiveCamera) PixelRay(x, y float32) Ray {
	pCamera := c.rasterToCamera.Point(Vector{x, y, 0})
	r := NewRay(Vector{0, 0, 0}, pCamera)
	return c.cameraToWorld.Ray(r)
}

func NewPerspectiveCamera(w, h uint, fov float32) *PerspectiveCamera {
	c := &PerspectiveCamera{
		projectiveCamera: projectiveCamera{
			cameraToScreen: perspective(fov, 1e-2, 1000.0),
		},
	}
	c.cameraTransforms(w, h)
	return c
}
