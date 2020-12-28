// +build !js,!tk

package main

func (ui *model) ApplyToggleTiles() {
}

func (ui *model) PostConfig() {
	if GameConfig.Small {
		UIHeight = 24
		UIWidth = 80
	}
}
