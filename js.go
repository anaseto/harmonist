// +build js

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"syscall/js"

	"github.com/anaseto/gruid"
	jsd "github.com/anaseto/gruid/drivers/js"
	"github.com/anaseto/gruid/ui"
)

var driver gruid.Driver

func initDriver() {
	dr := jsd.NewDriver(jsd.Config{
		TileManager: &monochromeTileManager{},
		AppCanvasId: "gamecanvas",
		AppDivId:    "gamediv",
	})
	driver = dr
}

func clearCache() {
	dr := driver.(*jsd.Driver)
	dr.ClearCache()
}

func (md *model) updateZoom() {
	// do nothing
}

func main() {
	initDriver()
	mainMenu := newMainMenu()
	err := initConfig()
	if err != nil {
		mainMenu.err = err
	}
	applyThemeConf()
	for {
		app := gruid.NewApp(gruid.AppConfig{
			Driver: driver,
			Model:  mainMenu,
		})
		if err := app.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
		switch mainMenu.action {
		case MainPlayGame:
			mainMenu.err = RunGame()
		case MainReplayGame:
			mainMenu.err = RunReplay()
		}
		mainMenu.action = MainMenuDefault
	}
}

type mainMenuAction int

const (
	MainMenuDefault mainMenuAction = iota
	MainPlayGame
	MainReplayGame
)

type mainMenu struct {
	grid   gruid.Grid
	menu   *ui.Menu
	errs   *ui.Label
	err    error
	action mainMenuAction
}

func newMainMenu() *mainMenu {
	gd := gruid.NewGrid(UIWidth, UIHeight)
	md := &mainMenu{grid: gd}
	style := ui.MenuStyle{
		Active: gruid.Style{}.WithFg(ColorYellow),
	}
	md.menu = ui.NewMenu(ui.MenuConfig{
		Grid: gruid.NewGrid(UIWidth/2, 3),
		Entries: []ui.MenuEntry{
			{Text: ui.Text("- (P)lay"), Keys: []gruid.Key{"p", "P"}},
			{Text: ui.Text("- (R)eplay"), Keys: []gruid.Key{"r", "R"}},
		},
		Style: style,
	})
	md.errs = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}.WithFg(ColorRed)))
	return md
}

func menuRange() gruid.Range {
	return gruid.NewRange(20, 18, UIWidth, UIHeight)
}

func (md *mainMenu) Update(msg gruid.Msg) gruid.Effect {
	if md.action != MainMenuDefault {
		return nil
	}
	md.menu.Update(menuRange().RelMsg(msg))
	switch md.menu.Action() {
	case ui.MenuMove:
	case ui.MenuInvoke:
		switch md.menu.Active() {
		case 0:
			md.action = MainPlayGame
			return gruid.End()
		case 1:
			md.action = MainReplayGame
			return gruid.End()
		}
	}
	return nil
}

func (md *mainMenu) Draw() gruid.Grid {
	md.grid.Fill(gruid.Cell{Rune: ' '})
	drawWelcome(md.grid)
	md.grid.Slice(menuRange()).Copy(md.menu.Draw())
	if md.err != nil {
		md.errs.SetText(md.err.Error())
		md.errs.Draw(md.grid.Slice(gruid.NewRange(10, 4, UIWidth, 6)))
	}
	return md.grid
}

const repit = "harmonistreplay"
const replock = "harmonistreplock"

func RunGame() error {
	gd := gruid.NewGrid(UIWidth, UIHeight)
	m := &model{gd: gd, g: &game{}}
	repw := &bytes.Buffer{}
	defer func() {
		if m.finished {
			SetItem(repit, repw.Bytes())
			RemoveItem(replock)
			RemoveSaveFile()
		} else {
			rl, err := GetItem(replock)
			if err != nil {
				log.Printf("replock: %v", err)
			}
			if rl != nil {
				// save ongoing game replay
				SetItem(repit, repw.Bytes())
			}
		}
	}()
	rl, err := GetItem(replock)
	if rl == nil {
		if err != nil {
			log.Printf("replock: %v", err)
		}
		// new game: remove old replay
		RemoveReplay()
	}
	replay, err := GetItem(repit)
	if err != nil {
		log.Printf("get replay: %v", err)
	} else if replay != nil {
		repw.Write(replay)
	}
	app := gruid.NewApp(gruid.AppConfig{
		Driver:      driver,
		Model:       m,
		FrameWriter: repw,
	})
	SetItem(replock, []byte("1"))
	return app.Start(context.Background())
}

func RunReplay() error {
	replay, err := GetItem(repit)
	if err != nil {
		return fmt.Errorf("replay loading: %v", err)
	}
	if replay == nil {
		return errors.New("no replay found")
	}
	repr := &bytes.Buffer{}
	repr.Write(replay)
	fd, err := gruid.NewFrameDecoder(repr)
	if err != nil {
		return fmt.Errorf("frame decoder: %v", err)
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
	return app.Start(context.Background())
}

// io compatibility functions

func DataDir() (string, error) {
	return "", nil
}

// GetItem retrieves a base64 encoded item from localStorage. It returns nil if
// the item does not exist in the storage. It returns an error if localStorage
// is not available, or an item existed but could not be decoded.
func GetItem(item string) ([]byte, error) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return nil, errors.New("localStorage not found")
	}
	save := storage.Call("getItem", item)
	if save.Type() != js.TypeString {
		// this means the item does not exist
		return nil, nil
	}
	s, err := base64.StdEncoding.DecodeString(save.String())
	if err != nil {
		return nil, err
	}
	return s, nil
}

// SetItem sets an item to a given value in the localStorage. The value will be
// base64 encoded.
func SetItem(item string, value []byte) error {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return errors.New("localStorage not found")
	}
	s := base64.StdEncoding.EncodeToString(value)
	storage.Call("setItem", item, s)
	return nil
}

// RemoveItem removes an item from localStorage.
func RemoveItem(item string) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		log.Print("localStorage not found")
	}
	storage.Call("removeItem", item)
}

const harmonistsave = "harmonistsave"
const harmonistconfig = "harmonistconfig"

func (g *game) Save() error {
	save, err := g.GameSave()
	if err != nil {
		return err
	}
	err = SetItem(harmonistsave, save)
	if err != nil {
		return err
	}
	return nil
}

func SaveConfig() error {
	conf, err := GameConfig.ConfigSave()
	if err != nil {
		return err
	}
	err = SetItem(harmonistconfig, conf)
	if err != nil {
		return err
	}
	return nil
}

func RemoveSaveFile() error {
	RemoveItem(harmonistsave)
	return nil
}

func RemoveReplay() {
	RemoveItem(replock)
	RemoveItem(repit)
}

func (g *game) Load() (bool, error) {
	s, err := GetItem(harmonistsave)
	if err != nil || s == nil {
		return false, nil
	}
	lg, err := g.DecodeGameSave(s)
	if err != nil {
		return false, err
	}
	if lg.Version != Version {
		return false, nil
	}
	*g = *lg
	return true, nil
}

func LoadConfig() (bool, error) {
	s, err := GetItem(harmonistconfig)
	if err != nil || s == nil {
		return false, err
	}
	c, err := DecodeConfigSave(s)
	if err != nil {
		return false, err
	}
	if c.Version != GameConfig.Version {
		return false, nil
	}
	GameConfig = *c
	return true, nil
}

func (g *game) WriteDump() error {
	pre := js.Global().Get("document").Call("getElementById", "dump")
	pre.Set("innerHTML", g.Dump())
	return nil
}
