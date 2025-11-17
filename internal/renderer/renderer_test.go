package renderer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/smithhmark/gotracer/internal/vector"
)

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

func TestLightHelperDiffuse(t *testing.T) {
	//l, n vector.SVector3, intensity float64) float64 {
	tests := []struct {
		l         vector.SVector3
		n         vector.SVector3
		intensity float64
		expected  float64
	}{
		{
			vector.SVector3{0, 5, 0},
			vector.SVector3{0, 1, 0},
			1.0,
			1.0,
		},
		{
			vector.SVector3{0, 5, 0},
			vector.SVector3{0, -1, 0},
			1.0,
			0.0,
		},
		{
			vector.SVector3{0, 5, 0},
			vector.SVector3{1, 1, 0}.Norm(),
			1.0,
			math.Sqrt(2) / 2.,
		},
		{
			vector.SVector3{0, -5, 0},
			vector.SVector3{1, 1, 0}.Norm(),
			1.0,
			0.0,
		},
	}
	for _, test := range tests {
		res := lightHelperDiffuse(test.l, test.n, test.intensity)
		if !almost(res, test.expected, 0.001) {
			t.Errorf("lightHelper(%v, %v, %v) != %v, got: %v",
				test.l, test.n, test.intensity, test.expected, res)
		}
	}
}

func TestComputeLight_diffuse(t *testing.T) {
	tests := []struct {
		sc       Scene
		point    vector.SVector3
		normal   vector.SVector3
		expected float64
	}{
		{
			Scene{},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0},
			0.0,
		},
		{
			Scene{Lighting: Illumination{Ambient: 0.2}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0},
			0.2,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				PointLight{
					Location:  vector.SVector3{0, 5, 0},
					Intensity: 1.0}}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0},
			1.0,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				PointLight{
					Location:  vector.SVector3{0, 5, 0},
					Intensity: .5}}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0},
			0.5,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				PointLight{
					Location:  vector.SVector3{0, 5, 0},
					Intensity: 1.0}}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, -1, 0},
			0.0,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				PointLight{
					Location:  vector.SVector3{0, 5, 0},
					Intensity: 1.0}}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{1, 1, 0}.Norm(),
			math.Sqrt(2) / 2.0,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				DirectionalLight{
					Direction: vector.SVector3{0, 1, 0},
					Intensity: 1},
			}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0}.Norm(),
			1,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				DirectionalLight{
					Direction: vector.SVector3{0, 5, 0},
					Intensity: .50}}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{1, 1, 0}.Norm(),
			0.5 * math.Sqrt(2) / 2.0,
		},
		{
			Scene{Lighting: Illumination{Positional: []LightSource{
				DirectionalLight{
					Direction: vector.SVector3{0, 0, 0},
					Intensity: 1},
			}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{0, 1, 0}.Norm(),
			0,
		},
		{
			Scene{Lighting: Illumination{
				Positional: []LightSource{
					PointLight{
						Location:  vector.SVector3{0, 5, 0},
						Intensity: .5},
					DirectionalLight{
						Direction: vector.SVector3{0, 1, 0},
						Intensity: .5},
				}}},
			vector.SVector3{0, 0, 0},
			vector.SVector3{1, 1, 0}.Norm(),
			math.Sqrt(2) / 2.0,
		},
	}
	specular := -1.0 // magic number to signal no specular lighting
	viewer := vector.SVector3{}

	for tno, test := range tests {
		rcvd := test.sc.ComputeLight(test.point, test.normal, viewer, specular)
		if !almost(rcvd, test.expected, 0.0001) {
			t.Errorf("test: %d:ComputeLight()  != %v, got: %v",
				tno, test.expected, rcvd)
		}
	}
}

func TestScaleColor(t *testing.T) {
	tests := []struct {
		r, g, b uint8
		i       float64
	}{
		{255, 255, 255, 1.0},
		{255, 255, 255, 0.0},
		{255, 255, 255, 0.5},
		{0, 255, 0, 1.0},
		{0, 255, 0, 0.0},
		{0, 255, 0, 0.50},
	}
	for _, test := range tests {
		input := color.RGBA{R: test.r, G: test.g, B: test.b, A: 0xff}
		expected := color.RGBA{
			R: uint8(test.i * float64(test.r)),
			G: uint8(test.i * float64(test.g)),
			B: uint8(test.i * float64(test.b)),
			A: 0xff}
		fmt.Printf("running: scaleColor(%v, %v)\n", input, test.i)
		rcvd := scaleColor(input, test.i)
		if rcvd != expected {
			t.Errorf("scaleColor(%v, %v) != %v, got:%v", input, test.i, expected, rcvd)
		}
	}
}
