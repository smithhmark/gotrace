package renderer

import (
	"github.com/smithhmark/gotracer/internal/vector"
	"image"
	"math"
	"testing"
)

func almost(l, r float64) bool {
	ep := 0.000001
	if math.Abs(l-r) < ep {
		return true
	}
	return false
}

func TestCanvasToViewport(t *testing.T) {
	tests := []struct {
		x, y int
		vp   vector.SVector3
		cs   image.Rectangle
		exp  vector.SVector3
	}{
		{
			x:   3,
			y:   2,
			vp:  vector.SVector3{1, 1, 1},
			cs:  image.Rect(0, 0, 10, 10),
			exp: vector.SVector3{(-0.5 + 3.0/10.), (.5 - 2.0/10.0), 1},
		},
		{
			x:   7,
			y:   8,
			vp:  vector.SVector3{1, 1, 1},
			cs:  image.Rect(0, 0, 10, 10),
			exp: vector.SVector3{(-0.5 + 7.0/10.), (.5 - 8.0/10.0), 1},
		},
	}
	for _, test := range tests {
		rcvd := canvasToViewport(test.x, test.y, test.vp, test.cs)
		if !rcvd.Almost(test.exp) {
			t.Errorf("canvasToViewport(%v, %v, %v, %d) != %v, got: %v",
				test.x, test.y, test.vp, test.cs, test.exp, rcvd)
		}
	}
}
