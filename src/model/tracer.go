package model

import (
	"math/rand"
	"time"
)

const (
	standardAlbedo   = 0.18
	MAX_RAY_DISTANCE = 1000000.0
	MAX_RAY_DEPTH    = 5
)

var BACKGROUND_COLOR Color
var BLACK = NewColor(0, 0, 0)

type Tracer interface {
	GetRayColor(Ray, *Scene, int) Color
}

type naiveRayTracer struct{}

func NewNaiveRayTracer() Tracer {
	return &naiveRayTracer{}
}

func (nrt *naiveRayTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	as := scene.AccelerationStructure

	hit, ok := as.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BACKGROUND_COLOR
	}

	point := PointFromRay(ray, hit.distance)
	si := &SurfaceInteraction{
		Point:  point,
		Normal: hit.normal,
		Object: hit.object,
		AS:     as,
		// already normalized
		Incident: ray.Direction,
		depth:    depth,
		tracer:   nrt,
	}

	color := NewColor(0, 0, 0)
	for _, light := range scene.Lights {
		if pointInShadow(light, point, as) {
			continue
		}
		facingRatio := hit.normal.Dot(ray.Direction.Times(-1))
		if facingRatio <= 0 {
			continue
		}

		objectColor := hit.object.GetColor(si, light)
		color = color.Add(objectColor)
	}
	return color
}

func pointInShadow(light Light, point Vector, as AccelerationStructure) bool {
	lightSegment := light.GetLightSegment(point)
	shadowRay := NewRay(point, lightSegment)
	maxDistance := lightSegment.Length()
	if _, ok := as.ClosestIntersection(shadowRay, maxDistance); ok {
		return true
	}
	return false
}

type pathTracer struct {
	random  *rand.Rand
	samples int
}

func NewPathTracer(n int) *pathTracer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &pathTracer{random: r, samples: n}
}

func (pt *pathTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	as := scene.AccelerationStructure

	hit, ok := as.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BLACK
	}

	point := PointFromRay(ray, hit.distance)
	si := &SurfaceInteraction{
		Point:  point,
		Normal: hit.normal,
		Object: hit.object,
		AS:     as,
		// already normalized
		Incident: ray.Direction,
		depth:    depth,
		tracer:   pt,
	}
	o := hit.object.(Triangle)

	emittedColor := NewColor(0, 0, 0)
	surfaceDiffuseColor := NewColor(0, 0, 0)
	if rad, ok := o.Material.(*RadiantMaterial); ok {
		emittedColor = rad.Color
		surfaceDiffuseColor = rad.Color
	}
	if diff, ok := o.Material.(*DiffuseMaterial); ok {
		facingRatio := si.Normal.Dot(si.Incident.Times(-1))
		surfaceDiffuseColor = diff.Color.Times(facingRatio)
	}
	if debug, ok := o.Material.(*PosFuncMat); ok {
		surfaceDiffuseColor = debug.GetColor(si, nil)
	}

	sumSampleColor := NewColor(0, 0, 0)
	for i := 0; i < pt.samples; i++ {
		// random new ray
		randomDirection := pt.randomInHemisphere(hit.normal)
		newRay := NewRay(point, randomDirection)
		costheta := hit.normal.Dot(randomDirection)
		recursiveColor := pt.GetRayColor(newRay, scene, depth+1)
		sampleColor := emittedColor.Add(recursiveColor.Times(2).Times(costheta).Product(surfaceDiffuseColor))
		sumSampleColor = sumSampleColor.Add(sampleColor)
	}
	return sumSampleColor.Times(1.0 / float64(pt.samples))
}

func (pt *pathTracer) randomInHemisphere(normal Vector) Vector {
	// uniform hemisphere sampling: pbrt 774
	// TODO: rotate the sample from ey to normal?
	/*
		z := random.Float64()
		det := 1 - z*z
		r := 0
		if det > 0 {
			r = math.Sqrt(det)
		}
		phi := 2 * math.Pi * random.Float64()
		v := Vector{r * math.Cos(phi), r * math.Sin(phi), z}
	*/

	// this is slow and dumb..
	for {
		randomVector := Vector{pt.random.Float64() - 0.5, pt.random.Float64() - 0.5, pt.random.Float64() - 0.5}.Normalize()
		if normal.Dot(randomVector) <= 0 {
			continue
		}
		return randomVector
	}
}
