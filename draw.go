package main

import (
	"fmt"
	"sort"
	"strings"
	//"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

const (
	UIWidth  = 80
	UIHeight = 24
)

var (
	DisableAnimations bool = false
	Xterm256Color          = false
	Terminal               = false
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
	AttrText gruid.AttrMask = iota
	AttrInMap
	AttrReverse
)

func (md *model) DrawWelcome() {
	tst := gruid.Style{}
	st := gruid.Style{}.WithAttrs(AttrInMap)
	stt := ui.StyledText{}.WithMarkups(map[rune]gruid.Style{
		't': tst.WithFg(ColorGreen), // title
		'W': st.WithFg(ColorViolet), // wall box
		'l': st.WithFg(ColorFgLOS).WithBg(ColorBgLOS),
		'L': st.WithFg(ColorFgLOSLight).WithBg(ColorBgLOS),
		'p': st.WithFg(ColorFgPlayer).WithBg(ColorBgLOS),
		'm': st.WithFg(ColorFgWanderingMonster),
		'P': st.WithFg(ColorFgPlace),
		'd': st.WithFg(ColorFgDark),
		's': st.WithFg(ColorFgMagicPlace),
		'T': st.WithFg(ColorFgTree),
		'b': st.WithFg(ColorFgBananas),
		'z': st.WithFg(ColorFgSleepingMonster),
		'r': st.WithFg(ColorFgWanderingMonster).WithBg(ColorBgLOS),
		'o': st.WithFg(ColorFgObject),
	})
	gd := md.gd
	rg := gd.Range()
	text := fmt.Sprintf("     Harmonist %s\n", Version)
	text += `@t───────────────────────
 @d#@l##@W###############@d### 
@d#.@L..@W#@t  HARMONIST  @W#@d.@b)@zt@d#
@d#.@pb@L.@W###############@d.## 
@d #@L...@r...@l#@d#@oπ@d.@P>@d##.....#  
@d @l#@L..@r.@mg@r..+@d..@mG@d..@P+@d.....#  
@l#@p@@@L.@l#@d≈@m♫@d..##@o☼@d.@o&@d##..@T♣@d.".##
@l#@L.@l#@d#≈≈≈..##@P+@d##..@mh@d."#.@s∩@d#
@l#@L.@l.@d##≈≈≈.........""""##
@t───────────────────────
`
	stt.WithText(text).Draw(gd.Slice(rg.Shift(20, 6, 0, 0)))
}

func (md *model) DrawKeysDescription(title string, actions []string) {
	md.pagerMode = modeHelpKeys
	md.mode = modePager
	if CustomKeys {
		title = fmt.Sprintf(" Default %s ", title)
	} else {
		title = fmt.Sprintf(" %s ", title)
	}
	md.pager.SetBox(&ui.Box{Title: ui.Text(title).WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	lines := []ui.StyledText{}
	for i := 0; i < len(actions)-1; i += 2 {
		stt := ui.StyledText{}
		if actions[i+1] != "" {
			stt = stt.WithTextf(" %-36s %s", actions[i], actions[i+1])
		} else {
			stt = stt.WithTextf(" %s ", actions[i]).WithStyle(gruid.Style{}.WithFg(ColorCyan))
		}
		if i%4 == 2 {
			stt = stt.WithStyle(stt.Style().WithBg(ColorBgLOS))
		}
		lines = append(lines, stt)
	}
	md.pager.SetLines(lines)
}

func (md *model) KeysHelp() {
	md.DrawKeysDescription("Basic Commands", []string{
		"Move/Jump", "arrows or wasd or hjkl or mouse left",
		"Wait a turn", "“.” or 5 or enter or mouse left on @",
		"Interact (Equip/Descend/Rest...)", "e",
		"Evoke/Zap magara", "v or z",
		"Inventory", "i",
		"Examine", "x or mouse hover",
		"Menu", "M",
		"Advanced Commands", "",
		"Save and Quit", "S",
		"View previous messages", "m",
		"Go to nearest stairs", "G",
		"Autoexplore (use with caution)", "o",
		"Write state statistics to file", "#",
		"Quit without saving", "Q",
		"Change settings and key bindings", "=",
	})
}

func (md *model) ExamineHelp() {
	md.DrawKeysDescription("Examine/Travel Commands", []string{
		"Move cursor", "arrows or wasd or hjkl or mouse hover",
		"Go to/select target", "“.” or enter or mouse left",
		"View target description", "v or mouse right",
		"Cycle through monsters", "+",
		"Cycle through stairs", ">",
		"Cycle through objects", "o",
		"Toggle exclude area from auto-travel", "e or mouse middle",
	})
}

const TextWidth = 72

func (md *model) WizardInfo() {
	//g := ui.st
	//ui.Clear()
	//b := &bytes.Buffer{}
	//fmt.Fprintf(b, "Monsters: %d (%d)\n", len(g.Monsters), g.MaxMonsters())
	//fmt.Fprintf(b, "Danger: %d (%d)\n", g.Danger(), g.MaxDanger())
	//ui.DrawText(b.String(), 0, 0)
	//ui.Flush()
	//ui.WaitForContinue(-1)
}

func (md *model) MapWidth() int {
	return DungeonWidth
}

func (md *model) MapHeight() int {
	return DungeonHeight
}

func (md *model) PositionDrawing(pos gruid.Point) (r rune, fgColor, bgColor gruid.Color) {
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
				fgColor = m.Color(g)
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

func (md *model) HPColor() rune {
	g := md.g
	hpColor := 'G'
	switch g.Player.HP + g.Player.HPbonus {
	case 1, 2:
		hpColor = 'C'
	case 3, 4:
		hpColor = 'W'
	}
	return hpColor
}

func (md *model) MPColor() rune {
	g := md.g
	mpColor := 'g'
	switch g.Player.MP {
	case 1, 2:
		mpColor = 'c'
	case 3, 4:
		mpColor = 'w'
	}
	return mpColor
}

func (md *model) updateStatus() {
	g := md.g
	var entries []ui.MenuEntry

	st := gruid.Style{}
	stt := ui.StyledText{}.WithMarkups(map[rune]gruid.Style{
		'G': st.WithFg(ColorFgHPok),
		'g': st.WithFg(ColorFgMPok),
		'W': st.WithFg(ColorFgHPwounded),
		'w': st.WithFg(ColorFgMPpartial),
		'C': st.WithFg(ColorFgHPcritical),
		'c': st.WithFg(ColorFgMPcritical),
		'x': st.WithFg(ColorFgStatusExpire),
		's': st.WithFg(ColorFgStatusGood),
		'o': st.WithFg(ColorFgStatusOther),
		'b': st.WithFg(ColorFgStatusBad),
		'B': st.WithFg(ColorCyan),
		'M': st.WithFg(ColorYellow).WithAttrs(AttrInMap),
	})
	// depth
	var depth string
	if g.Depth == -1 {
		depth = "D: Out! "
	} else {
		depth = fmt.Sprintf(" D:%d ", g.Depth)
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(depth), Disabled: true})

	// turns
	entries = append(entries, ui.MenuEntry{Text: stt.WithTextf("T: %d ", g.Turn), Disabled: true})

	// HP
	nWounds := g.Player.HPMax() - g.Player.HP - g.Player.HPbonus
	if nWounds <= 0 {
		nWounds = 0
	}
	hpColor := md.HPColor()
	hps := "HP:"
	hp := g.Player.HP
	if hp < 0 {
		hp = 0
	}
	if !GameConfig.ShowNumbers {
		hps = fmt.Sprintf("%s@%c%s@B%s@N%s ",
			hps,
			hpColor,
			strings.Repeat("♥", hp),
			strings.Repeat("♥", g.Player.HPbonus),
			strings.Repeat("♥", nWounds),
		)
	} else {
		if g.Player.HPbonus > 0 {
			hps = fmt.Sprintf("@%c%d+%d/%d@N ", hpColor, hp, g.Player.HPbonus, g.Player.HPMax())
		} else {
			hps = fmt.Sprintf("@%c%d/%d@N ", hpColor, hp, g.Player.HPMax())
		}
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(hps), Disabled: true})

	// MP
	MPspent := g.Player.MPMax() - g.Player.MP
	if MPspent <= 0 {
		MPspent = 0
	}
	mpColor := md.MPColor()
	mps := "MP:"
	if !GameConfig.ShowNumbers {
		mps = fmt.Sprintf("%s@%c%s@N%s ",
			mps,
			mpColor,
			strings.Repeat("♥", g.Player.MP),
			strings.Repeat("♥", MPspent),
		)
	} else {
		mps = fmt.Sprintf("@%c%d/%d@N ", mpColor, g.Player.MP, g.Player.MPMax())
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(mps), Disabled: true})

	// bananas
	bananas := fmt.Sprintf("@M)@N:%1d/%1d ", g.Player.Bananas, MaxBananas)
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(bananas), Disabled: true})

	// statuses TODO
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)

	if len(sts) > 0 {
		entries = append(entries, ui.MenuEntry{Text: stt.WithText("| "), Disabled: true})
	}
	for _, st := range sts {
		r := 'o'
		if st.Good() {
			r = 's'
			t := DurationTurn
			if g.Ev != nil && g.Player.Expire[st] >= g.Ev.Rank() && g.Player.Expire[st]-g.Ev.Rank() <= t {
				r = 'x'
			}
		} else if st.Bad() {
			r = 'b'
		}
		var sttext string
		if !st.Flag() {
			sttext = fmt.Sprintf("@%c%s(%d)@N ", r, st.Short(), g.Player.Statuses[st]/DurationStatusStep)
		} else {
			sttext = fmt.Sprintf("@%c%s@N ", r, st.Short())
		}
		entries = append(entries, ui.MenuEntry{Text: stt.WithText(sttext), Disabled: true})
	}

	//altBgEntries(entries)
	md.status.SetEntries(entries)
}

