// Package diskfree provides a status bar widget that display disk
// free in the users home directory.
//
// TODO home dir should be just the default
package meminfo

import (
	"image/color"
	"image/draw"
	"sync"
	"time"

	"github.com/guillermo/go.procmeminfo"
	"github.com/tv42/quobar"
	"github.com/tv42/quobar/draw/sparkline"
)

const (
	// TODO make these configurable
	step  = 1 * time.Second
	limit = 60 * 5
)

// TODO configurability
var (
	white  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	yellow = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	trafficLights = []sparkline.Threshold{
		{
			Max:   100,
			Color: white,
		},
		{
			Max:   50,
			Color: green,
		},
		{
			Max:   30,
			Color: yellow,
		},
		{
			Max:   20,
			Color: red,
		},
	}
)

// TODO more than just swap free, and find a way to share
// /proc/meminfo parsing

type Swap struct{}

func init() {
	quobar.Register(Swap{})
}

// New returns a new instance of the plugin.
func (Swap) New(state *quobar.State) (quobar.Drawer, error) {
	p := &swap{
		state: state,
		chart: sparkline.New(limit, state.Config.Foreground, trafficLights),
	}
	// TODO give sparkline a way to set min/max, init to 0..100
	// TODO shutdown mechanism
	go p.update()
	return p, nil
}

type swap struct {
	state *quobar.State

	mu    sync.Mutex
	chart *sparkline.Sparkline
	err   error
}

func (p *swap) update() {
	ticker := time.NewTicker(step)
	for {
		// this library is pretty bad. it has no state, but talks
		// about "Update". yet, i'd rather not write the parser myself
		// right now.
		mem := procmeminfo.MemInfo{}
		err := mem.Update()
		// this works even if the above errored
		total := mem["SwapTotal"]
		free := mem["SwapFree"]
		var percentage float32
		if total > 0 {
			t := float32(total)
			f := float32(free)
			percentage = 100 * f / t
		}

		p.mu.Lock()
		switch err {
		case nil:
			p.chart.Add(percentage)
		default:
			p.err = err
		}
		p.mu.Unlock()

		<-ticker.C
	}
}

func (p *swap) Draw(dst draw.Image) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.err != nil {
		return p.err
	}

	// TODO have a label or icon
	// TODO draw both sparkline and textual number
	p.chart.Draw(dst)
	return nil
}
