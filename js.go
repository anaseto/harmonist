// +build js

package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/drivers/js"
)

var driver gruid.Driver

func init() {
	dr := js.NewDriver(js.Config{
		TileManager: &monochromeTileManager{},
	})
	//dr.PreventQuit()
	driver = dr
}

func clearCache() {
	dr := driver.(*js.Driver)
	dr.ClearCache()
}
