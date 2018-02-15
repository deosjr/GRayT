package model

import "testing"

func Test4x4Multiply(t *testing.T) {
	for i, tt := range []struct {
		m1   matrix4x4
		m2   matrix4x4
		want matrix4x4
	}{
		{
			m1: matrix4x4{
				{5, 2, 6, 1},
				{0, 6, 2, 0},
				{3, 8, 1, 4},
				{1, 8, 5, 6},
			},
			m2: matrix4x4{
				{7, 5, 8, 0},
				{1, 8, 2, 6},
				{9, 4, 3, 8},
				{5, 3, 7, 9},
			},
			want: matrix4x4{
				{96, 68, 69, 69},
				{24, 56, 18, 52},
				{58, 95, 71, 92},
				{90, 107, 81, 142},
			},
		},
	} {
		got := tt.m1.multiply(tt.m2)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func Test4x4Determinant(t *testing.T) {
	for i, tt := range []struct {
		m1   matrix4x4
		want float64
	}{
		{
			m1: matrix4x4{
				{3, 2, -1, 4},
				{2, 1, 5, 7},
				{0, 5, 2, -6},
				{-1, 2, 1, 0},
			},
			want: -418.0,
		},
		{
			m1: matrix4x4{
				{1, 0, 2, -1},
				{0, 1, 0, -1},
				{0, 0, -6, 8},
				{0, 0, 0, 5},
			},
			want: -30.0,
		},
		{
			m1: matrix4x4{
				{3, 0, 2, -1},
				{1, 2, 0, -1},
				{4, 0, 6, -3},
				{5, 0, 2, 0},
			},
			want: 20.0,
		},
	} {
		got := tt.m1.determinant()
		if got != tt.want {
			t.Errorf("%d) got %v want %f", i, got, tt.want)
		}
	}
}

func Test4x4Inverse(t *testing.T) {
	for i, tt := range []struct {
		m1   matrix4x4
		want matrix4x4
	}{
		{
			m1: matrix4x4{
				{1, 1, 1, 0},
				{0, 3, 1, 2},
				{2, 3, 1, 0},
				{1, 0, 2, 1},
			},
			want: matrix4x4{
				{-3, -0.5, 1.5, 1},
				{1, 0.25, -0.25, -0.5},
				{3, 0.25, -1.25, -0.5},
				{-3, 0, 1, 1},
			},
		},
	} {
		got := tt.m1.inverse()
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}

func Test4x4Transpose(t *testing.T) {
	for i, tt := range []struct {
		m1   matrix4x4
		want matrix4x4
	}{
		{
			m1: matrix4x4{
				{1, 1, 1, 0},
				{0, 3, 1, 2},
				{2, 3, 1, 0},
				{1, 0, 2, 1},
			},
			want: matrix4x4{
				{1, 0, 2, 1},
				{1, 3, 3, 0},
				{1, 1, 1, 2},
				{0, 2, 0, 1},
			},
		},
	} {
		got := tt.m1.transpose()
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