func (md *model) ReadScroll() {
	sc, ok := md.g.Objects.Scrolls[md.g.Player.Pos]
	if !ok {
		md.g.PrintStyled("Error while reading message.", logError)
		return
	}
	md.g.Print("You read the message.")
	md.mode = modeSmallPager
	st := gruid.Style{}
	switch sc {
	case ScrollLore:
		md.smallPager.SetBox(&ui.Box{Title: ui.Text("Lore Message").WithStyle(st.WithFg(ColorCyan))})
		stts := []ui.StyledText{}
		text := ui.Text(sc.Text(md.g)).Format(56)
		for _, s := range strings.Split(text.Text(), "\n") {
			stts = append(stts, ui.Text(s))
		}
		md.smallPager.SetLines(stts)
		if !md.g.Stats.Lore[md.g.Depth] {
			md.g.StoryPrint("Read lore message")
		}
		md.g.Stats.Lore[md.g.Depth] = true
		if len(md.g.Stats.Lore) == 4 {
			AchLoreStudent.Get(md.g)
		}
		if len(md.g.Stats.Lore) == len(md.g.Params.Lore) {
			AchLoremaster.Get(md.g)
		}
	default:
		md.smallPager.SetBox(&ui.Box{Title: ui.Text("Story Message").WithStyle(st.WithFg(ColorCyan))})
		stts := []ui.StyledText{}
		text := ui.Text(sc.Text(md.g)).Format(56)
		for _, s := range strings.Split(text.Text(), "\n") {
			stts = append(stts, ui.Text(s))
		}
		md.smallPager.SetLines(stts)
	}
}
