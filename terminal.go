// +build !js,!tk

package main

func (ui *gameui) ApplyToggleTiles() {
}

func (ui *gameui) PostConfig() {
	if gameConfig.Small {
		UIHeight = 24
		UIWidth = 80
	}
}
