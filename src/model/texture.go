package model

import (
	"image"
	"math"
)

type Texture interface {
	GetColor(si *SurfaceInteraction) Color
}

type texture struct {
	// TODO: 2d uv struct, calculating uv on si?
	uvFunc func(si *SurfaceInteraction) Vector
	// (u,v) to (s,t) or possibly (x,y,z) coords
	mappingFunc func(u, v float32) Vector
	colorFunc   func(textureSpace Vector) Color
}

// by default, a texture is a mapping from
// (u,v) coords -> texture space -> color
func (t texture) GetColor(si *SurfaceInteraction) Color {
	uv := t.uvFunc(si)
	st := t.mappingFunc(uv.X, uv.Y)
	return t.colorFunc(st)
}

type ConstantTexture struct {
	Color Color
}

func NewConstantTexture(color Color) ConstantTexture {
	return ConstantTexture{
		Color: color,
	}
}

// override for efficiency but essentially a texture
// whose mappingFunc is identity and colorFunc returns a constant
func (ct ConstantTexture) GetColor(_ *SurfaceInteraction) Color {
	return ct.Color
}

type UVTexture struct {
	texture
}

func NewUVTexture(uvFunc func(*SurfaceInteraction) Vector) UVTexture {
	return UVTexture{
		texture: texture{
			uvFunc: uvFunc,
			mappingFunc: func(u, v float32) Vector {
				return Vector{u, v, 0}
			},
			colorFunc: func(st Vector) Color {
				return Color{st.X, st.Y, 0}
			},
		},
	}
}

func TriangleMeshUVFunc(si *SurfaceInteraction) Vector {
	tr := si.GetObject().(TriangleInMesh)
	p := si.UntransformedPoint
	l0, l1, l2 := tr.Barycentric(p)
	p0, p1, p2 := tr.PointIndices()
	uv0 := tr.Mesh.UV[p0]
	uv1 := tr.Mesh.UV[p1]
	uv2 := tr.Mesh.UV[p2]
	uv := uv0.Times(l0).Add(uv1.Times(l1)).Add(uv2.Times(l2))
	return uv
}

type ImageTexture struct {
	texture
	img image.Image
}

func NewImageTexture(img image.Image, uvFunc func(*SurfaceInteraction) Vector) ImageTexture {
	b := img.Bounds()
	minx := float32(b.Min.X)
	miny := float32(b.Min.Y)
	w := float32(b.Max.X) - minx
	h := float32(b.Max.Y) - miny
	return ImageTexture{
		texture: texture{
			uvFunc: uvFunc,
			// golang image bounds do not necessarily contain 0,0
			// Min.X, Min.Y is top left
			// Max.X-1, Max.Y-1 is bottom right
			// whilst u,v normally has 0,0 as bottom left
			// and 1,1 as top right ?
			mappingFunc: func(u, v float32) Vector {
				v = 1 - v // invert y axis
				return Vector{
					X: minx + w*u,
					Y: miny + h*v,
					Z: 0,
				}
			},
			colorFunc: func(st Vector) Color {
				r, g, b, _ := img.At(int(st.X), int(st.Y)).RGBA()
				r256 := (float32(r) / 65535)
				g256 := (float32(g) / 65535)
				b256 := (float32(b) / 65535)
				return Color{r256, g256, b256}
			},
		},
		img: img,
	}
}

type CheckerboardTexture struct {
	texture
	frequency int
}

// frequency determines mapping from (u,v) [0,1] space to (s,t) [0,2*f] space
// which makes the actual color function trivial
func NewCheckerboardTexture(f int, uvFunc func(*SurfaceInteraction) Vector) CheckerboardTexture {
	return CheckerboardTexture{
		texture: texture{
			uvFunc: uvFunc,
			mappingFunc: func(u, v float32) Vector {
				return Vector{2 * float32(f) * u, 2 * float32(f) * v, 0}
			},
			colorFunc: func(st Vector) Color {
				if (int(math.Floor(float64(st.X)))+int(math.Floor(float64(st.Y))))%2 == 0 {
					return NewColor(255, 255, 255)
				}
				return NewColor(0, 0, 0)
			},
		},
		frequency: f,
	}
}
