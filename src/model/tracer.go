package model

import (
	"math"
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

type whittedRayTracer struct{}

func NewWhittedRayTracer() Tracer {
	return whittedRayTracer{}
}

func (wrt whittedRayTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	si, ok := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BACKGROUND_COLOR
	}

	si.as = scene.AccelerationStructure
	si.depth = depth
	si.tracer = wrt

	color := NewColor(0, 0, 0)
	material := si.object.GetMaterial()
	for _, light := range scene.Lights {
		if pointInShadow(light, si.point, si.as) {
			continue
		}
		facingRatio := si.normal.Dot(si.incident.Times(-1))
		if facingRatio <= 0 {
			continue
		}

		var objectColor Color
		switch mat := material.(type) {
		case *RadiantMaterial:
			objectColor = mat.Color
		case *DiffuseMaterial:
			lightSegment := light.GetLightSegment(si.point)
			lightRatio := si.normal.Dot(lightSegment.Normalize())
			factors := standardAlbedo / math.Pi * light.Intensity(lightSegment.Length()) * lightRatio
			lightColor := light.Color().Times(factors)
			objectColor = mat.Color.Product(lightColor)
		case *ReflectiveMaterial:
			i := si.incident
			n := si.object.SurfaceNormal(si.point)
			reflection := i.Sub(n.Times(2 * i.Dot(n)))
			newRay := NewRay(si.point, reflection)
			// TODO: retain maxdistance for tracing
			objectColor = wrt.GetRayColor(newRay, scene, depth+1) //.Times(1 - standardAlbedo) // simulates nonperfect reflection
		case *PosFuncMat:
			objectColor = mat.GetColor(si)
		}

		//objectColor := si.object.GetColor(si, light)
		color = color.Add(objectColor.Times(facingRatio))
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
	random *rand.Rand
}

func NewPathTracer() Tracer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &pathTracer{random: r}
}

func (pt *pathTracer) GetRayColor(ray Ray, scene *Scene, depth int) Color {
	if depth == MAX_RAY_DEPTH {
		return BLACK
	}

	si, ok := scene.AccelerationStructure.ClosestIntersection(ray, MAX_RAY_DISTANCE)
	if !ok {
		return BLACK
	}

	si.as = scene.AccelerationStructure
	si.depth = depth
	si.tracer = pt

	surfaceDiffuseColor := si.object.GetColor(si)
	if si.object.IsLight() {
		return surfaceDiffuseColor
	}

	// random new ray
	randomDirection := randomInHemisphere(pt.random, si.normal)
	newRay := NewRay(si.point, randomDirection)
	cos := si.normal.Dot(randomDirection)
	recursiveColor := pt.GetRayColor(newRay, scene, depth+1)
	brdf := surfaceDiffuseColor.Times(1.0 / math.Pi)
	pdf := 1.0 / (2.0 * math.Pi)
	sampleColor := recursiveColor.Times(cos / pdf).Product(brdf)
	return sampleColor
}

// this is actually slower than the very naive method before..
func randomInHemisphere(random *rand.Rand, normal Vector) Vector {
	// uniform hemisphere sampling: pbrt 774
	// samples from hemisphere with z-axis = up direction
	z := random.Float64()
	det := 1 - z*z
	r := 0.0
	if det > 0 {
		r = math.Sqrt(det)
	}
	phi := 2 * math.Pi * random.Float64()
	v := Vector{r * math.Cos(phi), r * math.Sin(phi), z}

	ez := Vector{0, 0, 1}
	rotationVector := ez.Cross(normal)
	theta := math.Acos(ez.Dot(normal))
	return Rotate(theta, rotationVector).Vector(v)
}
