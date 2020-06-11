// +build js

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"runtime"
	"unicode/utf8"

	"syscall/js"
)

func main() {
	ui := &gameui{}
	err := ui.Init()
	if err != nil {
		log.Fatalf("nboohu: %v\n", err)
	}
	defer ui.Close()
	GameConfig.Tiles = true
	GameConfig.Version = Version
	LinkColors()
	GameConfig.DarkLOS = true
	ApplyDarkLOS()
	go func() {
		for {
			ui.ReqAnimFrame()
		}
	}()
	for {
		newGame(ui)
	}
}

func newGame(ui *gameui) {
	g := &game{}
	ui.g = g
	load, err := g.LoadConfig()
	if load && err != nil {
		log.Printf("Error loading config: %v\n", err)
		err = g.SaveConfig()
		if err != nil {
			log.Printf("Error resetting config: %v\n", err)
		}
	} else if load {
		CustomKeys = true
	}
	ApplyConfig()
	ui.PostConfig()
	if runtime.GOARCH != "wasm" {
		ui.DrawWelcome()
	} else {
		again := ui.HandleStartMenu()
		if again {
			return
		}
	}
	load, err = g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Printf("Error loading saved game… starting new game. (%v)", err)
	} else {
		ui.DrawBufferInit()
	}
	g.ui = ui
	g.EventLoop()
	ui.Clear()
	ui.DrawColoredText("Do you want to collect some more bananas today?\n\n───Click or press any key to play again───", 7, 5, ColorFg)
	ui.DrawText(SaveError, 0, 10)
	ui.Flush()
	ui.PressAnyKey()
}

func (ui *gameui) HandleStartMenu() (again bool) {
	l := ui.DrawWelcomeCommon()
	g := ui.g
	for {
		a := ui.StartMenu(l)
		switch a {
		case StartWatchReplay:
			err := g.LoadReplay()
			if err != nil {
				ui.ColorLine(l+1, ColorRed)
				ui.Flush()
				Sleep(AnimDurShort)
				log.Printf("Load replay: %v", err)
				return true
			}
			small := GameConfig.Small
			GameConfig.Small = true
			ui.ApplyToggleLayoutWithClear(false)
			ui.RestartDrawBuffers()
			ui.Replay()
			if small {
				GameConfig.Small = false
				ui.ApplyToggleLayoutWithClear(false)
			}
			return true
		default:
			return false
		}
	}
}

var SaveError string

type gameui struct {
	g         *game
	cursor    position
	display   js.Value
	cache     map[UICell]js.Value
	ctx       js.Value
	width     int
	height    int
	mousepos  position
	menuHover menu
	itemHover int
}

func (ui *gameui) InitElements() error {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		e.Call("preventDefault")
		return nil
	}), false)
	canvas.Call("setAttribute", "tabindex", "1")
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("imageSmoothingEnabled", false)
	ui.width = 16
	ui.height = 24
	canvas.Set("height", 24*UIHeight)
	canvas.Set("width", 16*UIWidth)
	ui.cache = make(map[UICell]js.Value)
	return nil
}

func (ui *gameui) Draw(cell UICell, x, y int) {
	var canvas js.Value
	if cv, ok := ui.cache[cell]; ok {
		canvas = cv
	} else {
		canvas = js.Global().Get("document").Call("createElement", "canvas")
		canvas.Set("width", 16)
		canvas.Set("height", 24)
		ctx := canvas.Call("getContext", "2d")
		ctx.Set("imageSmoothingEnabled", false)
		buf := getImage(cell).Pix
		ua := js.Global().Get("Uint8Array").New(js.ValueOf(len(buf)))
		js.CopyBytesToJS(ua, buf)
		ca := js.Global().Get("Uint8ClampedArray").New(ua)
		imgdata := js.Global().Get("ImageData").New(ca, 16, 24)
		ctx.Call("putImageData", imgdata, 0, 0)
		ui.cache[cell] = canvas
	}
	ui.ctx.Call("drawImage", canvas, x*ui.width, ui.height*y)
}

