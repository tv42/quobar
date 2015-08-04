// Command quobar is an X11 status bar.
package main

// You can copy this file and edit it, to customize quobar. This
// should only be needed if you have custom plugins.

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/tv42/quobar"
)

// Import your custom plugins here.
import (
	_ "github.com/tv42/quobar/plugin"
)

var defaultConfig = quobar.Config{
	HeightMillimeters: 3,
	FontPath:          "/usr/share/fonts/truetype/roboto/Roboto-Regular.ttf",
	Foreground:        color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff},
	Background:        color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s takes no arguments.\n", prog)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	if err := quobar.Main(defaultConfig); err != nil {
		log.Fatal(err)
	}
}
