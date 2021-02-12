// +build !js,!sdl

package main

import (
	"io/ioutil" // TODO: change when ioutil deprecated
	"log"
)

func (md *model) ApplyToggleTiles() {
}

func init() {
	// discard logged stuff when using the terminal, as the default writer
	// is not appropriate.
	log.SetOutput(ioutil.Discard)
}
