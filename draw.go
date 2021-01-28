package main

import (
	"fmt"
	//"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func (md *model) Draw() gruid.Grid {
	md.gd.Fill(gruid.Cell{Rune: ' '})
	switch md.mode {
	case modeDump:
		md.gd.Copy(md.pager.Draw())
		return md.gd
	case modeWelcome:
		return drawWelcome(md.gd)
	}
	// Draw map in all other cases, as it may be covered only partially by
	// other modes.
	md.drawMap(md.gd.Slice(md.gd.Range().Shift(0, 2, 0, -1)))
	md.log.StyledText = md.DrawLog()
	md.log.Draw(md.gd.Slice(md.gd.Range().Lines(0, 2)))
	if md.mp.ex.pos != InvalidPos {
		md.drawPosInfo()
	}
	switch md.mode {
	case modePager:
		md.gd.Copy(md.pager.Draw())
	case modeSmallPager:
		md.gd.Slice(gruid.NewRange(10, 2, UIWidth, UIHeight-1)).Copy(md.smallPager.Draw())
	case modeMenu:
		switch md.menuMode {
		case modeInventory, modeEquip, modeEvokation:
			md.gd.Copy(md.menu.Draw())
			md.description.Box = &ui.Box{Title: ui.Text("Description")}
			md.description.Draw(md.gd.Slice(md.gd.Range().Columns(UIWidth/2, UIWidth)))
		case modeGameMenu, modeSettings:
			md.gd.Copy(md.menu.Draw())
		case modeKeys, modeKeysChange:
			md.gd.Copy(md.keysMenu.Draw())
		}
	}
	md.gd.Slice(md.gd.Range().Line(UIHeight - 1)).Copy(md.status.Draw())
	return md.gd
}

func drawWelcome(gd gruid.Grid) gruid.Grid {
	tst := gruid.Style{}
	st := gruid.Style{}.WithAttrs(AttrInMap)
	stt := ui.StyledText{}.WithMarkups(map[rune]gruid.Style{
		't': tst.WithFg(ColorGreen), // title
		'W': st.WithFg(ColorViolet), // wall box
		'l': st.WithFg(ColorFgLOS).WithBg(ColorBgLOS),
		'L': st.WithFg(ColorFgLOSLight).WithBg(ColorBgLOS),
		'p': st.WithFg(ColorFgPlayer).WithBg(ColorBgLOS),
		'm': st.WithFg(ColorFgWanderingMonster),
		'M': st.WithFg(ColorFgWanderingMonster).WithBg(ColorBgLOS),
		'P': st.WithFg(ColorFgPlace),
		'd': st.WithFg(ColorFgDark),
		's': st.WithFg(ColorFgMagicPlace),
		'T': st.WithFg(ColorFgTree),
		'b': st.WithFg(ColorFgBananas),
		'z': st.WithFg(ColorFgSleepingMonster),
		'r': st.WithFg(ColorFgWanderingMonster).WithBg(ColorBgLOS),
		'o': st.WithFg(ColorFgObject),
	})
	rg := gd.Range()
	text := fmt.Sprintf("     Harmonist %s\n", Version)
	text += `@t───────────────────────
 @d#@l##@W###############@d### 
@d#.@L..@W#@t  HARMONIST  @W#@d.@b)@zt@d#
@d#.@pb@L.@W###############@d.## 
@d #@L...@r...@l#@d#@oπ@d.@P>@d##.....#  
@d @l#@L..@r.@Mg@r..+@d..@mG@d..@P+@d.....#  
@l#@p@@@L.@l#@d≈@m♫@d..##@o☼@d.@o&@d##..@T♣@d.".##
@l#@L.@l#@d#≈≈≈..##@P+@d##..@mh@d."#.@s∩@d#
@l#@L.@l.@d##≈≈≈.........""""##
@t───────────────────────
`
	stt.WithText(text).Draw(gd.Slice(rg.Shift(20, 6, 0, 0)))
	return gd
}

func (md *model) drawMap(gd gruid.Grid) {
	for i := range md.g.Dungeon.Cells {
		p := idxtopos(i)
		r, fg, bg := md.positionDrawing(p)
		attrs := AttrInMap
		if md.g.Highlight[p] || p == md.mp.ex.pos {
			attrs |= AttrReverse
		}
		gd.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg, Attrs: attrs}})
	}
}

