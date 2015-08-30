// Package sparkline draws small charts of timeseries data.
package sparkline

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/golang/freetype/raster"
)

type Sparkline struct {
	items      int
	data       []float32
	fg         color.Color
	thresholds []Threshold
}

func New(items int, fg color.Color, thresholds []Threshold) *Sparkline {
	s := &Sparkline{
		items:      items,
		fg:         fg,
		thresholds: thresholds,
		data:       make([]float32, 0, items),
	}
	return s
}

func (s *Sparkline) Add(n float32) {
	if len(s.data) < s.items {
		s.data = append(s.data, n)
		return
	}
	// could be fancy with a ringbuffer instead of copying
	copy(s.data, s.data[1:])
	s.data[len(s.data)-1] = n
}

func (s *Sparkline) Draw(dst draw.Image) {
	if len(s.data) == 0 {
		return
	}

	min := s.data[0]
	max := s.data[0]
	for _, n := range s.data[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	// try to be helpful when there's limited data
	if max < 1.0 {
		max = 1.0
	}
	if min > max-1.0 {
		min = max - 1.0
		if min < 0.0 {
			min = 0.0
		}
	}

	bounds := dst.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()

	tmp := image.NewRGBA(image.Rectangle{Max: image.Point{X: dx, Y: dy}})
	p := raster.NewRGBAPainter(tmp)

	r := raster.NewRasterizer(dx, dy)
	r.UseNonZeroWinding = true

	var q raster.Path
	q.Start(s.scale(0, s.data[0], min, max, dx, dy))
	for i, n := range s.data[1:] {
		pt := s.scale(i+1, n, min, max, dx, dy)
		q.Add1(pt)
	}
	const strokeWidth = raster.Fix32(5 << 8)
	r.AddStroke(q, strokeWidth, raster.RoundCapper, raster.RoundJoiner)
	p.SetColor(s.fg)
	r.Rasterize(p)

	r.Clear()
	q.Clear()
	headPt := s.scale(len(s.data)-1, s.data[len(s.data)-1], min, max, dx, dy)
	q.Start(headPt)
	// miniscule nudge so something actually is output
	q.Add1(headPt.Add(raster.Point{X: 1, Y: 1}))
	const headWidth = raster.Fix32(8 << 8)
	r.AddStroke(q, headWidth, raster.RoundCapper, raster.RoundJoiner)
	// TODO really decide between uint64 vs float32 vs uint32 etc
	//	value := uint64((s.data[len(s.data)-1] - min) / max * float32(^uint64(0)))
	value := uint64(s.data[len(s.data)-1])
	headColor := PickColor(s.thresholds, value)
	p.SetColor(headColor)
	r.Rasterize(p)
	draw.Draw(dst, bounds, tmp, image.ZP, draw.Over)
}

func (s *Sparkline) scale(idx int, n, min, max float32, dx, dy int) raster.Point {
	// 24.8 format, so shift by 8
	x := float32(idx) / float32(s.items) * float32(dx)
	y := (1.0 - ((n - min) / max)) * float32(dy)
	p := raster.Point{
		X: raster.Fix32(x)<<8 | (raster.Fix32(x*256) & (1<<8 - 1)),
		Y: raster.Fix32(y)<<8 | (raster.Fix32(y*256) & (1<<8 - 1)),
	}
	return p
}
