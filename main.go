// +build !js

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func main() {
	optVersion := flag.Bool("v", false, "print version number")
	optNoAnim := flag.Bool("n", false, "no animations")
	optReplay := flag.String("r", "", "path to replay file (_ means default location)")
	optLogFile := flag.String("o", "", "log to output file")
	opt16colors := new(bool)
	opt256colors := new(bool)
	optFullscreen := new(bool)
	if Terminal {
		opt16colors = flag.Bool("s", false, "use 16-color simple palette")
		opt256colors = flag.Bool("x", false, "use xterm 256-color palette (solarized approximation)")
	} else {
		optFullscreen = flag.Bool("F", false, "fullscreen")
	}
	flag.Parse()

	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optNoAnim {
		DisableAnimations = true
	}
	if runtime.GOOS != "windows" {
		Xterm256Color = true
	} else {
		Xterm256Color = false
		Only8Colors = true
	}
	if *opt256colors {
		Xterm256Color = true
		Only8Colors = false
	} else if *opt16colors {
		Xterm256Color = false
		Only8Colors = false
	}
	err := initConfig()
	if err != nil {
		log.Print(err)
	}
	applyThemeConf()
	initDriver(*optFullscreen)
	if *optReplay != "" {
		RunReplay(*optReplay)
	} else {
		RunGame(*optLogFile)
	}
}

func RunGame(logfile string) {
	gd := gruid.NewGrid(UIWidth, UIHeight)
	m := &model{gd: gd, g: &game{}}
	var repw io.WriteCloser
	dir, err := DataDir()
	defer func() {
		if repw != nil {
			repw.Close()
		}
		if m.finished && dir != "" {
			RemoveSaveFile()
			_, err := os.Stat(filepath.Join(dir, "replay.part"))
			if err != nil {
				log.Printf("no replay file: %v", err)
				return
			}
			if err := os.Rename(filepath.Join(dir, "replay.part"), filepath.Join(dir, "replay")); err != nil {
				log.Printf("writing replay file: %v", err)
			}
		}
	}()
	if err == nil {
		replay, err := os.OpenFile(filepath.Join(dir, "replay.part"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err == nil {
			repw = replay
		} else {
			log.Printf("writing to replay file: %v", err)
		}
	} else {
		log.Print(err)
	}

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("opening log file: %v", err)
		} else {
			defer f.Close()
			LogGame = true
			log.SetOutput(f)
		}
	}

	if !Tiles && !LogGame {
		log.SetOutput(ioutil.Discard)
	}
	app := gruid.NewApp(gruid.AppConfig{
		Driver:      driver,
		Model:       m,
		FrameWriter: repw,
	})
	err = app.Start(context.Background())
	if !Tiles && !LogGame {
		log.SetOutput(os.Stderr)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func RunReplay(file string) {
	if file == "_" {
		dir, err := DataDir()
		if err == nil {
			file = filepath.Join(dir, "replay")
		} else {
			log.Print(err)
		}
	}
	replay, err := os.Open(file)
	if err != nil {
		log.Fatalf("loading replay file: %v", err)
	}
	defer replay.Close()
	fd, err := gruid.NewFrameDecoder(replay)
	if err != nil {
		log.Printf("frame decoder: %v", err)
	}
	gd := gruid.NewGrid(UIWidth, UIHeight)
	rep := ui.NewReplay(ui.ReplayConfig{
		Grid:         gd,
		FrameDecoder: fd,
	})
	app := gruid.NewApp(gruid.AppConfig{
		Driver: driver,
		Model:  rep,
	})
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func subSig(ctx context.Context, msgs chan<- gruid.Msg) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)
	select {
	case <-ctx.Done():
	case <-sig:
		msgs <- gruid.MsgQuit{}
	}
}
