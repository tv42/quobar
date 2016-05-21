// Package battery provides a status bar widget that displays battery
// status.
//
// TODO provide time estimate
//
// TODO make it graphical, make it use color.
package battery

import (
	"bytes"
	"fmt"
	"image/draw"

	"github.com/distatus/battery"
	"github.com/tv42/quobar"
	"github.com/tv42/quobar/draw/truetype"
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
	buf := bytes.NewBufferString("bat:")
	for i, bat := range bats {
		if bat.Full == 0.0 {
			return fmt.Errorf("battery #%d full energy level is zero", i)
		}
		fmt.Fprintf(buf, " %.0f%%", bat.Current/bat.Full*100)
		switch bat.State {
		case battery.Charging:
			fmt.Fprint(buf, " AC")
		}
	}
	msg := buf.String()
	return p.font.Text(dst, msg, truetype.Foreground(p.state.Config.Foreground))
}
