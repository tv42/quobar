// Package truetype draws a textual message.
package truetype

// Truetype rendering is not simple. I honestly have no idea if I'm
// doing things right or wrong. I know I'm doing them less wrong than
// many code snippets I see online. Help is appreciated.
//
// This FAQ pretty much says every other approach than "try if it
// fits, and then try again slightly smaller" is doomed. Well, that's
// just crap, and I don't think we're gonna do that on every update.
// http://www.freetype.org/freetype2/docs/ft2faq.html#other-size

// Freetype resources
//
// Font metrics:
// http://www.freetype.org/freetype2/docs/tutorial/index.html
// and especially
// http://www.freetype.org/freetype2/docs/tutorial/step2.html
//
// Random things:
// https://www.allegro.cc/forums/thread/598093
// http://stackoverflow.com/questions/24030488/freetype-2-character-exact-size-and-exact-position/29263566

// Alternative libraries we could use:
//
// - https://github.com/golang/freetype
//
// Cannot get baseline information (via descender/ascender):
// https://github.com/golang/freetype/issues/15
//
// - http://godoc.org/j4k.co/freetype2
//
// I like the style but it doesn't expose nearly enough information.

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/azul3d/engine/native/freetype"
)

// Font is a Truetype font loaded and ready for use.
type Font struct {
	font *freetype.Font
}

// Open loads a Truetype font file into memory.
func Open(ttfPath string) (*Font, error) {
	buf, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		return nil, err
	}
	ft, err := freetype.Init()
	if err != nil {
		return nil, err
	}
	font, err := ft.Load(buf)
	if err != nil {
		return nil, err
	}
	return &Font{font: font}, nil
}

const (
	debugBaseline    = false
	debugGlyphHeight = false
)

// Text renders a text onto an image.
//
// TODO: support multiple lines
func (f *Font) Text(dst draw.Image, text string, opts ...Option) error {
	o := options{
		fg: color.White,
	}
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return err
		}
	}

	// X size 0 means preserve aspect ratio.
	//
	// Freetype docs say "You should not rely on the resulting glyphs
	// matching, or being constrained, to this pixel size." -- well,
	// we don't. We clip them to that bounding box, and don't crash or
	// corrupt memory. If it looks ugly, it's the font designer's
	// fault; switch to a better font.
	if err := f.font.SetSizePixels(0, dst.Bounds().Dy()); err != nil {
		return err
	}

	// Origin point for the baseline: left edge, lifted from the
	// bottom as much as the font descender says.
	//
	// TODO this results in a slightly ugly layout, with a lot of
	// negative space below the text, and less above; this is the
	// opposite of what one is supposed to do. The only fix, given
	// existing fonts with deep drops for "g" etc, is to just add
	// negative space above the text. Unfortunately, there seems to be
	// no good way to decide how much space to add; the state of the
	// art seems to be "expert human manually adjusts", which sucks.
	//
	// If we were to perform a similar calculation for the baseline,
	// starting from the drop and coming down by how much Ascender
	// tells us to, we'd end up in the same spot; that is of no use.
	pos := image.Point{
		X: dst.Bounds().Min.X,
		// Descender is in dimensionless units, of which there are
		// UnitsPerEm in the width & height of the em box, which is
		// the size we asked to fit in our image.
		//
		// Because of us truncating subpixel alignment to nearest
		// pixel, this might might end up being off by 1, but if we
		// add floor/ceil here, we'll just shuffle that lost pixel
		// between top and bottom of the image. Perhaps we should give
		// SetSizePixels a pre-shrunk height?
		Y: dst.Bounds().Max.Y + f.font.Descender*dst.Bounds().Dy()/f.font.UnitsPerEm,
	}

	if debugBaseline {
		draw.Draw(dst,
			image.Rectangle{
				Min: image.Point{X: dst.Bounds().Min.X, Y: pos.Y},
				Max: image.Point{X: dst.Bounds().Max.X, Y: pos.Y + 1},
			},
			image.NewUniform(color.RGBA{G: 0xFF, A: 0xFF}), image.ZP,
			draw.Over)
	}

	var prev rune
	for _, ch := range text {
		glyph, err := f.font.Load(f.font.Index(ch))
		if err != nil {
			return err
		}
		// All glyph measurements are in 26.6 format, and should be
		// pixel-aligned by now by the hinting process. Shift by 6 to
		// get pixels.
		i, err := glyph.Image()
		if err != nil {
			return err
		}
		if prev != 0 {
			x, _, err := f.font.Kerning(prev, ch)
			if err != nil {
				return err
			}
			pos.X += x >> 6
		}
		prev = ch
		to := image.Rectangle{pos, dst.Bounds().Max}
		to.Min.X += glyph.HMetrics.BearingX >> 6
		to.Min.Y -= glyph.HMetrics.BearingY >> 6

		if debugGlyphHeight {
			y := pos.Y - glyph.HMetrics.BearingY>>6
			draw.Draw(dst,
				image.Rectangle{
					Min: image.Point{X: to.Min.X, Y: y},
					Max: image.Point{X: to.Min.X + glyph.HMetrics.Advance>>6, Y: y + 1},
				},
				image.NewUniform(color.RGBA{R: 0xFF, B: 0xFF, A: 0xFF}), image.ZP,
				draw.Over)
		}

		draw.DrawMask(dst, to,
			image.NewUniform(o.fg), image.ZP,
			i, i.Bounds().Min,
			draw.Over)
		pos.X += glyph.HMetrics.Advance >> 6
	}
	return nil
}

// Option is a text rendering option that influences the output.
type Option option

// option is used to hide the implementation details of Option.
type option func(*options) error

type options struct {
	fg color.Color
}

// Foreground sets the text foreground color.
func Foreground(c color.Color) Option {
	return func(o *options) error {
		o.fg = c
		return nil
	}
}
