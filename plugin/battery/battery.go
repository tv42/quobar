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
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

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

const sysfsDir = "/sys/class/power_supply"

func readFloat(path string) (float64, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	s := string(buf)
	s = strings.TrimSuffix(s, "\n")
	return strconv.ParseFloat(s, 64)
}

func (p *batteries) Draw(dst draw.Image) error {
	f, err := os.Open(sysfsDir)
	if err != nil {
		return fmt.Errorf("cannot access power supply state: %v", err)
	}
	defer f.Close()
	buf := bytes.NewBufferString("bat:")
	for {
		fis, err := f.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot list batteries: %v", err)
		}
		for _, fi := range fis {
			if !strings.HasPrefix(fi.Name(), "BAT") {
				continue
			}
			energyNow, err := readFloat(path.Join(sysfsDir, fi.Name(), "energy_now"))
			if err != nil {
				return fmt.Errorf("cannot fetch battery %q current energy level: %v", fi.Name(), err)
			}
			energyFull, err := readFloat(path.Join(sysfsDir, fi.Name(), "energy_full"))
			if err != nil {
				return fmt.Errorf("cannot fetch battery %q full energy level: %v", fi.Name(), err)
			}
			if energyFull == 0.0 {
				return fmt.Errorf("battery %q full energy level is zero", fi.Name())
			}
			fmt.Fprintf(buf, " %.0f%%", energyNow/energyFull*100)
		}
	}
	msg := buf.String()
	return p.font.Text(dst, msg, truetype.Foreground(p.state.Config.Foreground))
}
