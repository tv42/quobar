package truetype_test

import (
	"image"
	"testing"

	"github.com/tv42/quobar/draw/truetype"
	"github.com/tv42/quobar/internal/approve"
)

const (
	roboto     = "/usr/share/fonts/truetype/roboto/hinted/Roboto-Regular.ttf"
	liberation = "/usr/share/fonts/truetype/liberation/LiberationSerif-Regular.ttf"
)

func makeFont(t testing.TB, ttf string) *truetype.Font {
	// TODO we should probably bundle fonts in testdata
	font, err := truetype.Open(ttf)
	if err != nil {
		t.Fatalf("cannot load font: %v", err)
	}
	return font
}

func TestSimple(t *testing.T) {
	t.Parallel()
	font := makeFont(t, roboto)
	bounds := image.Rect(0, 0, 1000, 300)
	// trigger false assumptions about 0,0 origin
	bounds = bounds.Add(image.Point{X: 10000, Y: 10000})
	img := image.NewRGBA(bounds)
	if err := font.Text(img, "Hello, world"); err != nil {
		t.Fatalf("text rendering error: %v", err)
	}
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}

func TestSmall(t *testing.T) {
	t.Parallel()
	font := makeFont(t, roboto)
	bounds := image.Rect(0, 0, 100, 32)
	// trigger false assumptions about 0,0 origin
	bounds = bounds.Add(image.Point{X: 10000, Y: 10000})
	img := image.NewRGBA(bounds)
	if err := font.Text(img, "Hello, world"); err != nil {
		t.Fatalf("text rendering error: %v", err)
	}
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}

func TestKerning(t *testing.T) {
	t.Parallel()
	// Need to use a font that actually has kerning. Roboto doesn't.
	font := makeFont(t, liberation)
	bounds := image.Rect(0, 0, 1000, 300)
	// trigger false assumptions about 0,0 origin
	bounds = bounds.Add(image.Point{X: 10000, Y: 10000})
	img := image.NewRGBA(bounds)
	if err := font.Text(img, "AV"); err != nil {
		t.Fatalf("text rendering error: %v", err)
	}
	if err := approve.Image(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}
