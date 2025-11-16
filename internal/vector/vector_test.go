package vector

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		left   SVector3
		right  SVector3
		result SVector3
	}{
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{0, 0, 0},
			result: SVector3{0, 0, 0},
		},
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{1, 1, 1},
			result: SVector3{1, 1, 1},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{0, 0, 0},
			result: SVector3{1, 1, 1},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{1, 1, 1},
			result: SVector3{2, 2, 2},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{1, 2, 4},
			result: SVector3{2, 3, 5},
		},
	}
	for _, test := range tests {
		l := test.left
		r := test.right
		rcvd := l.Add(r)
		if rcvd != test.result {
			t.Errorf("%v.Add(%v) != %v, got: %v", l, r, test.result, rcvd)

		}
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		left   SVector3
		right  SVector3
		result SVector3
	}{
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{0, 0, 0},
			result: SVector3{0, 0, 0},
		},
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{1, 1, 1},
			result: SVector3{-1, -1, -1},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{0, 0, 0},
			result: SVector3{1, 1, 1},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{1, 1, 1},
			result: SVector3{0, 0, 0},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  SVector3{1, 2, 4},
			result: SVector3{0, -1, -3},
		},
	}
	for _, test := range tests {
		l := test.left
		r := test.right
		rcvd := l.Sub(r)
		if rcvd != test.result {
			t.Errorf("%v.Sub(%v) != %v, got: %v", l, r, test.result, rcvd)

		}
	}
}

func TestScale(t *testing.T) {
	tests := []struct {
		left   SVector3
		right  float64
		result SVector3
	}{
		{
			left:   SVector3{0, 0, 0},
			right:  100.0,
			result: SVector3{0, 0, 0},
		},
		{
			left:   SVector3{1, 1, 1},
			right:  1.0,
			result: SVector3{1, 1, 1},
		},
		{
			left:   SVector3{1, 2, 3},
			right:  10.0,
			result: SVector3{10, 20, 30},
		},
	}
	for _, test := range tests {
		l := test.left
		r := test.right
		rcvd := l.Scale(r)
		if rcvd != test.result {
			t.Errorf("%v.Scale(%v) != %v, got: %v", l, r, test.result, rcvd)

		}
	}
}

func TestMag(t *testing.T) {
	tests := []struct {
		left   SVector3
		result float64
	}{
		{
			left:   SVector3{0, 0, 0},
			result: 0.0,
		},
		{
			left:   SVector3{3, 4, 0},
			result: 5.0,
		},
		{
			left:   SVector3{0, 5, 12},
			result: 13.0,
		},
	}
	for _, test := range tests {
		l := test.left
		rcvd := l.Mag()
		if rcvd != test.result {
			t.Errorf("%v.Mag() != %v, got: %v", l, test.result, rcvd)

		}
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		left   SVector3
		right  SVector3
		result float64
	}{
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{0, 0, 0},
			result: 0.0,
		},
		{
			left:   SVector3{0, 1, 0},
			right:  SVector3{0, 0, 1},
			result: 0.0,
		},
		{
			left:   SVector3{1, 0, 0},
			right:  SVector3{0, 1, 0},
			result: 0.0,
		},
	}
	for _, test := range tests {
		l := test.left
		r := test.right
		rcvd := l.Dot(r)
		if rcvd != test.result {
			t.Errorf("%v.Dot(%v) != %v, got: %v", l, r, test.result, rcvd)

		}
	}
}

func TestCross(t *testing.T) {
	tests := []struct {
		left   SVector3
		right  SVector3
		result SVector3
	}{
		{
			left:   SVector3{0, 0, 0},
			right:  SVector3{0, 0, 0},
			result: SVector3{0, 0, 0},
		},
		{ // example take from: https://www.cuemath.com/geometry/cross-product/
			left:   SVector3{3, 4, 5},
			right:  SVector3{7, 8, 9},
			result: SVector3{-4, 8, -4},
		},
		{ // example taken from:https://mathinsight.org/cross_product_examples
			left:   SVector3{3, -3, 1},
			right:  SVector3{4, 9, 2},
			result: SVector3{-15, -2, 39},
		},
	}
	for _, test := range tests {
		l := test.left
		r := test.right
		rcvd := l.Cross(r)
		if rcvd != test.result {
			t.Errorf("%v.Cross(%v) != %v, got: %v", l, r, test.result, rcvd)

		}
	}
}

func TestNorm(t *testing.T) {
	eqVecNorm := float64(0.57735026)
	tests := []struct {
		left   SVector3
		result SVector3
	}{
		{
			left:   SVector3{0, 0, 0},
			result: SVector3{0, 0, 0},
		},
		{
			left:   SVector3{1, 1, 1},
			result: SVector3{eqVecNorm, eqVecNorm, eqVecNorm},
		},
		{
			left:   SVector3{.1, .1, .1},
			result: SVector3{eqVecNorm, eqVecNorm, eqVecNorm},
		},
		{
			left:   SVector3{10, 10, 10},
			result: SVector3{eqVecNorm, eqVecNorm, eqVecNorm},
		},
		{
			left:   SVector3{10, 0, 0},
			result: SVector3{1, 0, 0},
		},
		{
			left:   SVector3{0, 10, 0},
			result: SVector3{0, 1, 0},
		},
		{
			left:   SVector3{0, 0, 10},
			result: SVector3{0, 0, 1},
		},
	}
	for _, test := range tests {
		l := test.left
		rcvd := l.Norm()
		if rcvd != test.result && !rcvd.Almost(test.result) {
			t.Errorf("%v.Norm() != %v, got: %v", l, test.result, rcvd)

		}
	}
}
