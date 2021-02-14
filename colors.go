package main

import (
	"github.com/anaseto/gruid"
)

const (
	Color16Background          gruid.Color = gruid.ColorDefault // background
	Color16BackgroundSecondary gruid.Color = 1 + 0
	Color16Foreground          gruid.Color = gruid.ColorDefault
	Color16ForegroundSecondary gruid.Color = 1 + 7
	Color16ForegroundEmph      gruid.Color = 1 + 15
	Color16Yellow              gruid.Color = 1 + 3
	Color16Orange              gruid.Color = 1 + 1 // red
	Color16Red                 gruid.Color = 1 + 9 // bright red
	Color16Magenta             gruid.Color = 1 + 5
	Color16Violet              gruid.Color = 1 + 12 // bright blue
	Color16Blue                gruid.Color = 1 + 4
	Color16Cyan                gruid.Color = 1 + 6
	Color16Green               gruid.Color = 1 + 2
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

func LinkColors() {
	ColorBg = Color16Background
	ColorBgDark = Color16Background
	ColorBgLOS = Color16BackgroundSecondary
	ColorFg = Color16Foreground
	ColorFgDark = Color16ForegroundSecondary
	ColorFgLOS = Color16ForegroundEmph
	ColorFgLOSLight = Color16Yellow
	ColorFgObject = Color16Yellow
	ColorFgTree = Color16Green
	ColorFgConfusedMonster = Color16Green
	ColorFgLignifiedMonster = Color16Yellow
	ColorFgParalysedMonster = Color16Cyan
	ColorFgExcluded = Color16Red
	ColorFgExplosionEnd = Color16Orange
	ColorFgExplosionStart = Color16Yellow
	ColorFgExplosionWallEnd = Color16Magenta
	ColorFgExplosionWallStart = Color16Violet
	ColorFgHPcritical = Color16Red
	ColorFgHPok = Color16Green
	ColorFgHPwounded = Color16Yellow
	ColorFgMPcritical = Color16Magenta
	ColorFgMPok = Color16Blue
	ColorFgMPpartial = Color16Violet
	ColorFgMagicPlace = Color16Cyan
	ColorFgMonster = Color16Red
	ColorFgPlace = Color16Magenta
	ColorFgPlayer = Color16Blue
	ColorFgBananas = Color16Yellow
	ColorFgSleepingMonster = Color16Violet
	ColorFgStatusBad = Color16Red
	ColorFgStatusGood = Color16Blue
	ColorFgStatusExpire = Color16Violet
	ColorFgStatusOther = Color16Yellow
	ColorFgWanderingMonster = Color16Orange
}

func ApplyDarkLOS() {
	ColorBgDark = Color16Background
	ColorBgLOS = Color16BackgroundSecondary
	ColorFgDark = Color16ForegroundSecondary
	if Only8Colors {
		ColorFgLOS = Color16Green
	} else {
		ColorFgLOS = Color16ForegroundEmph
	}
}

func ApplyLightLOS() {
	ColorBgDark = Color16BackgroundSecondary
	ColorBgLOS = Color16Background
	ColorFgDark = Color16Foreground
	if Only8Colors {
		ColorFgLOS = Color16Green
	} else {
		ColorFgLOS = Color16ForegroundEmph
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
