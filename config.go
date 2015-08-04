package quobar

import "image/color"

// Config contains configuration common to all plugins.
type Config struct {
	// Desired height of the status bar.
	HeightMillimeters float32

	// Truetype font file (*.ttf) to use by default
	FontPath string

	Foreground color.RGBA
	Background color.RGBA
}
