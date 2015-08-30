package sparkline_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/tv42/quobar/draw/sparkline"
	"github.com/tv42/quobar/internal/approve"
)

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
