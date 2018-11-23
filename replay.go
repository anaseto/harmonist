package main

import (
	"errors"
	"fmt"
	"os"
)

func Replay(file string) error {
	tui := &termui{}
	g := &game{}
	tui.g = g
	g.ui = tui
	err := g.LoadReplay()
	if err != nil {
		return fmt.Errorf("loading replay: %v", err)
	}
	err = tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()
	tui.PostInit()
	tui.DrawBufferInit()
	tui.Replay()
	tui.Close()
	return nil
}

func (ui *termui) Replay() {
	g := ui.g
	dl := g.DrawLog
	g.DrawLog = nil
	for i := 0; i < len(dl); i++ {
		df := dl[i]
		for _, dr := range df.Draws {
			x, y := ui.GetPos(dr.I)
			ui.SetGenCell(x, y, dr.Cell.R, dr.Cell.Fg, dr.Cell.Bg, dr.Cell.InMap)
		}
		ui.Flush()
		err := ui.HandleReplayKey()
		if err != nil {
			break
		}
	}
}

func (ui *termui) HandleReplayKey() error {
	for {
		e := ui.PollEvent()
		if e.interrupt {
			return errors.New("interrupted")
		}
		if e.key == "Q" {
			return errors.New("quit")
		}
		if e.key != "" || (e.mouse && e.button != -1) {
			return nil
		}
	}
}