func (ui *gameui) GetMousePos(evt js.Value) (int, int) {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	rect := canvas.Call("getBoundingClientRect")
	scaleX := canvas.Get("width").Float() / rect.Get("width").Float()
	scaleY := canvas.Get("height").Float() / rect.Get("height").Float()
	x := (evt.Get("clientX").Float() - rect.Get("left").Float()) * scaleX
	y := (evt.Get("clientY").Float() - rect.Get("top").Float()) * scaleY
	return (int(x) - 1) / ui.width, (int(y) - 1) / ui.height
}

// io compatibility functions

func (g *game) DataDir() (string, error) {
	return "", nil
}

func (g *game) Save() error {
	if runtime.GOARCH != "wasm" {
		return errors.New("Saving games is not available in the web html version.") // TODO remove when it works
	}
	save, err := g.GameSave()
	if err != nil {
		SaveError = err.Error()
		return err
	}
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		SaveError = "localStorage not found"
		return errors.New("localStorage not found")
	}
	s := base64.StdEncoding.EncodeToString(save)
	storage.Call("setItem", "nboohusave", s)
	SaveError = ""
	return nil
}

func (g *game) SaveConfig() error {
	if runtime.GOARCH != "wasm" {
		return nil
	}
	conf, err := GameConfig.ConfigSave()
	if err != nil {
		SaveError = err.Error()
		return err
	}
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		SaveError = "localStorage not found"
		return errors.New("localStorage not found")
	}
	s := base64.StdEncoding.EncodeToString(conf)
	storage.Call("setItem", "nboohuconfig", s)
	SaveError = ""
	return nil
}

func (g *game) SaveReplay() error {
	if runtime.GOARCH != "wasm" {
		return nil
	}
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		SaveError = "localStorage not found"
		return errors.New("localStorage not found")
	}
	data, err := g.EncodeDrawLog()
	if err != nil {
		return err
	}
	s := base64.StdEncoding.EncodeToString(data)
	storage.Call("setItem", "nboohureplay", s)
	SaveError = ""
	return nil
}

func (g *game) RemoveSaveFile() error {
	storage := js.Global().Get("localStorage")
	storage.Call("removeItem", "nboohusave")
	return nil
}

func (g *game) RemoveDataFile(file string) error {
	storage := js.Global().Get("localStorage")
	storage.Call("removeItem", file)
	return nil
}

func (g *game) Load() (bool, error) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return true, errors.New("localStorage not found")
	}
	save := storage.Call("getItem", "nboohusave")
	if save.Type() != js.TypeString || runtime.GOARCH != "wasm" {
		return false, nil
	}
	s, err := base64.StdEncoding.DecodeString(save.String())
	if err != nil {
		return true, err
	}
	lg, err := g.DecodeGameSave(s)
	if err != nil {
		return true, err
	}
	*g = *lg

	// // XXX: gob encoding works badly with gopherjs, it seems, some maps get broken

	return true, nil
}

func (g *game) LoadConfig() (bool, error) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return true, errors.New("localStorage not found")
	}
	conf := storage.Call("getItem", "nboohuconfig")
	if conf.Type() != js.TypeString || runtime.GOARCH != "wasm" {
		return false, nil
	}
	s, err := base64.StdEncoding.DecodeString(conf.String())
	if err != nil {
		return true, err
	}
	c, err := g.DecodeConfigSave(s)
	if err != nil {
		return true, err
	}
	if c.Version != GameConfig.Version {
		return true, errors.New("Version mismatch")
	}
	GameConfig = *c
	return true, nil
}

func (g *game) LoadReplay() error {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return errors.New("localStorage not found")
	}
	save := storage.Call("getItem", "nboohureplay")
	if save.Type() != js.TypeString || runtime.GOARCH != "wasm" {
		return errors.New("invalid storage")
	}
	data, err := base64.StdEncoding.DecodeString(save.String())
	if err != nil {
		return err
	}
	dl, err := g.DecodeDrawLog(data)
	if err != nil {
		return err
	}
	g.DrawLog = dl
	return nil
}

func (g *game) WriteDump() error {
	pre := js.Global().Get("document").Call("getElementById", "dump")
	pre.Set("innerHTML", g.Dump())
	err := g.SaveReplay()
	if err != nil {
		return fmt.Errorf("writing replay: %v", err)
	}
	return nil
}

// End of io compatibility functions

