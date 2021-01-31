package main

import (
	"github.com/anaseto/gruid"
)

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

	Color16Base03  gruid.Color = gruid.ColorDefault // background
	Color16Base02  gruid.Color = 8
	Color16Base01  gruid.Color = 10
	Color16Base00  gruid.Color = 11
	Color16Base0   gruid.Color = gruid.ColorDefault // foreground
	Color16Base1   gruid.Color = 14
	Color16Base2   gruid.Color = 7
	Color16Base3   gruid.Color = 15
	Color16Yellow  gruid.Color = 3
	Color16Orange  gruid.Color = 9
	Color16Red     gruid.Color = 1
	Color16Magenta gruid.Color = 5
	Color16Violet  gruid.Color = 13
	Color16Blue    gruid.Color = 4
	Color16Cyan    gruid.Color = 6
	Color16Green   gruid.Color = 2
)

// uicolors: http://ethanschoonover.com/solarized
var (
	ColorBase03  gruid.Color = Color16Base03
	ColorBase02  gruid.Color = Color16Base02
	ColorBase01  gruid.Color = Color16Base01
	ColorBase00  gruid.Color = Color16Base00
	ColorBase0   gruid.Color = Color16Base0
	ColorBase1   gruid.Color = Color16Base1
	ColorBase2   gruid.Color = Color16Base2
	ColorBase3   gruid.Color = Color16Base3
	ColorYellow  gruid.Color = Color16Yellow
	ColorOrange  gruid.Color = Color16Orange
	ColorRed     gruid.Color = Color16Red
	ColorMagenta gruid.Color = Color16Magenta
	ColorViolet  gruid.Color = Color16Violet
	ColorBlue    gruid.Color = Color16Blue
	ColorCyan    gruid.Color = Color16Cyan
	ColorGreen   gruid.Color = Color16Green
)

func (md *model) Map256ColorTo16(c gruid.Color) gruid.Color {
	switch c {
	case Color256Base03:
		return Color16Base03
	case Color256Base02:
		return Color16Base02
	case Color256Base01:
		return Color16Base01
	case Color256Base00:
		return Color16Base00
	case Color256Base0:
		return Color16Base0
	case Color256Base1:
		return Color16Base1
	case Color256Base2:
		return Color16Base2
	case Color256Base3:
		return Color16Base3
	case Color256Yellow:
		return Color16Yellow
	case Color256Orange:
		return Color16Orange
	case Color256Red:
		return Color16Red
	case Color256Magenta:
		return Color16Magenta
	case Color256Violet:
		return Color16Violet
	case Color256Blue:
		return Color16Blue
	case Color256Cyan:
		return Color16Cyan
	case Color256Green:
		return Color16Green
	default:
		return c
	}
}

func (md *model) Map16ColorTo256(c gruid.Color) gruid.Color {
	switch c {
	case Color16Base03:
		return Color256Base03
	case Color16Base02:
		return Color256Base02
	case Color16Base01:
		return Color256Base01
	case Color16Base00:
		return Color256Base00
	case Color16Base1:
		return Color256Base1
	case Color16Base2:
		return Color256Base2
	case Color16Base3:
		return Color256Base3
	case Color16Yellow:
		return Color256Yellow
	case Color16Orange:
		return Color256Orange
	case Color16Red:
		return Color256Red
	case Color16Magenta:
		return Color256Magenta
	case Color16Violet:
		return Color256Violet
	case Color16Blue:
		return Color256Blue
	case Color16Cyan:
		return Color256Cyan
	case Color16Green:
		return Color256Green
	default:
		return c
	}
}

var (
	ColorBg,
	ColorBgBorder,
	ColorBgDark,
	ColorBgLOS,
	ColorFg,
	ColorFgObject,
	ColorFgTree,
	ColorFgConfusedMonster,
	ColorFgLignifiedMonster,
	ColorFgParalysedMonster,
	ColorFgDark,
	ColorFgExcluded,
	ColorFgExplosionEnd,
	ColorFgExplosionStart,
	ColorFgExplosionWallEnd,
	ColorFgExplosionWallStart,
	ColorFgHPcritical,
	ColorFgHPok,
	ColorFgHPwounded,
	ColorFgLOS,
	ColorFgLOSLight,
	ColorFgMPcritical,
	ColorFgMPok,
	ColorFgMPpartial,
	ColorFgMagicPlace,
	ColorFgMonster,
	ColorFgPlace,
	ColorFgPlayer,
	ColorFgBananas,
	ColorFgSleepingMonster,
	ColorFgStatusBad,
	ColorFgStatusGood,
	ColorFgStatusExpire,
	ColorFgStatusOther,
	ColorFgWanderingMonster gruid.Color
)

