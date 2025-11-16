package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/smithhmark/gotracer/internal/vector"
)

func background() color.Color {
	return color.RGBA{R: 255, G: 255, B: 255, A: 0xff} // white
}

type Illumination struct {
	Ambient    float64
	Positional []LightSource
}

func (i Illumination) ComputeLight(p, n vector.SVector3) float64 {
	atPoint := i.Ambient

	for _, lightSource := range i.Positional {
		atPoint += lightSource.ContributeLight(p, n)
	}
	return atPoint
}

type PointLight struct {
	Location  vector.SVector3
	Intensity float64
}

func lightHelper(l, n vector.SVector3, intensity float64) float64 {
	dot := n.Dot(l)
	if dot > 0 {
		tmp := dot / (n.Mag() * l.Mag())
		return intensity * tmp
	}
	return 0
}

func (pl PointLight) ContributeLight(p, n vector.SVector3) float64 {
	l := pl.Location.Sub(p)
	return lightHelper(l, n, pl.Intensity)
}

type DirectionalLight struct {
	Direction vector.SVector3
	Intensity float64
}

func (dl DirectionalLight) ContributeLight(p, n vector.SVector3) float64 {
	l := dl.Direction
	return lightHelper(l, n, dl.Intensity)
}

type LightSource interface {
	ContributeLight(p, n vector.SVector3) float64
}

type Scene struct {
	Lighting Illumination
	Objects  []Intersectable
	Bg       color.Color
	Observer vector.SVector3
}

func testLighting() Illumination {
	var ill Illumination
	ill.Ambient = 0.2
	ill.Positional = append(ill.Positional,
		PointLight{
			Location:  vector.SVector3{2, 1, 0},
			Intensity: .6})
	ill.Positional = append(ill.Positional,
		DirectionalLight{
			Direction: vector.SVector3{1, 4, 4},
			Intensity: .2})
	return ill
}

func testSpheres() Scene {
	var sc Scene
	sc.Lighting = Illumination{Ambient: 1.0}

	sc.Bg = background()
	sc.Objects = append(sc.Objects,
		Sphere{
			center: vector.SVector3{0, -1, 3},
			radius: 1,
			shade:  color.RGBA{255, 0, 0, 0xff}})
	sc.Objects = append(sc.Objects,
		Sphere{
			center: vector.SVector3{2, 0, 4},
			radius: 1,
			shade:  color.RGBA{0, 0, 255, 0xff}})
	sc.Objects = append(sc.Objects,
		Sphere{
			center: vector.SVector3{-2, 0, 4},
			radius: 1,
			shade:  color.RGBA{0, 255, 0, 0xff}})

	return sc
}

func testSpheres2() Scene {
	sc := testSpheres()
	sc.Lighting = testLighting()
	sc.Objects = append(sc.Objects,
		Sphere{
			center: vector.SVector3{0, -5001, 0},
			radius: 5000,
			shade:  color.RGBA{255, 255, 0, 0xff}})
	return sc
}

type Sphere struct {
	center vector.SVector3
	radius float64
	shade  color.Color
}

func (s Sphere) Shade() color.Color {
	return s.shade
}

func (s Sphere) Center() vector.SVector3 {
	return s.center
}

func (s Sphere) Normal(o, d vector.SVector3) vector.SVector3 {
	return vector.SVector3{}
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
	Center() vector.SVector3
	Normal(o, d vector.SVector3) vector.SVector3
}

func Render() image.Image {
	sc := testSpheres2()
	vp := vector.SVector3{1, 1, 1} // viewport TODO:move into Scene
	//rec := image.Rect(0, 0, 256, 256)
	rec := image.Rect(0, 0, 1024, 1024)
	canvas := image.NewRGBA(rec)

	for x := range rec.Dx() {
		for y := range rec.Dy() {
			dir := canvasToViewport(x, y, vp, rec)
			c := traceRay(sc, dir, vp.Z(), math.Inf(1))
			canvas.Set(x, y, c)
		}
	}

	return canvas
}

func scaleColor(shade color.Color, intensity float64) color.Color {
	if intensity >= 1.0 {
		return shade
	}
	if intensity <= 0 {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0xff}
	}
	r, g, b, _ := shade.RGBA()
	scale := float64(0xffff)
	var nr float64 = intensity * float64(r)
	var ng float64 = intensity * float64(g)
	var nb float64 = intensity * float64(b)
	tmpC := color.RGBA{
		R: uint8(nr / scale * 255),
		G: uint8(ng / scale * 255),
		B: uint8(nb / scale * 255),
		A: 0xff}
	return tmpC
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

	obj := scene.Objects[closest_i]

	//ep := .00001
	p := scene.Observer.Add(d.Scale(closest_dist))
	n := p.Sub(obj.Center()).Norm()
	//p = p.Add(n.Scale(ep)) // epsilon offset to try reducing artifacts
	shade := obj.Shade()
	lightLevel := scene.Lighting.ComputeLight(p, n)
	return scaleColor(shade, lightLevel)
	//return color.RGBA{R: 0, G: 0, B: 200, A: 0xff}
}

func almost(x, y, ep float64) bool {
	if math.Abs(x-y) < ep {
		return true
	}
	return false
}

func canvasToViewport(x, y int, vp vector.SVector3, rec image.Rectangle) vector.SVector3 {
	leftEdge := -1 * vp.X() / 2
	topEdge := vp.Y() / 2
	newx := leftEdge + float64(x)/float64(rec.Dx())*vp.X()
	newy := topEdge - float64(y)/float64(rec.Dy())*vp.Y()
	ret := vector.SVector3{newx, newy, vp.Z()}
	return ret
}
