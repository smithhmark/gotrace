package renderer

import (
	//"fmt"
	"image"
	"image/color"
	"math"

	"github.com/smithhmark/gotracer/internal/vector"
)

func Render() image.Image {
	observer := vector.SVector3{}
	vp := vector.SVector3{1, 1, 1}
	rec := image.Rect(0, 0, 256, 256)
	canvas := image.NewRGBA(rec)

	for x := range rec.Dx() {
		for y := range rec.Dy() {
			dir := canvasToViewport(x, y, vp, rec)
			c := traceRay(observer, dir, vp.Z(), float32(math.Inf(1)))
			canvas.Set(x, y, c)
		}
		//fmt.Print(".")
	}

	return canvas
}

func traceRay(o, d vector.SVector3, near, far float32) color.Color {
	return color.RGBA{R: 0, G: 0, B: 200, A: 0xff}
}

func canvasToViewport(x, y int, vp vector.SVector3, rec image.Rectangle) vector.SVector3 {
	xScale := vp.X() / float32(rec.Dx())
	yScale := vp.Y() / float32(rec.Dy())
	ret := vector.SVector3{float32(x) * xScale, float32(y) * yScale, vp.Z()}
	return ret
}
