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
	//dr.SetScale(2.0, 2.0)
	//dr.PreventQuit()
	driver = dr
}
