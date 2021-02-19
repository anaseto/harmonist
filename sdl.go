// +build sdl

package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid-sdl"
)

var driver gruid.Driver
var isFullscreen bool

func initDriver(fullscreen bool) {
	isFullscreen = fullscreen
	dr := sdl.NewDriver(sdl.Config{
		TileManager: &monochromeTileManager{},
		Fullscreen:  fullscreen,
	})
	driver = dr
}

func (md *model) updateZoom() {
	if !isFullscreen {
		dr := driver.(*sdl.Driver)
		dr.SetScale(1+0.25*float32(md.zoomlevel), 1+0.25*float32(md.zoomlevel))
	}
}

func clearCache() {
	dr := driver.(*sdl.Driver)
	dr.ClearCache()
}
