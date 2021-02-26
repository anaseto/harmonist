// +build sdl

package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid-sdl"
)

var driver gruid.Driver
var isFullscreen bool

func initDriver(fullscreen bool) {
	isFullscreen = fullscreen
	icon, err := base64pngToRGBA(TileImgs["favicon"])
	if err != nil {
		log.Printf("decoding window icon: %v", err)
	}
	dr := sdl.NewDriver(sdl.Config{
		TileManager: &monochromeTileManager{},
		Fullscreen:  fullscreen,
		WindowTitle: "Harmonist",
		WindowIcon:  icon,
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
