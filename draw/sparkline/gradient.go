package sparkline

import (
	"image/color"
)

// blend R/G/B values. a and b are in the range [0, 0xFFFF],
// percentage is [0.0, 1.0].
func blend(a, b uint32, percentage float32) uint16 {
	return uint16(float32(a)*(1.0-percentage) + float32(b)*percentage)
}

// Gradient picks a color that's somewhere in between a and b,
// depending on percentage. 0.0 results in a, 1.0 results in b.
func Gradient(a, b color.Color, percentage float32) color.Color {
	switch {
	case percentage < 0.0:
		percentage = 0.0
	case percentage > 1.0:
		percentage = 1.0
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	r := color.RGBA64{
		R: blend(ar, br, percentage),
		G: blend(ag, bg, percentage),
		B: blend(ab, bb, percentage),
		A: blend(aa, ba, percentage),
	}
	return r
}
