// +build !js

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {
	optSolarized := flag.Bool("s", false, "Use true 16-color solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	optCenteredCamera := flag.Bool("c", false, "centered camera")
	color8 := false
	if runtime.GOOS == "windows" {
		color8 = true
	}
	opt8colors := flag.Bool("o", color8, "use only 8-color palette")
	opt256colors := flag.Bool("x", !color8, "use xterm 256-color palette (solarized approximation)")
	optNoAnim := flag.Bool("n", false, "no animations")
	optReplay := flag.String("r", "", "path to replay file")
	flag.Parse()
	if *optSolarized {
		SolarizedPalette()
	} else if color8 && !*opt256colors || !color8 && *opt8colors {
		SolarizedPalette()
		Simple8ColorPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optReplay != "" {
		err := Replay(*optReplay)
		if err != nil {
			log.Printf("harmonist: replay: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *optCenteredCamera {
		CenteredCamera = true
	}
	if *optNoAnim {
		DisableAnimations = true
	}

	ui := &gameui{}
	g := &game{}
	ui.g = g
	if CenteredCamera {
		UIWidth = 80
	}
	err := ui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "harmonist: %v\n", err)
		os.Exit(1)
	}
	defer ui.Close()

	LinkColors()
	GameConfig.DarkLOS = true
	GameConfig.Version = Version

	load, err := g.LoadConfig()
	var cfgerrstr string
	if load && err != nil {
		cfgerrstr = err.Error()
	} else if load {
		CustomKeys = true
	}
	ApplyConfig()
	ui.PostConfig()
	ui.DrawWelcome()
	load, err = g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.PrintfStyled("Error: %v", logError, err)
		g.PrintStyled("Could not load saved game… starting new game.", logError)
	} else {
		ui.DrawBufferInit()
	}
	if cfgerrstr != "" {
		g.PrintStyled(cfgerrstr, logError)
	}
	g.ui = ui
	g.EventLoop()
}
