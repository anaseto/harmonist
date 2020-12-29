// +build !sdl,!js

package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/drivers/tcell"
	tc "github.com/gdamore/tcell/v2"
)

var driver gruid.Driver

func init() {
	st := styler{}
	dr := tcell.NewDriver(tcell.Config{StyleManager: st})
	//dr.PreventQuit()
	driver = dr
}

// styler implements the tcell.StyleManager interface.
type styler struct{}

func (sty styler) GetStyle(st gruid.Style) tc.Style {
	st := tc.StyleDefault
	st = st.Foreground(tc.ColorValid + tc.Color(cst.Fg)).Background(tc.ColorValid + tc.Color(cst.Bg))
	return st
}
