package model

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{
			u:    Vector{0, 0, 0},
			v:    Vector{0, 0, 0},
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{1, 1, 1},
			v:    Vector{0, 0, 0},
			want: Vector{1, 1, 1},
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			v:    Vector{63.7, -15, 5},
			want: Vector{105.7, -11.86, 1053.234},
		},
	} {
		got := tt.u.Add(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestSub(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{
			u:    Vector{0, 0, 0},
			v:    Vector{0, 0, 0},
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{1, 1, 1},
			v:    Vector{0, 0, 0},
			want: Vector{1, 1, 1},
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			v:    Vector{63.5, -15, 5},
			want: Vector{-21.5, 18.14, 1043.234},
		},
	} {
		got := tt.u.Sub(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestTimes(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		f    float32
		want Vector
	}{
		{
			u:    Vector{0, 0, 0},
			f:    0.0,
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{1, 1, 1},
			f:    0.0,
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{-42, 3.14, 1048.234},
			f:    3.14,
			want: Vector{-131.88, 9.8596, 3291.45476},
		},
	} {
		got := tt.u.Times(tt.f)
		if !compareVectors(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestDotProduct(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		v    Vector
		want float32
	}{
		{
			u:    Vector{0, 0, 0},
			v:    Vector{0, 0, 0},
			want: 0.0,
		},
		{
			u:    Vector{1, 1, 1},
			v:    Vector{0, 0, 0},
			want: 0.0,
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			v:    Vector{63.5, -15, 5},
			want: 7861.07,
		},
	} {
		got := tt.u.Dot(tt.v)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestLength(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		want float32
	}{
		{
			u:    Vector{0, 0, 0},
			want: 0.0,
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			want: 1049.0797769264261,
		},
	} {
		got := tt.u.Length()
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestNormalize(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		want Vector
	}{
		{
			u:    Vector{0, 0, 0},
			want: Vector{0, 0, 0}, // TODO: NaN! (division by zero), exception?
		},
		{
			u:    Vector{0, 0, 1},
			want: Vector{0, 0, 1},
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			want: Vector{0.040035086867321754, 0.0029930993515092934, 0.9991937916019084},
		},
	} {
		got := tt.u.Normalize()
		if !compareVectors(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestCrossProduct(t *testing.T) {
	for i, tt := range []struct {
		u    Vector
		v    Vector
		want Vector
	}{
		{
			u:    Vector{0, 0, 0},
			v:    Vector{0, 0, 0},
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{1, 1, 1},
			v:    Vector{0, 0, 0},
			want: Vector{0, 0, 0},
		},
		{
			u:    Vector{42, 3.14, 1048.234},
			v:    Vector{63.5, -15, 5},
			want: Vector{15739.21, 66352.859, -829.39},
		},
	} {
		got := tt.u.Cross(tt.v)
		if !compareVectors(got, tt.want) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func TestRay(t *testing.T) {
	for i, tt := range []struct {
		o    Vector
		d    Vector
		want Ray
	}{
		{
			o: Vector{0, 0, 0},
			d: Vector{0, 0, 0},
			want: Ray{
				Origin:    Vector{0, 0, 0},
				Direction: Vector{0, 0, 0},
			},
		},
		{
			o: Vector{1, 1, 1},
			d: Vector{1, 1, 1},
			want: Ray{
				Origin:    Vector{1, 1, 1},
				Direction: Vector{0.5773502691896258, 0.5773502691896258, 0.5773502691896258},
			},
		},
		{
			o: Vector{63.5, -15, 5},
			d: Vector{42, 3.14, 1048.234},
			want: Ray{
				Origin:    Vector{63.5, -15, 5},
				Direction: Vector{0.040035086867321754, 0.0029930993515092934, 0.9991937916019084},
			},
		},
	} {
		got := NewRay(tt.o, tt.d)
		if !compareVectors(got.Origin, tt.want.Origin) || !compareVectors(got.Direction, tt.want.Direction) {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

// floating point precision
func compareVectors(u, v Vector) bool {
	for _, d := range Dimensions {
		if !compareFloat32(u.Get(d), v.Get(d)) {
			return false
		}
	}
	return true
}

func compareFloat32(a, b float32) bool {
	if math.Abs(float64(a-b)) > 0.001 {
		return false
	}
	return true
}
