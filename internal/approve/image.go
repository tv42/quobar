package approve

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
)

func colorEq(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

func Image(img image.Image) error {
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
