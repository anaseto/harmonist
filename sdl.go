// +build sdl

package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/drivers/sdl"
)

var driver gruid.Driver

func init() {
	dr := sdl.NewDriver(sdl.Config{
		TileManager: &monochromeTileManager{},
	})
	//dr.SetScale(1.5, 1.5)
	//dr.PreventQuit()
	driver = dr
}
