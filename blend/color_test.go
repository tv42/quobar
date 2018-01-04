package blend_test

import (
	"image/color"
	"testing"

	"github.com/tv42/quobar/blend"
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

func TestPickColorKnown(t *testing.T) {
	for idx, test := range []struct {
		input  uint64
		expect color.Color
	}{
		{9000, white},
		{100, red},
		{0, red},
		{1024, white},
		{768, green},
		{512, yellow},
		{256, red},
	} {
		// if it goes through mixing, it's converted to RGBA64
		if g, e := blend.PickColor(trafficLights, test.input), color.RGBA64Model.Convert(test.expect); g != e {
			t.Errorf("#%d: value %d wrong color: %v != %v", idx, test.input, g, e)
		}
	}
}

func TestPickColorEmpty(t *testing.T) {
	if g, e := blend.PickColor(nil, 42), color.Black; g != e {
		t.Errorf("wrong color for empty slice: %v != %v", g, e)
	}
}
