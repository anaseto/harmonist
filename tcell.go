// +build !sdl,!js

package main

import (
	"runtime"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/drivers/tcell"
	tc "github.com/gdamore/tcell/v2"
)

const Tiles = false

func (md *model) ApplyToggleTiles() {
	// do nothing
}

func (md *model) updateZoom() {
	// do nothing
}

var driver gruid.Driver
var color8 bool

func initDriver(fullscreen bool) {
	st := styler{}
	dr := tcell.NewDriver(tcell.Config{StyleManager: st})
	//dr.PreventQuit()
	driver = dr
	Terminal = true
	if runtime.GOOS == "windows" {
		color8 = true
	}
}

// styler implements the tcell.StyleManager interface.
type styler struct{}

func (sty styler) GetStyle(cst gruid.Style) tc.Style {
	st := tc.StyleDefault
	if Xterm256Color {
		cst.Fg = map16ColorTo256(cst.Fg, true)
		cst.Bg = map16ColorTo256(cst.Bg, false)
		st = st.Background(tc.ColorValid + tc.Color(cst.Bg)).Foreground(tc.ColorValid + tc.Color(cst.Fg))
	} else {
		if !GameConfig.DarkLOS {
			cst.Fg = map16ColorToLight(cst.Fg)
			cst.Bg = map16ColorToLight(cst.Bg)
		}
		if color8 {
			cst.Fg = map16ColorTo8Color(cst.Fg)
			cst.Bg = map16ColorTo8Color(cst.Bg)
		}
		fg := tc.Color(cst.Fg)
		bg := tc.Color(cst.Bg)
		if cst.Bg == gruid.ColorDefault {
			st = st.Background(tc.ColorDefault)
		} else {
			st = st.Background(tc.ColorValid + bg - 1)
		}
		if cst.Fg == gruid.ColorDefault {
			st = st.Foreground(tc.ColorDefault)
		} else {
			st = st.Foreground(tc.ColorValid + fg - 1)
		}
	}
	if cst.Attrs&AttrReverse != 0 {
		st = st.Reverse(true)
	}
	return st
}

func map16ColorTo8Color(c gruid.Color) gruid.Color {
	if c >= 1+8 {
		c -= 8
	}
	return c
}

func map16ColorToLight(c gruid.Color) gruid.Color {
	switch c {
	case ColorBackgroundSecondary:
		return ColorForegroundEmph
	case ColorForegroundSecondary:
		return ColorBackgroundSecondary
	case ColorForegroundEmph:
		return 1 + 0
	default:
		return c
	}
}

// xterm solarized colors: http://ethanschoonover.com/solarized
const (
	Color256Base03  gruid.Color = 234
	Color256Base02  gruid.Color = 235
	Color256Base01  gruid.Color = 240
	Color256Base00  gruid.Color = 241 // for dark on light background
	Color256Base0   gruid.Color = 244
	Color256Base1   gruid.Color = 245
	Color256Base2   gruid.Color = 254
	Color256Base3   gruid.Color = 230
	Color256Yellow  gruid.Color = 136
	Color256Orange  gruid.Color = 166
	Color256Red     gruid.Color = 160
	Color256Magenta gruid.Color = 125
	Color256Violet  gruid.Color = 61
	Color256Blue    gruid.Color = 33
	Color256Cyan    gruid.Color = 37
	Color256Green   gruid.Color = 64
)

func map16ColorTo256(c gruid.Color, fg bool) gruid.Color {
	switch c {
	case ColorBackground:
		if fg {
			if GameConfig.DarkLOS {
				return Color256Base0
			}
			return Color256Base00
		}
		if GameConfig.DarkLOS {
			return Color256Base03
		}
		return Color256Base3
	case ColorBackgroundSecondary:
		if GameConfig.DarkLOS {
			return Color256Base02
		}
		return Color256Base2
	case ColorForegroundEmph:
		if GameConfig.DarkLOS {
			return Color256Base1
		}
		return Color256Base01
	case ColorForegroundSecondary:
		if GameConfig.DarkLOS {
			return Color256Base01
		}
		return Color256Base1
	case ColorYellow:
		return Color256Yellow
	case ColorOrange:
		return Color256Orange
	case ColorRed:
		return Color256Red
	case ColorMagenta:
		return Color256Magenta
	case ColorViolet:
		return Color256Violet
	case ColorBlue:
		return Color256Blue
	case ColorCyan:
		return Color256Cyan
	case ColorGreen:
		return Color256Green
	default:
		return c
	}
}

func clearCache() {
}
