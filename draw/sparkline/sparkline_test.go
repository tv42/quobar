package sparkline_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/tv42/quobar/blend"
	"github.com/tv42/quobar/draw/sparkline"
	"github.com/tv42/quobar/internal/approve"
)

var (
	white  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	yellow = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
)

var trafficLights = []blend.Threshold{
	{
		Max:   1024,
		Color: white,
	},
	{
		Max:   768,
		Color: green,
	},
	{
		Max:   512,
		Color: yellow,
	},
	{
		Max:   256,
		Color: red,
	},
}

func TestSimple(t *testing.T) {
	t.Parallel()
	gray := color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	s := sparkline.New(10, gray, trafficLights)
	for _, n := range []float32{
		1000,
		800,
		751,
		454,
		300,
		599,
		564,
		901,
		1500,
		1300,
		500,
		270,
	} {
		s.Add(n)
	}
	bounds := image.Rect(0, 0, 1000, 300)
	// trigger false assumptions about 0,0 origin
	bounds = bounds.Add(image.Point{X: 10000, Y: 10000})
	img := image.NewRGBA(bounds)
	s.Draw(img)
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}

func TestSmall(t *testing.T) {
	t.Parallel()
	gray := color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	s := sparkline.New(10, gray, trafficLights)
	for _, n := range []float32{
		1000,
		800,
		751,
		454,
		300,
		599,
		564,
		901,
		1500,
		1300,
		500,
		270,
	} {
		s.Add(n)
	}
	bounds := image.Rect(0, 0, 100, 32)
	// trigger false assumptions about 0,0 origin
	bounds = bounds.Add(image.Point{X: 10000, Y: 10000})
	img := image.NewRGBA(bounds)
	s.Draw(img)
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}