func LinkColors() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgDark = ColorBase03
	ColorBgLOS = ColorBase02
	ColorFg = ColorBase0
	ColorFgDark = ColorBase01
	ColorFgLOS = ColorBase0
	ColorFgLOSLight = ColorBase1
	ColorFgObject = ColorYellow
	ColorFgTree = ColorGreen
	ColorFgConfusedMonster = ColorGreen
	ColorFgLignifiedMonster = ColorYellow
	ColorFgParalysedMonster = ColorCyan
	ColorFgExcluded = ColorRed
	ColorFgExplosionEnd = ColorOrange
	ColorFgExplosionStart = ColorYellow
	ColorFgExplosionWallEnd = ColorMagenta
	ColorFgExplosionWallStart = ColorViolet
	ColorFgHPcritical = ColorRed
	ColorFgHPok = ColorGreen
	ColorFgHPwounded = ColorYellow
	ColorFgMPcritical = ColorMagenta
	ColorFgMPok = ColorBlue
	ColorFgMPpartial = ColorViolet
	ColorFgMagicPlace = ColorCyan
	ColorFgMonster = ColorRed
	ColorFgPlace = ColorMagenta
	ColorFgPlayer = ColorBlue
	ColorFgBananas = ColorYellow
	ColorFgSleepingMonster = ColorViolet
	ColorFgStatusBad = ColorRed
	ColorFgStatusGood = ColorBlue
	ColorFgStatusExpire = ColorViolet
	ColorFgStatusOther = ColorYellow
	ColorFgWanderingMonster = ColorOrange
}

func ApplyDarkLOS() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgDark = ColorBase03
	ColorBgLOS = ColorBase02
	ColorFgDark = ColorBase01
	ColorFg = ColorBase0
	if Only8Colors {
		ColorFgLOS = ColorGreen
		ColorFgLOSLight = ColorYellow
	} else {
		ColorFgLOS = ColorBase0
		//ColorFgLOSLight = ColorBase1
		ColorFgLOSLight = ColorYellow
	}
}

func ApplyLightLOS() {
	if Only8Colors {
		ApplyDarkLOS()
		ColorBgLOS = ColorBase2
		ColorFgLOS = ColorBase00
	} else {
		ColorBg = ColorBase3
		ColorBgBorder = ColorBase2
		ColorBgDark = ColorBase3
		ColorBgLOS = ColorBase2
		ColorFgDark = ColorBase1
		ColorFgLOS = ColorBase00
		ColorFg = ColorBase00
	}
}

func SolarizedPalette() {
	ColorBase03 = Color16Base03
	ColorBase02 = Color16Base02
	ColorBase01 = Color16Base01
	ColorBase00 = Color16Base00
	ColorBase0 = Color16Base0
	ColorBase1 = Color16Base1
	ColorBase2 = Color16Base2
	ColorBase3 = Color16Base3
	ColorYellow = Color16Yellow
	ColorOrange = Color16Orange
	ColorRed = Color16Red
	ColorMagenta = Color16Magenta
	ColorViolet = Color16Violet
	ColorBlue = Color16Blue
	ColorCyan = Color16Cyan
	ColorGreen = Color16Green
}

const (
	Black gruid.Color = iota
	Maroon
	Green
	Olive
	Navy
	Purple
	Teal
	Silver
)

func Map16ColorTo8Color(c gruid.Color) gruid.Color {
	switch c {
	case Color16Base03:
		return Black
	case Color16Base02:
		return Black
	case Color16Base01:
		return Silver
	case Color16Base00:
		return Black
	case Color16Base1:
		return Silver
	case Color16Base2:
		return Silver
	case Color16Base3:
		return Silver
	case Color16Yellow:
		return Olive
	case Color16Orange:
		return Purple
	case Color16Red:
		return Maroon
	case Color16Magenta:
		return Purple
	case Color16Violet:
		return Teal
	case Color16Blue:
		return Navy
	case Color16Cyan:
		return Teal
	case Color16Green:
		return Green
	default:
		return c
	}
}

var Only8Colors bool

func Simple8ColorPalette() {
	Only8Colors = true
}

const (
	AttrInMap gruid.AttrMask = 1 + iota
	AttrReverse
)
