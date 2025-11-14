package vector

type SVector3 [3]float32

func (l SVector3) Add(r SVector3) SVector3 {
	var res = SVector3{l[0] + r[0], l[1] + r[1], l[2] + r[2]}
	//res := SVector3{l[0] + r[0], l[1] + r[1], l[2] + r[2]}
	return res
}

func (l SVector3) Sub(r SVector3) SVector3 {
	var res = SVector3{l[0] - r[0], l[1] - r[1], l[2] - r[2]}
	//res := SVector3{l[0] + r[0], l[1] + r[1], l[2] + r[2]}
	return res
}
