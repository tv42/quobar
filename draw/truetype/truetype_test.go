package truetype_test

import (
	"errors"
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/tv42/quobar/draw/truetype"
)

const (
	roboto     = "/usr/share/fonts/truetype/roboto/Roboto-Regular.ttf"
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

func colorEq(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

func approve(img image.Image) error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("unknown caller")
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Errorf("caller address not known: %v", pc)
	}
	name := fn.Name()
	idx := strings.LastIndex(name, ".")
	if idx == -1 {
		return fmt.Errorf("cannot determine package: %q", name)
	}
	importPath := name[:idx]
	importPath = strings.TrimSuffix(importPath, "_test")
	testName := name[idx+1:]
	pkg, err := build.Import(importPath, ".", build.FindOnly)
	if err != nil {
		return err
	}
	dir := path.Join(pkg.Dir, "testdata")

	newName := path.Join(dir, testName+".new.png")
	newF, err := os.Create(newName)
	if err != nil {
		return fmt.Errorf("cannot open PNG for saving: %v", err)
	}
	defer newF.Close()
	if err := png.Encode(newF, img); err != nil {
		return fmt.Errorf("cannot save PNG: %v: %v", newName, err)
	}
	if err := newF.Close(); err != nil {
		return fmt.Errorf("cannot finish saving PNG: %v", err)
	}

	goodName := path.Join(dir, testName+".good.png")
	goodF, err := os.Open(goodName)
	if err != nil {
		return fmt.Errorf("cannot open good file: %v", err)
	}
	defer goodF.Close()
	good, err := png.Decode(goodF)
	if err != nil {
		return fmt.Errorf("cannot load good PNG: %v: %v", goodName, err)
	}

	if g, e := good.Bounds().Size(), img.Bounds().Size(); !g.Eq(e) {
		return fmt.Errorf("size mismatch: %v != %v", g, e)
	}

	off := img.Bounds().Min.Sub(good.Bounds().Min)
	for y := good.Bounds().Min.Y; y < good.Bounds().Max.Y; y++ {
		for x := good.Bounds().Min.X; x < good.Bounds().Max.X; x++ {
			colorGood := good.At(x, y)
			p := image.Point{X: x, Y: y}.Add(off)
			colorGot := img.At(p.X, p.Y)
			if !colorEq(colorGood, colorGot) {
				return fmt.Errorf("pixel difference at %v", p)
			}
		}
	}
	return nil
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
	if err := approve(img); err != nil {
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
	if err := approve(img); err != nil {
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
	if err := approve(img); err != nil {
		t.Fatalf("not approved: %v", err)
	}
}
