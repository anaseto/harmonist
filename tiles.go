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

//func (c gruid.Color) String() string {
//color := "#002b36"
//switch c {
//case 0:
//color = "#073642"
//case 1:
//color = "#dc322f"
//case 2:
//color = "#859900"
//case 3:
//color = "#b58900"
//case 4:
//color = "#268bd2"
//case 5:
//color = "#d33682"
//case 6:
//color = "#2aa198"
//case 7:
//color = "#eee8d5"
//case 8:
//color = "#002b36"
//case 9:
//color = "#cb4b16"
//case 10:
//color = "#586e75"
//case 11:
//color = "#657b83"
//case 12:
//color = "#839496"
//case 13:
//color = "#6c71c4"
//case 14:
//color = "#93a1a1"
//case 15:
//color = "#fdf6e3"
//}
//return color
//}

func ColorToRGBA(c gruid.Color, fg bool) color.Color {
	cl := color.RGBA{}
	opaque := uint8(255)
	switch c {
	case Color16Base02:
		cl = color.RGBA{7, 54, 66, opaque}
	case Color16Red:
		cl = color.RGBA{220, 50, 47, opaque}
	case Color16Green:
		cl = color.RGBA{133, 153, 0, opaque}
	case Color16Yellow:
		cl = color.RGBA{181, 137, 0, opaque}
	case Color16Blue:
		cl = color.RGBA{38, 139, 210, opaque}
	case Color16Magenta:
		cl = color.RGBA{211, 54, 130, opaque}
	case Color16Cyan:
		cl = color.RGBA{42, 161, 152, opaque}
	case Color16Base2:
		cl = color.RGBA{238, 232, 213, opaque}
	case Color16Base03: // DefaultColor (TODO: improve this code)
		cl = color.RGBA{0, 43, 54, opaque}
		if fg {
			cl = color.RGBA{131, 148, 150, opaque}
		}
	case Color16Orange:
		cl = color.RGBA{203, 75, 22, opaque}
	case Color16Base01:
		cl = color.RGBA{88, 110, 117, opaque}
	case Color16Base00:
		cl = color.RGBA{101, 123, 131, opaque}
	case Color16Violet:
		cl = color.RGBA{108, 113, 196, opaque}
	case Color16Base1:
		cl = color.RGBA{147, 161, 161, opaque}
	case Color16Base3:
		cl = color.RGBA{253, 246, 227, opaque}
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

func (tm *monochromeTileManager) GetImage(gc gruid.Cell) *image.RGBA {
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
