package model

type Camera struct {
	Ray
	ViewDistance float64
	View         View
	Image        Image
}

func NewCamera(o, d Vector, vd float64, w, h uint) Camera {
	img := newImage(w, h)
	r := NewRay(o, d)
	return Camera{
		Ray:          r,
		ViewDistance: vd,
		View:         newView(r, vd, h, w),
		Image:        img,
	}
}

type View struct {
	ulhc    Vector // upper left hand corner
	xVector Vector // unit vector from ULHC to URHC
	yVector Vector // unit vector from ULHC to LLHC
}

func newView(r Ray, vd float64, h, w uint) View {

	n := r.Direction
	c := r.Origin.Add(n.Times(vd))

	// TODO: waarom de fuq kan ik y zo kiezen?
	// wat zijn randgevallen en wanneer gaat dit mis?
	y := Vector{0, 1, 0}
	xv := y.Cross(n).Normalize()
	yv := xv.Cross(n).Normalize()

	// if len(H)=1 then len(W)=aspectRatio
	aspectRatio := float64(w) / float64(h)

	ulhc := c.Add(xv.Times(-aspectRatio / 2.0)).Add(yv.Times(-0.5))

	return View{
		ulhc:    ulhc,
		xVector: xv,
		yVector: yv,
	}
}

// 3d translation of 2d point on view
// TODO: use midpoint of pixel instead of ULHC
// TODO: aspectRatio meenemen, circel vervormt nu
func (c Camera) PixelVector(x, y int) Vector {
	v := c.View
	xlen := float64(x) / float64(c.Image.width) * v.xVector.Length()
	ylen := float64(y) / float64(c.Image.height) * v.yVector.Length()
	xDisplacement := v.xVector.Times(xlen)
	yDisplacement := v.yVector.Times(ylen)
	return v.ulhc.Add(xDisplacement).Add(yDisplacement)
}
