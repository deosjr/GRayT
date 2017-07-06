package model

type Color struct {
	r, g, b float64
}

func NewColor(r, g, b uint8) Color {
	return Color{
		r: float64(r) / 255,
		g: float64(g) / 255,
		b: float64(b) / 255,
	}
}

func (c Color) Add(d Color) Color {
	return Color{
		r: c.r + d.r,
		g: c.g + d.g,
		b: c.b + d.b,
	}
}

func (c Color) Times(f float64) Color {
	return Color{
		r: f * c.r,
		g: f * c.g,
		b: f * c.b,
	}
}

func (c Color) Product(d Color) Color {
	return Color{
		r: c.r * d.r,
		g: c.g * d.g,
		b: c.b * d.b,
	}
}

func float64touint8(f float64) uint8 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 255
	}
	return uint8(f * 255)
}

// RGB vector in [0,1] -> to [0,255]
func (c Color) R() uint8 {
	return float64touint8(c.r)
}

func (c Color) G() uint8 {
	return float64touint8(c.g)
}

func (c Color) B() uint8 {
	return float64touint8(c.b)
}
