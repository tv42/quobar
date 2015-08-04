// Package clock provides a status bar widget that displays the
// current time.
package clock

import (
	"image/draw"
	"time"

	"github.com/tv42/quobar"
	"github.com/tv42/quobar/draw/truetype"
)

// Clock shows the current time.
type Clock struct{}

func init() {
	quobar.Register(Clock{})
}

// New returns a new instance of the plugin.
func (Clock) New(state *quobar.State) (quobar.Drawer, error) {
	font, err := truetype.Open(state.Config.FontPath)
	if err != nil {
		return nil, err
	}
	p := &clock{
		state: state,
		font:  font,
	}
	return p, nil
}

// Default time format. See http://golang.org/pkg/time/#pkg-constants
// for the syntax.
const DefaultFormat = `Mon Jan 2 15:04`

type clock struct {
	state  *quobar.State
	font   *truetype.Font
	Format string
}

func (p *clock) Draw(dst draw.Image) error {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}
	msg := time.Now().Format(format)
	return p.font.Text(dst, msg, truetype.Foreground(p.state.Config.Foreground))
}
