package quobar

import (
	"fmt"
	"image/draw"
	"reflect"
)

// State contains runtime state visible to the plugins.
type State struct {
	// Global configuration, as loaded from the configuration file.
	Config Config

	// The detected resolution of the display.
	Resolution Resolution
}

// Plugin is implemented by all status bar widget plugins.
//
// TODO layout negotiation, something like flexbox
type Plugin interface {
	// New returns a new instance of this plugin.
	//
	// In addition to methods required by the interface, the returned
	// value must support JSON unmarshaling.
	New(*State) (Drawer, error)
}

// Drawer is a widget that knows how to draw itself.
type Drawer interface {
	Draw(img draw.Image) error
}

type registration struct {
	first bool
	Plugin
}

var plugins = map[string]registration{}

// Register a plugin. The plugin name will be derived from the import
// path and the local type name.
//
// As all Register errors are programmer errors, and it is intended to
// be used at init time, Register panics on errors.
func Register(p Plugin) {
	// indirect because we don't want "*foo" for pointer types
	typ := reflect.Indirect(reflect.ValueOf(p)).Type()
	pkg := typ.PkgPath()
	name := pkg + "#" + typ.Name()
	if _, ok := plugins[name]; ok {
		panic(fmt.Errorf("Plugin registered twice: %q", name))
	}
	plugins[name] = registration{Plugin: p}
	// register the first plugin of the package as a default
	if _, ok := plugins[pkg]; !ok {
		plugins[pkg] = registration{Plugin: p, first: true}
	}
}
