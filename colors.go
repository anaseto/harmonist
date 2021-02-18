package main

import (
	"github.com/anaseto/gruid"
)

// Thoses are the colors of the main palette. They are given 16-palette color
// numbers compatible with terminals, though they are then mapped to more
// precise colors depending on options and the driver. Dark colorscheme is
// assumed by default, but it can be changed in configuration.
const (
	ColorBackground          gruid.Color = gruid.ColorDefault // background
	ColorBackgroundSecondary gruid.Color = 1 + 0              // black
	ColorForeground          gruid.Color = gruid.ColorDefault
	ColorForegroundSecondary gruid.Color = 1 + 7  // white
	ColorForegroundEmph      gruid.Color = 1 + 15 // bright white
	ColorYellow              gruid.Color = 1 + 3
	ColorOrange              gruid.Color = 1 + 1 // red
	ColorRed                 gruid.Color = 1 + 9 // bright red
	ColorMagenta             gruid.Color = 1 + 5
	ColorViolet              gruid.Color = 1 + 12 // bright blue
	ColorBlue                gruid.Color = 1 + 4
	ColorCyan                gruid.Color = 1 + 6
	ColorGreen               gruid.Color = 1 + 2
)

var (
	ColorBg,
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

func init() {
	ColorBg = ColorBackground
	ColorBgDark = ColorBackground
	ColorBgLOS = ColorBackgroundSecondary
	ColorFg = ColorForeground
	ColorFgDark = ColorForegroundSecondary
	ColorFgLOS = ColorForegroundEmph
	ColorFgLOSLight = ColorYellow
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

var Only8Colors bool

func ApplyDarkLOS() {
	ColorBgDark = ColorBackground
	ColorBgLOS = ColorBackgroundSecondary
	ColorFgDark = ColorForegroundSecondary
	if Only8Colors && !Tiles {
		ColorFgLOS = ColorGreen
	} else {
		ColorFgLOS = ColorForegroundEmph
	}
}

func ApplyLightLOS() {
	ColorBgDark = ColorBackgroundSecondary
	ColorBgLOS = ColorBackground
	ColorFgDark = ColorForeground
	if Only8Colors && !Tiles {
		ColorFgLOS = ColorGreen
	} else {
		ColorFgLOS = ColorForegroundEmph
	}
}

const (
	AttrInMap gruid.AttrMask = 1 + iota
	AttrReverse
)
