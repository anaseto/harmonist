// +build js sdl

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"

	"github.com/anaseto/gruid"
)

const Tiles = true

func init() {
	settingsActions = append(settingsActions, ActionToggleTiles)
}

func (md *model) ApplyToggleTiles() {
	GameConfig.Tiles = !GameConfig.Tiles
	err := SaveConfig()
	if err != nil {
		md.g.Printf("Error saving config changes: %v", err)
	}
	clearCache()
}

func ColorToRGBA(c gruid.Color, fg bool) color.Color {
	cl := color.RGBA{}
	opaque := uint8(255)
	switch c {
	case ColorBackgroundSecondary:
		if GameConfig.DarkLOS {
			cl = color.RGBA{7, 54, 66, opaque}
		} else {
			cl = color.RGBA{238, 232, 213, opaque}
		}
	case ColorRed:
		cl = color.RGBA{220, 50, 47, opaque}
	case ColorGreen:
		cl = color.RGBA{133, 153, 0, opaque}
	case ColorYellow:
		cl = color.RGBA{181, 137, 0, opaque}
	case ColorBlue:
		cl = color.RGBA{38, 139, 210, opaque}
	case ColorMagenta:
		cl = color.RGBA{211, 54, 130, opaque}
	case ColorCyan:
		cl = color.RGBA{42, 161, 152, opaque}
	case ColorOrange:
		cl = color.RGBA{203, 75, 22, opaque}
	case ColorViolet:
		cl = color.RGBA{108, 113, 196, opaque}
	case ColorForegroundEmph:
		if GameConfig.DarkLOS {
			cl = color.RGBA{147, 161, 161, opaque}
		} else {
			cl = color.RGBA{88, 110, 117, opaque}
		}
	case ColorForegroundSecondary:
		if GameConfig.DarkLOS {
			cl = color.RGBA{88, 110, 117, opaque}
		} else {
			cl = color.RGBA{147, 161, 161, opaque}
		}
	default:
		if GameConfig.DarkLOS {
			cl = color.RGBA{0, 43, 54, opaque}
			if fg {
				cl = color.RGBA{131, 148, 150, opaque}
			}
		} else {
			cl = color.RGBA{253, 246, 227, opaque}
			if fg {
				cl = color.RGBA{101, 123, 131, opaque}
			}
		}
	}
	return cl
}

var TileImgs map[string][]byte

var MapNames = map[rune]string{
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'☻':  "dreaming",
	'♪':  "music1",
	'♫':  "footsteps",
	'#':  "wall",
	'@':  "player",
	'§':  "fog",
	'♣':  "simella",
	'+':  "door",
	'.':  "ground",
	'"':  "foliage",
	'•':  "tick",
	'●':  "rock",
	'×':  "times",
	',':  "comma",
	'}':  "rbrace",
	'%':  "percent",
	':':  "colon",
	'\\': "backslash",
	'~':  "tilde",
	'*':  "asterisc",
	'—':  "hbar",
	'/':  "slash",
	'|':  "vbar",
	'∞':  "kill",
	' ':  "space",
	'[':  "lbracket",
	']':  "rbracket",
	')':  "rparen",
	'(':  "lparen",
	'>':  "stairs",
	'!':  "potion",
	';':  "semicolon",
	'∩':  "stone",
	'_':  "stone",
	'&':  "barrel",
	'☼':  "light",
	'π':  "table",
	'Π':  "holedwall",
	'?':  "scroll",
	'Δ':  "portal",
	'Ξ':  "barrier",
	'=':  "amulet",
	'Θ':  "window",
	'≈':  "water",
	'◊':  "chasm",
	'^':  "rubble",
	'○':  "nolight",
	'‗':  "queenrock",
}

var LetterNames = map[rune]string{
	'(':  "lparen",
	')':  "rparen",
	'@':  "player",
	'{':  "lbrace",
	'}':  "rbrace",
	'[':  "lbracket",
	']':  "rbracket",
	'♪':  "music1",
	'♫':  "music2",
	'•':  "tick",
	'♣':  "simella",
	' ':  "space",
	'!':  "exclamation",
	'?':  "interrogation",
	',':  "comma",
	':':  "colon",
	';':  "semicolon",
	'\'': "quote",
	'—':  "longhyphen",
	'-':  "hyphen",
	'|':  "pipe",
	'/':  "slash",
	'\\': "backslash",
	'%':  "percent",
	'┐':  "boxne",
	'┤':  "boxe",
	'│':  "vbar",
	'┘':  "boxse",
	'┌':  "boxnw",
	'└':  "boxsw",
	'─':  "hbar",
	'►':  "arrow",
	'×':  "times",
	'.':  "dot",
	'#':  "hash",
	'"':  "quotes",
	'+':  "plus",
	'“':  "lquotes",
	'”':  "rquotes",
	'=':  "equal",
	'>':  "gt",
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'§':  "fog",
	'●':  "rock",
	'~':  "tilde",
	'*':  "asterisc",
	'∞':  "kill",
	'☻':  "dreaming",
	'…':  "dots",
	'∩':  "stone",
	'_':  "stone",
	'♥':  "heart",
	'&':  "barrel",
	'☼':  "light",
	'π':  "table",
	'Π':  "holedwall",
	'←':  "larrow",
	'↓':  "darrow",
	'→':  "rarrow",
	'↑':  "uarrow",
	'Δ':  "portal",
	'«':  "ldiag",
	'»':  "rdiag",
	'Ξ':  "barrier",
	'Θ':  "window",
	'≈':  "water",
	'◊':  "chasm",
	'^':  "rubble",
	'○':  "nolight",
	'‗':  "queenrock",
}

type monochromeTileManager struct{}

func (tm *monochromeTileManager) TileSize() gruid.Point {
	return gruid.Point{16, 24}
}

func (tm *monochromeTileManager) GetImage(gc gruid.Cell) image.Image {
	var pngImg []byte
	hastile := false
	if gc.Style.Attrs&AttrInMap != 0 && GameConfig.Tiles {
		pngImg = TileImgs["map-notile"]
		if im, ok := TileImgs["map-"+string(gc.Rune)]; ok {
			pngImg = im
			hastile = true
		} else if im, ok := TileImgs["map-"+MapNames[gc.Rune]]; ok {
			pngImg = im
			hastile = true
		}
	}
	if !hastile {
		pngImg = TileImgs["map-notile"]
		if im, ok := TileImgs["letter-"+string(gc.Rune)]; ok {
			pngImg = im
		} else if im, ok := TileImgs["letter-"+LetterNames[gc.Rune]]; ok {
			pngImg = im
		}
	}
	buf := make([]byte, len(pngImg))
	base64.StdEncoding.Decode(buf, pngImg) // TODO: check error
	br := bytes.NewReader(buf)
	img, err := png.Decode(br)
	if err != nil {
		log.Printf("Rune %s: could not decode png: %v", string(gc.Rune), err)
	}
	rect := img.Bounds()
	rgbaimg := image.NewRGBA(rect)
	draw.Draw(rgbaimg, rect, img, rect.Min, draw.Src)
	bgc := ColorToRGBA(gc.Style.Bg, false)
	fgc := ColorToRGBA(gc.Style.Fg, true)
	if gc.Style.Attrs&AttrReverse != 0 {
		fgc, bgc = bgc, fgc
	}
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c := rgbaimg.At(x, y)
			r, _, _, _ := c.RGBA()
			if r == 0 {
				rgbaimg.Set(x, y, bgc)
			} else {
				rgbaimg.Set(x, y, fgc)
			}
		}
	}
	return rgbaimg
}