func (md *model) positionDrawing(pos gruid.Point) (r rune, fgColor, bgColor gruid.Color) {
	g := md.g
	m := g.Dungeon
	c := m.Cell(pos)
	fgColor = ColorFg
	bgColor = ColorBg
	if !c.Explored && (!g.Wizard || g.WizardMode == WizardNormal) {
		r = ' '
		bgColor = ColorBgDark
		if g.HasNonWallExploredNeighbor(pos) {
			r = '¤'
			fgColor = ColorFgDark
		}
		if mons, ok := g.LastMonsterKnownAt[pos]; ok && !mons.Seen {
			r = '☻'
			fgColor = ColorFgSleepingMonster
		}
		if g.Noise[pos] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
		} else if g.NoiseIllusion[pos] {
			r = '♪'
			fgColor = ColorFgMagicPlace
		}
		return
	}
	if g.Wizard && g.WizardMode != WizardNormal {
		if !c.Explored && g.HasNonWallExploredNeighbor(pos) && g.WizardMode == WizardSeeAll {
			r = '¤'
			fgColor = ColorFgDark
			bgColor = ColorBgDark
			return
		}
		if c.T == WallCell {
			if len(g.Dungeon.CardinalNonWallNeighbors(pos)) == 0 {
				r = ' '
				return
			}
		}
	}
	if g.Player.Sees(pos) && !(g.Wizard && g.WizardMode == WizardMap) {
		fgColor = ColorFgLOS
		bgColor = ColorBgLOS
	} else {
		fgColor = ColorFgDark
		bgColor = ColorBgDark
	}
	if g.ExclusionsMap[pos] && c.T.IsPlayerPassable() {
		fgColor = ColorFgExcluded
	}
	if trkn, okTrkn := g.TerrainKnowledge[pos]; okTrkn && (!g.Wizard || g.WizardMode == WizardNormal) {
		c.T = trkn
	}
	var fgTerrain gruid.Color
	switch {
	case c.CoversPlayer():
		r, fgTerrain = c.Style(g, pos)
		if pos == g.Player.Pos {
			fgColor = ColorFgPlayer
		} else if fgTerrain != ColorFgLOS {
			fgColor = fgTerrain
		}
		if _, ok := g.MagicalBarriers[pos]; ok {
			fgColor = ColorFgMagicPlace
		}
	case pos == g.Player.Pos && !(g.Wizard && g.WizardMode == WizardMap):
		r = '@'
		fgColor = ColorFgPlayer
	default:
		// TODO: maybe some wrong knowledge issues
		r, fgTerrain = c.Style(g, pos)
		if fgTerrain != ColorFgLOS {
			fgColor = fgTerrain
		}
		if g.MonsterTargLOS != nil {
			if g.MonsterTargLOS[pos] {
				fgColor = ColorFgWanderingMonster
			}
		} else if g.MonsterLOS[pos] {
			fgColor = ColorFgWanderingMonster
		}
		if cld, ok := g.Clouds[pos]; ok && g.Player.Sees(pos) {
			r = '§'
			if cld == CloudFire {
				fgColor = ColorFgWanderingMonster
			} else if cld == CloudNight {
				fgColor = ColorFgSleepingMonster
			}
		}
		if g.Player.Sees(pos) || (g.Wizard && g.WizardMode == WizardSeeAll) {
			m := g.MonsterAt(pos)
			if m.Exists() {
				r = m.Kind.Letter()
				fgColor = m.color(g)
			}
		} else if (!g.Wizard || g.WizardMode == WizardNormal) && g.Noise[pos] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
		} else if g.NoiseIllusion[pos] {
			r = '♪'
			fgColor = ColorFgMagicPlace
		} else if mons, ok := g.LastMonsterKnownAt[pos]; (!g.Wizard || g.WizardMode == WizardNormal) && ok {
			if !mons.Seen {
				r = '☻'
				fgColor = ColorFgWanderingMonster
			} else {
				r = mons.Kind.Letter()
				if mons.LastSeenState == Resting {
					fgColor = ColorFgSleepingMonster
				} else if mons.Kind.Peaceful() {
					fgColor = ColorFgPlayer
				} else {
					fgColor = ColorFgWanderingMonster
				}
			}
		}
		if fgColor == ColorFgLOS && g.Illuminated[idx(pos)] && c.IsIlluminable() {
			fgColor = ColorFgLOSLight
		}
	}
	return
}
