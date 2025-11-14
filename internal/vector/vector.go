package vector

import (
	"math"
)

type SVector3 [3]float32

func almostEqual(a, b, epsilon float32) bool {
	dif := float64(a - b)
	return math.Abs(dif) < float64(epsilon)
}

func (l SVector3) Almost(r SVector3) bool {
	ep := float32(0.000001)
	return almostEqual(l[0], r[0], ep) && almostEqual(l[1], r[1], ep) && almostEqual(l[2], r[2], ep)
}

func (l SVector3) Add(r SVector3) SVector3 {
	var res = SVector3{l[0] + r[0], l[1] + r[1], l[2] + r[2]}
	return res
}

func (l SVector3) Sub(r SVector3) SVector3 {
	var res = SVector3{l[0] - r[0], l[1] - r[1], l[2] - r[2]}
	return res
}

func (l SVector3) Scale(k float32) SVector3 {
	var res = SVector3{l[0] * k, l[1] * k, l[2] * k}
	return res
}

func (l SVector3) Dot(r SVector3) float32 {
	var res float32
	res = l[0]*r[0] + l[1]*r[1] + l[2]*r[2]
	return res
}

func (l SVector3) Mag() float32 {
	var res float64
	res = float64(l[0])*float64(l[0]) + float64(l[1])*float64(l[1]) + float64(l[2])*float64(l[2])
	return float32(math.Sqrt(res))
}

func (l SVector3) Cross(r SVector3) SVector3 {
	x := l[1]*r[2] - l[2]*r[1]
	y := l[0]*r[2] - l[2]*r[0]
	z := l[0]*r[1] - l[1]*r[0]
	return SVector3{x, -1 * y, z}
}

func (l SVector3) Norm() SVector3 {
	mag := l.Mag()
	if mag == 0.0 {
		return SVector3{0, 0, 0}
	}
	mag = 1 / mag
	return l.Scale(mag)
}
