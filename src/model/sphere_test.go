package model

import "testing"

func TestSphereIntersect(t *testing.T) {
	for i, tt := range []struct {
		s         Sphere
		r         Ray
		want      float32
		wantTruth bool
	}{
		{
			s: Sphere{
				Center: Vector{0, 0, 0},
				Radius: 0.0,
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, 0},
			},
			want:      0.0,
			wantTruth: false,
		},
		{
			s: Sphere{
				Center: Vector{206.155, 0, 0},
				Radius: 50.0,
			},
			r: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{157.648, 150.48, 0}.Normalize(),
			},
			want:      0.0,
			wantTruth: false,
		},
		{
			s: Sphere{
				Center: Vector{1, 2, 3},
				Radius: 3.0,
			},
			r: Ray{
				Origin:    Vector{10, 10, 10},
				Direction: Vector{-1, -1, -1}.Normalize(),
			},
			want:      11.210655149486414,
			wantTruth: true,
		},
	} {
		got, found := tt.s.Intersect(tt.r)
		if !found && tt.wantTruth == false {
			continue
		}
		if (!found && tt.wantTruth == true) || (found && tt.wantTruth == false) {
			t.Errorf("%d) incorrect bool value; want %v", i, tt.wantTruth)
			continue
		}
		if !compareFloat32(got.distance, tt.want) {
			t.Errorf("%d) got %v want %v", i, got.distance, tt.want)
		}
	}
}

func TestSphereSurfaceNormal(t *testing.T) {
	for i, tt := range []struct {
		s    Sphere
		v    Vector
		want Vector
	}{
		{
			s: Sphere{
				Center: Vector{0, 0, 0},
				Radius: 0.0,
			},
			v:    Vector{0, 0, 0},
			want: Vector{0, 0, 0},
		},
		{
			s: Sphere{
				Center: Vector{0, 0, 0},
				Radius: 30.0,
			},
			v:    Vector{5, 3.14, 6.1},
			want: Vector{0.5889710504096182, 0.36987381965724025, 0.7185446814997342},
		},
	} {
		got := tt.s.SurfaceNormal(tt.v)
		if !compareVectors(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
