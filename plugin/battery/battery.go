// Package battery provides a status bar widget that displays battery
// status.
//
// TODO provide time estimate
//
// TODO make it graphical
package battery

import (
	"bytes"
	"fmt"
	"image/color"
	"image/draw"

	"github.com/distatus/battery"
	"github.com/tv42/quobar"
	"github.com/tv42/quobar/blend"
	"github.com/tv42/quobar/draw/truetype"
)

func percentKludge(f float64) uint64 {
	return uint64(f * 100 * 1000)
}

// TODO configurability
var (
	yellow = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	criticality = []blend.Threshold{
		{
			// values are percentages times 1000

			Max: percentKludge(0.4),
			// this will get overwritten with configured color
			//
			// TODO that's ugly, clean up when adding configurability
			Color: color.White,
		},
		{
			Max:   percentKludge(0.4),
			Color: yellow,
		},
		{
			Max:   percentKludge(0.3),
			Color: red,
		},
	}
)

// Batteries shows the percentage of remaining energy in all batteries.
type Batteries struct{}

func init() {
	quobar.Register(Batteries{})
}

// New returns a new instance of the plugin.
func (Batteries) New(state *quobar.State) (quobar.Drawer, error) {
	font, err := truetype.Open(state.Config.FontPath)
	if err != nil {
		return nil, err
	}
	p := &batteries{
		state: state,
		font:  font,
	}
	return p, nil
}

type batteries struct {
	state *quobar.State
	font  *truetype.Font
}

func (p *batteries) Draw(dst draw.Image) error {
	// TODO handle >1 battery
	bats, err := battery.GetAll()
	if err != nil {
		return fmt.Errorf("cannot access power supply state: %v", err)
	}
	lowest := 1.0
	buf := bytes.NewBufferString("bat:")
	for i, bat := range bats {
		if bat.Full == 0.0 {
			return fmt.Errorf("battery #%d full energy level is zero", i)
		}
		remaining := bat.Current / bat.Full
		if remaining < lowest {
			lowest = remaining
		}
		fmt.Fprintf(buf, " %.0f%%", remaining*100)
		switch bat.State {
		case battery.Charging:
			fmt.Fprint(buf, " AC")
		}
	}
	msg := buf.String()

	criticality[0].Color = p.state.Config.Foreground
	c := blend.PickColor(criticality, percentKludge(lowest))
	return p.font.Text(dst, msg, truetype.Foreground(c))
}
