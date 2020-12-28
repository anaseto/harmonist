// +build sdl

package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/drivers/sdl"
)

var driver gruid.Driver

func init() {
	t, err := getTileDrawer()
	if err != nil {
		log.Fatal(err)
	}
	dr := sdl.NewDriver(sdl.Config{
		TileManager: t,
	})
	//dr.SetScale(2.0, 2.0)
	//dr.PreventQuit()
	driver = dr
}
