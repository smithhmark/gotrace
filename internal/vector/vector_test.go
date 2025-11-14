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
