// +build !js,!tk

package main

func (md *model) ApplyToggleTiles() {
}

func (md *model) PostConfig() {
	if GameConfig.Small {
		UIHeight = 24
		UIWidth = 80
	}
}
