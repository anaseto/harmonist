// +build !js

package main

import (
	"bytes"
	"context"
	"log"

	"github.com/anaseto/gruid"
)

func main() {
	gd := gruid.NewGrid(80, 24)
	m := &model{gd: gd, g: &game{}}
	framebuf := &bytes.Buffer{} // for compressed recording

	// define new application
	app := gruid.NewApp(gruid.AppConfig{
		Driver:      driver,
		Model:       m,
		FrameWriter: framebuf,
	})

	// start application
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

//func main() {
//optSolarized := flag.Bool("s", false, "Use true 16-color solarized palette")
//optVersion := flag.Bool("v", false, "print version number")
//optCenteredCamera := flag.Bool("c", false, "centered camera")
//color8 := false
//if runtime.GOOS == "windows" {
//color8 = true
//}
//opt8colors := flag.Bool("o", color8, "use only 8-color palette")
//opt256colors := flag.Bool("x", !color8, "use xterm 256-color palette (solarized approximation)")
//optNoAnim := flag.Bool("n", false, "no animations")
//optReplay := flag.String("r", "", "path to replay file")
//flag.Parse()
//if *optSolarized {
//SolarizedPalette()
//} else if color8 && !*opt256colors || !color8 && *opt8colors {
//SolarizedPalette()
//Simple8ColorPalette()
//}
//if *optVersion {
//fmt.Println(Version)
//os.Exit(0)
//}
//if *optReplay != "" {
//err := Replay(*optReplay)
//if err != nil {
//log.Printf("harmonist: replay: %v\n", err)
//os.Exit(1)
//}
//os.Exit(0)
//}
//if *optCenteredCamera {
//CenteredCamera = true
//}
//if *optNoAnim {
//DisableAnimations = true
//}

//ui := &model{}
//g := &state{}
//ui.st = g
//if CenteredCamera {
//UIWidth = 80
//}
//err := ui.Init()
//if err != nil {
//fmt.Fprintf(os.Stderr, "harmonist: %v\n", err)
//os.Exit(1)
//}
//defer ui.Close()

//LinkColors()
//GameConfig.DarkLOS = true
//GameConfig.Version = Version

//load, err := g.LoadConfig()
//var cfgerrstr string
//var cfgreseterr string
//if load && err != nil {
//cfgerrstr = fmt.Sprintf("Error loading config: %s", err.Error())
//err = g.SaveConfig()
//if err != nil {
//cfgreseterr = fmt.Sprintf("Error resetting config: %s", err.Error())
//}
//} else if load {
//CustomKeys = true
//}
//ApplyConfig()
//ui.PostConfig()
//ui.DrawWelcome()
//load, err = g.Load()
//if !load {
//g.InitLevel()
//} else if err != nil {
//g.InitLevel()
//g.PrintfStyled("Error: %v", logError, err)
//g.PrintStyled("Could not load saved state… starting new state.", logError)
//} else {
//ui.DrawBufferInit()
//}
//if cfgerrstr != "" {
//g.PrintStyled(cfgerrstr, logError)
//}
//if cfgreseterr != "" {
//g.PrintStyled(cfgreseterr, logError)
//}
//g.ui = ui
//g.EventLoop()
//}
