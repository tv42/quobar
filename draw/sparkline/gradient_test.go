package sparkline_test

import (
	"image/color"
	"testing"

	"github.com/tv42/quobar/draw/sparkline"
)

type colorTest struct {
	a, b       color.Color
	percentage float32
	result     color.Color
}

func TestGradientEdgeCases(t *testing.T) {
	var (
		color1 = color.RGBA{R: 12, G: 34, B: 56, A: 87}
		color2 = color.RGBA{R: 98, G: 76, B: 54, A: 32}
	)
	var tests = []colorTest{
		{
			a:          color1,
			b:          color2,
			percentage: 0.0,
			result:     color1,
		},
		{
			a:          color1,
			b:          color2,
			percentage: -0.1,
			result:     color1,
		},
		{
			a:          color1,
			b:          color2,
			percentage: 1.0,
			result:     color2,
		},
		{
			a:          color1,
			b:          color2,
			percentage: 1.1,
			result:     color2,
		},
		{
			a:          color1,
			b:          color2,
			percentage: 9000,
			result:     color2,
		},
	}
	for idx, test := range tests {
		c := sparkline.Gradient(test.a, test.b, test.percentage)
		if g, e := c, color.RGBA64Model.Convert(test.result); g != e {
			t.Errorf("#%d: %f mix of %v..%v: %v != %v", idx, test.percentage, test.a, test.b, g, e)
		}
	}
}

func TestGradientMixing(t *testing.T) {
	var tests = []colorTest{
		{
			a:          color.RGBA64{R: 1000, G: 2000, B: 3000, A: 4000},
			b:          color.RGBA64{R: 10000, G: 20000, B: 30000, A: 40000},
			percentage: 0.5,
			result:     color.RGBA64{R: 5500, G: 11000, B: 16500, A: 22000},
		},
		{
			a:          color.RGBA64{R: 1000, G: 2000, B: 3000, A: 4000},
			b:          color.RGBA64{R: 10000, G: 20000, B: 30000, A: 40000},
			percentage: 0.9,
			result:     color.RGBA64{R: 9100, G: 18200, B: 27300, A: 36400},
		},
	}
	for idx, test := range tests {
		c := sparkline.Gradient(test.a, test.b, test.percentage)
		if g, e := c, test.result; g != e {
			t.Errorf("#%d: %f mix of %v..%v: %v != %v", idx, test.percentage, test.a, test.b, g, e)
		}
	}
}
