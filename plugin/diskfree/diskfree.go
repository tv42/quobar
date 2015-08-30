// Package diskfree provides a status bar widget that displays disk
// free in the users home directory.
//
// TODO home dir should be just the default
package diskfree

import (
	"errors"
	"image/color"
	"image/draw"
	"os"
	"sync"
	"syscall"
	"time"

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
			Max:   10 * 1024,
			Color: white,
		},
		{
			Max:   5 * 1024,
			Color: green,
		},
		{
			Max:   2 * 1024,
			Color: yellow,
		},
		{
			Max:   1 * 1024,
			Color: red,
		},
	}
)

type DiskFree struct{}

func init() {
	quobar.Register(DiskFree{})
}

// New returns a new instance of the plugin.
func (DiskFree) New(state *quobar.State) (quobar.Drawer, error) {
	path := os.Getenv("HOME")
	if path == "" {
		return nil, errors.New("HOME not set in environment")
	}
	p := &diskFree{
		state: state,
		path:  path,
		chart: sparkline.New(limit, state.Config.Foreground, trafficLights),
	}
	// seed with a baseline so the graph draws high, to begin with
	p.chart.Add(0.0)
	// TODO shutdown mechanism
	go p.update()
	return p, nil
}

type diskFree struct {
	state *quobar.State
	path  string

	mu    sync.Mutex
	chart *sparkline.Sparkline
	err   error
}

func (p *diskFree) update() {
	ticker := time.NewTicker(step)
	var st syscall.Statfs_t
	for {
		err := syscall.Statfs(p.path, &st)

		p.mu.Lock()
		switch err {
		case nil:
			megabytes := float32(st.Bfree) * float32(st.Bsize) / 1024 / 1024
			p.chart.Add(megabytes)
		default:
			p.err = err
		}
		p.mu.Unlock()

		<-ticker.C
	}
}

func (p *diskFree) Draw(dst draw.Image) error {
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
