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

func (sc Scene) ComputeLight(p, n, v vector.SVector3, s float64) float64 {
	atPoint := sc.Lighting.Ambient

	for _, lightSource := range sc.Lighting.Positional {
		mn, mx := lightSource.GetShadowSearchRange()
		d := lightSource.DirectionFrom(p)
		blockingObjId, _ := closestIntersection(sc, p, d, mn, mx)
		if blockingObjId < 0 {
			// if there is object blocking this ligheSource
			atPoint += lightSource.ContributeLight(p, n, v, s)
		}
	}
	return atPoint
}

type PointLight struct {
	Location  vector.SVector3
	Intensity float64
}

func lightHelperDiffuse(l, n vector.SVector3, intensity float64) float64 {
	dot := n.Dot(l)
	if dot > 0 {
		tmp := dot / (n.Mag() * l.Mag())
		return intensity * tmp
	}
	return 0
}

func lightHelperSpecular(l, n, v vector.SVector3, intensity, specular float64) float64 {
	if specular != -1 {
		r := n.Scale(2.0 * n.Dot(l)).Sub(l)
		dot := r.Dot(v)
		if dot > 0 {
			tmp := dot / (r.Mag() * v.Mag())
			tmp = math.Pow(tmp, specular)
			return intensity * tmp
		}
	}
	return 0

}

func (pl PointLight) DirectionFrom(p vector.SVector3) vector.SVector3 {
	return pl.Location.Sub(p)
}
func (pl PointLight) GetShadowSearchRange() (float64, float64) {
	return 0.0001, 1.0
}

func (pl PointLight) ContributeLight(p, n, v vector.SVector3, s float64) float64 {
	l := pl.DirectionFrom(p)
	i := pl.Intensity

	return lightHelperDiffuse(l, n, i) + lightHelperSpecular(l, n, v, i, s)
}

type DirectionalLight struct {
	Direction vector.SVector3
	Intensity float64
}

func (dl DirectionalLight) DirectionFrom(p vector.SVector3) vector.SVector3 {
	return dl.Direction
}

func (dl DirectionalLight) GetShadowSearchRange() (float64, float64) {
	return 0.0001, math.Inf(1)
}

func (dl DirectionalLight) ContributeLight(p, n, v vector.SVector3, s float64) float64 {
	l := dl.DirectionFrom(p)
	i := dl.Intensity
	return lightHelperDiffuse(l, n, i) + lightHelperSpecular(l, n, v, i, s)
}

type LightSource interface {
	ContributeLight(p, n, v vector.SVector3, s float64) float64
	GetShadowSearchRange() (float64, float64)
	DirectionFrom(p vector.SVector3) vector.SVector3
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
			center:   vector.SVector3{0, -1, 3},
			radius:   1,
			shade:    color.RGBA{255, 0, 0, 0xff},
			specular: 500})
	sc.Objects = append(sc.Objects,
		Sphere{
			center:   vector.SVector3{2, 0, 4},
			radius:   1,
			shade:    color.RGBA{0, 0, 255, 0xff},
			specular: 500})
	sc.Objects = append(sc.Objects,
		Sphere{
			center:   vector.SVector3{-2, 0, 4},
			radius:   1,
			shade:    color.RGBA{0, 255, 0, 0xff},
			specular: 10})

	return sc
}

func testSpheres2() Scene {
	sc := testSpheres()
	sc.Lighting = testLighting()
	sc.Objects = append(sc.Objects,
		Sphere{
			center:   vector.SVector3{0, -5001, 0},
			radius:   5000,
			shade:    color.RGBA{255, 255, 0, 0xff},
			specular: 1000})
	return sc
}

type Sphere struct {
	center   vector.SVector3
	radius   float64
	shade    color.Color
	specular float64
}

func (s Sphere) Shade() color.Color {
	return s.shade
}

func (s Sphere) Specular() float64 {
	return s.specular
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
	Specular() float64
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

func closestIntersection(scene Scene, o, d vector.SVector3, near, far float64) (int, float64) {
	closest_dist := math.Inf(1)
	var closest_i int = -1 // index into scene
	for idx, ible := range scene.Objects {
		t1, t2 := ible.Intersect(o, d)
		if t1 >= near && t1 < far && t1 < closest_dist {
			closest_dist = t1
			closest_i = idx
		}
		if t2 >= near && t2 < far && t2 < closest_dist {
			closest_dist = t2
			closest_i = idx
		}
	}
	return closest_i, closest_dist
}
func traceRay(scene Scene, d vector.SVector3, near, far float64) color.Color {
	closest_i, closest_dist := closestIntersection(scene, scene.Observer, d, near, far)
	if closest_i < 0 {
		return scene.Bg
	}

	obj := scene.Objects[closest_i]

	//ep := .00001
	p := scene.Observer.Add(d.Scale(closest_dist))
	n := p.Sub(obj.Center()).Norm()
	//p = p.Add(n.Scale(ep)) // epsilon offset to try reducing artifacts
	shade := obj.Shade()
	spec := obj.Specular()
	lightLevel := scene.ComputeLight(p, n, d.Scale(-1), spec)
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