func (ui *gameui) Init() error {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	gamediv := js.Global().Get("document").Call("getElementById", "gamediv")
	js.Global().Get("document").Call(
		"addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			if !e.Get("ctrlKey").Bool() && !e.Get("metaKey").Bool() {
				e.Call("preventDefault")
			} else {
				return nil
			}
			s := e.Get("key").String()
			if s == "F11" {
				screenfull := js.Global().Get("screenfull")
				if screenfull.Get("enabled").Bool() {
					screenfull.Call("toggle", gamediv)
				}
				return nil
			}
			if s == "Unidentified" {
				s = e.Get("code").String()
			}
			if len(InCh) < cap(InCh) {
				InCh <- uiInput{key: s}
			}
			return nil
		}))
	canvas.Call(
		"addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			x, y := ui.GetMousePos(e)
			if len(InCh) < cap(InCh) {
				InCh <- uiInput{mouse: true, mouseX: x, mouseY: y, button: e.Get("button").Int()}
			}
			return nil
		}))
	canvas.Call(
		"addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if CenteredCamera {
				return nil
			}
			e := args[0]
			x, y := ui.GetMousePos(e)
			if x != ui.mousepos.X || y != ui.mousepos.Y {
				ui.mousepos.X = x
				ui.mousepos.Y = y
				if len(InCh) < cap(InCh) {
					InCh <- uiInput{mouse: true, mouseX: x, mouseY: y, button: -1}
				}
			}
			return nil
		}))
	ui.menuHover = -1
	ui.InitElements()
	SolarizedPalette()
	ui.HideCursor()
	settingsActions = append(settingsActions, toggleTiles)
	return nil
}

var InCh chan uiInput
var Interrupt chan bool

func init() {
	InCh = make(chan uiInput, 5)
	Interrupt = make(chan bool)
	Flushdone = make(chan bool)
	ReqFrame = make(chan bool)
}

func (ui *gameui) Close() {
	// nothing to do
}

func (ui *gameui) Flush() {
	ReqFrame <- true
	<-Flushdone
}

func (ui *gameui) ReqAnimFrame() {
	<-ReqFrame
	js.Global().Get("window").Call("requestAnimationFrame",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} { ui.FlushCallback(args[0]); return nil }))
}

func (ui *gameui) ApplyToggleLayoutWithClear(clear bool) {
	GameConfig.Small = !GameConfig.Small
	if GameConfig.Small {
		if clear {
			ui.Clear()
			ui.Flush()
		}
		UIHeight = 24
		UIWidth = 80
	} else {
		UIHeight = 26
		if CenteredCamera {
			UIWidth = 80
		} else {
			UIWidth = 100
		}
	}
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Set("height", 24*UIHeight)
	canvas.Set("width", 16*UIWidth)
	ui.g.DrawBuffer = make([]UICell, UIWidth*UIHeight)
	ui.cache = make(map[UICell]js.Value)
	if clear {
		ui.Clear()
	}
}

func (ui *gameui) ApplyToggleLayout() {
	ui.ApplyToggleLayoutWithClear(true)
}

var Flushdone chan bool
var ReqFrame chan bool

func (ui *gameui) FlushCallback(t js.Value) {
	ui.DrawLogFrame()
	for _, cdraw := range ui.g.DrawLog[len(ui.g.DrawLog)-1].Draws {
		cell := cdraw.Cell
		ui.Draw(cell, cdraw.X, cdraw.Y)
	}
	Flushdone <- true
}

func (ui *gameui) PollEvent() (in uiInput) {
	select {
	case in = <-InCh:
	case in.interrupt = <-Interrupt:
	}
	switch in.key {
	case "Escape", "Space":
		in.key = "\x1b"
	case "Enter", "\r", "\n":
		in.key = "."
	case "ArrowLeft":
		in.key = "4"
	case "ArrowRight":
		in.key = "6"
	case "ArrowUp", "BackSpace":
		in.key = "8"
	case "ArrowDown":
		in.key = "2"
	case "PageUp":
		in.key = "u"
	case "PageDown":
		in.key = "d"
	case "Numpad5", "Delete":
		in.key = "5"
	default:
		if utf8.RuneCountInString(in.key) != 1 {
			in.key = ""
		}
	}
	return in
}
