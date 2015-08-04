package placeholder

import (
	"fmt"
	"image/draw"

	"github.com/tv42/quobar"
	"github.com/tv42/quobar/draw/truetype"
)

// Placeholder shows the pixel size of the widget.
//
// It is mostly useful for debugging.
type Placeholder struct{}

func init() {
	quobar.Register(Placeholder{})
}

// New returns a new instance of the plugin.
func (Placeholder) New(state *quobar.State) (quobar.Drawer, error) {
	font, err := truetype.Open(state.Config.FontPath)
	if err != nil {
		return nil, err
	}
	p := &placeholder{
		state: state,
		font:  font,
	}
	return p, nil
}

type placeholder struct {
	state *quobar.State
	font  *truetype.Font
}

func (p *placeholder) Draw(dst draw.Image) error {
	msg := fmt.Sprintf("%dx%d", dst.Bounds().Dx(), dst.Bounds().Dy())
	// TODO centered
	return p.font.Text(dst, msg, truetype.Foreground(p.state.Config.Foreground))
}
