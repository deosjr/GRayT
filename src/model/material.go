package model

type Material interface {
	GetColor(Vector) Color
}

type DiffuseMaterial struct {
	Color Color
}

func (dm *DiffuseMaterial) GetColor(p Vector) Color {
	return dm.Color
}

// temporary material to play around with
type PosFuncMat struct {
	Func func(Vector) Color
}

func (pfm *PosFuncMat) GetColor(p Vector) Color {
	return pfm.Func(p)
}
