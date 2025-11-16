package renderer

import (
	//"fmt"
	"image"
	"image/color"
	"math"

	"github.com/smithhmark/gotracer/internal/vector"
)

type Scene struct {
	Objects  []Intersectable
	Bg       color.Color
	Observer vector.SVector3
}

func testSpheres() Scene {
	var sc Scene
	sc.Bg = background()
	sc.Objects = append(sc.Objects, Sphere{center: vector.SVector3{0, -1, 3}, radius: 1, shade: color.RGBA{255, 0, 0, 0xff}})
	sc.Objects = append(sc.Objects, Sphere{center: vector.SVector3{2, 0, 4}, radius: 1, shade: color.RGBA{0, 0, 255, 0xff}})
	sc.Objects = append(sc.Objects, Sphere{center: vector.SVector3{-2, 0, 4}, radius: 1, shade: color.RGBA{0, 255, 0, 0xff}})

	return sc
}

type Sphere struct {
	center vector.SVector3
	radius float64
	shade  color.Color
}

func background() color.Color {
	return color.RGBA{R: 255, G: 255, B: 255, A: 0xff} // white
}

func (s Sphere) Shade() color.Color {
	return s.shade
}

func (s Sphere) Intersect(o, d vector.SVector3) (float64, float64) {
	co := o.Sub(s.center)
	a := d.Dot(d)
	b := 2 * co.Dot(d)
	c := co.Dot(co) - s.radius*s.radius

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return math.Inf(1), math.Inf(1)
	}

	t1 := (-1*b + math.Sqrt(float64(discriminant))) / (2 * a)
	t2 := (-1*b - math.Sqrt(float64(discriminant))) / (2 * a)
	return t1, t2

}

type Intersectable interface {
	Intersect(o, d vector.SVector3) (float64, float64)
	Shade() color.Color
}

func Render() image.Image {
	sc := testSpheres()
	vp := vector.SVector3{1, 1, 1} // viewport TODO:move into Scene
	rec := image.Rect(0, 0, 256, 256)
	canvas := image.NewRGBA(rec)

	for x := range rec.Dx() {
		for y := range rec.Dy() {
			dir := canvasToViewport(x, y, vp, rec)
			c := traceRay(sc, dir, vp.Z(), math.Inf(1))
			canvas.Set(x, y, c)
		}
		//fmt.Print(".")
	}

	return canvas
}

func traceRay(scene Scene, d vector.SVector3, near, far float64) color.Color {
	closest_dist := math.Inf(1)
	var closest_i int = -1 // index into scene

	for idx, ible := range scene.Objects {
		t1, t2 := ible.Intersect(scene.Observer, d)
		if t1 >= near && t1 < far && t1 < closest_dist {
			closest_dist = t1
			closest_i = idx
		}
		if t2 >= near && t2 < far && t2 < closest_dist {
			closest_dist = t2
			closest_i = idx
		}
	}
	if closest_i < 0 {
		return scene.Bg
	}
	return scene.Objects[closest_i].Shade()
	//return color.RGBA{R: 0, G: 0, B: 200, A: 0xff}
}

func canvasToViewport(x, y int, vp vector.SVector3, rec image.Rectangle) vector.SVector3 {
	leftEdge := -1 * vp.X() / 2
	topEdge := vp.Y() / 2
	newx := leftEdge + float64(x)/float64(rec.Dx())*vp.X()
	newy := topEdge - float64(y)/float64(rec.Dy())*vp.Y()
	ret := vector.SVector3{newx, newy, vp.Z()}
	return ret
}
